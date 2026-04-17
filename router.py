from fastapi import FastAPI, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession

import crud
import schemas
import scheduling
from database import get_db

app = FastAPI()


# Task Categories
@app.post("/task-categories/", response_model=schemas.TaskCategoryResponse, status_code=status.HTTP_201_CREATED)
async def create_task_category(category: schemas.TaskCategoryCreate, db: AsyncSession = Depends(get_db)):
    return await crud.create_task_category(db=db, category=category)


@app.get("/task-categories/", response_model=list[schemas.TaskCategoryResponse])
async def read_task_categories(skip: int = 0, limit: int = 100, db: AsyncSession = Depends(get_db)):
    return await crud.get_task_categories(db=db, skip=skip, limit=limit)


@app.get("/task-categories/{category_id}", response_model=schemas.TaskCategoryResponse)
async def read_task_category(category_id: int, db: AsyncSession = Depends(get_db)):
    db_category = await crud.get_task_category(db=db, category_id=category_id)
    if db_category is None:
        raise HTTPException(status_code=404, detail="Task Category not found")
    return db_category


@app.patch("/task-categories/{category_id}", response_model=schemas.TaskCategoryResponse)
async def update_task_category(category_id: int, category: schemas.TaskCategoryUpdate, db: AsyncSession = Depends(get_db)):
    db_category = await crud.update_task_category(db=db, category_id=category_id, category_update=category)
    if db_category is None:
        raise HTTPException(status_code=404, detail="Task Category not found")
    return db_category


@app.delete("/task-categories/{category_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_task_category(category_id: int, db: AsyncSession = Depends(get_db)):
    success = await crud.delete_task_category(db=db, category_id=category_id)
    if not success:
        raise HTTPException(status_code=404, detail="Task Category not found")


# Hard Routines
@app.post("/hard-routines/", response_model=schemas.HardRoutineResponse, status_code=status.HTTP_201_CREATED)
async def create_hard_routine(routine: schemas.HardRoutineCreate, db: AsyncSession = Depends(get_db)):
    return await crud.create_hard_routine(db=db, routine=routine)


@app.get("/hard-routines/", response_model=list[schemas.HardRoutineResponse])
async def read_hard_routines(skip: int = 0, limit: int = 100, db: AsyncSession = Depends(get_db)):
    return await crud.get_hard_routines(db=db, skip=skip, limit=limit)


@app.get("/hard-routines/{routine_id}", response_model=schemas.HardRoutineResponse)
async def read_hard_routine(routine_id: int, db: AsyncSession = Depends(get_db)):
    db_routine = await crud.get_hard_routine(db=db, hard_routine_id=routine_id)
    if db_routine is None:
        raise HTTPException(status_code=404, detail="Hard Routine not found")
    return db_routine


@app.patch("/hard-routines/{routine_id}", response_model=schemas.HardRoutineResponse)
async def update_hard_routine(routine_id: int, routine: schemas.HardRoutineUpdate, db: AsyncSession = Depends(get_db)):
    db_routine = await crud.update_hard_routine(db=db, hard_routine_id=routine_id, routine_update=routine)
    if db_routine is None:
        raise HTTPException(status_code=404, detail="Hard Routine not found")
    return db_routine


@app.delete("/hard-routines/{routine_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_hard_routine(routine_id: int, db: AsyncSession = Depends(get_db)):
    success = await crud.delete_hard_routine(db=db, hard_routine_id=routine_id)
    if not success:
        raise HTTPException(status_code=404, detail="Hard Routine not found")


# Tasks
@app.post("/tasks/", response_model=schemas.TaskResponse, status_code=status.HTTP_201_CREATED)
async def create_task(task: schemas.TaskCreate, db: AsyncSession = Depends(get_db)):
    return await crud.create_task(db=db, task=task)


@app.get("/tasks/", response_model=list[schemas.TaskResponse])
async def read_tasks(skip: int = 0, limit: int = 100, db: AsyncSession = Depends(get_db)):
    return await crud.get_tasks(db=db, skip=skip, limit=limit)


@app.get("/tasks/{task_id}", response_model=schemas.TaskResponse)
async def read_task(task_id: int, db: AsyncSession = Depends(get_db)):
    db_task = await crud.get_task(db=db, task_id=task_id)
    if db_task is None:
        raise HTTPException(status_code=404, detail="Task not found")
    return db_task


@app.patch("/tasks/{task_id}", response_model=schemas.TaskResponse)
async def update_task(task_id: int, task: schemas.TaskUpdate, db: AsyncSession = Depends(get_db)):
    db_task = await crud.update_task(db=db, task_id=task_id, task_update=task)
    if db_task is None:
        raise HTTPException(status_code=404, detail="Task not found")
    return db_task


@app.delete("/tasks/{task_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_task(task_id: int, db: AsyncSession = Depends(get_db)):
    success = await crud.delete_task(db=db, task_id=task_id)
    if not success:
        raise HTTPException(status_code=404, detail="Task not found")

# User-Profile endpoints
@app.post("/user/", response_model=schemas.UserProfileCreate)
async def create_user(user: schemas.UserProfileCreate, db: AsyncSession = Depends(get_db)):
    return await crud.create_user(db, user)

# Other endpoints

@app.get("/dashboard/today", response_model=schemas.DashboardResponse)
async def get_dashboard_today(db: AsyncSession = Depends(get_db)):
    """
    Returns unified timeline data, stats, and unscheduled pool for today.
    """
    return await scheduling.get_dashboard_today(db=db)

# Advanced Scheduling Endpoints

@app.post("/tasks/wrap", response_model=list[schemas.TaskResponse], status_code=status.HTTP_201_CREATED)
async def wrap_task_endpoint(task: schemas.TaskCreate, db: AsyncSession = Depends(get_db)):
    """
    Schedules a new task by wrapping it around existing Hard Routines and Scheduled Tasks. Automatically splits the task if it encounters obstacles.
    """
    return await scheduling.wrap_task(db=db, task_create=task)


@app.post("/tasks/shift", response_model=list[schemas.TaskResponse])
async def shift_tasks_endpoint(shift_req: schemas.ShiftTasksRequest, db: AsyncSession = Depends(get_db)):
    """
    Shifts all scheduled tasks forward. Will automatically wrap (split) tasks
    if the shift causes them to collide with Hard Routines or other shifted tasks.
    Raises a 400 error if tasks overflow past the user's sleep time.
    """
    return await scheduling.shift_tasks(db=db, shift_from_time=shift_req.shift_from_time, shift_amount_minutes=shift_req.shift_amount_minutes)
