package texts

import "strings"

var (
	MorningMessages = []string{
		"доброе утро", "доброго утра", "добрый день", "добрый вечер", "доброй ночи", "утречко", "ночечко", "проснул", "открыл глаза",
	}
	SpokiMessages    = []string{"сладких снов", "спокойной ночи", "до завтра"}
	PhasmaMessages   = []string{"фасма", "фазма", "фазму", "фасму", "фазмой", "фасмой", "фазме", "фасме", "фазмы", "фасмы", "phasma", "phasmaphobia", "призрак"}
	SickMessages     = []string{"заболел", "заболела", "заболело", "заболели", "болею", "болит"}
	PotterMessages   = []string{"гарри", "поттер", "гарик", "гаррик", "потник", "потер", "гари", "хогвартс", "хогварт", "хогвардс", "хогвард", "гаррипоттер"}
	ValorantMessages = []string{"валорант", "валик", "валарант"}
	LegionEmojis     = []string{"🇱", "🇪", "🇬", "🇮", "🇴", "🇳"}
	CovenEmojis      = []string{"🇨", "🇴", "🇻", "🇪", "🇳"}
)

func ContainsAny(content string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(content, n) {
			return true
		}
	}
	return false
}

func CeramicsMessage(customEmoji string) string {
	return "внезапная " + customEmoji + " перекличка " + customEmoji + " ебучих " + customEmoji + " керамических " + customEmoji + " изделий " + customEmoji +
		" внезапная " + customEmoji + " перекличка " + customEmoji + " ебучих " + customEmoji + " керамических " + customEmoji + " изделий " + customEmoji
}
