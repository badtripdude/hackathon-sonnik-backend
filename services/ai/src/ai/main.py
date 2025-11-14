import os
from io import BytesIO
from typing import Optional

from elevenlabs import ElevenLabs
from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.responses import StreamingResponse, JSONResponse
from openai import OpenAI

from ai.infra.config import settings
from ai.shemas import GptResponse, GptRequest, AsrResponse, TtsRequest

openai_client = OpenAI(
    api_key=os.getenv("OPENAI_API_KEY"),
)

eleven_client = ElevenLabs(
    api_key=os.getenv("ELEVENLABS_API_KEY"),
)
app = FastAPI(title="AI", version="0.1.0",
              )


@app.post("/ai/asr", response_model=AsrResponse)
async def asr_whisper(
        audio: UploadFile = File(..., description="Аудиофайл (wav/mp3/m4a и т.п.)"),
        language: Optional[str] = None,
        prompt: Optional[str] = None,
):
    """
    Speech-To-Text через OpenAI Whisper.
    """
    if not os.getenv("OPENAI_API_KEY"):
        raise HTTPException(status_code=500, detail="OPENAI_API_KEY is not set")

    try:
        # Читаем файл в память
        audio_bytes = await audio.read()
        audio_file = BytesIO(audio_bytes)
        audio_file.name = audio.filename or "audio_input"

        # Вызов OpenAI Whisper
        transcription = openai_client.audio.transcriptions.create(
            model=settings.openai_whisper_model,
            file=audio_file,
            language=language,
            prompt=prompt,
            response_format="json",
        )

        return AsrResponse(
            text=transcription.text,
            language=language,
            model=settings.openai_whisper_model,
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"ASR error: {e}")


# ---------- ElevenLabs TTS: /ai/tts ----------

@app.post("/ai/tts")
async def tts_elevenlabs(req: TtsRequest):
    """
    Text-To-Speech через ElevenLabs.
    Отдаём аудио стримом (audio/mpeg).
    """
    if not os.getenv("ELEVENLABS_API_KEY"):
        raise HTTPException(status_code=500, detail="ELEVENLABS_API_KEY is not set")

    voice_id = req.voice_id or settings.elevenlabs_voice_id
    if not voice_id:
        raise HTTPException(status_code=400, detail="voice_id is required (no default set)")

    try:
        # ElevenLabs Python SDK возвращает bytes
        audio_bytes: bytes = eleven_client.text_to_speech.convert(
            voice_id=voice_id,
            model_id=req.model_id,
            text=req.text,
            output_format=req.output_format,
        )

        filename = "tts_output.mp3"
        media_type = "audio/mpeg"

        def iterfile():
            yield audio_bytes

        headers = {
            "x-voice-id": voice_id,
            "x-model-id": req.model_id,
        }

        return StreamingResponse(
            iterfile(),
            media_type=media_type,
            headers=headers,
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"TTS error: {e}")


# ---------- GPT Text-To-Text: /ai/gpt ----------

@app.post("/ai/gpt", response_model=GptResponse)
async def gpt_text(req: GptRequest):
    """
    Text-To-Text через OpenAI GPT (Responses API).
    """

    model = req.model or settings.openai_gpt_model

    try:
        # Используем Responses API
        response = openai_client.responses.create(
            model=model,
            input=[{"role": m.role, "content": m.content} for m in req.messages],
        )

        # Вытаскиваем текст
        output_text = ""
        if response.output and response.output[0].content:
            # Обычно там один элемент text
            for part in response.output[0].content:
                if part.type == "output_text":
                    output_text += part.text

        usage = getattr(response, "usage", None)

        return GptResponse(
            text=output_text.strip(),
            model=model,
            usage_input_tokens=getattr(usage, "input_tokens", None) if usage else None,
            usage_output_tokens=getattr(usage, "output_tokens", None) if usage else None,
            usage_total_tokens=getattr(usage, "total_tokens", None) if usage else None,
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"GPT error: {e}")


@app.get("/ai/health")
async def health():
    return JSONResponse({"status": "ok"})
