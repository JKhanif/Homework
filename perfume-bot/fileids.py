import asyncio
import aiohttp
from aiogram import Bot
from aiogram.types import BufferedInputFile

BOT_TOKEN = "8621070234:AAFPllQjIW8vukY5ZmS4a3PW1Uh9XdDo3RA"
CHAT_ID = 568523564

urls = [
"https://pcdn.goldapple.ru/p/p/7195500011/web/696d674d61696e5064708ddc3c997df0d79.jpg",
"https://basket-12.wbbasket.ru/vol1900/part190081/190081400/images/big/1.webp",
"https://avatars.mds.yandex.net/get-mpic/13448948/2a00000199f1a8e580f9af5fc7930935037d/orig",
"https://usmall.ru/image/1371/89/24/d1a50f7b0f6231739063e040b2984bdf.jpeg",
"https://avatars.mds.yandex.net/get-mpic/5254153/2a0000018e9ff972b8c42b42d4cb9c6ddde1/orig",
"https://basket-29.wbbasket.ru/vol5668/part566890/566890166/images/big/1.webp",
"https://basket-26.wbbasket.ru/vol4721/part472112/472112553/images/big/1.webp",
"https://pcdn.goldapple.ru/p/p/7380100006/web/696d674d61696e5064708ddc3e31a923812.jpg",
"https://basket-20.wbbasket.ru/vol3341/part334105/334105700/images/big/1.webp",
"https://basket-38.wbbasket.ru/vol8383/part838392/838392201/images/big/1.webp",
"https://basket-19.wbbasket.ru/vol3123/part312344/312344845/images/big/1.webp",
"https://basket-17.wbbasket.ru/vol2824/part282412/282412120/images/big/1.webp",
"https://basket-15.wbbasket.ru/vol2303/part230385/230385177/images/big/1.webp",
"https://i.letu.ru/common/img/pim/2025/08/AUX_36282651-ed9c-4f97-8ac3-1ae46f7ead4c.jpg",
"https://pcdn.goldapple.ru/p/p/26250800001/web/696d674d61696e5064705f38346363646261383162336534393231396634316562363037356134383261388ddcbd17acbb6c1.jpg",
"https://basket-31.wbbasket.ru/vol6433/part643335/643335640/images/big/1.webp",
"https://avatars.mds.yandex.net/get-mpic/12411215/2a00000199c044053bc56539c695dba187ae/orig",
"https://basket-29.wbbasket.ru/vol5765/part576514/576514760/images/big/1.webp",
"https://fimgs.net/photogram/p1200/9f/im/6a8I6I9wOKCT60kA.jpg"
]

async def main():
    bot = Bot(token=BOT_TOKEN)

    async with aiohttp.ClientSession() as session:
        for url in urls:
            try:
                async with session.get(url, headers={
                    "User-Agent": "Mozilla/5.0"
                }) as r:

                    if r.status != 200:
                        print(f"SKIP {url} status={r.status}")
                        continue

                    data = await r.read()

                photo = BufferedInputFile(data, filename="image.jpg")

                msg = await bot.send_photo(
                    chat_id=CHAT_ID,
                    photo=photo
                )

                print(url, "->", msg.photo[-1].file_id)

            except Exception as e:
                print("ERROR", url, e)

    await bot.session.close()

asyncio.run(main())