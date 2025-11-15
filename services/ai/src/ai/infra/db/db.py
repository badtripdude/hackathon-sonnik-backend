import os

from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
from sqlalchemy.orm import sessionmaker, declarative_base

from ai.infra.config import settings


DATABASE_URL = settings.database_url

engine = create_async_engine(DATABASE_URL, echo=False, future=True)

async_session_maker = sessionmaker(
    bind=engine,
    class_=AsyncSession,
    expire_on_commit=False,
)

Base = declarative_base()


async def get_db() -> AsyncSession:
    async with async_session_maker() as session:
        yield session
