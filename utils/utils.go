package utils

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var Prefix string = os.Getenv("BOT_PREFIX")
var Version string = "1.0.0"

var ActiveGuilds map[string]*discordgo.Guild
var ActiveGuildsMutex = sync.RWMutex{}

func MemberHasPermission(s *discordgo.Session, guildID string, userID string, permission int64) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}

func ExecuteCommand(args []string) (out string, e error) {
	cmd := exec.Command(args[0], args[1:]...)
	output, e := cmd.CombinedOutput()
	if e != nil {
		fmt.Printf("ExecuteCommand: error: %v\n", e)
		return "", e
	}
	return string(output), nil
}
