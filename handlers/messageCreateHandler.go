package handlers

import (
	"strings"

	"github.com/acidicneko/nekosan/commands"
	"github.com/acidicneko/nekosan/utils"
	"github.com/bwmarrin/discordgo"
)

func MessageCreateHandler(bot *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID == bot.State.User.ID { // ignore messages from the bot itself
		return
	}
	if len(event.Content) < len(utils.Prefix)+1 || event.Content[0:len(utils.Prefix)+1] != utils.Prefix+" " {
		return
	}
	args := strings.Split(event.Content, " ")
	if args[0] == utils.Prefix {
		commands.HandleCommands(args[1:], bot, event)
	}
}
