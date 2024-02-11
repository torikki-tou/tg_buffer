import os
import sys

import aiogram
import aiogram.types
import asyncio
import pika


def main():
    connection = pika.BlockingConnection(
        pika.ConnectionParameters(host="rabbit"),
    )
    channel = connection.channel()

    bot = aiogram.Bot(os.getenv("TOKEN"))
    dp = aiogram.Dispatcher()

    def callback(ch, method, properties, body):
        print(f" [x] Received {body.decode()}")
        upd = aiogram.types.update.Update.model_validate_json(body)
        try:
            asyncio.run(dp.feed_raw_update(bot, upd.model_dump()))
        except BaseException as e:
            print(e)

    channel.basic_consume(
        queue="updates",
        on_message_callback=callback,
        auto_ack=True,
    )

    print(" [*] Waiting for messages. To exit press CTRL+C")
    channel.start_consuming()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("Interrupted")
        try:
            sys.exit(0)
        except SystemExit:
            os._exit(0)
