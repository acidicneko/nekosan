package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/acidicneko/nekosan/commands"
	"github.com/acidicneko/nekosan/console"
	"github.com/acidicneko/nekosan/handlers"
	"github.com/bwmarrin/discordgo"
)

var (
	startConsole bool
	startLogin   bool
)

func init() {
	flag.BoolVar(&startConsole, "console", false, "Open console only. No login.")
	flag.BoolVar(&startLogin, "login", false, "Login to Discord only. No console launched.")
	flag.Parse()
}

func main() {
	if !startConsole && !startLogin {
		fmt.Println("No launch flag provided!\nProvide atleast one launch flag\n\t" +
			"-console: Open console only. No login.\n\t-login: Login to Discord only. No console launched.")
		os.Exit(1)
	}
	bot, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Printf("Error creating bot\n%v\n", err)
		os.Exit(1)
	}

	bot.AddHandler(handlers.MessageCreateHandler)

	// Required INTENTS
	bot.Identify.Intents = discordgo.IntentsGuildMessages

	commands.InitCommands()

	if startLogin {
		err = bot.Open()
		if err != nil {
			fmt.Println("error opening connection", err)
			return
		}
		bot.UpdateGameStatus(0, "with your life")
		fmt.Println("NekoSan logged in into Discord!")
		fmt.Println("Bot is now running. Type \"exit\" or press Ctrl+C to stop.")
	}

	if startConsole {
		console.InitConsole(bot)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Stopping Bot...")
	bot.Close()
}
