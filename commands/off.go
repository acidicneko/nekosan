package commands

import (
	"fmt"
	"os"

	"github.com/acidicneko/nekosan/utils"
	"github.com/bwmarrin/discordgo"
)

var off = cmdlet{
	Name:  "off",
	Usage: "turn off the Bot",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		if event.Author.ID != utils.ActiveGuilds[event.GuildID].OwnerID {
			_, err := bot.ChannelMessageSend(event.ChannelID, "You don't have the permission to turn off the bot!")
			if err != nil {
				fmt.Println("Error while sending messages!")
			}
			return
		}
		fmt.Println("Connection close request recieved from remote!")
		_, err := bot.ChannelMessageSend(event.ChannelID, "Going offline! Bye Bye!")
		if err != nil {
			fmt.Println("Error while sending messages!")
		}
		bot.Close()
		os.Exit(0)
	},
}
