package commands

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var off = cmdlet{
	Name:  "off",
	Usage: "turn off the Bot",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		fmt.Println("Connection close request recieved from remote!")
		_, err := bot.ChannelMessageSend(event.ChannelID, "Going offline! Bye Bye!")
		if err != nil {
			fmt.Println("Error while sending messages!")
		}
		bot.Close()
		os.Exit(0)
	},
}
