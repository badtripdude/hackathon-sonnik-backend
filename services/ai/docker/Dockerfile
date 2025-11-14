FROM ghcr.io/astral-sh/uv:python3.11-bookworm AS builder
ARG INSTALL_DEV=false
WORKDIR /app

COPY pyproject.toml uv.lock ./
RUN --mount=type=cache,target=/root/.cache/uv \
    if [ "$INSTALL_DEV" = "true" ]; then \
        uv sync --frozen --extra dev --no-install-project --no-editable ; \
    else \
        uv sync --frozen --no-install-project --no-editable ; \
    fi

COPY . .

FROM python:3.11.9-slim-bookworm

WORKDIR /app
COPY --from=builder /app /app

ENV PATH="/app/.venv/bin:$PATH"

ENV PYTHONPATH=/app/src