package commands

import (
	"errors"
	"fmt"

	"github.com/acidicneko/nekosan/player"
	"github.com/bwmarrin/discordgo"
)

func findUserVoiceState(bot *discordgo.Session, guildID string, userID string) (*discordgo.VoiceState, error) {
	var guild *discordgo.Guild = nil
	for _, guild = range bot.State.Guilds {
		if guild.ID == guildID {
			break
		}
	}
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			return vs, nil
		}
	}
	return nil, errors.New("could not find user's voice state")
}

var join = cmdlet{
	Name:  "join",
	Usage: "Join user's voice channel",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		// check if user isn't in any voice state
		vs, _ := findUserVoiceState(bot, event.GuildID, event.Author.ID)
		if vs == nil {
			_, err := bot.ChannelMessageSend(event.ChannelID, "You are not in any voice channel. Join a voice channel first!")
			if err != nil {
				fmt.Println("Error while sending messages!")
			}
			return
		}
		//check if bot is already in a voice state in current guild
		if vs, _ := bot.State.VoiceState(event.GuildID, bot.State.User.ID); vs != nil {
			_, err := bot.ChannelMessageSend(event.ChannelID, "Already in a voice channel!")
			if err != nil {
				fmt.Println("Error while sending messages!")
			}
			return
		}
		vc, _ := bot.ChannelVoiceJoin(event.GuildID, vs.ChannelID, false, false)
		player.GuildAudioManagers[event.GuildID] = &player.GuildAudioManager{
			VoiceConn: vc,
			BotStatus: player.Resting,
		}
	},
}
