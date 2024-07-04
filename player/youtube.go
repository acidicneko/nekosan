package player

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/acidicneko/nekosan/utils"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func FindYTSong(query string) (url string, err error) {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YT_API_KEY")))
	if err != nil {
		log.Printf("Couldn't create new search service\n")
	}
	call := youtubeService.Search.List([]string{"id", "snippet"}).Q(query).MaxResults(1)
	response, err := call.Do()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	Id := response.Items[0].Id.VideoId
	return Id, err
}

func FindYTSongYTDLP(query string) (url string, err error) {
	yt_dlp_args := []string{"yt-dlp", "ytsearch:1\"" + query + "\"", "--get-id", "--flat-playlist", "--no-check-certificate"}
	result, err := utils.ExecuteCommand(yt_dlp_args)
	if err != nil {
		log.Printf("Error searching query: %s\n", query)
		return "", err
	}
	result = strings.TrimSuffix(result, "\n")
	return result, nil
}
