package commands

import (
	"fmt"

	"github.com/acidicneko/nekosan/player"
	"github.com/bwmarrin/discordgo"
)

var leave = cmdlet{
	Name:  "leave",
	Usage: "Leave user voice state",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		// check if bot isn't in any voice state
		if vs, _ := bot.State.VoiceState(event.GuildID, bot.State.User.ID); vs == nil {
			_, err := bot.ChannelMessageSend(event.ChannelID, "Not in any voice channel!")
			if err != nil {
				fmt.Println("Error while sending messages!")
			}
			return
		}
		// disconnect the bot and delete the GuildAudioManager associated to the guild
		player.GuildAudioManagers[event.GuildID].VoiceConn.Disconnect()
		delete(player.GuildAudioManagers, event.GuildID)
		/*voiceConnectionsMutex.Lock()
		requiredVoiceConnection, ok := voiceConnections[event.GuildID]
		voiceConnectionsMutex.Unlock()
		if !ok {
			bot.ChannelMessageSend(event.ChannelID, "Failed to disconnect the Voice connection in current guild!")
			return
		}
		requiredVoiceConnection.Disconnect()
		voiceConnectionsMutex.Lock()
		delete(voiceConnections, event.GuildID)
		voiceConnectionsMutex.Unlock()*/
	},
}
