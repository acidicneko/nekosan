package handlers

import (
	"fmt"

	"github.com/acidicneko/nekosan/utils"
	"github.com/bwmarrin/discordgo"
)

func GuildCreateHandler(bot *discordgo.Session, event *discordgo.GuildCreate) {
	channels := event.Guild.Channels
	for i := 0; i < len(channels); i++ {
		if channels[i].Name == "server-logs" {
			bot.ChannelMessageSend(channels[i].ID, "NekoSan is online!")
		}
	}
	utils.ActiveGuildsMutex.Lock()
	utils.ActiveGuilds[event.Guild.ID] = event.Guild
	utils.ActiveGuildsMutex.Unlock()

	fmt.Printf("Guild Loaded: %s\n", event.Guild.Name)
}
