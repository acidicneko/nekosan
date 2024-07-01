package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var usage = cmdlet{
	Name:  "usage",
	Usage: "print usage information about given command",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		if len(args) < 1 {
			bot.ChannelMessageSend(event.ChannelID, "`usage: no command specified`")
			return
		}
		requiredCmd, ok := cmdMap[args[0]]
		if !ok {
			bot.ChannelMessageSend(event.ChannelID, fmt.Sprintf("`usage: %s: no such command`", args[0]))
			return
		}
		embed := &discordgo.MessageEmbed{
			Title: "Command usage info",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    event.Author.Username,
				IconURL: event.Author.AvatarURL(""),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Command Name",
					Value:  requiredCmd.Name,
					Inline: false,
				},
				{
					Name:   "Usage",
					Value:  requiredCmd.Usage,
					Inline: false,
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
