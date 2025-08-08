package bot

import (
	"fmt"
	rand "math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/texts"
)

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

	specialOdds := map[int]float64{1: 99.0, 7: 70.0, 13: 25.0, 40: 1.5}
	var chance int
	const rollMax = 120
	if specialChance, ok := specialOdds[amount]; ok {
		chance = rollMax - int(float64(rollMax)*(specialChance/100.0))
	} else {
		baseMultiplier := 3.0
		switch {
		case amount > 30:
			baseMultiplier = 3.3
		case amount > 20:
			baseMultiplier = 3.2
		case amount > 10:
			baseMultiplier = 3.1
		}
		chance = int(float64(amount) * baseMultiplier)
	}
	if chance < 1 {
		chance = 1
	} else if chance > rollMax {
		chance = rollMax
	}

	roll := rand.IntN(rollMax) + 1
	success := roll >= chance

	if success {
		successMessages := []string{
			fmt.Sprintf("<@%s> смог осилить %d литров пива! 🍺", m.Author.ID, amount),
			fmt.Sprintf("<@%s> успешно справился с %d литрами! Это достойно уважения! 🍻", m.Author.ID, amount),
			fmt.Sprintf("<@%s> выпил %d литров, пивной монстр на свободе! 🍻🦹", m.Author.ID, amount),
			fmt.Sprintf("<@%s> залпом поглотил %d литров и выглядит, как чемпион! 🏆", m.Author.ID, amount),
			fmt.Sprintf("<@%s> выпил %d литров пива и готов к новым свершениям! 🍻🚀", m.Author.ID, amount),
			fmt.Sprintf("<@%s> справился с %d литрами пива! Не плохо! 🍺", m.Author.ID, amount),
		}
		successMessage := successMessages[rand.IntN(len(successMessages))]
		s.ChannelMessageSend(m.ChannelID, successMessage)
		if amount == 40 {
			s.ChannelMessageSend(m.ChannelID, "https://media.giphy.com/media/gPbhyNB9Vpde0/giphy.gif?cid=790b7611u68bncsm51wuk8e8whzjalqm9r0gi2mpqxaiqpr3&ep=v1_gifs_search&rid=giphy.gif&ct=g")
			time.Sleep(1 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "Невероятно!!!!!! 40 литров!!!!!!!! Ты, наверное, из пивного королевства! 🍻👑")
			time.Sleep(3 * time.Second)
			s.ChannelMessageSend(m.ChannelID, "https://media.giphy.com/media/Zw3oBUuOlDJ3W/giphy.gif?cid=790b7611rwi3azyed54indak41tqabn2pga0fbqr5da2z44d&ep=v1_gifs_search&rid=giphy.gif&ct=g")
			return
		}
		if rand.IntN(100) < 50 {
			gifs := []string{
				"https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExZGt0bGtuZHphOTg1bHo2b3BwYW5sZG00Y3U1MHN6amY5aGl2aDdodSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/lTGLOH7ml3poQ6JoFg/giphy.gif",
				"https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExYzJpaTcxZTYzeW1zN3Jhc2VxbjR0YndqZWVjb3Btb3AxZzJuZDk0aSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/yB9T6y9k1GQSkZZp9v/giphy.gif",
				"https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExYWw5NXNyaDQ0Ymh0ejg5NzgzY3Y2cm5ndXllaHVpdTJrZ2tiYmFwaSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/xQG0wbo9A3WHC/giphy-downsized-large.gif",
				"https://media3.giphy.com/media/v1.Y2lkPTc5MGI3NjExMjd2ZTVsZmtvd2F2aTR1ZXJ5ZG5yM2EybzV5OWltMmJzdWttcWsxMyZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/YooCD0Y2fw1C6VFBwl/giphy.gif",
				"https://media.giphy.com/media/26tP21xUQnOCIIoFi/giphy.gif?cid=790b7611iyvxpdr8q647v1zbgay9muul2t1u1y0vjyzm4fg8&ep=v1_gifs_search&rid=giphy.gif&ct=g",
			}
			gif := gifs[rand.IntN(len(gifs))]
			s.ChannelMessageSend(m.ChannelID, gif)
		}
		return
	}

	drankAmount := roll / 3
	if drankAmount < 1 {
		drankAmount = 1
	}
	var muteDuration time.Duration
	switch {
	case amount >= 40:
		muteDuration = 10 * time.Minute
	case amount >= 30:
		muteDuration = 5 * time.Minute
	case amount >= 20:
		muteDuration = 3 * time.Minute
	case amount >= 10:
		muteDuration = 2 * time.Minute
	default:
		muteDuration = 1 * time.Minute
	}
	muteUntil := time.Now().Add(muteDuration)
	_ = s.GuildMemberTimeout(m.GuildID, m.Author.ID, &muteUntil)
	var durationText string
	if muteDuration.Minutes() < 1 {
		durationText = fmt.Sprintf("%d секунд", int(muteDuration.Seconds()))
	} else {
		durationText = fmt.Sprintf("%d минут", int(muteDuration.Minutes()))
	}
	failureMessages := []string{
		fmt.Sprintf("<@%s> не смог осилить даже %d литров пива и облевал весь пол! Кто это убирать будет?! 🤢🤮", m.Author.ID, drankAmount),
		fmt.Sprintf("<@%s> попытался выпить %d литр, но потерпел неудачу и свалился под стол! 😵", m.Author.ID, drankAmount),
		fmt.Sprintf("<@%s> проиграл борьбу на %d литрах пива и отправляется в бан на %s! 😴", m.Author.ID, drankAmount, durationText),
		fmt.Sprintf("<@%s> взял на себя слишком много! %d литр пива уже оказался выше его сил! 🥴", m.Author.ID, drankAmount),
		fmt.Sprintf("<@%s> был слишком уверен в себе и перепил. %d литров — не шутка! 🤢", m.Author.ID, drankAmount),
	}
	failureMessage := failureMessages[rand.IntN(len(failureMessages))]
	s.ChannelMessageSend(m.ChannelID, failureMessage)
	if rand.IntN(100) < 50 {
		gifs := []string{
			"https://media.giphy.com/media/6S9cWuMVtjfPz1GYqK/giphy.gif?cid=ecf05e47f7cas4uugmw9k7whhb5fx06n7zlpzwwgcjw482n4&ep=v1_gifs_search&rid=giphy.gif&ct=g",
			"https://media.giphy.com/media/zrj0yPfw3kGTS/giphy.gif?cid=ecf05e47f7cas4uugmw9k7whhb5fx06n7zlpzwwgcjw482n4&ep=v1_gifs_search&rid=giphy.gif&ct=g",
			"https://media.giphy.com/media/2CvuL80h6YTbq/giphy.gif?cid=ecf05e47f7cas4uugmw9k7whhb5fx06n7zlpzwwgcjw482n4&ep=v1_gifs_search&rid=giphy.gif&ct=g",
			"https://media.giphy.com/media/RqbkeCZGgipSo/giphy.gif?cid=ecf05e47afa5rztdshpog9jf8m2ecm4ecw8pn38ihu8qxypn&ep=v1_gifs_search&rid=giphy.gif&ct=g",
		}
		gif := gifs[rand.IntN(len(gifs))]
		s.ChannelMessageSend(m.ChannelID, gif)
	}
}

func handlePenis(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author.ID
	if len(m.Mentions) != 0 {
		if member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID); err == nil {
			user = member.User.ID
		}
	}
	response := penisCommand()
	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>\n%s", user, response), m.Reference())
}

func handleBoobs(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author.ID
	if len(m.Mentions) != 0 {
		if member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID); err == nil {
			user = member.User.ID
		}
	}
	response := boobsCommand()
	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>\n%s", user, response), m.Reference())
}

func handleGay(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Mentions) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Пожалуйста, упомяни пользователя для проверки гейства!")
		return
	}
	userID := m.Mentions[0].ID
	if rand.IntN(10) == 0 {
		userID = m.Author.ID
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, а может быть ты, моя голубая луна???!!!", userID))
	}
	gayMessage(s, m, userID)
}

func handlePiskiMultiple(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author.ID
	users := make([]string, 0)
	if len(m.Mentions) != 0 {
		if rand.IntN(10) == 0 {
			s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, а вот и нет, писька это ты!!!", user), m.Reference())
			return
		}
		for _, mention := range m.Mentions {
			if member, err := s.GuildMember(m.GuildID, mention.ID); err == nil {
				users = append(users, member.User.ID)
			}
		}
	}
	s.ChannelMessageSendReply(m.ChannelID, piskaMessage(users), m.Reference())
}

func handleRoll(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author.ID
	roll := rand.IntN(20) + 1
	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, ты выкинул %d", user, roll), m.Reference())
}

func handleMagicBall(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Мой ответ: %s", texts.MagicBallMessages[rand.IntN(len(texts.MagicBallMessages))]), m.Reference())
}

func handlePiskaSingle(s *discordgo.Session, m *discordgo.MessageCreate) {
	user := m.Author.ID
	if len(m.Mentions) != 0 {
		if rand.IntN(10) == 0 {
			s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, а вот и нет, писька это ты!!!", user), m.Reference())
			return
		}
		if member, err := s.GuildMember(m.GuildID, m.Mentions[0].ID); err == nil {
			user = member.User.ID
		}
	}
	piskaProc := rand.IntN(101)
	if piskaProc == 100 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s>, ты просто прекрасная писька на ВСЕ 100%%", user), m.Reference())
		return
	}
	if piskaProc == 0 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Извини, <@%s>, но ты совсем не писька (0%%), приходи когда описюнеешь", user), m.Reference())
		return
	}
	if rand.IntN(2) == 0 && piskaProc > 50 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> настоящая писька на %d%%, вот тебе цитата: %s", user, piskaProc, texts.QuotesPublic[rand.IntN(len(texts.QuotesPublic))]), m.Reference())
		return
	}
	if piskaProc > 50 {
		s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> писька на %d%%, молодец, так держать!", user, piskaProc), m.Reference())
		return
	}
	s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("<@%s> писька на %d%%, но нужно еще вырасти!", user, piskaProc), m.Reference())
}
