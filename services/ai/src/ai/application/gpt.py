from ai.infra.config import settings
from ai.infra.db.models import Message


async def call_gpt(openai_client, messages: list[Message], model=None) -> str:
    if not model:
        model = settings.openai_gpt_model
    dialog_history = []
    dialog_history.append({
        "role": "system",
        "content": "Ты — ассистент-сонник. Задай сначала 2–3 уточняющих вопроса, если описание сна короткое, "
                   "затем дай интерпретацию в 2–3 блоках: «Факты сна», «Психологическая трактовка», «Практические "
                   "выводы»."
    })
    dialog_history.extend([{"role": m.role, "content": m.content} for m in messages])
    response = openai_client.responses.create(
        model=model,
        input=dialog_history,
    )

    # Вытаскиваем текст
    output_text = ""
    if response.output and response.output[0].content:
        # Обычно там один элемент text
        for part in response.output[0].content:
            if part.type == "output_text":
                output_text += part.text
    return output_text.strip()
