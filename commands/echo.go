package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var echo = cmdlet{
	Name:  "echo",
	Usage: "return user defined string",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		s := strings.Join(args, " ")
		bot.ChannelMessageDelete(event.ChannelID, event.Reference().MessageID)
		_, err := bot.ChannelMessageSend(event.ChannelID, s)
		if err != nil {
			fmt.Println("Error while sending messages!")
		}
	},
}
