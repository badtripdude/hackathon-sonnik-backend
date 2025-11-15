import os
from io import BytesIO
from typing import Optional
from uuid import UUID

from elevenlabs import ElevenLabs
from fastapi import FastAPI, UploadFile, File, HTTPException, Depends, status
from fastapi.responses import StreamingResponse, JSONResponse
from openai import OpenAI

from ai.application.gpt import call_gpt
from ai.bootstrap import get_chat_repo
from ai.infra.config import settings
from ai.infra.db.db import get_db
from ai.infra.repo.chat_repo import ChatRepository
from ai.shemas import GptResponse, GptRequest, AsrResponse, TtsRequest, InterpretRequest, InterpretResponse, \
    MessageCreate, MessageOut, FolderOut, FolderCreate, FolderRename, ChatOut, ChatCreate, ChatRename

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
    if not os.getenv("ELEVENLABS_API_KEY"):
        raise HTTPException(status_code=500, detail="ELEVENLABS_API_KEY is not set")

    voice_id = req.voice_id or settings.elevenlabs_voice_id
    if not voice_id:
        raise HTTPException(status_code=400, detail="voice_id is required (no default set)")

    try:
        raw_audio = eleven_client.text_to_speech.convert(
            voice_id=voice_id,
            model_id=req.model_id,
            text=req.text,
            output_format=req.output_format,
        )

        # если convert() возвращает генератор байтов:
        def iter_audio():
            for chunk in raw_audio:
                # гарантируем, что это bytes
                if isinstance(chunk, str):
                    chunk = chunk.encode("utf-8")
                yield chunk

        headers = {
            "x-voice-id": voice_id,
            "x-model-id": req.model_id,
        }

        return StreamingResponse(
            iter_audio(),
            media_type="audio/mpeg",
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


@app.post("/ai/interpret", response_model=InterpretResponse)
async def interpret(
        req: InterpretRequest,
        db=Depends(get_db),
        # user_id: uuid.UUID = Depends(...),  # подставь свою зависимость из auth
):
    user_id = req.user_id
    repo = ChatRepository(db)

    # если нет chat_id — создаем новый чат
    if req.chat_id is None:
        chat = await repo.create_chat(
            user_id=user_id,
            title="Новый сон",  # потом можно переименовать по содержанию
        )
        chat_id = chat.id
    else:
        chat = await repo.get_chat(chat_id=req.chat_id, user_id=user_id)
        if not chat:
            raise HTTPException(status_code=404, detail="Chat not found")
        chat_id = chat.id

    # сохраняем сообщение пользователя
    await repo.add_message(
        chat_id=chat_id,
        role="user",
        content=req.message,
        user_id=user_id,
    )

    # забираем историю для контекста
    history = await repo.get_last_messages(chat_id=chat_id, limit=20)

    # тут собираешь messages для LLM (system + history + новый user)
    # llm_messages = build_llm_messages(...)
    # gpt_response = await call_gpt(llm_messages)

    gpt_reply_text = await call_gpt(openai_client, history, model=settings.openai_gpt_model)

    # сохраняем ответ ассистента
    await repo.add_message(
        chat_id=chat_id,
        role="assistant",
        content=gpt_reply_text,
        user_id=None,
    )

    return InterpretResponse(
        chat_id=chat_id,
        reply=gpt_reply_text,
    )


# ---------- FOLDERS ----------

@app.get("/folders", response_model=list[FolderOut])
async def list_folders(
        user_id: UUID,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    GET /ai/folders?user_id=...
    """
    folders = await repo.list_folders(user_id=user_id)
    return folders


@app.post("/folders", response_model=FolderOut, status_code=status.HTTP_201_CREATED)
async def create_folder(
        body: FolderCreate,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    POST /ai/folders
    {
      "user_id": "...",
      "title": "Мой промпт-инженеринг"
    }
    """
    folder = await repo.create_folder(
        user_id=body.user_id,
        title=body.title,
    )
    return folder


@app.patch("/folders/{folder_id}", response_model=FolderOut)
async def rename_folder(
        folder_id: UUID,
        body: FolderRename,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    PATCH /ai/folders/{folder_id}
    {
      "user_id": "...",
      "title": "Новое имя"
    }
    """
    folder = await repo.rename_folder(
        folder_id=folder_id,
        user_id=body.user_id,
        new_title=body.title,
    )
    if folder is None:
        raise HTTPException(status_code=404, detail="Folder not found")
    return folder


# ---------- CHATS ----------

@app.get("/chats", response_model=list[ChatOut])
async def list_chats(
        user_id: UUID,
        folder_id: UUID | None = None,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    GET /ai/chats?user_id=...&folder_id=...
    """
    chats = await repo.list_chats(
        user_id=user_id,
        folder_id=folder_id,
    )
    return chats


@app.post("/chats", response_model=ChatOut, status_code=status.HTTP_201_CREATED)
async def create_chat(
        body: ChatCreate,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    POST /ai/chats
    {
      "user_id": "...",
      "title": "Новый чат",
      "folder_id": null
    }
    """
    chat = await repo.create_chat(
        user_id=body.user_id,
        title=body.title,
        folder_id=body.folder_id,
    )
    return chat


@app.patch("/chats/{chat_id}", response_model=ChatOut)
async def rename_chat(
        chat_id: UUID,
        body: ChatRename,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    PATCH /ai/chats/{chat_id}
    {
      "user_id": "...",
      "title": "Новое имя чата"
    }
    """
    chat = await repo.rename_chat(
        chat_id=chat_id,
        user_id=body.user_id,
        new_title=body.title,
    )
    if chat is None:
        raise HTTPException(status_code=404, detail="Chat not found")
    return chat


# ---------- MESSAGES ----------

@app.get("/chats/{chat_id}/messages", response_model=list[MessageOut])
async def list_messages(
        chat_id: UUID,
        user_id: UUID,
        limit: int = 20,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    GET /ai/chats/{chat_id}/messages?user_id=...&limit=20
    (user_id используем для проверки владения чатом)
    """
    chat = await repo.get_chat(chat_id=chat_id, user_id=user_id)
    if chat is None:
        raise HTTPException(status_code=404, detail="Chat not found")

    messages = await repo.get_last_messages(chat_id=chat_id, limit=limit)
    return messages


@app.post(
    "/chats/{chat_id}/messages",
    response_model=MessageOut,
    status_code=status.HTTP_201_CREATED,
)
async def add_message(
        chat_id: UUID,
        body: MessageCreate,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    POST /ai/chats/{chat_id}/messages
    {
      "user_id": "...",
      "role": "user" | "assistant" | "system",
      "content": "текст сообщения"
    }
    """
    chat = await repo.get_chat(chat_id=chat_id, user_id=body.user_id)
    if chat is None:
        raise HTTPException(status_code=404, detail="Chat not found")

    msg = await repo.add_message(
        chat_id=chat_id,
        role=body.role,
        content=body.content,
        user_id=body.user_id,
    )
    return msg


@app.delete("/folders/{folder_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_folder(
        folder_id: UUID,
        user_id: UUID,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    DELETE /ai/folders/{folder_id}?user_id=...
    """
    ok = await repo.delete_folder(folder_id=folder_id, user_id=user_id)
    if not ok:
        raise HTTPException(status_code=404, detail="Folder not found")
    # 204 — без тела
    return


@app.delete("/chats/{chat_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_chat(
        chat_id: UUID,
        user_id: UUID,
        repo: ChatRepository = Depends(get_chat_repo),
):
    """
    DELETE /ai/chats/{chat_id}?user_id=...
    """
    ok = await repo.delete_chat(chat_id=chat_id, user_id=user_id)
    if not ok:
        raise HTTPException(status_code=404, detail="Chat not found")
    return


@app.get("/ai/health")
async def health():
    return JSONResponse({"status": "ok"})
