package bot

import (
	"fmt"
	rand "math/rand/v2"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func addGuildEmojiOrFallback(s *discordgo.Session, m *discordgo.MessageCreate, emojiID, fallback string) {
	emoji, err := s.GuildEmoji(m.GuildID, emojiID)
	if err != nil {
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, fallback)
		return
	}
	_ = s.MessageReactionAdd(m.ChannelID, m.ID, emoji.APIName())
}

func reactMany(s *discordgo.Session, m *discordgo.MessageCreate, emojis []string) {
	for _, v := range emojis {
		_ = s.MessageReactionAdd(m.ChannelID, m.ID, v)
		time.Sleep(100 * time.Millisecond)
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
	size := rand.IntN(21)
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
	messageContent := strings.Repeat("🌈", rainbowCount) + "\n" + result
	sentMessage, err := s.ChannelMessageSend(m.ChannelID, messageContent)
	if err != nil {
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
		_ = s.MessageReactionAdd(m.ChannelID, sentMessage.ID, emoji)
	}
}

// poll moved from main
func poll(session *discordgo.Session, m *discordgo.MessageCreate) {
	users, err := session.GuildMembers(m.GuildID, "", 300)
	if err != nil || len(users) < 3 {
		return
	}
	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	users = users[:3]
	poll := &discordgo.MessageEmbed{
		Title:  "Кто сегодня писька??? 🤔🤔🤔",
		Color:  0x00ff00,
		Fields: []*discordgo.MessageEmbedField{{Name: "1", Value: getNick(users[0]), Inline: true}, {Name: "2", Value: getNick(users[1]), Inline: true}, {Name: "3", Value: getNick(users[2]), Inline: true}},
	}
	pollMessage, err := session.ChannelMessageSendEmbed(m.ChannelID, poll)
	if err != nil {
		return
	}
	reactions := []string{"1️⃣", "2️⃣", "3️⃣"}
	for _, v := range reactions {
		_ = session.MessageReactionAdd(pollMessage.ChannelID, pollMessage.ID, v)
	}
	time.Sleep(30 * time.Minute)
	pollResults, err := session.ChannelMessage(pollMessage.ChannelID, pollMessage.ID)
	if err != nil {
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
	var winner *discordgo.Member
	switch mostVotedOption {
	case "1️⃣":
		winner = users[0]
	case "2️⃣":
		winner = users[1]
	case "3️⃣":
		winner = users[2]
	}
	if winner != nil {
		_, _ = session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Поздравляем, <@%s>, ты сегодня писька! 🎉🎉🎉", winner.User.ID))
	}
}
