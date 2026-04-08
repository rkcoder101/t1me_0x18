from datetime import datetime, timedelta

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from fastapi import HTTPException, status

import models
import schemas


# Task Category CRUD
async def create_task_category(db: AsyncSession, category: schemas.TaskCategoryCreate) -> models.TaskCategory:
    db_category = models.TaskCategory(**category.model_dump())
    db.add(db_category)
    await db.commit()
    await db.refresh(db_category)
    return db_category


async def get_task_categories(db: AsyncSession, skip: int = 0, limit: int = 100) -> list[models.TaskCategory]:
    result = await db.execute(select(models.TaskCategory).offset(skip).limit(limit))
    return list(result.scalars().all())


async def get_task_category(db: AsyncSession, category_id: int) -> models.TaskCategory | None:
    return await db.get(models.TaskCategory, category_id)


async def update_task_category(db: AsyncSession, category_id: int, category_update: schemas.TaskCategoryUpdate) -> models.TaskCategory | None:
    db_category = await db.get(models.TaskCategory, category_id)
    if not db_category:
        return None
    update_data = category_update.model_dump(exclude_unset=True)
    for key, value in update_data.items():
        setattr(db_category, key, value)
    await db.commit()
    await db.refresh(db_category)
    return db_category


async def delete_task_category(db: AsyncSession, category_id: int) -> bool:
    db_category = await db.get(models.TaskCategory, category_id)
    if not db_category:
        return False
    await db.delete(db_category)
    await db.commit()
    return True


# Hard Routine CRUD


async def create_hard_routine(db: AsyncSession, routine: schemas.HardRoutineCreate) -> models.HardRoutine:
    db_routine = models.HardRoutine(**routine.model_dump())
    db.add(db_routine)
    await db.commit()
    await db.refresh(db_routine)
    return db_routine


async def get_hard_routines(db: AsyncSession, skip: int = 0, limit: int = 100) -> list[models.HardRoutine]:
    result = await db.execute(select(models.HardRoutine).offset(skip).limit(limit))
    return list(result.scalars().all())


async def get_hard_routnine(db: AsyncSession, hard_routine_id: int = 0) -> models.HardRoutine | None:
    return await db.get(models.HardRoutine, hard_routine_id)


# Task CRUD and Scheduling Logic


def _time_overlaps(start1: datetime, duration1: int, start2: datetime, duration2: int) -> bool:
    end1 = start1 + timedelta(minutes=duration1)
    end2 = start2 + timedelta(minutes=duration2)
    return max(start1, start2) < min(end1, end2)


async def check_schedule_conflicts(db: AsyncSession, scheduled_start: datetime, estimated_duration: int, exclude_task_id: int | None = None) -> list[str]:
    warnings = []
    scheduled_date_val = scheduled_start.date()

    query = select(models.Task).filter(models.Task.scheduled_date == scheduled_date_val)
    if exclude_task_id:
        query = query.filter(models.Task.id != exclude_task_id)

    tasks_on_day = (await db.execute(query)).scalars().all()

    for task in tasks_on_day:
        task_end = task.scheduled_start + timedelta(minutes=task.estimated_duration)
        if _time_overlaps(scheduled_start, estimated_duration, task.scheduled_start, task.estimated_duration):
            warnings.append(f"Task overlaps with existing task '{task.title}' ({task.scheduled_start.strftime('%H:%M')} - {task_end.strftime('%H:%M')}).")

    # Check overlaps with Hard Routines
    weekday_str = scheduled_start.strftime("%A").lower()
    routines = (await db.execute(select(models.HardRoutine).filter(models.HardRoutine.is_active == True))).scalars().all()

    for routine in routines:
        if weekday_str in [w.value for w in routine.weekdays]:
            # Convert routine start_time to datetime on the same day
            routine_start = datetime.combine(scheduled_date_val, routine.start_time)
            routine_end = routine_start + timedelta(minutes=routine.duration)
            if _time_overlaps(scheduled_start, estimated_duration, routine_start, routine.duration):
                warnings.append(f"Task overlaps with hard routine '{routine.name}'. You might not be able to do the hard routine at all you will miss it, r u sure?")

    return warnings


async def create_task(db: AsyncSession, task: schemas.TaskCreate) -> models.Task:
    task_data = task.model_dump(exclude_unset=True)

    if task_data.get("category_id"):
        category = await get_task_category(db, task_data["category_id"])
        if category:
            if "flexibility" not in task_data:
                task_data["flexibility"] = category.scheduling_flexibility
            if "energy_required" not in task_data:
                task_data["energy_required"] = category.energy_required

    if "flexibility" not in task_data:
        task_data["flexibility"] = models.Flexibility.M
    if "energy_required" not in task_data:
        task_data["energy_required"] = models.Energy.M

    # Check overlaps
    warnings = await check_schedule_conflicts(db, task_data["scheduled_start"], task_data["estimated_duration"])
    if warnings:
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail={
                "message": "Scheduling conflict detected. Please reschedule or use force=true to bypass.",
                "warnings": warnings,
            },
        )

    # Calculate scheduled_date
    task_data["scheduled_date"] = task_data["scheduled_start"].date()

    db_task = models.Task(**task_data)
    db.add(db_task)
    await db.commit()
    await db.refresh(db_task)
    return db_task


async def get_tasks(db: AsyncSession, skip: int = 0, limit: int = 100) -> list[models.Task]:
    result = await db.execute(select(models.Task).offset(skip).limit(limit))
    return list(result.scalars().all())


async def update_task(db: AsyncSession, task_id: int, task_update: schemas.TaskUpdate) -> models.Task | None:
    db_task = await db.get(models.Task, task_id)

    if not db_task:
        return None

    update_data = task_update.model_dump(exclude_unset=True)

    # Check overlap if schedule changes
    new_start = update_data.get("scheduled_start", db_task.scheduled_start)
    new_duration = update_data.get("estimated_duration", db_task.estimated_duration)

    warnings = await check_schedule_conflicts(db, new_start, new_duration, exclude_task_id=task_id)
    if warnings:
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail={
                "message": "Scheduling conflict detected",
                "warnings": warnings,
            },
        )

    if "scheduled_start" in update_data:
        update_data["scheduled_date"] = new_start.date()

    for key, value in update_data.items():
        setattr(db_task, key, value)

    await db.commit()
    await db.refresh(db_task)
    return db_task


async def delete_task(db: AsyncSession, task_id: int) -> bool:
    db_task = await db.get(models.Task, task_id)
    if not db_task:
        return False
    await db.delete(db_task)
    await db.commit()
    return True
