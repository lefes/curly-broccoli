package main

import (
	"fmt"
	rand "math/rand/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/lefes/curly-broccoli/jokes"
	"github.com/lefes/curly-broccoli/pkg/logging"
	"github.com/lefes/curly-broccoli/pkg/weather"
	"github.com/lefes/curly-broccoli/quotes"
)

func handleRaceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if raceInProgress {
		s.ChannelMessageSend(m.ChannelID, "Гонка уже идет! Дождитесь окончания текущей гонки.")
		return
	}

	raceInProgress = true
	s.ChannelMessageSend(m.ChannelID, "Заезд начинается! Напишите !го, чтобы присоединиться. У вас есть 1 минута.")

	time.AfterFunc(1*time.Minute, func() {
		startRace(s, m)
	})
}

func handleJoinRaceCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !raceInProgress {
		s.ChannelMessageSend(m.ChannelID, "Сейчас нет активной гонки. Напишите !гонка, чтобы начать новую.")
		return
	}

	raceMutex.Lock()
	defer raceMutex.Unlock()

	if _, exists := raceParticipants[m.Author.ID]; exists {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, ты уже участвуешь в заезде!", m.Author.ID))
		return
	}

	emoji := raceEmojis[rand.IntN(len(raceEmojis))]
	raceParticipants[m.Author.ID] = emoji
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> присоединился к гонке как %s!", m.Author.ID, emoji))
}

func startRace(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(raceParticipants) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Недостаточно участников для начала гонки. Гонка отменена.")
		raceInProgress = false
		raceParticipants = make(map[string]string)
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Гонка начнется через 30 секунд! 🏁")
	time.Sleep(30 * time.Second)

	raceTrack := make(map[string]int)
	trackLength := 20

	for id := range raceParticipants {
		raceTrack[id] = 0
	}

	raceMessageContent := buildRaceMessage(raceTrack, raceParticipants, trackLength)
	raceMessage, err := s.ChannelMessageSend(m.ChannelID, raceMessageContent)
	if err != nil {
		fmt.Println("error sending race message:", err)
		return
	}

	winner := ""
	for winner == "" {
		time.Sleep(1 * time.Second)

		for id := range raceParticipants {
			raceTrack[id] += rand.IntN(3)
			if raceTrack[id] >= trackLength {
				raceTrack[id] = trackLength
				winner = id
				break
			}
		}

		updatedRaceMessageContent := buildRaceMessage(raceTrack, raceParticipants, trackLength)
		_, err := s.ChannelMessageEdit(m.ChannelID, raceMessage.ID, updatedRaceMessageContent)
		if err != nil {
			fmt.Println("error editing race message:", err)
			return
		}
	}

	winnerMessage := fmt.Sprintf("🎉 Победитель гонки: <@%s> %s! Поздравляем! 🎉", winner, raceParticipants[winner])
	_, err = s.ChannelMessageSend(m.ChannelID, winnerMessage)
	if err != nil {
		fmt.Println("error sending winner message:", err)
	}

	raceInProgress = false
	raceParticipants = make(map[string]string)
}

func buildRaceMessage(raceTrack map[string]int, raceParticipants map[string]string, trackLength int) string {
	raceMessage := "🏁 Гонка в процессе: 🏁\n\n"
	longestName := 0

	for name := range raceParticipants {
		if len(name) >= longestName {
			longestName = len(name)
		}
	}

	for id, emoji := range raceParticipants {
		track := strings.Repeat("-", raceTrack[id]) + emoji + strings.Repeat("-", trackLength-raceTrack[id])
		raceMessage += fmt.Sprintf("<@%s>: ", id) + strings.Repeat(" ", longestName-len(id)) + fmt.Sprintf("%s\n", track)
	}

	return raceMessage
}

func handleBeerCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Использование: !пиво <число от 1 до 40>")
		return
	}

	amount, err := strconv.Atoi(args[1])
	if err != nil || amount < 1 || amount > 40 {
		s.ChannelMessageSend(m.ChannelID, "Пожалуйста, введите число от 1 до 40.")
		return
	}

	chance := (amount * 3)
	roll := rand.IntN(120) + 1
	fmt.Printf("Date: %s, Author: %s, Amount: %d, Chance: %d, Roll: %d\n", time.Now().Format("2006-01-02 15:04:05"), m.Author.Username, amount, chance, roll)

	successMessages := []string{
		fmt.Sprintf("<@%s> смог осилить %d литров пива! 🍺", m.Author.ID, amount),
		fmt.Sprintf("<@%s> успешно справился с %d литрами! Это достойно уважения! 🍻", m.Author.ID, amount),
		fmt.Sprintf("<@%s> выпил %d литров, пивной монстр на свободе! 🍻🦹", m.Author.ID, amount),
		fmt.Sprintf("<@%s> залпом поглотил %d литров и выглядит, как чемпион! 🏆", m.Author.ID, amount),
		fmt.Sprintf("<@%s> выпил %d литров пива и готов к новым свершениям! 🍻🚀", m.Author.ID, amount),
		fmt.Sprintf("<@%s> справился с %d литрами пива! Не плохо! 🍺", m.Author.ID, amount),
		fmt.Sprintf("<@%s> выпил %d литров пива и готов к новым подвигам! 🍻🚀", m.Author.ID, amount),
	}

	failureMessages := []string{
		fmt.Sprintf("<@%s> не смог осилить даже %d литров пива и облевал весь пол! Кто это убирать будет?! 🤢🤮", m.Author.ID, roll/3),
		fmt.Sprintf("<@%s> попытался выпить %d литр, но потерпел неудачу и свалился под стол! 😵", m.Author.ID, roll/3),
		fmt.Sprintf("<@%s> проиграл борьбу на %d литрах пива и отправляется в бан на %s! 😴", m.Author.ID, roll/3, getMuteDuration(amount)),
		fmt.Sprintf("<@%s> взял на себя слишком много! %d литр пива уже оказался выше его сил! 🥴", m.Author.ID, roll/3),
		fmt.Sprintf("<@%s> был слишком уверен в себе и перепил. %d литров — не шутка! 🤢", m.Author.ID, roll/3),
		fmt.Sprintf("<@%s> свалился под весом %d литров пива и отправляется в тайм-аут! 😵", m.Author.ID, roll/3),
	}

	if roll >= chance {
		successMessage := successMessages[rand.IntN(len(successMessages))]
		s.ChannelMessageSend(m.ChannelID, successMessage)

		if amount == 40 {
			s.ChannelMessageSend(m.ChannelID, "https://media.giphy.com/media/gPbhyNB9Vpde0/giphy.gif?cid=790b7611u68bncsm51wuk8e8whzjalqm9r0gi2mpqxaiqpr3&ep=v1_gifs_search&rid=giphy.gif&ct=g")
			time.Sleep(1 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "Невероятно!!!!!! 40 литров!!!!!!!! Ты, наверное, из пивного королевства! 🍻👑")
			time.Sleep(5 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "https://media.giphy.com/media/Zw3oBUuOlDJ3W/giphy.gif?cid=790b7611rwi3azyed54indak41tqabn2pga0fbqr5da2z44d&ep=v1_gifs_search&rid=giphy.gif&ct=g")
			return
		}

		if rand.IntN(100) < 50 { // 50% шанс показать GIF
			gif := gifs[rand.IntN(len(gifs))]
			s.ChannelMessageSend(m.ChannelID, gif)
		}

	} else {
		failureMessage := failureMessages[rand.IntN(len(failureMessages))]
		muteDuration := getMuteDuration(amount)
		muteUntil := time.Now().Add(muteDuration)

		err = s.GuildMemberTimeout(m.GuildID, m.Author.ID, &muteUntil)
		if err != nil {
			fmt.Println("Error muting member:", err)
		}

		s.ChannelMessageSend(m.ChannelID, failureMessage)

		if rand.IntN(100) < 50 { // 50% шанс показать GIF
			gif := gifs[rand.IntN(len(gifs))]
			s.ChannelMessageSend(m.ChannelID, gif)
		}
	}
}

func getMuteDuration(amount int) time.Duration {
	switch {
	case amount >= 40:
		return 10 * time.Minute
	case amount >= 30:
		return 5 * time.Minute
	case amount >= 20:
		return 3 * time.Minute
	case amount >= 10:
		return 2 * time.Minute
	default:
		return 1 * time.Minute
	}
}

func poll(session *discordgo.Session, m *discordgo.MessageCreate) {
	users, err := session.GuildMembers(m.GuildID, "", 300)
	if err != nil {
		fmt.Println("error getting users,", err)
		return
	}

	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	users = users[:3]

	poll := &discordgo.MessageEmbed{
		Title: "Кто сегодня писька??? 🤔🤔🤔",
		Color: 0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "1",
				Value:  getNick(users[0]),
				Inline: true,
			},
			{
				Name:   "2",
				Value:  getNick(users[1]),
				Inline: true,
			},
			{
				Name:   "3",
				Value:  getNick(users[2]),
				Inline: true,
			},
		},
	}

	pollMessage, err := session.ChannelMessageSendEmbed(m.ChannelID, poll)
	if err != nil {
		fmt.Println("error sending poll,", err)
		return
	}

	reactions := []string{"1️⃣", "2️⃣", "3️⃣"}

	for _, v := range reactions {
		err := session.MessageReactionAdd(pollMessage.ChannelID, pollMessage.ID, v)
		if err != nil {
			fmt.Println("error adding reaction,", err)
			return
		}
	}

	time.Sleep(30 * time.Minute)

	pollResults, err := session.ChannelMessage(pollMessage.ChannelID, pollMessage.ID)
	if err != nil {
		fmt.Println("error getting poll results,", err)
		return
	}

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

	_, err = session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Поздравляем, <@%s>, ты сегодня писька! 🎉🎉🎉", winner.User.ID))
	if err != nil {
		fmt.Println("error congratulating the winner,", err)
	}

}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	Token = os.Getenv("TOKEN")
	if Token == "" {
		panic("You need to set the TOKEN environment variable.")
	}

	weatherApiKey = os.Getenv("WEATHER_API_KEY")
	if weatherApiKey == "" {
		panic("You need to set the WEATHER_API_KEY environment variable.")
	}
}

func getNick(member *discordgo.Member) string {
	if member.Nick == "" {
		return member.User.Username
	}
	return member.Nick
}

func piskaMessage(users []string) string {
	var message string
	message += "🤔🤔🤔"
	for _, user := range users {
		// #nosec G404 -- This is a false positive
		piskaProc := rand.IntN(101)
		switch {
		case piskaProc == 0:
			message += fmt.Sprintf("\nИзвини, <@%s>, но ты совсем не писька (0%%), приходи когда описюнеешь", user)
		case piskaProc == 100:
			message += fmt.Sprintf("\n<@%s>, ты просто прекрасная писька на ВСЕ 100%%", user)
		case piskaProc >= 50:
			message += fmt.Sprintf("\n<@%s> писька на %d%%, молодец, так держать!", user, piskaProc)
		default:
			message += fmt.Sprintf("\n<@%s> писька на %d%%, но нужно еще вырасти", user, piskaProc)
		}
	}
	return message
}

func penisCommand() string {
	size := rand.IntN(30) + 1
	shaft := strings.Repeat("=", size)
	penis := fmt.Sprintf("8%s>", shaft)

	var message string
	switch size {
	case 1:
		message = "Обладатель микроскопического стручка! Не грусти, бро, зато ты король клитора!"
	case 30:
		message = "Святые угодники! У тебя там баобаб вырос? Поздравляем, теперь ты главный калибр эскадры!"
	default:
		message = fmt.Sprintf("Размер: %d см", size)
	}

	return fmt.Sprintf("```\n%s\n```\n%s", penis, message)
}

func boobsCommand() string {
	// Генерируем размер груди от 0 до 20
	size := rand.IntN(21)

	// Строим визуальное представление груди
	leftBoob := "(" + strings.Repeat(" ", size/4) + "." + strings.Repeat(" ", size/4) + ")"
	rightBoob := "(" + strings.Repeat(" ", size/4) + "." + strings.Repeat(" ", size/4) + ")"
	boobs := leftBoob + " " + rightBoob

	var message string
	switch size {
	case 0:
		message = "Ноль? Не беда! Главное — душевная глубина."
	case 20:
		message = "Это не грудь, это просто обоюдоострый инструмент соблазна!"
	case 1, 2:
		message = "Ну, это почти незаметно, но всегда можно подсунуть носок!"
	case 3, 4, 5:
		message = "Мал да удал! Кто-то явно фанат японских аниме."
	case 6, 7, 8:
		message = "Пока что скромно, но всё впереди. Кстати, push-up никто не отменял!"
	case 9, 10, 11:
		message = "Средний размер — идеальный баланс! Завидую тому, кто будет с этим работать."
	case 12, 13, 14:
		message = "Ого, это уже что-то серьезное. Тебе точно нужно больше топов и меньше gravity."
	case 15, 16, 17:
		message = "Вот это да, пышные формы! С такой грудью можно смело идти на кастинг к Victoria's Secret."
	case 18, 19:
		message = "Невероятно! Это не просто размер — это целое событие! Скоро тебе нужен будет поддерживающий персонал."
	default:
		message = fmt.Sprintf("Размер: %d ", size)
	}

	return fmt.Sprintf("```\n%s\n```\n%s", boobs, message)
}

func gayMessage(s *discordgo.Session, m *discordgo.MessageCreate, user string) {
	// Генерация процента гейства
	gayProc := rand.IntN(101)
	var result string
	var rainbowCount int

	switch {
	case gayProc == 0:
		result = fmt.Sprintf("<@%s>, у тебя пока 0%% GaYства. Не сдавайся! 🥺", user)
		rainbowCount = 1
	case gayProc == 100:
		result = fmt.Sprintf("<@%s>, ты просто совершенство! 400%% GaYства! 🌈✨🦄💖🌟", user)
		rainbowCount = 20
	case gayProc >= 61:
		result = fmt.Sprintf("<@%s>, у тебя %d%% GaYства! Держись, радужный воин! 💃✨", user, gayProc)
		rainbowCount = 15
	case gayProc >= 21:
		result = fmt.Sprintf("<@%s>, у тебя %d%% GaYства. Попробуй танцевать под Lady Gaga! 💃🎶", user, gayProc)
		rainbowCount = 10
	default:
		result = fmt.Sprintf("<@%s>, у тебя %d%% GaYства. Нужно больше блесток и радуг! ✨🌈", user, gayProc)
		rainbowCount = 5
	}

	messageContent := fmt.Sprint(strings.Repeat("🌈", rainbowCount), "\n", result)

	sentMessage, err := s.ChannelMessageSend(m.ChannelID, messageContent)
	if err != nil {
		fmt.Println("error sending message:", err)
		return
	}

	var reactions []string
	switch {
	case gayProc == 0:
		reactions = []string{}
	case gayProc == 100:
		reactions = []string{"🌈", "✨", "🦄", "💖"}
	case gayProc >= 61:
		reactions = []string{"🌈", "✨", "🦄"}
	case gayProc >= 21:
		reactions = []string{"🌈", "✨"}
	default:
		reactions = []string{"🌈"}
	}

	for _, emoji := range reactions {
		err := s.MessageReactionAdd(m.ChannelID, sentMessage.ID, emoji)
		if err != nil {
			fmt.Println("error adding reaction:", err)
		}
	}
}

func main() {

	logging.InitLogger()
	mainLogger := logging.GetLogger("main")
	weather.InitWeatherLogger()

	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		mainLogger.Error("Error creating Discord session:", err)
		return
	}

	weatherLogger := logging.GetLogger("weather")
	weatherApiBaseUrl := "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"
	weatherCommandRe := regexp.MustCompile(`^!(weather|погода)(?:\s+([\p{L}\s]+))?(?:\s+(\d+))?$`)

	quote := quotes.New()

	session.Identify.Intents = discordgo.IntentsGuildMessages

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {

		weatherMathes := weatherCommandRe.FindStringSubmatch(m.Content)
		if len(weatherMathes) > 0 {
			err := weather.HandleWeatherMessage(s, m, weatherApiKey, weatherApiBaseUrl, weatherMathes)
			if err != nil {
				weatherLogger.Error("Error handling weather message:", err)
			}
		}

		if m.Author.ID == s.State.User.ID {
			return
		}
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

		if m.Content == "!гонка" {
			handleRaceCommand(s, m)
		} else if m.Content == "!го" {
			handleJoinRaceCommand(s, m)
		} else if strings.HasPrefix(m.Content, "!пиво") {
			handleBeerCommand(s, m)
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
				mainLogger.Error("error reacting to message,", err)
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
				mainLogger.Error("error reacting to message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "легион") {
			for _, v := range legionEmojis {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, v)
				time.Sleep(100 * time.Millisecond)
				if err != nil {
					mainLogger.Error("error reacting to message,", err)
				}
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "ковен") || strings.Contains(strings.ToLower(m.Content), "сестры") || strings.Contains(strings.ToLower(m.Content), "сёстры") {
			for _, v := range covenEmojis {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, v)
				time.Sleep(100 * time.Millisecond)
				if err != nil {
					mainLogger.Error("error reacting to message,", err)
				}
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "спасибо") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Это тебе спасибо! 😎😎😎", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "бобр") || strings.Contains(strings.ToLower(m.Content), "бобер") || strings.Contains(strings.ToLower(m.Content), "курва") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Kurwa bóbr. Ja pierdolę, Jakie bydlę jebane 🦫🦫🦫", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "привет") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Привет!", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "пиф") && strings.ContainsAny(strings.ToLower(m.Content), "паф") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Пиф-паф!🔫🔫🔫", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		} else if strings.Contains(strings.ToLower(m.Content), "pif") && strings.ContainsAny(strings.ToLower(m.Content), "paf") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Pif-paf!🔫🔫🔫", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "алкаш") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Эй мальчик, давай обмен,я же вижу что ты алкаш (c) Чайок", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "дед инсайд") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Глисты наконец-то померли?", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я гей") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я лесбиянка") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я би") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я натурал") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я иногда тоже!", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "понедельник") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "День тяжелый 😵‍💫", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		for _, v := range sickMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				_, err := s.ChannelMessageSendReply(m.ChannelID, "Скорее выздоравливай и больше не болей! 😍", m.Reference())
				if err != nil {
					mainLogger.Error("error sending message,", err)
				}
			}
		}

		for _, v := range phasmaMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, "👻")
				if err != nil {
					mainLogger.Error("error reacting to message,", err)
				}
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!пенис") {
			user := m.Author.ID
			if len(m.Mentions) != 0 {
				member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
				if err == nil {
					user = member.User.ID
				}
			}

			response := penisCommand()
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>\n%s", user, response), m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!бубс") {
			// Выбираем пользователя: упомянутого или автора
			user := m.Author.ID
			if len(m.Mentions) != 0 {
				member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
				if err == nil {
					user = member.User.ID
				}
			}

			response := boobsCommand()
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>\n%s", user, response), m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "полчаса") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "полчаса, полчаса - не вопрос. Не ответ полчаса, полчаса (c) Чайок", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "керамика") {
			customEmoji := "<:PotFriend:1271815662695743590>"
			response := fmt.Sprintf("внезапная %s перекличка %s ебучих %s керамических %s изделий %s внезапная %s перекличка %s ебучих %s керамических %s изделий %s",
				customEmoji, customEmoji, customEmoji, customEmoji, customEmoji,
				customEmoji, customEmoji, customEmoji, customEmoji, customEmoji)
			_, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if m.Content == "!голосование" {
			go poll(s, m)
		}

		if strings.Contains(strings.ToLower(m.Content), "!quote") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, quote.GetRandom(), m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "!academia") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, quote.GetRandomAcademia(), m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!медведь") {
			user := m.Author.ID
			if len(m.Mentions) != 0 {
				//#nosec G404 -- This is a false positive
				if rand.IntN(10) == 0 {
					_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, кажется медведь прямо сейчас завалит тебя 🐻🐻🐻", user), m.Reference())
					if err != nil {
						mainLogger.Error("error sending message,", err)
					}
					return
				}
				member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
				if err == nil {
					user = member.User.ID
				}
			}

			//#nosec G404 -- This is a false positive
			medvedProc := rand.IntN(101)
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, завалишь медведя с %d%% вероятностью 🐻", user, medvedProc), m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!ролл") || strings.HasPrefix(strings.ToLower(m.Content), "!d20") {
			user := m.Author.ID
			//#nosec G404 -- This is a false positive
			roll := rand.IntN(20) + 1
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, ты выкинул %d", user, roll), m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!писька") {
			user := m.Author.ID
			if len(m.Mentions) != 0 {
				//#nosec G404 -- This is a false positive
				if rand.IntN(10) == 0 {
					_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, а вот и нет, писька это ты!!!", user), m.Reference())
					if err != nil {
						mainLogger.Error("error sending message,", err)
					}
					return
				}
				member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID)
				if err == nil {
					user = member.User.ID
				}
			}

			//#nosec G404 -- This is a false positive
			piskaProc := rand.IntN(101)

			if piskaProc == 100 {
				_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, ты просто прекрасная писька на ВСЕ 100%%", user), m.Reference())
				if err != nil {
					mainLogger.Error("error sending message,", err)
				}
				return
			}

			if piskaProc == 0 {
				_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Извини, <@%s>, но ты совсем не писька (0%%), приходи когда описюнеешь", user), m.Reference())
				if err != nil {
					mainLogger.Error("error sending message,", err)
				}
				return
			}

			//#nosec G404 -- This is a false positive
			if rand.IntN(2) == 0 && piskaProc > 50 {
				//#nosec G404 -- This is a false positive
				_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> настоящая писька на %d%%, вот тебе цитата: %s", user, piskaProc, quotesPublic[rand.IntN(len(quotesPublic))]), m.Reference())
				if err != nil {
					mainLogger.Error("error sending message,", err)
				}
				return
			}

			if piskaProc > 50 {
				_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> писька на %d%%, молодец, так держать!", user, piskaProc), m.Reference())
				if err != nil {
					mainLogger.Error("error sending message,", err)
				}
				return
			}

			//#nosec G404 -- This is a false positive
			_, err = s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> писька на %d%%, но нужно еще вырасти!", user, piskaProc), m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!гей") {
			var userID string

			if len(m.Mentions) > 0 {
				userID = m.Mentions[0].ID
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Пожалуйста, упомяни пользователя для проверки гейства!")
				if err != nil {
					mainLogger.Error("error sending message:", err)
				}
				return
			}

			if rand.IntN(10) == 0 {
				userID = m.Author.ID
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, а может быть ты, моя голубая луна???!!!", userID))
				if err != nil {
					mainLogger.Error("error sending message:", err)
				}
			}

			gayMessage(s, m, userID)
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!письки") {
			user := m.Author.ID
			users := make([]string, 0)
			if len(m.Mentions) != 0 {
				// nosec G404 -- This is a false positive
				if rand.IntN(10) == 0 {
					_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, а вот и нет, писька это ты!!!", user), m.Reference())
					if err != nil {
						mainLogger.Error("error sending message,", err)
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
				mainLogger.Error("error sending message,", err)
			}
			return

		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!анекдот") {
			joke, err := jokes.GetJoke()
			if err != nil {
				mainLogger.Error("error getting joke,", err)
				return
			}
			_, err = s.ChannelMessageSendReply(m.ChannelID, joke, m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		for _, v := range potterMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, "🧙")
				if err != nil {
					mainLogger.Error("error reacting message,", err)
				}
			}
		}

		for _, v := range valorantMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, "🔥")
				if err != nil {
					mainLogger.Error("error reacting message,", err)
				}
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я писюн") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже писюн!!!", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я писька") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже писька!!!", m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!шар") {
			//#nosec G404 -- This is a false positive
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Мой ответ: %s", magicBallMessages[rand.IntN(len(magicBallMessages))]), m.Reference())
			if err != nil {
				mainLogger.Error("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "демон") {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "👹")
			if err != nil {
				mainLogger.Error("error reacting message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "клоун") {
			err := s.MessageReactionAdd(m.ChannelID, m.ID, "🤡")
			if err != nil {
				mainLogger.Error("error reacting message,", err)
			}
		}

	})

	err = session.Open()
	if err != nil {
		mainLogger.Error("error opening connection,", err)
		return
	}

	mainLogger.Info("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan struct{})

}
