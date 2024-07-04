package commands

import (
	"fmt"
	"strings"

	"github.com/acidicneko/nekosan/player"
	"github.com/bwmarrin/discordgo"
)

var play = cmdlet{
	Name:  "play",
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
		AudioManager, ok := player.GuildAudioManagers[event.GuildID]
		if !ok {
			bot.ChannelMessageSend(event.ChannelID, `Corrupted connection detected.\n
				Try removing the bot and joining again to your voice channel`)
			return
		}
		query := ""
		for _, k := range args {
			query += k
			query += " "
		}
		query = strings.TrimPrefix(query, " ")
		//url, e := player.FindYTSong(query)
		msg, _ := bot.ChannelMessageSend(event.ChannelID, fmt.Sprintf("Searching YT for: `%s`", query))
		url, e := player.FindYTSongYTDLP(query)
		if e != nil {
			bot.ChannelMessageSend(event.ChannelID, fmt.Sprintf("Failed to fetch query from YT: `%s`", query))
			return
		}
		song, err := player.GetSongInfo(url)
		if err != nil {
			bot.ChannelMessageSend(event.ChannelID, "Error while fetching the song!")
			return
		}
		bot.ChannelMessageDelete(event.ChannelID, msg.ID)
		AudioManager.Enqueue(bot, event, song)
		if AudioManager.BotStatus == player.Resting {
			AudioManager.PlaySong(bot, event)
		}
	},
}
