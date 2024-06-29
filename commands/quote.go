package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var quote = cmdlet{
	Name:  "quote",
	Usage: "provide a motivational quote to cheer you up",
	Run: func(args []string, bot *discordgo.Session, event *discordgo.MessageCreate) {
		response, err := http.Get("https://zenquotes.io/api/random")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()
		b, error := io.ReadAll(response.Body)
		if error != nil {
			fmt.Println(err)
			return
		}
		jsonString := strings.TrimSuffix(strings.TrimPrefix(string(b), "[ "), " ]")
		var data map[string]interface{}
		err1 := json.Unmarshal([]byte(jsonString), &data)
		if err1 != nil {
			fmt.Printf("could not unmarshal json: %s\n", err1)
			return
		}
		_, err2 := bot.ChannelMessageSend(event.ChannelID, fmt.Sprintf("`%s`\n     ~ %s", data["q"], data["a"]))
		if err2 != nil {
			fmt.Println("Error while sending messages!")
		}
	},
}
