# AI Agent Instructions for t1me-0x18

This document (`AGENTS.md`) contains standard operating procedures, architectural guidelines, code style rules, and common commands for AI agents operating within this repository. 

## 1. Tech Stack Overview

- **Language:** Python >= 3.12
- **Framework:** FastAPI
- **ORM / Database:** SQLAlchemy 2.0 (Async mode) + asyncpg + PostgreSQL
- **Migrations:** Alembic
- **Package Manager:** uv

Additional project context:

- This repository is for an AI-powered developer time scheduler called `t1me-0x18`.
- Key features: Pomodoro timer, dashboard for recent todos, productivity metrics, AI auto-scheduling with manual overrides and dynamic rescheduling.
- UI plan: a Charm (Golang) based TUI will be developed for local desktop usage; backend remains FastAPI + async SQLAlchemy.
- Database models live in `models.py` and use PostgreSQL-specific types where appropriate (ARRAY for weekdays). When adding new models, prefer explicit __tablename__ names and include relationships where it makes sense.

## 2. Common Commands

All commands should be executed via `uv` or inside the activated virtual environment (`.venv`).

### Running the Application
- **Dev Server:** `uv run uvicorn main:app --reload` (or `uv run fastapi dev main.py`)
- **Prod Server:** `uv run uvicorn main:app --host 0.0.0.0 --port 8000`

### Build / Dependency Management
- **Sync Dependencies:** `uv sync`
- **Add Dependency:** `uv add <package_name>`
- **Add Dev Dependency:** `uv add --dev <package_name>`


### Database & Migrations
- **Migrations** Do not generate or run migrations, let the developer handle them.

---

## 3. Code Style & Architecture Guidelines

### Typing & Type Hinting
- **Strict Typing:** All functions, methods, and variables must be fully typed.
- **Function Signatures:** Keep function signatures in a single line.
- **Python 3.12 Features:** Use modern Python typing features:
  - Use `X | Y` instead of `Union[X, Y]`.
  - Use built-in generics like `list[str]` and `dict[str, int]` instead of `List` and `Dict` from `typing`.
  - Prefer the `type Alias = ...` syntax for type aliases.
- **Return Types:** Always specify return types, e.g., `def do_something() -> str:`. Use `-> None:` for functions with no return.

### Naming Conventions
- **Classes:** PascalCase (e.g., `DatabaseConfig`, `UserCreate`).
- **Functions, Variables, Attributes:** snake_case (e.g., `get_user_by_email`, `total_count`).
- **Constants & Enums:** UPPER_SNAKE_CASE (e.g., `MAX_RETRY_ATTEMPTS`, `JWT_SECRET`).
- **Database Models:** Singular nouns for class names (e.g., `User`, `Post`), plural for `__tablename__` (e.g., `users`, `posts`).
- **Pydantic Schemas:** Use clear suffixes like `Create`, `Update`, `Read`, `Response` to distinguish from SQLAlchemy models (e.g., `UserCreate`, `UserResponse`).

### Formatting & Linting (Ruff)
- **Line Length:** 88 characters (Black standard).
- **Quotes:** Use double quotes `"` for strings. Single quotes `'` are acceptable for internal dictionary keys but standard is double.
- **Docstrings:** Use Google-style docstrings or standard standard triple-quote `"""` descriptions for complex logic. Avoid useless docstrings (e.g. `"""Gets user"""` on `def get_user()`).

### Imports
- **Ordering:** Handled automatically by Ruff (`I` rules). Grouping is:
  1. Standard library imports (e.g., `import os`, `import asyncio`)
  2. Third-party imports (e.g., `from fastapi import FastAPI`)
  3. First-party / Local module imports (e.g., `from database import get_db`)
- **Absolute vs Relative:** Prefer absolute imports starting from the root of the project over deep relative imports (`..`).

### Database (SQLAlchemy 2.0 Async)
- **Async Driver:** Use `asyncpg` for PostgreSQL (e.g., `postgresql+asyncpg://`).
- **Engine & Sessions:** Use `create_async_engine` and `async_sessionmaker`.
- **Query 2.0 Style:** NEVER use legacy `session.query(Model)`. 
  - Instead, use `sqlalchemy.select`, `update`, `delete`, and `insert`.
  - Example: `result = await session.get(User,user_id)`
  - Use `result.scalars().first()` or `result.scalars().all()`.
- **No Blocking Operations:** Ensure all database I/O is awaited. Avoid blocking the main event loop.

### API & FastAPI Best Practices
- **Dependency Injection:** Use FastAPI's `Depends()` heavily, especially for database sessions (`get_db`), current user resolution, and common query parameters.
- **Routing:** Separate routers into logical domains (e.g., `routers/users.py`, `routers/auth.py`) and include them in `main.py` using `app.include_router()`.
- **Status Codes:** Explicitly define `status_code` in `@router.post()` and `@router.delete()` decorators (e.g., `status.HTTP_201_CREATED`, `status.HTTP_204_NO_CONTENT`).

### Error Handling
- **Fail Fast:** Validate inputs early and fail fast.
- **HTTP Exceptions:** Raise FastAPI's `HTTPException` for expected client errors (400, 401, 403, 404).
- **Do Not Mask Errors:** Avoid blanket `except Exception as e:` blocks unless explicitly logging the trace and re-raising. If catching a specific exception, handle it locally or wrap it in a custom exception.
- **Logging:** Use standard Python `logging` or a structured logger (like `structlog` or `loguru`). Log meaningful context (user_id, action, etc.) on `ERROR` and `WARNING`.

## 4. AI Agent Guidelines
- **Formatting of code** Do not format the part of the codebase which is not touched by your change.
- **Assumptions:** Never assume the presence of a library not found in `pyproject.toml` or `uv.lock`. If necessary, propose adding it first.
- **Code Context:** Review the surrounding project architecture before suggesting sweeping refactors. Mimic the existing structure (e.g., domain-driven vs. layered architecture).
