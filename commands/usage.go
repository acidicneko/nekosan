package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var usage = cmdlet{
	Name:  "usage",
	Usage: "print usage information about given command",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		requiredCmd, ok := cmdMap[args[0]]
		if !ok {
			bot.ChannelMessageSend(event.ChannelID, fmt.Sprintf("usage: %s: no such command", args[0]))
			return
		}
		msg := fmt.Sprintf("Command name: `%s`\nUsage: `%s`", requiredCmd.Name, requiredCmd.Usage)
		_, err := bot.ChannelMessageSend(event.ChannelID, msg)
		if err != nil {
			fmt.Println("Error while sending messages!")
		}
	},
}
