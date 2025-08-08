package bot

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/internal/logging"
	"github.com/lefes/curly-broccoli/internal/texts"
	"github.com/lefes/curly-broccoli/pkg/jokes"
	"github.com/lefes/curly-broccoli/pkg/weather"
)

type Router struct {
	weatherClient *weather.Client
}

func NewRouter(weatherClient *weather.Client) *Router {
	return &Router{weatherClient: weatherClient}
}

func (r *Router) Handler() func(s *discordgo.Session, m *discordgo.MessageCreate) {
	log := logging.Get("router")
	weatherLogger := logging.Get("weather")
	weatherCommandRe := regexp.MustCompile(`^!(weather|погода)(?:\s+([\p{L}\s]+))?(?:\s+(\d+))?$`)

	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		// Weather
		if matches := weatherCommandRe.FindStringSubmatch(m.Content); len(matches) > 0 {
			if err := weather.HandleWeatherMessage(*r.weatherClient, s, m, matches); err != nil {
				weatherLogger.Error("handle weather:", err)
			}
			return
		}

		contentLower := strings.ToLower(m.Content)

		// Commands
		switch {
		case strings.HasPrefix(contentLower, "!пиво"):
			handleBeerCommand(s, m)
			return
		case strings.HasPrefix(contentLower, "!голосование"):
			go poll(s, m)
			return
		case strings.HasPrefix(contentLower, "!пенис"):
			handlePenis(s, m)
			return
		case strings.HasPrefix(contentLower, "!бубс"):
			handleBoobs(s, m)
			return
		case strings.HasPrefix(contentLower, "!гей"):
			handleGay(s, m)
			return
		case strings.HasPrefix(contentLower, "!письки"):
			handlePiskiMultiple(s, m)
			return
		case strings.HasPrefix(contentLower, "!писька"):
			handlePiskaSingle(s, m)
			return
		case strings.HasPrefix(contentLower, "!ролл") || strings.HasPrefix(contentLower, "!d20"):
			handleRoll(s, m)
			return
		case strings.HasPrefix(contentLower, "!шар"):
			handleMagicBall(s, m)
			return
		case strings.HasPrefix(contentLower, "!анекдот"):
			joke, err := jokes.GetJoke()
			if err == nil {
				s.ChannelMessageSendReply(m.ChannelID, joke, m.Reference())
			} else {
				log.Error("get joke:", err)
			}
			return
		}

		// Triggers
		if texts.ContainsAny(contentLower, texts.MorningMessages) {
			addGuildEmojiOrFallback(s, m, "1016631674106294353", "🫠")
		}
		if texts.ContainsAny(contentLower, texts.SpokiMessages) {
			addGuildEmojiOrFallback(s, m, "1016631826338566144", "😴")
		}
		if strings.Contains(contentLower, "легион") {
			reactMany(s, m, texts.LegionEmojis)
		}
		if strings.Contains(contentLower, "ковен") || strings.Contains(contentLower, "сестры") || strings.Contains(contentLower, "сёстры") {
			reactMany(s, m, texts.CovenEmojis)
		}
		if strings.Contains(contentLower, "спасибо") {
			s.ChannelMessageSendReply(m.ChannelID, "Это тебе спасибо! 😎😎😎", m.Reference())
		}
		if strings.Contains(contentLower, "бобр") || strings.Contains(contentLower, "бобер") || strings.Contains(contentLower, "курва") {
			s.ChannelMessageSendReply(m.ChannelID, "Kurwa bóbr. Ja pierdolę, Jakie bydlę jebane 🦫🦫🦫", m.Reference())
		}
		if strings.Contains(contentLower, "привет") {
			s.ChannelMessageSendReply(m.ChannelID, "Привет!", m.Reference())
		}
		if strings.Contains(contentLower, "пиф") && strings.ContainsAny(contentLower, "паф") {
			s.ChannelMessageSendReply(m.ChannelID, "Пиф-паф!🔫🔫🔫", m.Reference())
		} else if strings.Contains(contentLower, "pif") && strings.ContainsAny(contentLower, "paf") {
			s.ChannelMessageSendReply(m.ChannelID, "Pif-paf!🔫🔫🔫", m.Reference())
		}
		if strings.Contains(contentLower, "алкаш") {
			s.ChannelMessageSendReply(m.ChannelID, "Эй мальчик, давай обмен,я же вижу что ты алкаш (c) Чайок", m.Reference())
		}
		if strings.Contains(contentLower, "дед инсайд") {
			s.ChannelMessageSendReply(m.ChannelID, "Глисты наконец-то померли?", m.Reference())
		}
		if strings.Contains(contentLower, "я гей") {
			s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
		}
		if strings.Contains(contentLower, "я лесбиянка") {
			s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
		}
		if strings.Contains(contentLower, "я би") {
			s.ChannelMessageSendReply(m.ChannelID, "Я тоже!", m.Reference())
		}
		if strings.Contains(contentLower, "я натурал") {
			s.ChannelMessageSendReply(m.ChannelID, "Я иногда тоже!", m.Reference())
		}
		if strings.Contains(contentLower, "понедельник") {
			s.ChannelMessageSendReply(m.ChannelID, "День тяжелый 😵‍💫", m.Reference())
		}
		if texts.ContainsAny(contentLower, texts.SickMessages) {
			s.ChannelMessageSendReply(m.ChannelID, "Скорее выздоравливай и больше не болей! 😍", m.Reference())
		}
		if texts.ContainsAny(contentLower, texts.PhasmaMessages) {
			s.MessageReactionAdd(m.ChannelID, m.ID, "👻")
		}
		if texts.ContainsAny(contentLower, texts.PotterMessages) {
			s.MessageReactionAdd(m.ChannelID, m.ID, "🧙")
		}
		if texts.ContainsAny(contentLower, texts.ValorantMessages) {
			s.MessageReactionAdd(m.ChannelID, m.ID, "🔥")
		}
		if strings.Contains(contentLower, "я писюн") {
			s.ChannelMessageSendReply(m.ChannelID, "Я тоже писюн!!!", m.Reference())
		}
		if strings.Contains(contentLower, "я писька") {
			s.ChannelMessageSendReply(m.ChannelID, "Я тоже писька!!!", m.Reference())
		}
		if strings.Contains(contentLower, "демон") {
			s.MessageReactionAdd(m.ChannelID, m.ID, "👹")
		}
		if strings.Contains(contentLower, "клоун") {
			s.MessageReactionAdd(m.ChannelID, m.ID, "🤡")
		}
		if strings.Contains(contentLower, "керамика") {
			customEmoji := "<:PotFriend:1271815662695743590>"
			response := texts.CeramicsMessage(customEmoji)
			s.ChannelMessageSendReply(m.ChannelID, response, m.Reference())
		}
	}
}
