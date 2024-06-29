package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type cmdlet struct {
	Name  string
	Usage string
	Run   func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate)
}

var cmdMap map[string]cmdlet

func InitCommands() {
	cmdMap = make(map[string]cmdlet)
	registerCommand(echo)
	registerCommand(version)
	registerCommand(usage)
	registerCommand(off)
	registerCommand(quote)
}

func registerCommand(command cmdlet) {
	cmdMap[command.Name] = command
}

func HandleCommands(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
	cmdCalled, ok := cmdMap[args[0]]
	if !ok {
		errMsg := fmt.Sprintf("%s: command not found", args[0])
		bot.ChannelMessageSend(event.ChannelID, errMsg)
		return
	}
	cmdCalled.Run(args[1:], bot, event)
}
