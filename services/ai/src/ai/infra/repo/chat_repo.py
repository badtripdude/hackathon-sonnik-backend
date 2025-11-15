import uuid
from typing import Sequence

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession

from ai.infra.db.models import Folder, Chat, Message


class ChatRepository:
    def __init__(self, db: AsyncSession):
        self.db = db

    # ---------- FOLDERS ----------

    async def create_folder(
        self,
        user_id: uuid.UUID,
        title: str,
    ) -> Folder:
        folder = Folder(user_id=user_id, title=title)
        self.db.add(folder)
        await self.db.commit()
        await self.db.refresh(folder)
        return folder

    async def list_folders(self, user_id: uuid.UUID) -> Sequence[Folder]:
        stmt = select(Folder).where(Folder.user_id == user_id).order_by(Folder.created_at)
        result = await self.db.execute(stmt)
        return result.scalars().all()

    async def get_folder(self, folder_id: uuid.UUID, user_id: uuid.UUID) -> Folder | None:
        stmt = select(Folder).where(
            Folder.id == folder_id,
            Folder.user_id == user_id,
        )
        result = await self.db.execute(stmt)
        return result.scalar_one_or_none()

    async def delete_folder(
        self,
        folder_id: uuid.UUID,
        user_id: uuid.UUID,
    ) -> bool:
        """
        Удаляет папку пользователя.
        Ожидаем, что в БД стоят каскадные связи на чаты/сообщения
        (ON DELETE CASCADE) либо folder_id nullable.
        """
        folder = await self.get_folder(folder_id=folder_id, user_id=user_id)
        if folder is None:
            return False

        await self.db.delete(folder)
        await self.db.commit()
        return True
    # ---------- CHATS ----------

    async def create_chat(
        self,
        user_id: uuid.UUID,
        title: str,
        folder_id: uuid.UUID | None = None,
    ) -> Chat:
        chat = Chat(
            user_id=user_id,
            title=title,
            folder_id=folder_id,
        )
        self.db.add(chat)
        await self.db.commit()
        await self.db.refresh(chat)
        return chat

    async def list_chats(
        self,
        user_id: uuid.UUID,
        folder_id: uuid.UUID | None = None,
    ) -> Sequence[Chat]:
        stmt = select(Chat).where(Chat.user_id == user_id)

        if folder_id is not None:
            stmt = stmt.where(Chat.folder_id == folder_id)

        stmt = stmt.order_by(Chat.updated_at.desc())
        result = await self.db.execute(stmt)
        return result.scalars().all()

    async def get_chat(self, chat_id: uuid.UUID, user_id: uuid.UUID) -> Chat | None:
        stmt = select(Chat).where(
            Chat.id == chat_id,
            Chat.user_id == user_id,
        )
        result = await self.db.execute(stmt)
        return result.scalar_one_or_none()

    async def delete_chat(
        self,
        chat_id: uuid.UUID,
        user_id: uuid.UUID,
    ) -> bool:
        """
        Удаляет чат пользователя.
        Ожидаем, что сообщения удалятся каскадно.
        """
        chat = await self.get_chat(chat_id=chat_id, user_id=user_id)
        if chat is None:
            return False

        await self.db.delete(chat)
        await self.db.commit()
        return True
    # ---------- MESSAGES ----------

    async def add_message(
        self,
        chat_id: uuid.UUID,
        role: str,
        content: str,
        user_id: uuid.UUID | None = None,
    ) -> Message:
        message = Message(
            chat_id=chat_id,
            user_id=user_id,
            role=role,
            content=content,
        )
        self.db.add(message)

        # обновляем updated_at у чата
        chat = await self.get_chat(chat_id=chat_id, user_id=user_id) if user_id else None
        # если не хочешь лишний запрос — можно сделать update по chat_id

        await self.db.commit()
        await self.db.refresh(message)
        return message

    async def get_last_messages(
        self,
        chat_id: uuid.UUID,
        limit: int = 20,
    ) -> Sequence[Message]:
        stmt = (
            select(Message)
            .where(Message.chat_id == chat_id)
            .order_by(Message.created_at.desc())
            .limit(limit)
        )
        result = await self.db.execute(stmt)
        messages = result.scalars().all()
        # вернем в хронологическом порядке
        return list(reversed(messages))

    async def rename_chat(
        self,
        chat_id: uuid.UUID,
        user_id: uuid.UUID,
        new_title: str,
    ) -> Chat | None:
        chat = await self.get_chat(chat_id=chat_id, user_id=user_id)
        if chat is None:
            return None
        chat.title = new_title
        self.db.add(chat)
        await self.db.commit()
        await self.db.refresh(chat)
        return chat