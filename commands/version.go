package commands

import (
	"fmt"
	"runtime"
	"time"

	"github.com/acidicneko/nekosan/utils"
	"github.com/bwmarrin/discordgo"
)

var version = cmdlet{
	Name:  "version",
	Usage: "print bot version information",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		embed := &discordgo.MessageEmbed{
			Title: "NekoSan :black_cat:",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    bot.State.User.Username,
				IconURL: bot.State.User.AvatarURL(""),
			},
			Description: "Made with :brown_heart: by Acidicneko",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Version",
					Value:  "v" + utils.Version,
					Inline: false,
				},
				{
					Name:   "Written in",
					Value:  runtime.Version(),
					Inline: false,
				},
				{
					Name:  "Build Date",
					Value: time.Now().Format(time.RFC3339),
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://raw.githubusercontent.com/acidicneko/nekosan/main/assets/cat.png",
			},
			Color:     0xed2939,                        // Red
			Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "NekoSan",
				IconURL: bot.State.User.AvatarURL(""),
			},
		}

		_, err := bot.ChannelMessageSendEmbed(event.ChannelID, embed)
		if err != nil {
			fmt.Println("Error while sending messages!")
		}
	},
}
