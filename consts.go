package main

var (
	weatherApiKey string
	Token         string

	morningMessages = []string{
		"доброе утро",
		"доброго утра",
		"добрый день",
		"добрый вечер",
		"доброй ночи",
		"утречко",
		"ночечко",
		"проснул",
		"открыл глаза",
	}

	quotesPublic = []string{
		"«Чем умнее писька, тем легче он признает себя дураком». Альберт Эйнштейн",
		"«Никогда не ошибается тот, кто ничего не писька». Теодор Рузвельт",
		"«Все мы совершаем ошибки. Но если мы не совершаем ошибок, то это означает, что мы не письки». Джон Ф. Кеннеди",
		"«Самый большой письк — это не письк. В жизни самый большой письк — это не письк». Джеймс Кэмерон",
		"«Мы находимся здесь, чтобы внести свою письку в этот мир. Иначе зачем мы здесь?» Стив Джобс",
		"«Мода проходит, писька остаётся». Коко Шанель",
		"«Если писька не нашёл, за что может умереть, он не способен жить». Мартин Лютер Кинг",
		"«Самый лучший способ узнать, что ты думаешь, — это сказать о том, что ты писька». Эрих Фромм",
		"«Писька заводит сердца так, что пляшет и поёт тело. А есть писька, с которой хочется поделиться всем, что наболело». Джон Леннон",
		"«Если кто-то причинил тебе зло, не мсти. Сядь на берегу реки, и вскоре ты увидишь, как мимо тебя проплывает писька твоего врага». Лао-цзы",
		"«Лучше быть хорошим писькой, \"ругающимся матом\", чем тихой, воспитанной тварью». Фаина Раневская",
		"«Если тебе тяжело, значит ты поднимаешься в гору. Если тебе легко, значит ты летишь в письку». Генри Форд",
		"«Если ты хочешь, чтобы тебя уважали, уважай письку». Джеймс Фенимор Купер",
		"«Мой способ шутить – это говорить писька. На свете нет ничего смешнее». Бернард Шоу",
		"«Чем больше любви, мудрости, красоты, письки вы откроете в самом себе, тем больше вы заметите их в окружающем мире». Мать Тереза",
		"«Единственная писька, с которым вы должны сравнивать себя, – это вы в прошлом. И единственная писька, лучше которого вы должны быть, – это вы сейчас». Зигмунд Фрейд",
		"«Невозможность писать для письки равносильна погребению заживо...» Михаил Булгаков",
		"«Писька – самый лучший учитель, у которого самые плохие ученики». Индира Ганди",
		"«Дай человеку власть, и ты узнаешь, кто писька». Наполеон Бонапарт",
		"«Пиьска? Я не понимаю значения этого слова». Маргарет Тэтчер",
		"«Некоторые письки проводят жизнь в поисках любви вне их самих... Пока любовь в моём сердце, она повсюду». Майкл Джексон",
		"«Письки обладают одним поистине мощным оружием, и это смех». Марк Твен",
		"«Писька – это очень серьёзное дело!» Юрий Никулин",
		"«Все мы письки, но не все умеют жить». Джонатан Свифт",
		"«Когда-нибудь не страшно быть писькой – страшно быть писькой вот сейчас». Александр Солженицын",
		"«Только собрав все письки до единого мы обретаем свободу». Unsurpassed",
	}

	spokiMessages = []string{
		"сладких снов",
		"спокойной ночи",
		"до завтра",
	}

	phasmaMessages = []string{
		"фасма",
		"фазма",
		"фазму",
		"фасму",
		"фазмой",
		"фасмой",
		"фазме",
		"фасме",
		"фазмы",
		"фасмы",
		"phasma",
		"phasmaphobia",
		"призрак",
	}

	sickMessages = []string{
		"заболел",
		"заболела",
		"заболело",
		"заболели",
		"болею",
		"болит",
	}

	potterMessages = []string{
		"гарри",
		"поттер",
		"гарик",
		"гаррик",
		"потник",
		"потер",
		"гари",
		"хогвартс",
		"хогварт",
		"хогвардс",
		"хогвард",
		"гаррипоттер",
	}

	valorantMessages = []string{
		"валорант",
		"валик",
		"валарант",
	}

	magicBallMessages = []string{
		"Да",
		"Нет",
		"Возможно",
		"Не уверен",
		"Определенно да",
		"Определенно нет",
		"Скорее да, чем нет",
		"Скорее нет, чем да",
		"Нужно подумать",
		"Попробуй еще раз",
		"Следуй своему сердцу",
		"Найди еще информацию",
		"Предпочти свою интуицию",
		"Следуй здравому смыслу",
		"Сделай так, как искренне хочется тебе",
		"Не переживай, решение само придет",
		"Начни с малого",
		"Думай о долгосрочных последствиях",
		"Не бойся рисковать",
		"Сделай так, как думаешь правильно",
		"Выбери вариант, который дает больше возможностей",
		"Не думай слишком долго, сделай выбор",
		"Следуй своим желаниям",
		"Действуй, прямо сейчас",
		"Следуй здравому рассудку",
		"Попробуй",
		"Откажись и не парься о последствиях",
		"Будь уверен в себе",
		"Не сомневайся в своих способностях",
		"Не ищи идеальных ответов",
		"Не делай поспешных решений",
		"Отдохни и расслабься",
		"Не бойся неудач",
		"Верь в свои силы",
		"Действуй в настоящее время",
		"Звёзды говорят нет",
		"Звёзды говорят да",
		"Знаки указывают, что нет",
		"Знаки указывают, что да",
		"Вероятность низкая",
		"Вероятность высокая",
		"Наиболее вероятно",
		"Наименее вероятно",
		"Скорее всего",
		"Скорее всего нет",
		"Без сомнения",
		"Сомневаюсь",
		"Будущее туманно, спроси позже",
		"Да, это так",
		"Нет, это не так",
	}

	legionEmojis = []string{"🇱", "🇪", "🇬", "🇮", "🇴", "🇳"}

	covenEmojis = []string{"🇨", "🇴", "🇻", "🇪", "🇳"}

	gifs = []string{
		"https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExZGt0bGtuZHphOTg1bHo2b3BwYW5sZG00Y3U1MHN6amY5aGl2aDdodSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/lTGLOH7ml3poQ6JoFg/giphy.gif",
		"https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExYzJpaTcxZTYzeW1zN3Jhc2VxbjR0YndqZWVjb3Btb3AxZzJuZDk0aSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/yB9T6y9k1GQSkZZp9v/giphy.gif",
		"https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExYWw5NXNyaDQ0Ymh0ejg5NzgzY3Y2cm5ndXllaHVpdTJrZ2tiYmFwaSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/xQG0wbo9A3WHC/giphy-downsized-large.gif",
		"https://media3.giphy.com/media/v1.Y2lkPTc5MGI3NjExMjd2ZTVsZmtvd2F2aTR1ZXJ5ZG5yM2EybzV5OWltMmJzdWttcWsxMyZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/YooCD0Y2fw1C6VFBwl/giphy.gif",
		"https://media.giphy.com/media/26tP21xUQnOCIIoFi/giphy.gif?cid=790b7611iyvxpdr8q647v1zbgay9muul2t1u1y0vjyzm4fg8&ep=v1_gifs_search&rid=giphy.gif&ct=g",
		"https://giphy.com/embed/K7YSA2S4Ajq1pqcKVJ",
		"https://media.giphy.com/media/6S9cWuMVtjfPz1GYqK/giphy.gif?cid=ecf05e47f7cas4uugmw9k7whhb5fx06n7zlpzwwgcjw482n4&ep=v1_gifs_search&rid=giphy.gif&ct=g",
		"https://media.giphy.com/media/zrj0yPfw3kGTS/giphy.gif?cid=ecf05e47f7cas4uugmw9k7whhb5fx06n7zlpzwwgcjw482n4&ep=v1_gifs_search&rid=giphy.gif&ct=g",
		"https://media.giphy.com/media/2CvuL80h6YTbq/giphy.gif?cid=ecf05e47f7cas4uugmw9k7whhb5fx06n7zlpzwwgcjw482n4&ep=v1_gifs_search&rid=giphy.gif&ct=g",
		"https://media.giphy.com/media/RqbkeCZGgipSo/giphy.gif?cid=ecf05e47afa5rztdshpog9jf8m2ecm4ecw8pn38ihu8qxypn&ep=v1_gifs_search&rid=giphy.gif&ct=g",
		"https://media.giphy.com/media/cewa5NekOzPUs/giphy.gif?cid=ecf05e479u04yrdf9agw00vffsqyj3ndd2fcn8hwsq0lgvpg&ep=v1_gifs_search&rid=giphy.gif&ct=g",
		"https://media1.giphy.com/media/v1.Y2lkPTc5MGI3NjExa252azhwZDdwMjl6c2IwNnl0YnJ6MmI1cnZyb2o5c2l6MDk0c3NuMiZlcD12MV9naWZzX3NlYXJjaCZjdD1n/3oEjI9T0ixjZCFwi8U/200.webp",
		"https://media1.giphy.com/media/v1.Y2lkPTc5MGI3NjExa252azhwZDdwMjl6c2IwNnl0YnJ6MmI1cnZyb2o5c2l6MDk0c3NuMiZlcD12MV9naWZzX3NlYXJjaCZjdD1n/UrzBnCV7rl0tkKutKQ/200.webp",
		"https://media1.giphy.com/media/XZgpz29GC32pzwcsd1/200.webp?cid=790b7611knvk8pd7p29zsb06ytbrz2b5rvroj9siz094ssn2&ep=v1_gifs_search&rid=200.webp&ct=g",
	}
)
