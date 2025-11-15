from ai.infra.config import settings
from ai.infra.db.models import Message


async def call_gpt(openai_client, messages: list[Message], model=None) -> str:
    if not model:
        model = settings.openai_gpt_model
    response = openai_client.responses.create(
        model=model,
        input=[{"role": m.role, "content": m.content} for m in messages],
    )

    # Вытаскиваем текст
    output_text = ""
    if response.output and response.output[0].content:
        # Обычно там один элемент text
        for part in response.output[0].content:
            if part.type == "output_text":
                output_text += part.text
    return output_text.strip()
