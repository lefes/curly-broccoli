package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token string = ""
)

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
}

func main() {
	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	morningMessages := []string{
		"доброе утро",
		"доброго утра",
		"добрый день",
		"добрый вечер",
		"доброй ночи",
		"утро",
		"утречко",
		"день",
		"днечко",
		"вечер",
		"вечечко",
		"ночь",
		"ночечко",
		"morning",
		"evening",
		"night",
		"day",
		"good morning",
		"good evening",
		"good night",
		"good day",
		"проснул",
		"открыл глаза",
	}

	spokiMessages := []string{
		"спок",
		"сладких снов",
		"спокойной ночи",
		"до завтра",
		"спать",
		"дрем",
		"кемар",
		"сплю",
		"пока",
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

	legionEmojis := []string{"🇱", "🇪", "🇬", "🇮", "🇴", "🇳"}

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

		// Checking on spasibo message
		if strings.Contains(strings.ToLower(m.Content), "спасибо") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Это тебе спасибо! 😎😎😎", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "привет" message
		if strings.Contains(strings.ToLower(m.Content), "привет") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Привет, друг!", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
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

		// Checking on "понедельник" message
		if strings.Contains(strings.ToLower(m.Content), "понедельник") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "День тяжелый 😵‍💫", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "заболел" message
		if strings.Contains(strings.ToLower(m.Content), "заболел") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "Скорее выздоравливай и больше не болей! 😍", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
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

	})

	err = session.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan struct{})

}
