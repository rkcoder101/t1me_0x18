import re
from datetime import datetime, timedelta, timezone
from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from fastapi import HTTPException, status
from sqlalchemy.exc import IntegrityError

import models
import schemas
from crud import get_task_category


def calculate_free_blocks(start_time: datetime, duration: int, obstacles: list[tuple[datetime, datetime]]) -> list[tuple[datetime, int]]:
    """
    Given a start time, desired duration, and a list of occupied time blocks (obstacles),
    returns a list of (start_time, chunk_duration_in_minutes) representing where the task can be scheduled.
    """
    obstacles.sort(key=lambda x: x[0])
    merged_obstacles = []
    for obs in obstacles:
        if not merged_obstacles:
            merged_obstacles.append(obs)
        else:
            last_obs = merged_obstacles[-1]
            if obs[0] <= last_obs[1]:
                merged_obstacles[-1] = (last_obs[0], max(last_obs[1], obs[1]))
            else:
                merged_obstacles.append(obs)

    current_time = start_time
    remaining_duration = duration
    parts = []

    while remaining_duration > 0:
        next_obs = None
        for obs_start, obs_end in merged_obstacles:
            if obs_end > current_time:
                next_obs = (obs_start, obs_end)
                break

        if not next_obs or next_obs[0] >= current_time + timedelta(minutes=remaining_duration):
            parts.append((current_time, remaining_duration))
            remaining_duration = 0
            break
        else:
            obs_start, obs_end = next_obs
            if obs_start > current_time:
                gap_duration = int((obs_start - current_time).total_seconds() / 60)
                if gap_duration > 0:
                    parts.append((current_time, gap_duration))
                    remaining_duration -= gap_duration
            current_time = max(current_time, obs_end)

    return parts


async def _get_hard_routine_obstacles(db: AsyncSession, target_date: datetime.date) -> list[tuple[datetime, datetime]]:
    """Fetches all active Hard Routines for a given date and returns them as datetime obstacle blocks."""
    obstacles = []
    weekday_str = target_date.strftime("%a").lower()
    routines = (await db.execute(select(models.HardRoutine).filter(models.HardRoutine.is_active))).scalars().all()
    for routine in routines:
        if weekday_str in [w.value for w in routine.weekdays]:
            r_start = datetime.combine(target_date, routine.start_time)
            if r_start.tzinfo is None:
                r_start = r_start.replace(tzinfo=timezone.utc)
            r_end = r_start + timedelta(minutes=routine.duration)
            obstacles.append((r_start, r_end))
    return obstacles


async def _get_scheduled_task_obstacles(db: AsyncSession, target_date: datetime.date, before_time: datetime | None = None) -> list[tuple[datetime, datetime]]:
    """Fetches scheduled tasks for a given date. If before_time is set, only fetches tasks scheduled before that time."""
    query = select(models.Task).filter(models.Task.scheduled_date == target_date).filter(models.Task.status == models.Status.scheduled)
    if before_time:
        query = query.filter(models.Task.scheduled_start < before_time)

    tasks = (await db.execute(query)).scalars().all()
    obstacles = []
    for pt in tasks:
        p_start = pt.scheduled_start
        p_end = p_start + timedelta(minutes=pt.estimated_duration)
        obstacles.append((p_start, p_end))
    return obstacles


async def _get_sleep_boundary(db: AsyncSession, target_date: datetime.date) -> datetime:
    """Calculates the sleep boundary datetime for a given date based on DailySchedule or UserProfile defaults."""
    daily_schedule = await db.get(models.DailySchedule, target_date)
    if daily_schedule:
        sleep_time = daily_schedule.sleep_start
    else:
        user_profile = await db.get(models.UserProfile, 1)
        if not user_profile:
            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="User profile not found. Cannot determine sleep bounds.")
        sleep_time = user_profile.default_sleep_start

    sleep_dt = datetime.combine(target_date, sleep_time)
    if sleep_dt.tzinfo is None:
        sleep_dt = sleep_dt.replace(tzinfo=timezone.utc)

    # Handle if sleep time crosses midnight
    if sleep_time.hour < 12:
        sleep_dt += timedelta(days=1)

    return sleep_dt


async def wrap_task(db: AsyncSession, task_create: schemas.TaskCreate) -> list[models.Task]:
    """
    Schedules a new task by wrapping it around existing Hard Routines AND Scheduled Tasks.
    If another task is currently in progress, it pauses it.
    Returns a list of the created Task chunks.
    """
    now = datetime.now(timezone.utc)

    # Pause the currently running task (if any)
    in_progress_tasks = (await db.execute(select(models.Task).filter(models.Task.status == models.Status.in_progress))).scalars().all()

    for running_task in in_progress_tasks:
        if running_task.last_started_at:
            elapsed_minutes = int((now - running_task.last_started_at).total_seconds() / 60)
            running_task.actual_duration = (running_task.actual_duration or 0) + elapsed_minutes
        running_task.last_started_at = None
        running_task.status = models.Status.paused
        db.add(running_task)

    # Gather obstacles (Hard Routines & Scheduled Tasks)
    requested_start = task_create.scheduled_start
    target_date = requested_start.date()

    obstacles = await _get_hard_routine_obstacles(db, target_date)
    obstacles.extend(await _get_scheduled_task_obstacles(db, target_date))

    # 3. Gap-Filling Loop using the helper
    parts = calculate_free_blocks(requested_start, task_create.estimated_duration, obstacles)

    # 4. Create Task Records in DB
    created_tasks = []
    base_task_data = task_create.model_dump()
    parent_task_id = None

    # Prepare category defaults
    if base_task_data.get("category_id"):
        category = await get_task_category(db, base_task_data["category_id"])
        if category:
            if base_task_data.get("flexibility") is None:
                base_task_data["flexibility"] = category.scheduling_flexibility
            if base_task_data.get("energy_required") is None:
                base_task_data["energy_required"] = category.energy_required

    if base_task_data.get("flexibility") is None:
        base_task_data["flexibility"] = models.Flexibility.M
    if base_task_data.get("energy_required") is None:
        base_task_data["energy_required"] = models.Energy.M

    for i, (start_time, duration) in enumerate(parts):
        new_task_data = base_task_data.copy()
        new_task_data["scheduled_start"] = start_time
        new_task_data["estimated_duration"] = duration
        new_task_data["scheduled_date"] = start_time.date()

        if len(parts) > 1:
            new_task_data["title"] = f"{task_create.title} (Part {i + 1})"

        if i > 0 and parent_task_id is not None:
            new_task_data["parent_task_id"] = parent_task_id

        db_task = models.Task(**new_task_data)
        db.add(db_task)
        await db.flush()  # To get the ID generated

        if i == 0 and len(parts) > 1:
            parent_task_id = db_task.id

        created_tasks.append(db_task)

    await db.commit()
    for task in created_tasks:
        await db.refresh(task)

    return created_tasks


async def shift_tasks(db: AsyncSession, shift_from_time: datetime, shift_amount_minutes: int) -> list[models.Task]:
    """
    Shifts all scheduled tasks that begin on or after `shift_from_time` by `shift_amount_minutes`.
    If pushing tasks forward causes them to overlap with Hard Routines or other already-shifted tasks,
    they will wrap (split).
    If the shift causes any task to spill past the user's `sleep_start`, it raises an HTTPException.
    """
    if shift_amount_minutes <= 0:
        return []

    target_date = shift_from_time.date()

    # 1. Fetch user's sleep_start for the day to check for overflow
    sleep_dt = await _get_sleep_boundary(db, target_date)

    # 2. Gather Hard Routines (Immovable Obstacles) and Prior Tasks
    base_obstacles = await _get_hard_routine_obstacles(db, target_date)
    base_obstacles.extend(await _get_scheduled_task_obstacles(db, target_date, before_time=shift_from_time))

    # 3. Gather Scheduled Tasks to Shift
    query = (
        select(models.Task)
        .where(models.Task.scheduled_date == target_date)
        .where(models.Task.status == models.Status.scheduled)
        .where(models.Task.scheduled_start >= shift_from_time)
        .order_by(models.Task.scheduled_start)
    )
    result = await db.execute(query)
    tasks_to_shift = result.scalars().all()
    if not tasks_to_shift:
        return []

    # 4. Simulate the Shift to check for overflow AND build the new allocation
    # We maintain a dynamic obstacle list that grows as we place shifted tasks
    dynamic_obstacles = list(base_obstacles)

    # To hold dicts of how to mutate each task: {"task": db_task, "parts": [(start, duration)]}
    allocations = []

    # Track the last end time of any placed task so subsequent tasks are pushed along
    latest_placed_time = shift_from_time

    for task in tasks_to_shift:
        # The new requested start is either its original start + shift,
        # or immediately after the latest task we just pushed.
        new_requested_start = max(task.scheduled_start + timedelta(minutes=shift_amount_minutes), latest_placed_time)

        parts = calculate_free_blocks(new_requested_start, task.estimated_duration, dynamic_obstacles)

        # Check if the final part exceeds sleep time
        final_part_end = parts[-1][0] + timedelta(minutes=parts[-1][1])
        if final_part_end > sleep_dt:
            # We must notify user
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail=f"Shifting tasks by {shift_amount_minutes} mins causes '{task.title}' to overflow past sleep time ({sleep_dt.strftime('%H:%M')}). Please reduce task durations or shift amount.",
            )

        allocations.append({"task": task, "parts": parts})

        # Add these new parts as obstacles for the NEXT task in the loop
        for p_start, p_dur in parts:
            dynamic_obstacles.append((p_start, p_start + timedelta(minutes=p_dur)))
            latest_placed_time = max(latest_placed_time, p_start + timedelta(minutes=p_dur))

    # 5. Apply the allocations to the DB
    updated_and_created = []

    for alloc in allocations:
        original_task = alloc["task"]
        parts = alloc["parts"]

        # Clean title if it already has (Part X)
        clean_title = re.sub(r"\s*\(Part \d+\)$", "", original_task.title)

        if len(parts) == 1:
            # Simple shift, no split
            original_task.scheduled_start = parts[0][0]
            db.add(original_task)
            updated_and_created.append(original_task)
        else:
            # Task Split!
            # Update original task to be Part 1
            original_task.title = f"{clean_title} (Part 1)"
            original_task.scheduled_start = parts[0][0]
            original_task.estimated_duration = parts[0][1]
            db.add(original_task)
            await db.flush()  # Ensure it's in DB to act as parent
            updated_and_created.append(original_task)

            parent_id = original_task.parent_task_id or original_task.id

            # Create Part 2+
            for i, (p_start, p_dur) in enumerate(parts[1:], start=2):
                new_task_data = {
                    "title": f"{clean_title} (Part {i})",
                    "description": original_task.description,
                    "category_id": original_task.category_id,
                    "flexibility": original_task.flexibility,
                    "energy_required": original_task.energy_required,
                    "scheduled_start": p_start,
                    "estimated_duration": p_dur,
                    "scheduled_date": p_start.date(),
                    "priority": original_task.priority,
                    "status": models.Status.scheduled,
                    "parent_task_id": parent_id,
                }
                new_task = models.Task(**new_task_data)
                db.add(new_task)
                await db.flush()
                updated_and_created.append(new_task)

    await db.commit()
    for t in updated_and_created:
        await db.refresh(t)

    return updated_and_created

async def get_dashboard_today(db: AsyncSession) -> schemas.DashboardResponse:
    now = datetime.now(timezone.utc)
    today = now.date()

    # 1. Fetch tasks scheduled for today
    tasks_query = (
        select(models.Task, models.TaskCategory.name.label("category_name"))
        .outerjoin(models.TaskCategory, models.Task.category_id == models.TaskCategory.id)
        .where(models.Task.scheduled_date == today)
    )
    try:
        result = await db.execute(tasks_query)
    except IntegrityError as e:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=f"{e}")

    all_tasks = result.all()

    timeline_items = []
    stats_done = 0
    stats_remaining = 0
    stats_overdue = 0

    for task_obj, cat_name in all_tasks:
        task = task_obj

        if task.status == models.Status.COMPLETED:
            stats_done += 1
        else:
            stats_remaining += 1
            if task.scheduled_start + timedelta(minutes=task.estimated_duration) < now:
                stats_overdue += 1

        timeline_items.append(
            schemas.TimelineItem(
                id=task.id,
                type=schemas.TimelineItemType.TASK,
                title=task.title,
                start_time=task.scheduled_start,
                duration=task.estimated_duration,
                category_name=cat_name,
                status=task.status,
                priority=task.priority,
            )
        )

    # 2. Fetch active hard routines for today
    today_weekday = models.Weekday(now.strftime("%a").lower())
    routines_query = select(models.HardRoutine).where(models.HardRoutine.is_active.is_(True)).where(models.HardRoutine.weekdays.any(today_weekday))
    result = await db.execute(routines_query)
    routines = result.scalars().all()

    for routine in routines:
        # Create a localized start time for today
        dt = datetime.combine(today, routine.start_time).replace(tzinfo=timezone.utc)
        timeline_items.append(schemas.TimelineItem(id=routine.id, type=schemas.TimelineItemType.ROUTINE, title=routine.name, start_time=dt, duration=routine.duration))

    # 3. Sort by start_time
    timeline_items.sort(key=lambda x: x.start_time)

    # 4. Inject Gaps
    final_timeline = []
    for i, item in enumerate(timeline_items):
        final_timeline.append(item)
        if i < len(timeline_items) - 1:
            next_item = timeline_items[i + 1]
            end_time = item.start_time + timedelta(minutes=item.duration)
            if next_item.start_time > end_time:
                gap_duration = int((next_item.start_time - end_time).total_seconds() / 60)
                final_timeline.append(schemas.TimelineItem(type=schemas.TimelineItemType.GAP, title=f"·· {gap_duration} min free ··", start_time=end_time, duration=gap_duration))

    return schemas.DashboardResponse(date=today, timeline=final_timeline, stats_done=stats_done, stats_remaining=stats_remaining, stats_overdue=stats_overdue)
