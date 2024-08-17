package main

import (
	"fmt"
	rand "math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/lefes/curly-broccoli/jokes"
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

	initialMessage := "🏁 **Гонка начинается!** 🏁\n\n"
	for id := range raceParticipants {
		initialMessage += fmt.Sprintf("<@%s> %s на старте 🏎️💨\n", id, raceParticipants[id])
	}
	raceMessage, err := s.ChannelMessageSend(m.ChannelID, initialMessage)
	if err != nil {
		fmt.Println("error sending message:", err)
		return
	}

	// Инициализация трека
	raceTrack := make(map[string]int)
	for id := range raceParticipants {
		raceTrack[id] = 0
	}

	// Запуск гонки
	winner := ""
	trackLength := 20
	for winner == "" {
		time.Sleep(1 * time.Second)
		raceStatus := "```🏁 Гонка в процессе 🏁\n\n"
		for id, emoji := range raceParticipants {
			raceTrack[id] += rand.IntN(3)
			if raceTrack[id] >= trackLength {
				raceTrack[id] = trackLength
				winner = id
				break
			}
			progress := strings.Repeat("—", raceTrack[id])
			emptySpace := strings.Repeat("—", trackLength-raceTrack[id])
			raceStatus += fmt.Sprintf("🚦 |%s%s%s|\n", progress, emoji, emptySpace)
		}
		raceStatus += "```"

		_, err := s.ChannelMessageEdit(m.ChannelID, raceMessage.ID, raceStatus)
		if err != nil {
			fmt.Println("error editing message:", err)
			return
		}
	}

	finalMessage := fmt.Sprintf("🎉 **Победитель гонки:** <@%s> %s! Поздравляем! 🏆🎉", winner, raceParticipants[winner])
	s.ChannelMessageSend(m.ChannelID, finalMessage)

	raceInProgress = false
	raceParticipants = make(map[string]string)
}

func handleBeerCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(m.Content, " ")
	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Использование: !пиво <число от 1 до 20>")
		return
	}

	amount, err := strconv.Atoi(args[1])
	if err != nil || amount < 1 || amount > 20 {
		s.ChannelMessageSend(m.ChannelID, "Пожалуйста, введите число от 1 до 20.")
		return
	}

	successChance := 100 - (amount * 5)
	if successChance < 5 {
		successChance = 5
	}

	roll := rand.IntN(100) + 1

	if roll <= successChance {
		var successMessage string
		if amount == 20 {
			successMessage = fmt.Sprintf("<@%s> выпил %d литров пива и остался жив?! 🎉🍻\n\n", m.Author.ID, amount)
			s.ChannelMessageSend(m.ChannelID, successMessage)
			s.ChannelMessageSend(m.ChannelID, "https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExejN4bjU1cTc1NDRodXU1OGd1NTExNTZheXRwOTdkaHNycWwyMTdtZyZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/qiSGGu0d2Dgac/giphy.gif")
		} else {
			successMessage = fmt.Sprintf("<@%s> успешно выпил %d литров пива! 🍺\n\n", m.Author.ID, amount)
			s.ChannelMessageSend(m.ChannelID, successMessage)
		}
	} else {
		var failureMessage string
		if amount == 20 {
			failureMessage = fmt.Sprintf("<@%s> не смог осилить %d литров пива и отправляется в бессознательное состояние на 5 минут! 🍺😴\n\n", m.Author.ID, amount)
			s.ChannelMessageSend(m.ChannelID, failureMessage)
			s.ChannelMessageSend(m.ChannelID, "https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExd3Rqb3NycG0xZTRqNHZoamgybmVmOGRvYTcyamViNGJ6ZGM0YjA1MSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/7bx7ZHokGnofm/giphy-downsized-large.gif")
		} else if amount >= 15 {
			failureMessage = fmt.Sprintf("<@%s> не осилил %d литров пива. Похоже, ты не подготовился к настоящей пьянке. Спокойной ночи на 5 минут! 🍺😴\n\n", m.Author.ID, amount)
			s.ChannelMessageSend(m.ChannelID, failureMessage)
		} else if amount >= 10 {
			failureMessage = fmt.Sprintf("<@%s> не смог выпить %d литров пива. Немного больше тренировки и получится! Мут на 5 минут! 🍻😴\n\n", m.Author.ID, amount)
			s.ChannelMessageSend(m.ChannelID, failureMessage)
		} else {
			failureMessage = fmt.Sprintf("<@%s> не справился с %d литрами пива. Надо больше тренироваться! Мут на 5 минут. 🍺😴\n\n", m.Author.ID, amount)
			s.ChannelMessageSend(m.ChannelID, failureMessage)
		}

		muteDuration := 5 * time.Minute
		muteUntil := time.Now().Add(muteDuration)
		err := s.GuildMemberTimeout(m.GuildID, m.Author.ID, &muteUntil)
		if err != nil {
			fmt.Println("Error muting member:", err)
			return
		}
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
	// Add reactions to the poll

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

	// Congratulate the winner
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

// Функция для команды пенис
func penisCommand(s *discordgo.Session, m *discordgo.MessageCreate) string {
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

func gayMessage(s *discordgo.Session, m *discordgo.MessageCreate, user string) {
	var message strings.Builder
	message.WriteString("\U0001F3F3\U0000FE0F\u200D\U0001F308\U0001F308\U0001F3F3\U0000FE0F\n")

	gayProc := rand.IntN(101)
	var result string

	switch {
	case gayProc == 0:
		result = fmt.Sprintf("<@%s>, у тебя пока 0%% GaYства. Не сдавайся! 🥺", user)
	case gayProc == 100:
		message.WriteString(strings.Repeat("🌈", 15))
		result = fmt.Sprintf("<@%s>, ты просто совершенство! 400%% GaYства! %s", user, strings.Join([]string{"🌈", "✨", "🦄", "💖", "🌟"}, " "))
	case gayProc >= 50:
		message.WriteString(strings.Repeat("🌈", 10))
		result = fmt.Sprintf("<@%s>, у тебя %d%% гейства! Держись, радужный воин! 💃✨", user, gayProc)
	default:
		message.WriteString(strings.Repeat("🌈", 5))
		result = fmt.Sprintf("<@%s>, у тебя %d%% гейства. Попробуй танцевать под Lady Gaga! 💃🎶", user, gayProc)
	}

	message.WriteString(result + "\n")

	message.WriteString(strings.Repeat("\U0001F308", 10) + "\n" + "\U0001F3F3\U0000FE0F\u200D\U0001F308\U0001F308\U0001F3F3\U0000FE0F")

	s.ChannelMessageSend(m.ChannelID, message.String())

	for _, emoji := range rainbowEmojis {
		time.Sleep(200 * time.Millisecond)
		s.MessageReactionAdd(m.ChannelID, m.ID, emoji)
	}

	if gayProc >= 50 {
		animatedMessage := "🌈 "
		for i := 0; i < 5; i++ {
			animatedMessage += strings.Repeat("🌈", i+1)
			_, err := s.ChannelMessageSend(m.ChannelID, animatedMessage)
			if err != nil {
				fmt.Println("error sending animated message:", err)
			}
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func main() {

	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	quote := quotes.New()

	session.Identify.Intents = discordgo.IntentsGuildMessages

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
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

		if strings.HasPrefix(m.Content, "!гонка") {
			handleRaceCommand(s, m)
		} else if strings.HasPrefix(m.Content, "!го") {
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

		if strings.Contains(strings.ToLower(m.Content), "легион") {
			for _, v := range legionEmojis {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, v)
				time.Sleep(100 * time.Millisecond)
				if err != nil {
					fmt.Println("error reacting to message,", err)
				}
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "ковен") || strings.Contains(strings.ToLower(m.Content), "сестры") || strings.Contains(strings.ToLower(m.Content), "сёстры") {
			for _, v := range covenEmojis {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, v)
				time.Sleep(100 * time.Millisecond)
				if err != nil {
					fmt.Println("error reacting to message,", err)
				}
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "спасибо") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Это тебе спасибо! 😎😎😎", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "бобр") || strings.Contains(strings.ToLower(m.Content), "бобер") || strings.Contains(strings.ToLower(m.Content), "курва") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Kurwa bóbr. Ja pierdolę, Jakie bydlę jebane 🦫🦫🦫", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "привет") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Привет!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

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

		if strings.Contains(strings.ToLower(m.Content), "алкаш") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Эй мальчик, давай обмен,я же вижу что ты алкаш (c) Чайок", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "дед инсайд") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Глисты наконец-то померли?", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я гей") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я лесбиянка") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я би") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "я натурал") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Я иногда тоже!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "понедельник") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "День тяжелый 😵‍💫", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		for _, v := range sickMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				_, err := s.ChannelMessageSendReply(m.ChannelID, "Скорее выздоравливай и больше не болей! 😍", m.Reference())
				if err != nil {
					fmt.Println("error sending message,", err)
				}
			}
		}

		for _, v := range phasmaMessages {
			if strings.Contains(strings.ToLower(m.Content), v) {
				err := s.MessageReactionAdd(m.ChannelID, m.ID, "👻")
				if err != nil {
					fmt.Println("error reacting to message,", err)
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

			response := penisCommand(s, m)
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>\n%s", user, response), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "полчаса") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "полчаса, полчаса - не вопрос. Не ответ полчаса, полчаса (c) Чайок", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.Contains(strings.ToLower(m.Content), "керамика") {
			customEmoji := "<:PotFriend:1271815662695743590>"
			response := fmt.Sprintf("внезапная %s перекличка %s ебучих %s керамических %s изделий %s внезапная %s перекличка %s ебучих %s керамических %s изделий %s",
				customEmoji, customEmoji, customEmoji, customEmoji, customEmoji,
				customEmoji, customEmoji, customEmoji, customEmoji, customEmoji)
			_, err := s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
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
				if rand.IntN(10) == 0 {
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
			medvedProc := rand.IntN(101)
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, завалишь медведя с %d%% вероятностью 🐻", user, medvedProc), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!ролл") || strings.HasPrefix(strings.ToLower(m.Content), "!d20") {
			user := m.Author.ID
			//#nosec G404 -- This is a false positive
			roll := rand.IntN(20) + 1
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, ты выкинул %d", user, roll), m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		if strings.HasPrefix(strings.ToLower(m.Content), "!писька") {
			user := m.Author.ID
			if len(m.Mentions) != 0 {
				//#nosec G404 -- This is a false positive
				if rand.IntN(10) == 0 {
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
			piskaProc := rand.IntN(101)

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
			if rand.IntN(2) == 0 && piskaProc > 50 {
				//#nosec G404 -- This is a false positive
				_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> настоящая писька на %d%%, вот тебе цитата: %s", user, piskaProc, quotesPublic[rand.IntN(len(quotesPublic))]), m.Reference())
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

		if strings.HasPrefix(strings.ToLower(m.Content), "!гей") {
			var userID string

			if len(m.Mentions) > 0 {
				userID = m.Mentions[0].ID
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Пожалуйста, упомяни пользователя для проверки гейства!")
				if err != nil {
					fmt.Println("error sending message:", err)
				}
				return
			}

			if rand.IntN(10) == 0 {
				userID = m.Author.ID
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, а может быть ты, моя голубая луна???!!!", userID))
				if err != nil {
					fmt.Println("error sending message:", err)
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
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Мой ответ: %s", magicBallMessages[rand.IntN(len(magicBallMessages))]), m.Reference())
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
