from enum import Enum
from sqlalchemy import Boolean, Column, Date, DateTime, ForeignKey, Integer, String, Text, Time, CheckConstraint
from sqlalchemy import Enum as SAEnum
from sqlalchemy.dialects.postgresql import ARRAY
from sqlalchemy.orm import relationship
from database import Base


class Weekday(str, Enum):
    mon = "mon"
    tue = "tue"
    wed = "wed"
    thu = "thu"
    fri = "fri"
    sat = "sat"
    sun = "sun"


class Flexibility(str, Enum):
    L = "L"
    M = "M"
    H = "H"


class Energy(str, Enum):
    L = "L"
    M = "M"
    H = "H"


class Status(str, Enum):
    scheduled = "scheduled"
    in_progress = "in-progress"
    completed = "completed"
    skipped = "skipped"
    cancelled = "cancelled"


class TaskCategory(Base):
    __tablename__ = "task_categories"

    id = Column(Integer, primary_key=True)
    name = Column(String, unique=True, nullable=False)
    scheduling_flexibility = Column(SAEnum(Flexibility, name="flexibility_enum"), default=Flexibility.M)
    energy_required = Column(SAEnum(Energy, name="energy_enum"), default=Energy.M)
    needs_focus_block = Column(Boolean, default=False)

    tasks = relationship("Task", back_populates="category")


class HardRoutine(Base):
    __tablename__ = "hard_routines"

    id = Column(Integer, primary_key=True)
    name = Column(String, unique=True, nullable=False)
    weekdays = Column(ARRAY(SAEnum(Weekday, name="weekday_enum")), nullable=False)
    start_time = Column(Time(timezone=True), nullable=False)
    duration = Column(Integer, nullable=False)
    is_active = Column(Boolean, default=True)


class Task(Base):
    __tablename__ = "tasks"
    __table_args__ = (CheckConstraint("priority >= 1 AND priority <= 5", name="priority_range"),)

    id = Column(Integer, primary_key=True)
    title = Column(String, nullable=False)
    description = Column(Text)
    category_id = Column(Integer, ForeignKey("task_categories.id"), nullable=True)
    flexibility = Column(SAEnum(Flexibility, name="flexibility_enum"), nullable=False, default=Flexibility.M)
    energy_required = Column(SAEnum(Energy, name="energy_enum"), nullable=False, default=Energy.M)
    scheduled_start = Column(DateTime(timezone=True), nullable=False)
    estimated_duration = Column(Integer, nullable=False)
    scheduled_date = Column(Date, nullable=False)
    actual_start = Column(DateTime(timezone=True), nullable=True)
    actual_duration = Column(Integer, nullable=True)
    actual_date = Column(Date, nullable=True)
    priority = Column(Integer, nullable=False, default=3)
    status = Column(SAEnum(Status, name="status_enum"), nullable=False, default=Status.scheduled)

    category = relationship("TaskCategory", back_populates="tasks")
