package commands

import (
	"fmt"

	"github.com/acidicneko/nekosan/player"
	"github.com/bwmarrin/discordgo"
)

var skip = cmdlet{
	Name:  "skip",
	Usage: "Play given video from YouTube",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		if vs, _ := bot.State.VoiceState(event.GuildID, bot.State.User.ID); vs == nil {
			_, err := bot.ChannelMessageSend(event.ChannelID, "Not in any voice channel!")
			if err != nil {
				fmt.Println("Error while sending messages!")
			}
			return
		}
		vs, _ := findUserVoiceState(bot, event.GuildID, event.Author.ID)
		if vs == nil {
			_, err := bot.ChannelMessageSend(event.ChannelID, "You are not in any voice channel. Join a voice channel first!")
			if err != nil {
				fmt.Println("Error while sending messages!")
			}
			return
		}
		AudioManager := player.GuildAudioManagers[event.GuildID]
		if AudioManager.BotStatus == player.Resting {
			bot.ChannelMessageSend(event.ChannelID, "Nothing is playing right now.")
			return
		}
		AudioManager.Skip(bot, event)
	},
}
