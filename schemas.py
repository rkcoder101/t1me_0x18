from datetime import date, datetime, time, timezone
from pydantic import BaseModel, Field, ConfigDict, field_validator

from models import Energy, Flexibility, Status, Weekday

# Task Category Schemas


class TaskCategoryBase(BaseModel):
    name: str
    scheduling_flexibility: Flexibility = Flexibility.M
    energy_required: Energy = Energy.M
    needs_focus_block: bool = False


class TaskCategoryCreate(TaskCategoryBase):
    pass


class TaskCategoryUpdate(BaseModel):
    name: str | None = None
    scheduling_flexibility: Flexibility | None = None
    energy_required: Energy | None = None
    needs_focus_block: bool | None = None


class TaskCategoryResponse(TaskCategoryBase):
    id: int

    model_config = ConfigDict(from_attributes=True)


# Hard Routine Schemas
class HardRoutineBase(BaseModel):
    name: str
    weekdays: set[Weekday] = Field(max_length=7)
    start_time: time
    duration: int = Field(gt=0)
    is_active: bool = True


class HardRoutineCreate(HardRoutineBase):
    pass


class HardRoutineUpdate(BaseModel):
    name: str | None = None
    weekdays: set[Weekday] | None = Field(None, max_length=7)
    start_time: time | None = None
    duration: int | None = Field(None, gt=0)
    is_active: bool | None = None


class HardRoutineResponse(HardRoutineBase):
    id: int

    model_config = ConfigDict(from_attributes=True)


# Task Schemas
class TaskBase(BaseModel):
    title: str
    description: str | None = None
    category_id: int | None = None
    flexibility: Flexibility | None = None
    energy_required: Energy | None = None
    scheduled_start: datetime
    estimated_duration: int = Field(gt=0)
    priority: int = Field(3, ge=1, le=5)


class TaskCreate(TaskBase):
    @field_validator("scheduled_start")
    @classmethod
    def check_start_time(cls, v: datetime) -> datetime:
        now = datetime.now(timezone.utc)
        if v < now:
            raise ValueError("Scheduled start time cannot be in the past.")
        return v


class TaskUpdate(BaseModel):
    title: str | None = None
    description: str | None = None
    category_id: int | None = None
    flexibility: Flexibility | None = None
    energy_required: Energy | None = None
    scheduled_start: datetime | None = None
    estimated_duration: int | None = None
    priority: int | None = Field(None, ge=1, le=5)
    status: Status | None = None
    actual_start: datetime | None = None
    actual_end: datetime | None = None
    actual_duration: int | None = None

    @field_validator("scheduled_start")
    @classmethod
    def check_start_time(cls, v: datetime | None) -> datetime | None:
        if v is not None:
            now = datetime.now(timezone.utc)
            if v < now:
                raise ValueError("Scheduled start time cannot be in the past.")
        return v


class TaskSegmentResponse(BaseModel):
    id: int
    task_id: int
    start_time: datetime
    end_time: datetime | None = None
    duration: int | None = None

    model_config = ConfigDict(from_attributes=True)


class TaskResponse(TaskBase):
    id: int
    scheduled_date: date
    actual_start: datetime | None = None
    actual_duration: int | None = None
    actual_date: date | None = None
    status: Status
    segments: list[TaskSegmentResponse] = []

    model_config = ConfigDict(from_attributes=True)


# User Profile Schemas
class UserProfileBase(BaseModel):
    default_work_start: time
    default_sleep_start: time
    timezone: str = "UTC"


class UserProfileCreate(UserProfileBase):
    pass


class UserProfileUpdate(BaseModel):
    default_work_start: time | None = None
    default_sleep_start: time | None = None
    timezone: str | None = None


class UserProfileResponse(UserProfileBase):
    id: int

    model_config = ConfigDict(from_attributes=True)


# Daily Schedule Schemas
class DailyScheduleBase(BaseModel):
    date: date
    work_start: time
    sleep_start: time


class DailyScheduleCreate(DailyScheduleBase):
    pass


class DailyScheduleUpdate(BaseModel):
    work_start: time | None = None
    sleep_start: time | None = None


class DailyScheduleResponse(DailyScheduleBase):
    model_config = ConfigDict(from_attributes=True)
