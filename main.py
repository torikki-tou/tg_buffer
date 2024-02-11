import aiogram
import aiogram.types as types
import asyncio
import aio_pika as pika
import aio_pika.message as message
import time
import random

bot = aiogram.Bot("123:qwerty") # Fake token
dp = aiogram.Dispatcher()


async def main():
    con = await pika.connect("amqp://guest:guest@rabbit")
    channel = await con.channel()

    queue = await channel.declare_queue("updates", durable=True)

    await queue.consume(callback, no_ack=True)

    print(" [*] Waiting for messages. To exit press CTRL+C")
    await asyncio.Future()


async def callback(msg: message.AbstractIncomingMessage):
    print(f" [x] Received {msg.body.decode()}")
    upd = aiogram.types.update.Update.model_validate_json(msg.body)
    try:
        await dp.feed_raw_update(bot, upd.model_dump())
    except BaseException as e:
        print(e)

@dp.message()
async def cmd_start(_: types.Message):
    t = time.time()
    await asyncio.sleep(random.randint(1, 5))
    print(t)


if __name__ == "__main__":
    asyncio.run(main())
