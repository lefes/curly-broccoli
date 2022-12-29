package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/lefes/curly-broccoli/quotes"
)

var (
	Token string = ""
)

type Quote interface {
	GetRandomAcademia() string
	GetRandom() string
}

func poll(session *discordgo.Session, m *discordgo.MessageCreate) {
	// Randomly create a poll with 3 options in the channel
	// Take 3 person from the channel
	users, err := session.GuildMembers(m.GuildID, "", 300)
	if err != nil {
		fmt.Println("error getting users,", err)
		return
	}

	// Get 3 random users
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(users), func(i, j int) { users[i], users[j] = users[j], users[i] })
	users = users[:3]

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
	_, err = session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Поздравляем, %s, ты сегодня писька! 🎉🎉🎉", getNick(winner)))
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
}

func getNick(member *discordgo.Member) string {
	if member.Nick == "" {
		return member.User.Username
	}
	return member.Nick
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

	sickMessages := []string{
		"заболел",
		"заболела",
		"заболело",
		"заболели",
		"болею",
		"болит",
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

		if strings.Contains(strings.ToLower(m.Content), "!писька") {
			rand.Seed(time.Now().UnixNano())
			//#nosec G404 -- This is a false positive
			_, err := s.ChannelMessageSendReply(m.ChannelID, fmt.Sprintf("Ты Писька на %d%%", rand.Intn(101)), m.Reference())
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
