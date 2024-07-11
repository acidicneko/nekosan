package player

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/acidicneko/nekosan/utils"
	YT "github.com/kkdai/youtube/v2"
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

func findAudioFormat(formats YT.FormatList) *YT.Format {
	var audioFormat *YT.Format
	var audioFormats YT.FormatList = formats.Type("audio")

	if len(audioFormats) > 0 {
		audioFormats.Sort()
		audioFormat = &audioFormats[0]
	}

	return audioFormat
}

func GetSongInfo(url string) (*Song, error) {
	client := YT.Client{}
	sng, err := client.GetVideo(url)
	if err != nil {
		log.Printf("Error while retrieving song %v\n", err)
		return nil, err
	}
	downloadURL, _ := client.GetStreamURL(sng, findAudioFormat(sng.Formats))
	return &Song{
		Name:        sng.Title,
		Author:      sng.Author,
		FullUrl:     url,
		DownloadUrl: downloadURL,
		Duration:    sng.Duration,
		ID:          sng.ID,
	}, err
}

func GetPlaylistInfo(url string) (*YT.Playlist, error) {
	client := YT.Client{}
	playlist, err := client.GetPlaylist(url)
	if err != nil {
		return nil, err
	}
	return playlist, nil
}
