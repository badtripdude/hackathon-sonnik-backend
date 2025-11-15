import uuid
from datetime import datetime
from typing import List, Optional


from pydantic import BaseModel, Field


class GptMessage(BaseModel):
    role: str = Field(..., description="system|user|assistant")
    content: str


class GptRequest(BaseModel):
    messages: List[GptMessage]
    model: Optional[str] = Field(default=None, description="Перекрыть дефолтную модель GPT")


class GptResponse(BaseModel):
    text: str
    model: str
    usage_input_tokens: Optional[int] = None
    usage_output_tokens: Optional[int] = None
    usage_total_tokens: Optional[int] = None


class TtsRequest(BaseModel):
    text: str
    voice_id: Optional[str] = Field(
        default=None,
        description="ID голоса ElevenLabs. Если не передан — берём дефолт из ENV",
    )
    model_id: Optional[str] = Field(
        default="eleven_multilingual_v2",
        description="Модель ElevenLabs TTS",
    )
    # формат можешь поменять под свои нужды
    output_format: str = Field(
        default="mp3_44100_128",
        description="Формат аудио ElevenLabs (см. доку)",
    )


class TtsResponseMetadata(BaseModel):
    voice_id: str
    model_id: str
    content_type: str
    filename: str


class AsrResponse(BaseModel):
    text: str
    language: Optional[str] = None
    model: str

class InterpretRequest(BaseModel):
    user_id: uuid.UUID | None = None
    chat_id: uuid.UUID | None = None
    message: str


class InterpretResponse(BaseModel):
    chat_id: uuid.UUID
    reply: str


from datetime import datetime
from uuid import UUID
from pydantic import BaseModel


class FolderCreate(BaseModel):
    user_id: UUID
    title: str


class FolderRename(BaseModel):
    user_id: UUID
    title: str


class FolderOut(BaseModel):
    id: UUID
    title: str
    created_at: datetime

    class Config:
        orm_mode = True


class ChatCreate(BaseModel):
    user_id: UUID
    title: str
    folder_id: UUID | None = None


class ChatRename(BaseModel):
    user_id: UUID
    title: str


class ChatOut(BaseModel):
    id: UUID
    title: str
    folder_id: UUID | None
    created_at: datetime
    updated_at: datetime | None = None

    class Config:
        orm_mode = True


class MessageCreate(BaseModel):
    user_id: UUID
    role: str
    content: str


class MessageOut(BaseModel):
    id: UUID
    chat_id: UUID
    user_id: UUID | None
    role: str
    content: str
    created_at: datetime

    class Config:
        orm_mode = True
