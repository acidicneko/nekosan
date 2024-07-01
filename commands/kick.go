package commands

import (
	"fmt"

	"github.com/acidicneko/nekosan/utils"
	"github.com/bwmarrin/discordgo"
)

var kick = cmdlet{
	Name:  "kick",
	Usage: "kick the specified members",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		if status, _ := utils.MemberHasPermission(bot, event.GuildID, event.Author.ID,
			discordgo.PermissionKickMembers|discordgo.PermissionAdministrator); !status {
			_, err := bot.ChannelMessageSend(event.ChannelID, "You don't have permission to kick members.")
			if err != nil {
				fmt.Println("Error while sending messages!")
			}
			return
		}
		for _, mention := range event.Mentions {
			bot.ChannelMessageSend(event.ChannelID, fmt.Sprintf("kicking: %s", mention.Username))
			err := bot.GuildMemberDelete(event.GuildID, mention.ID)
			if err != nil {
				fmt.Printf("%v\n", err)
				_, err1 := bot.ChannelMessageSend(event.ChannelID, "Couldn't kick the specified member")
				if err1 != nil {
					fmt.Println("Error while sending messages!")
				}
			}
		}
	},
}
