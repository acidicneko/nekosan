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

		if len(args) == 0 {
			_, err := bot.ChannelMessageSend(event.ChannelID, "No queue command specified!")
			if err != nil {
				fmt.Println("Error while sending messages!")
			}
			return
		}
		AudioManager := player.GuildAudioManagers[event.GuildID]
		if AudioManager == nil {
			bot.ChannelMessageSend(event.ChannelID, "No AudioManager found for current guild")
			return
		}
		if args[0] == "list" {
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
		} else if args[0] == "remove" {
			if len(args) < 2 {
				bot.ChannelMessageSend(event.ChannelID, "No track position specified to be removed!")
				return
			}
			s, _ := strconv.Atoi(args[1])
			bot.ChannelMessageSend(event.ChannelID, fmt.Sprintf("Removing song `%s` from queue.", AudioManager.QueueList[s]))
			AudioManager.QueueList = append(AudioManager.QueueList[:s-1], AudioManager.QueueList[s:]...)
		}
	},
}
