package commands

import (
	"fmt"

	"github.com/acidicneko/nekosan/utils"
	"github.com/bwmarrin/discordgo"
)

var version = cmdlet{
	Name:  "version",
	Usage: "print bot version information",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		msg := fmt.Sprintf("NekoSan Discord Bot, v%s, written in Golang!\nCrafted with :hearts:, by Acidicneko",
			utils.Version)
		_, err := bot.ChannelMessageSend(event.ChannelID, msg)
		if err != nil {
			fmt.Println("Error while sending messages!")
		}
	},
}
