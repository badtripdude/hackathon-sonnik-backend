from fastapi import Depends
from sqlalchemy.ext.asyncio import AsyncSession

from ai.infra.db.db import get_db
from ai.infra.repo.chat_repo import ChatRepository


async def get_chat_repo(db: AsyncSession = Depends(get_db)) -> ChatRepository:
    return ChatRepository(db=db)