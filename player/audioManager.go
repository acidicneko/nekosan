package player

import (
	"io"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type Song struct {
	Name        string
	Author      string
	FullUrl     string
	DownloadUrl string
	Duration    time.Duration
	ID          string
}

type Status int32

const (
	Resting Status = 0
	Playing Status = 1
	Paused  Status = 2
	Err     Status = 3
)

type GuildAudioManager struct {
	VoiceConn            *discordgo.VoiceConnection
	Queue                *Song
	QueueList            []*Song
	SkipInterrupt        chan bool
	CurrentStream        *dca.StreamingSession
	CurrentEncodeSession *dca.EncodeSession
	BotStatus            Status
}

// TODO: This should have a mutex
var GuildAudioManagers = make(map[string]*GuildAudioManager)

func (mb *GuildAudioManager) PlaySong(session *discordgo.Session, event *discordgo.MessageCreate) {
	song := mb.Dequeue()
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	encodingSession, err := dca.EncodeFile(song.DownloadUrl, options)
	mb.CurrentEncodeSession = encodingSession
	if err != nil {
		log.Println("Error encoding from yt url")
		log.Println(err.Error())
		return
	}
	defer encodingSession.Cleanup()
	time.Sleep(250 * time.Millisecond)
	err = mb.VoiceConn.Speaking(true)

	if err != nil {
		log.Println("Error connecting to discord voice")
		log.Println(err.Error())
	}
	mb.BotStatus = Playing
	embed := &discordgo.MessageEmbed{
		Title: ":notes: Now Playing",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    event.Author.Username,
			IconURL: event.Author.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  song.Name,
				Inline: false,
			},
			{
				Name:   "By",
				Value:  song.Author,
				Inline: false,
			},
			{
				Name:   "Duration",
				Value:  song.Duration.String(),
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://img.youtube.com/vi/" + song.ID + "/hqdefault.jpg",
		},
		Color:     0x5e81ac,
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "NekoSan",
			IconURL: session.State.User.AvatarURL(""),
		},
	}
	session.ChannelMessageSendEmbed(event.ChannelID, embed)
	done := make(chan error)
	stream := dca.NewStream(encodingSession, mb.VoiceConn, done)
	mb.CurrentStream = stream

	select {
	case err := <-done:
		log.Println("Song done")
		if err != nil && err != io.EOF {
			mb.BotStatus = Err
			log.Println(err.Error())
			return
		}
		mb.VoiceConn.Speaking(false)
		break
	case <-mb.SkipInterrupt:
		mb.VoiceConn.Speaking(false)
		return
	}
	mb.VoiceConn.Speaking(false)

	if len(mb.QueueList) == 0 {
		time.Sleep(250 * time.Millisecond)
		log.Println("Audio done")
		mb.Stop()
		mb.BotStatus = Resting
		return
	}

	time.Sleep(250 * time.Millisecond)
	log.Println("Play next in queue")
	go mb.PlaySong(session, event)
}

func (mb *GuildAudioManager) Enqueue(session *discordgo.Session, event *discordgo.MessageCreate, song *Song) {
	embed := &discordgo.MessageEmbed{
		Title: ":notes: Track added to queue",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    event.Author.Username,
			IconURL: event.Author.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Name",
				Value:  song.Name,
				Inline: false,
			},
			{
				Name:   "By",
				Value:  song.Author,
				Inline: false,
			},
			{
				Name:   "Duration",
				Value:  song.Duration.String(),
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://img.youtube.com/vi/" + song.ID + "/hqdefault.jpg",
		},
		Color:     0x88c0d0,
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "NekoSan",
			IconURL: session.State.User.AvatarURL(""),
		},
	}
	session.ChannelMessageSendEmbed(event.ChannelID, embed)
	mb.QueueList = append(mb.QueueList, song)
	mb.Queue = song
}

func (mb *GuildAudioManager) Dequeue() *Song {
	song := mb.QueueList[0]
	mb.QueueList = mb.QueueList[1:]
	return song
}

func (mb *GuildAudioManager) Stop() {
	mb.VoiceConn.Disconnect()
	mb.VoiceConn = nil
	mb.BotStatus = Resting
}

func (mb *GuildAudioManager) Skip(session *discordgo.Session, event *discordgo.MessageCreate) {
	if len(mb.QueueList) == 0 {
		mb.Stop()
	} else {
		if len(mb.SkipInterrupt) == 0 {
			session.ChannelMessageSend(event.ChannelID, "Skipping current song.")
			mb.CurrentEncodeSession.Cleanup()
		}
	}
}
