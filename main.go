package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"encoding/json"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/lefes/curly-broccoli/jokes"
	"github.com/lefes/curly-broccoli/quotes"
)

var (
	Token   string = ""
	counter        = 0
)

// Структура для хранения количества смертей
type DeathCounter struct {
	Count int `json:"count"`
}

// Путь к файлу для сохранения счетчика
const deathCountFile = "death_count.json"

func poll(session *discordgo.Session, m *discordgo.MessageCreate) {
	// Randomly create a poll with 3 options in the channel
	// Take 3 person from the channel
	users, err := session.GuildMembers(m.GuildID, "", 300)
	if err != nil {
		fmt.Println("error getting users,", err)
		return
	}

	// Get 3 random users
	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	users = users[:3]

	// Загрузка счетчика смертей из файла
	counter := loadDeathCount()

	// Create a poll
	poll := &discordgo.MessageEmbed{
		Title: "Кто сегодня писька??? 🤔🤔🤔",
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "1",
				Value:  getNick(users[0]),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "2",
				Value:  getNick(users[1]),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "3",
				Value:  getNick(users[2]),
				Inline: true,
			},
		},
	}

	// Send the poll
	pollMessage, err := session.ChannelMessageSendEmbed(m.ChannelID, poll)
	if err != nil {
		fmt.Println("error sending poll,", err)
		return
	}

	reactions := []string{"1️⃣", "2️⃣", "3️⃣"}
	// Add reactions to the poll

	for _, v := range reactions {
		err := session.MessageReactionAdd(pollMessage.ChannelID, pollMessage.ID, v)
		if err != nil {
			fmt.Println("error adding reaction,", err)
			return
		}
	}

	// Wait for 2 hours
	time.Sleep(30 * time.Minute)

	// Get the poll results
	pollResults, err := session.ChannelMessage(pollMessage.ChannelID, pollMessage.ID)
	if err != nil {
		fmt.Println("error getting poll results,", err)
		return
	}

	// Get the most voted option
	var mostVotedOption string
	var mostVotedOptionCount int
	for _, v := range pollResults.Reactions {
		if v.Count > mostVotedOptionCount {
			mostVotedOption = v.Emoji.Name
			mostVotedOptionCount = v.Count
		}
	}

	// Get the winner
	var winner *discordgo.Member
	switch mostVotedOption {
	case "1️⃣":
		winner = users[0]
	case "2️⃣":
		winner = users[1]
	case "3️⃣":
		winner = users[2]
	}

	// Congratulate the winner
	_, err = session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Поздравляем, <@%s>, ты сегодня писька! 🎉🎉🎉", winner.User.ID))
	if err != nil {
		fmt.Println("error congratulating the winner,", err)
	}

}

func init() {
	// Load dotenv
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	if Token == "" {
		flag.StringVar(&Token, "token", "", "token")
		flag.Parse()
	}
	if Token == "" {
		Token = os.Getenv("TOKEN")
		if Token == "" {
			panic("You need to input the token.")
		}
	}
	rand.Seed(time.Now().UnixNano())

}

func getNick(member *discordgo.Member) string {
	if member.Nick == "" {
		return member.User.Username
	}
	return member.Nick
}

func piskaMessage(users []string) string {
	var message string
	rand.Seed(time.Now().UnixNano())
	message += "🤔🤔🤔"
	for _, user := range users {
		// #nosec G404 -- This is a false positive
		piskaProc := rand.Intn(101)
		if piskaProc == 0 {
			message += fmt.Sprintf("\nИзвини, <@%s>, но ты совсем не писька (0%%), приходи когда описюнеешь", user)
		} else if piskaProc == 100 {
			message += fmt.Sprintf("\n<@%s>, ты просто прекрасная писька на ВСЕ 100%%", user)
		} else if piskaProc >= 50 {
			message += fmt.Sprintf("\n<@%s> писька на %d%%, молодец, так держать!", user, piskaProc)
		} else {
			message += fmt.Sprintf("\n<@%s> писька на %d%%, но нужно еще вырасти", user, piskaProc)
		}
	}
	return message
}

// Загрузка счетчика смертей из файла
func loadDeathCount() int {
	data, err := ioutil.ReadFile(deathCountFile)
	if err != nil {
		// Если файл не существует или есть ошибка чтения, начинаем с 0
		return 0
	}

	var counter DeathCounter
	err = json.Unmarshal(data, &counter)
	if err != nil {
		fmt.Println("Ошибка чтения счетчика смертей:", err)
		return 0
	}

	return counter.Count
}

// Сохранение счетчика смертей в файл
func saveDeathCount(count int) {
	counter := DeathCounter{Count: count}
	data, err := json.MarshalIndent(counter, "", "  ")
	if err != nil {
		fmt.Println("Ошибка сериализации счетчика смертей:", err)
		return
	}

	err = ioutil.WriteFile(deathCountFile, data, 0644)
	if err != nil {
		fmt.Println("Ошибка записи счетчика смертей:", err)
	}
}

func main() {
	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Create interface for quotes
	quote := quotes.New()

	morningMessages := []string{
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

	quotesPublic := []string{
		"«Чем умнее писька, тем легче он признает себя дураком». Альберт Эйнштейн",
		"«Никогда не ошибается тот, кто ничего не писька». Теодор Рузвельт",
		"«Все мы совершаем ошибки. Но если мы не совершаем ошибок, то это означает, что мы не письки». Джон Ф. Кеннеди",
		"«Самый большой письк — это не письк. В жизни самый большой письк — это не письк». Джеймс Кэмерон",
		"«Мы находимся здесь, чтобы внести свою письку в этот мир. Иначе зачем мы здесь?» Стив Джобс",
		"«Мода проходит, писька остаётся». Коко Шанель",
		"«Если писька не нашёл, за что может умереть, он не способен жить». Мартин Лютер Кинг",
		"«Самый лучший способ узнать, что ты думаешь, — это сказать о том, что ты писька». Эрих Фромм",
		"«Писька заводит сердца так, что пляшет и поёт тело. А есть писька, с которой хочется поделиться всем, что наболело». Джон Леннон",
		"«Если кто-то причинил тебе зло, не мсти. Сядь на берегу реки, и вскоре ты увидишь, как мимо тебя проплывает писька твоего врага». Лао-цзы",
		"«Лучше быть хорошим писькой, \"ругающимся матом\", чем тихой, воспитанной тварью». Фаина Раневская",
		"«Если тебе тяжело, значит ты поднимаешься в гору. Если тебе легко, значит ты летишь в письку». Генри Форд",
		"«Если ты хочешь, чтобы тебя уважали, уважай письку». Джеймс Фенимор Купер",
		"«Мой способ шутить – это говорить писька. На свете нет ничего смешнее». Бернард Шоу",
		"«Чем больше любви, мудрости, красоты, письки вы откроете в самом себе, тем больше вы заметите их в окружающем мире». Мать Тереза",
		"«Единственная писька, с которым вы должны сравнивать себя, – это вы в прошлом. И единственная писька, лучше которого вы должны быть, – это вы сейчас». Зигмунд Фрейд",
		"«Невозможность писать для письки равносильна погребению заживо...» Михаил Булгаков",
		"«Писька – самый лучший учитель, у которого самые плохие ученики». Индира Ганди",
		"«Дай человеку власть, и ты узнаешь, кто писька». Наполеон Бонапарт",
		"«Пиьска? Я не понимаю значения этого слова». Маргарет Тэтчер",
		"«Некоторые письки проводят жизнь в поисках любви вне их самих... Пока любовь в моём сердце, она повсюду». Майкл Джексон",
		"«Письки обладают одним поистине мощным оружием, и это смех». Марк Твен",
		"«Писька – это очень серьёзное дело!» Юрий Никулин",
		"«Все мы письки, но не все умеют жить». Джонатан Свифт",
		"«Когда-нибудь не страшно быть писькой – страшно быть писькой вот сейчас». Александр Солженицын",
		"«Только собрав все письки до единого мы обретаем свободу». Unsurpassed",
	}

	spokiMessages := []string{
		"сладких снов",
		"спокойной ночи",
		"до завтра",
	}

	phasmaMessages := []string{
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

	sickMessages := []string{
		"заболел",
		"заболела",
		"заболело",
		"заболели",
		"болею",
		"болит",
	}

	potterMessages := []string{
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

	valorantMessages := []string{
		"валорант",
		"валик",
		"валарант",
	}

	magicBallMessages := []string{
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

	legionEmojis := []string{"🇱", "🇪", "🇬", "🇮", "🇴", "🇳"}

	covenEmojis := []string{"🇨", "🇴", "🇻", "🇪", "🇳"}

	session.Identify.Intents = discordgo.IntentsGuildMessages

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		// Checking on spoki and morning event
		morning := false
		for _, v := range morningMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				morning = true
			}
		}

		spoki := false
		for _, v := range spokiMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				spoki = true
			}
		}

		if morning {
			emoji, err := session.GuildEmoji(m.GuildID, "1016631674106294353")
			if err != nil {
				emoji = &discordgo.Emoji{
					Name: "🫠",
				}
			}
			err = session.MessageReactionAdd(m.ChannelID, m.ID, emoji.APIName())
			if err != nil {
				fmt.Println("error reacting to message,", err)
			}
		}

		if spoki {
			emoji, err := session.GuildEmoji(m.GuildID, "1016631826338566144")
			if err != nil {
				emoji = &discordgo.Emoji{
					Name: "😴",
				}
			}
			err = session.MessageReactionAdd(m.ChannelID, m.ID, emoji.APIName())
			if err != nil {
				fmt.Println("error reacting to message,", err)
			}
		}

		// Checking on LEGION event
		if strings.Contains(strings.ToLower(m.Content), "легион") {
			for _, v := range legionEmojis {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, v)
				time.Sleep(100 * time.Millisecond)
				if err != nil {
					fmt.Println("error reacting to message,", err)
				}
			}
		}

		
		// Checking on COVEN event
		if strings.Contains(strings.ToLower(m.Content), "ковен") || strings.Contains(strings.ToLower(m.Content), "сестры") || strings.Contains(strings.ToLower(m.Content), "сёстры") {
			for _, v := range covenEmojis {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, v)
				time.Sleep(100 * time.Millisecond)
				if err != nil {
					fmt.Println("error reacting to message,", err)
				}
			}
		}

		// Checking on spasibo message
		if strings.Contains(strings.ToLower(m.Content), "спасибо") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Это тебе спасибо! 😎😎😎", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}
		
		// Checking on бобр message
		if strings.Contains(strings.ToLower(m.Content), "бобр") || strings.Contains(strings.ToLower(m.Content), "бобер") || strings.Contains(strings.ToLower(m.Content) "курва" {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Kurwa bóbr. Ja pierdolę, Jakie bydlę jebane 🦫🦫🦫", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on бобр message
		if strings.Contains(strings.ToLower(m.Content), "бобр") || strings.Contains(strings.ToLower(m.Content), "бобер") || strings.Contains(strings.ToLower(m.Content) "курва" {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Kurwa bóbr. Ja pierdolę, Jakie bydlę jebane 🦫🦫🦫", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}
		
		// Checking on "привет" message
		if strings.Contains(strings.ToLower(m.Content), "привет") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Привет!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Проверка, что сообщение от целевого пользователя и содержит слово "умер"
		if m.Author.ID == "850043154207604736" && strings.Contains(strings.ToLower(m.Content), "умер") {
			counter++
			response := fmt.Sprintf("Сволочи, они убили @%s %d раз(а) 💀🔫", m.Author.Username, counter)
			_, err := s.ChannelMessageSend(m.ChannelID, response)
			if err != nil {
				fmt.Println("Ошибка отправки сообщения:", err)
				return
			}

			// Сохранение счетчика смертей в файл
			saveDeathCount(counter)
		}

		// Checking on "пиф-паф" message
		if strings.Contains(strings.ToLower(m.Content), "пиф") && strings.ContainsAny(strings.ToLower(m.Content), "паф") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Пиф-паф!🔫🔫🔫", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		} else if strings.Contains(strings.ToLower(m.Content), "pif") && strings.ContainsAny(strings.ToLower(m.Content), "paf") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Pif-paf!🔫🔫🔫", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "алкаш" message
		if strings.Contains(strings.ToLower(m.Content), "алкаш") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Эй мальчик, давай обмен,я же вижу что ты алкаш (c) Чайок", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "дед инсайд" message
		if strings.Contains(strings.ToLower(m.Content), "дед инсайд") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Глисты наконец-то померли?", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "я гей" message
		if strings.Contains(strings.ToLower(m.Content), "я гей") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "я лесбиянка" message
		if strings.Contains(strings.ToLower(m.Content), "я лесбиянка") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "я би" message
		if strings.Contains(strings.ToLower(m.Content), "я би") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "я натурал" message
		if strings.Contains(strings.ToLower(m.Content), "я натурал") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я иногда тоже!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "понедельник" message
		if strings.Contains(strings.ToLower(m.Content), "понедельник") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "День тяжелый 😵‍💫", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "заболел" message
		for _, v := range sickMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				_, err := s.ChannelMessageSendReply(m.ChannelID, "Скорее выздоравливай и больше не болей! 😍", m.Reference())
				if err != nil {
					fmt.Println("error sending message,", err)
				}
			}
		}

		// Checking on "фазма" message
		for _, v := range phasmaMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, "👻")
				if err != nil {
					fmt.Println("error reacting to message,", err)
				}
			}
		}

		// Checking on "полчаса" message
		if strings.Contains(strings.ToLower(m.Content), "полчаса") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "полчаса, полчаса - не вопрос. Не ответ полчаса, полчаса (c) Чайок", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "!голосование") {
			go poll(s, m)
		}

		if strings.Contains(strings.ToLower(m.Content), "!quote") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, quote.GetRandom(), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "!academia") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, quote.GetRandomAcademia(), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!медведь") {
			user := m.Author.ID
			if len(m.Mentions) != 0 {
				//#nosec G404 -- This is a false positive
				if rand.Intn(10) == 0 {
					_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, кажется медведь прямо сейчас завалит тебя 🐻🐻🐻", user), m.Reference())
					if err != nil {
						fmt.Println("error sending message,", err)
					}
					return
				}
				member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
				if err == nil {
					user = member.User.ID
				}
			}

			//#nosec G404 -- This is a false positive
			medvedProc := rand.Intn(101)
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, завалишь медведя с %d%% вероятностью 🐻", user, medvedProc), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!ролл") || strings.HasPrefix(strings.ToLower(m.Content), "!d20") {
			user := m.Author.ID
			//#nosec G404 -- This is a false positive
			roll := rand.Intn(20) + 1
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, ты выкинул %d", user, roll), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!писька") {
			user := m.Author.ID
			if len(m.Mentions) != 0 {
				//#nosec G404 -- This is a false positive
				if rand.Intn(10) == 0 {
					_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, а вот и нет, писька это ты!!!", user), m.Reference())
					if err != nil {
						fmt.Println("error sending message,", err)
					}
					return
				}
				member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
				if err == nil {
					user = member.User.ID
				}
			}

			//#nosec G404 -- This is a false positive
			piskaProc := rand.Intn(101)

			if piskaProc == 100 {
				_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, ты просто прекрасная писька на ВСЕ 100%%", user), m.Reference())
				if err != nil {
					fmt.Println("error sending message,", err)
				}
				return
			}

			if piskaProc == 0 {
				_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Извини, <@%s>, но ты совсем не писька (0%%), приходи когда описюнеешь", user), m.Reference())
				if err != nil {
					fmt.Println("error sending message,", err)
				}
				return
			}

			//#nosec G404 -- This is a false positive
			if rand.Intn(2) == 0 && piskaProc > 50 {
				//#nosec G404 -- This is a false positive
				_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> настоящая писька на %d%%, вот тебе цитата: %s", user, piskaProc, quotesPublic[rand.Intn(len(quotesPublic))]), m.Reference())
				if err != nil {
					fmt.Println("error sending message,", err)
				}
				return
			}

			if piskaProc > 50 {
				_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> писька на %d%%, молодец, так держать!", user, piskaProc), m.Reference())
				if err != nil {
					fmt.Println("error sending message,", err)
				}
				return
			}

			//#nosec G404 -- This is a false positive
			_, err = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> писька на %d%%, но нужно еще вырасти!", user, piskaProc), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!письки") {
			user := m.Author.ID
			users := make([]string, 0)
			if len(m.Mentions) != 0 {
				//#nosec G404 -- This is a false positive
				if rand.Intn(10) == 0 {
					_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, а вот и нет, писька это ты!!!", user), m.Reference())
					if err != nil {
						fmt.Println("error sending message,", err)
					}
					return
				}
				for _, mention := range m.Mentions {
					member, err := s.GuildMember(m.GuildID, mention.ID)
					if err == nil {
						users = append(users, member.User.ID)
					}
				}
			}

			_, err := s.ChannelMessageSendReply(m.ChannelID, piskaMessage(users), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
			return

		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!анекдот") {
			joke, err := jokes.GetJoke()
			if err != nil {
				fmt.Println("error getting joke,", err)
				return
			}
			_, err = s.ChannelMessageSendReply(m.ChannelID, joke, m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		for _, v := range potterMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, "🧙")
				if err != nil {
					fmt.Println("error reacting message,", err)
				}
			}
		}

		for _, v := range valorantMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, "🔥")
				if err != nil {
					fmt.Println("error reacting message,", err)
				}
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я писюн") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже писюн!!!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я писька") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже писька!!!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!шар") {
			//#nosec G404 -- This is a false positive
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Мой ответ: %s", magicBallMessages[rand.Intn(len(magicBallMessages))]), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "демон") {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "👹")
			if err != nil {
				fmt.Println("error reacting message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "клоун") {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "🤡")
			if err != nil {
				fmt.Println("error reacting message,", err)
			}
		}

	})

	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan struct{})

}
