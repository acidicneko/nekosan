package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var list = cmdlet{
	Name:  "list",
	Usage: "list all the available commands and their usage",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		embed := &discordgo.MessageEmbed{
			Title: "Available Commands",
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
		for key, val := range cmdMap {
			tempField := &discordgo.MessageEmbedField{
				Name:   key,
				Value:  val.Usage,
				Inline: true,
			}
			embed.Fields = append(embed.Fields, tempField)
		}
		_, err := bot.ChannelMessageSendEmbed(event.ChannelID, embed)
		if err != nil {
			fmt.Printf("Error while sending messages!")
		}
	},
}
