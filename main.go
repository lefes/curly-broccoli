package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/lefes/curly-broccoli/jokes"
	"github.com/lefes/curly-broccoli/quotes"
)

type DeathCounter struct {
	Count int `json:"count"`
	mu    sync.Mutex
}

var deathCounter DeathCounter
var Token string
var (
	raceParticipants = make(map[string]string)
	raceEmojis       = []string{"🐶", "🐱", "🐭", "🐹", "🐰", "🦊", "🐻", "🐼", "🐨", "🐯", "🦁", "🐮", "🐷", "🐸", "🐵", "🐔", "🐧", "🐦", "🐤", "🦆", "🦅", "🦉", "🦇", "🐺", "🐗", "🐴", "🦄", "🐝", "🐛", "🦋", "🐌", "🐞", "🐜", "🦟", "🦗", "🕷", "🦂", "🐢", "🐍", "🦎", "🦖", "🦕", "🐙", "🦑", "🦐", "🦞", "🦀", "🐡", "🐠", "🐟", "🐬", "🐳", "🐋", "🦈", "🐊", "🐅", "🐆", "🦓", "🦍", "🦧", "🐘", "🦛", "🦏", "🐪", "🐫", "🦒", "🦘", "🐃", "🐂", "🐄", "🐎", "🐖", "🐏", "🐑", "🦙", "🐐", "🦌", "🐕", "🐩", "🦮", "🐕‍🦺", "🐈", "🐈‍⬛", "🐓", "🦃", "🦚", "🦜", "🦢", "🦩", "🕊", "🐇", "🦝", "🦨", "🦡", "🦦", "🦥", "🐁", "🐀", "🐿", "🦔", "🐾", "🚗", "🚕", "🚙", "🚌", "🚎", "🏎", "🚓", "🚑", "🚒", "🚐", "🛻", "🚚", "🚛", "🚜", "🦯", "🦽", "🦼", "🛴", "🚲", "🛵", "🏍", "🛺", "🚔", "🚍", "🚘", "🚖", "🚡", "🚠", "🚟", "🚃", "🚋", "🚞", "🚝", "🚄", "🚅", "🚈", "🚂", "🚆", "🚇", "🚊", "🚉", "✈", "🛫", "🛬", "🛩", "💺", "🛰", "🚀", "🛸"}
	raceInProgress   bool
	raceMutex        sync.Mutex
)
var muteDuration time.Duration
var raceMessage *discordgo.Message

const counterFile = "death_counter.json"

func (dc *DeathCounter) save() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	data, err := json.Marshal(dc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(counterFile, data, 0644)
}

func (dc *DeathCounter) load() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	data, err := ioutil.ReadFile(counterFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist, start with 0
		}
		return err
	}
	return json.Unmarshal(data, dc)
}

func (dc *DeathCounter) increment() int {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.Count++
	return dc.Count
}

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

	emoji := raceEmojis[rand.Intn(len(raceEmojis))]
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

	// Отправляем сообщение о начале гонки
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
	trackLength := 20 // Длина трека в символах
	for winner == "" {
		time.Sleep(1 * time.Second)
		raceStatus := "```🏁 Гонка в процессе 🏁\n\n"
		for id, emoji := range raceParticipants {
			raceTrack[id] += rand.Intn(3) // Случайный прогресс каждого участника
			if raceTrack[id] >= trackLength {
				raceTrack[id] = trackLength // Ограничиваем прогресс максимальной длиной трека
				winner = id
				break
			}
			// Прогресс участника
			progress := strings.Repeat("—", raceTrack[id])
			// Остаток трека
			emptySpace := strings.Repeat("—", trackLength-raceTrack[id])
			// Формируем строку с эмодзи участника на текущей позиции
			raceStatus += fmt.Sprintf("🚦 |%s%s%s|\n", progress, emoji, emptySpace)
		}
		raceStatus += "```"

		// Редактируем сообщение, чтобы обновить статус гонки
		_, err := s.ChannelMessageEdit(m.ChannelID, raceMessage.ID, raceStatus)
		if err != nil {
			fmt.Println("error editing message:", err)
			return
		}
	}

	// Сообщение о победителе
	finalMessage := fmt.Sprintf("🎉 **Победитель гонки:** <@%s> %s! Поздравляем! 🏆🎉", winner, raceParticipants[winner])
	s.ChannelMessageSend(m.ChannelID, finalMessage)

	// Сброс состояния гонки
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

	// Расчет вероятности успешного питья пива
	successChance := 100 - (amount * 5)
	if successChance < 5 {
		successChance = 5 // Минимальный шанс 5%
	}

	roll := rand.Intn(100) + 1

	if roll <= successChance {
		// Успешное питье
		var successMessage string
		if amount == 20 {
			successMessage = fmt.Sprintf("<@%s> выпил %d литров пива и остался жив?! 🎉🍻\n\n", m.Author.ID, amount)
			// Отправляем GIF анимацию для максимального количества литров
			s.ChannelMessageSend(m.ChannelID, successMessage)
			s.ChannelMessageSend(m.ChannelID, "https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExejN4bjU1cTc1NDRodXU1OGd1NTExNTZheXRwOTdkaHNycWwyMTdtZyZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/qiSGGu0d2Dgac/giphy.gif") // Замените ссылку на подходящий GIF
		} else {
			successMessage = fmt.Sprintf("<@%s> успешно выпил %d литров пива! 🍺\n\n", m.Author.ID, amount)
			s.ChannelMessageSend(m.ChannelID, successMessage)
		}
	} else {
		// Неудача, применяем мут
		var failureMessage string
		if amount == 20 {
			failureMessage = fmt.Sprintf("<@%s> не смог осилить %d литров пива и отправляется в бессознательное состояние на 5 минут! 🍺😴\n\n", m.Author.ID, amount)
			// Отправляем GIF анимацию для провала
			s.ChannelMessageSend(m.ChannelID, failureMessage)
			s.ChannelMessageSend(m.ChannelID, "https://i.giphy.com/media/v1.Y2lkPTc5MGI3NjExd3Rqb3NycG0xZTRqNHZoamgybmVmOGRvYTcyamViNGJ6ZGM0YjA1MSZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/7bx7ZHokGnofm/giphy-downsized-large.gif") // Замените ссылку на подходящий GIF
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

		// Применяем мут
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

	// Create a poll
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
	message := "🤔🤔🤔"
	for _, user := range users {
		piskaProc := rand.Intn(101)
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
	size := rand.Intn(30) + 1
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
	// Начинаем сообщение с радужных эмодзи
	var message strings.Builder
	message.WriteString("🏳️‍🌈🌈🏳️‍🌈\n")

	// Генерация процента гейства
	gayProc := rand.Intn(101)
	var result string

	switch {
	case gayProc == 0:
		result = fmt.Sprintf("<@%s>, у тебя пока 0%% GaYства. Не сдавайся! 🥺", user)
	case gayProc == 100:
		// Генерация радужных эмодзи в зависимости от процента
		message.WriteString(strings.Repeat("🌈", 15))
		result = fmt.Sprintf("<@%s>, ты просто совершенство! 400%% GaYства! %s", user, strings.Join([]string{"🌈", "✨", "🦄", "💖", "🌟"}, " "))
	case gayProc >= 50:
		// Генерация радужных эмодзи в зависимости от процента
		message.WriteString(strings.Repeat("🌈", 10))
		result = fmt.Sprintf("<@%s>, у тебя %d%% гейства! Держись, радужный воин! 💃✨", user, gayProc)
	default:
		// Генерация радужных эмодзи в зависимости от процента
		message.WriteString(strings.Repeat("🌈", 5))
		result = fmt.Sprintf("<@%s>, у тебя %d%% гейства. Попробуй танцевать под Lady Gaga! 💃🎶", user, gayProc)
	}

	// Добавляем результат в сообщение
	message.WriteString(result + "\n")

	// Завершаем сообщение радужными эмодзи
	message.WriteString(strings.Repeat("🌈", 10) + "\n" + "🏳️‍🌈🌈🏳️‍🌈")

	// Отправляем сообщение
	s.ChannelMessageSend(m.ChannelID, message.String())

	// Добавляем случайные реакции с радужными эмодзи
	rainbowEmojis := []string{"🌈", "✨", "🦄", "💖", "🌟", "💅", "🎉", "💃", "🕺", "🎶"}
	for _, emoji := range rainbowEmojis {
		time.Sleep(200 * time.Millisecond) // Пауза перед добавлением следующей реакции
		s.MessageReactionAdd(m.ChannelID, m.ID, emoji)
	}

	// "Анимация" с последовательной отправкой радужных сообщений, если процент гейства больше 50
	if gayProc >= 50 {
		animatedMessage := "🌈 "
		for i := 0; i < 5; i++ {
			animatedMessage += strings.Repeat("🌈", i+1)
			_, err := s.ChannelMessageSend(m.ChannelID, animatedMessage)
			if err != nil {
				fmt.Println("error sending animated message:", err)
			}
			time.Sleep(300 * time.Millisecond) // Пауза между сообщениями
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	session, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	if err := deathCounter.load(); err != nil {
		fmt.Println("Error loading death counter:", err)
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
		if strings.Contains(strings.ToLower(m.Content), "бобр") || strings.Contains(strings.ToLower(m.Content), "бобер") || strings.Contains(strings.ToLower(m.Content), "курва") {
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
			count := deathCounter.increment()
			response := fmt.Sprintf("Сволочи, они убили @%s %d раз(а) 💀🔫", m.Author.Username, count)
			_, err := s.ChannelMessageSend(m.ChannelID, response)
			if err != nil {
				fmt.Println("Ошибка отправки сообщения:", err)
				return
			}
			if err := deathCounter.save(); err != nil {
				fmt.Println("Ошибка сохранения счетчика смертей:", err)
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

		// Checking on "писька" message
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

		// Checking on "полчаса" message
		if strings.Contains(strings.ToLower(m.Content), "полчаса") {
			_, err := s.ChannelMessageSendReply(m.ChannelID, "полчаса, полчаса - не вопрос. Не ответ полчаса, полчаса (c) Чайок", m.Reference())
			if err != nil {
				fmt.Println("error sending message,", err)
			}
		}

		// Checking on "полчаса" message
		if strings.Contains(strings.ToLower(m.Content), "керамика") {
			// Замените на ваше кастомное эмодзи
			customEmoji := "<:PotFriend:1271815662695743590>" // замените на ваш эмодзи
			// Формируем строку с кастомными эмодзи
			response := fmt.Sprintf("внезапная %s перекличка %s ебучих %s керамических %s изделий %s внезапная %s перекличка %s ебучих %s керамических %s изделий %s",
				customEmoji, customEmoji, customEmoji, customEmoji, customEmoji,
				customEmoji, customEmoji, customEmoji, customEmoji, customEmoji)
			// Отправляем сообщение с ответом
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

		if strings.HasPrefix(strings.ToLower(m.Content), "!гей") {
			var userID string

			// Проверяем, есть ли упомянутые пользователи
			if len(m.Mentions) > 0 {
				// Если есть упомянутые пользователи, выбираем первого из них
				userID = m.Mentions[0].ID
			} else {
				// Если нет упомянутых пользователей, ничего не делаем или отправляем сообщение об ошибке
				_, err := s.ChannelMessageSend(m.ChannelID, "Пожалуйста, упомяни пользователя для проверки гейства!")
				if err != nil {
					fmt.Println("error sending message:", err)
				}
				return
			}

			// Маленький шанс, что сообщение будет отправлено обратно автору команды
			if rand.Intn(10) == 0 {
				userID = m.Author.ID
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>, а может быть ты, моя голубая луна???!!!", userID))
				if err != nil {
					fmt.Println("error sending message:", err)
				}
			}

			// Генерируем и отправляем сообщение
			gayMessage(s, m, userID)
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
