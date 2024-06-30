package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
		quote := fmt.Sprintf("%s", data["q"])
		author := fmt.Sprintf("%s", data["a"])

		embed := &discordgo.MessageEmbed{
			Title: "Get up you dead soul!",
			Author: &discordgo.MessageEmbedAuthor{
				Name:    event.Author.Username,
				IconURL: event.Author.AvatarURL(""),
			},
			Color:       0x00ff00, // Green
			Description: quote,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "By",
					Value:  author,
					Inline: false,
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "https://c4.wallpaperflare.com/wallpaper/186/380/857/your-name-sky-stars-kimi-no-na-wa-wallpaper-preview.jpg",
			},
			Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "NekoSan",
				IconURL: bot.State.User.AvatarURL(""),
			},
		}
		bot.ChannelMessageSendEmbed(event.ChannelID, embed)

	},
}
