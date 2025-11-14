from __future__ import annotations

import os

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    openai_api_key: str = Field(alias="OPENAI_API_KEY")
    elevenlabs_api_key: str = Field(alias="ELEVENLABS_API_KEY")
    openai_gpt_model: str = Field(alias="OPENAI_GPT_MODEL", default="gpt-4.1-mini")
    openai_whisper_model: str = Field(alias="OPENAI_WHISPER_MODEL", default="whisper-1")
    elevenlabs_voice_id: str = Field(alias="ELEVENLABS_VOICE_ID", default="")
    # --- Logging ---
    logging_config: str = Field(alias="LOGGING_CONFIG", default=".logging.yml")
    log_level: str = Field(alias="LOG_LEVEL", default="INFO")

    model_config = SettingsConfigDict(
        env_file=(
            ".env",
            f".env.{os.getenv('APP_ENV', 'dev')}",
            ".env.local",
            f".env.{os.getenv('APP_ENV', 'dev')}.local",
        ),
        env_file_encoding="utf-8",
        extra="ignore",
    )


settings = Settings()
