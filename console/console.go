package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/acidicneko/nekosan/utils"
	"github.com/bwmarrin/discordgo"
)

func execute(args []string, bot *discordgo.Session) {
	if args[0] == "about" {
		fmt.Println("Nekosan. A discord bot written in Golang! v", utils.Version)
	} else if args[0] == "exit" {
		fmt.Println("Stopping Bot...")
		bot.Close()
		fmt.Println("Closing console...")
		os.Exit(0)
	} else if args[0] == "logout" {
		fmt.Println("Logging out from Discord.\nNOTE: Console will keep running! Type \"exit\" to exit console.")
		bot.Close()
	} else if args[0] == "login" {
		fmt.Println("Logging in into Discord...")
		err := bot.Open()
		if err != nil {
			fmt.Println("error opening connection", err)
			return
		}
		bot.UpdateGameStatus(0, "with your life")
		fmt.Println("NekoSan logged in into Discord!")
	} else if args[0] == "settoken" {
		if len(args) != 2 {
			fmt.Println("settoken: wrong number of arguments!")
			return
		}
		os.Setenv("BOT_TOKEN", args[1])
	} else if args[0] == "setprefix" {
		if len(args) != 2 {
			fmt.Println("setprefix: wrong number of arguments!")
			return
		}
		utils.Prefix = args[1]
		fmt.Println("Bot prefix set to:", args[1])
	} else {
		fmt.Printf("%s: unknown command\n", args[0])
	}
}

func InitConsole(bot *discordgo.Session) {
	for {
		fmt.Printf("[Console] > ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()
		args := strings.Split(line, " ")
		execute(args, bot)
	}
}
