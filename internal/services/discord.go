package services

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lefes/curly-broccoli/config"
)

type DiscordService struct {
	s  *discordgo.Session
	dc *config.DiscordService
}

func NewDiscordService(dc *config.DiscordService) *DiscordService {
	session, err := discordgo.New("Bot " + dc.BotToken)
	if err != nil {
		panic(err)
	}
	return &DiscordService{s: session, dc: dc}
}

func (d *DiscordService) Open() (*discordgo.Session, error) {
	err := d.s.Open()
	if err != nil {
		return nil, err
	}
	return d.s, nil
}

/* func (r *DiscordService) GetAllUsers(guildID string) (*domain.DiscordMembers, error) { */
/* allMembers := &domain.DiscordMembers{} */
/* lastUserID := "" */

/* for { */
/* members, err := r.s.GuildMembers(guildID, lastUserID, 1000) */
/* if err != nil { */
/* log.Printf("Error fetching members: %v", err) */
/* return nil, err */
/* } */

/* if len(members) == 0 { */
/* break */
/* } */

/* allMembers.Members = append(allMembers.Members, members...) */

/* lastUserID = members[len(members)-1].User.ID */
/* } */

/* return allMembers, nil */
/* } */

/* func (r *DiscordService) GetDiscordIdByUsername(username string) (string, error) { */
/* user, err := r.s.User(username) */
/* if err != nil { */
/* return "", err */
/* } */

/* return user.ID, nil */
/* } */
