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