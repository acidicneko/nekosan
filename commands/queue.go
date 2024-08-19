package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/acidicneko/nekosan/player"
	"github.com/bwmarrin/discordgo"
)

var queue = cmdlet{
	Name:  "queue",
	Usage: "List songs in queue currently",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		AudioManager := player.GuildAudioManagers[event.GuildID]
		if AudioManager == nil {
			bot.ChannelMessageSend(event.ChannelID, "No AudioManager found for current guild")
			return
		}
		embed := &discordgo.MessageEmbed{
			Title: "Songs in Queue",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    event.Author.Username,
				IconURL: event.Author.AvatarURL(""),
			},
			Color:     0xed2939,                        // Red
			Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "NekoSan",
				IconURL: bot.State.User.AvatarURL(""),
			},
		}
		for val, song := range AudioManager.QueueList {
			tempField := &discordgo.MessageEmbedField{
				Name: strconv.Itoa(val+1) + ". " + song.Name + " - " + song.Duration.String(),
			}
			embed.Fields = append(embed.Fields, tempField)
		}
		_, err := bot.ChannelMessageSendEmbed(event.ChannelID, embed)
		if err != nil {
			fmt.Printf("Error while sending messages!")
		}

	},
}
