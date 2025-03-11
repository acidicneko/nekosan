package player

import (
	"encoding/binary"
	"io"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
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
	channels  int = 2                   // 1 for mono, 2 for stereo
	frameRate int = 48000               // Audio sampling rate
	frameSize int = 960                 // uint16 size of each audio frame
	maxBytes  int = (frameSize * 2) * 2 // Max size of opus data
)

const (
	Resting Status = 0
	Playing Status = 1
	Paused  Status = 2
	Err     Status = 3
)

type GuildAudioManager struct {
	VoiceConn     *discordgo.VoiceConnection
	Queue         *Song
	QueueList     []*Song
	SkipInterrupt chan bool
	StopPlayback  chan bool
	BotStatus     Status
}

// TODO: This should have a mutex
var GuildAudioManagers = make(map[string]*GuildAudioManager)

func (mb *GuildAudioManager) PlaySong(session *discordgo.Session, event *discordgo.MessageCreate) {
	if len(mb.QueueList) == 0 {
		log.Println("Queue is empty")
		mb.BotStatus = Resting
		return
	}

	song := mb.Dequeue()
	mb.BotStatus = Playing

	// Start speaking
	err := mb.VoiceConn.Speaking(true)
	if err != nil {
		log.Println("Error setting speaking status:", err)
		mb.BotStatus = Err
		return
	}

	ytdlp := exec.Command(
		"yt-dlp",
		"--no-playlist",
		"--force-generic-extractor",
		"--youtube-skip-dash-manifest",
		"--no-check-certificate",
		"-f", "bestaudio",
		"-o", "-",
		"https://www.youtube.com/watch?v="+song.FullUrl,
	)

	ffmpeg := exec.Command(
		"ffmpeg",
		"-i", "pipe:0",
		"-f", "s16le",
		"-ar", strconv.Itoa(frameRate),
		"-ac", strconv.Itoa(channels),
		"pipe:1",
	)

	ytdlpout, err := ytdlp.StdoutPipe()
	if err != nil {
		log.Println("Error creating yt-dlp stdout pipe:", err)
		mb.BotStatus = Err
		return
	}
	ffmpeg.Stdin = ytdlpout

	ffmpegout, err := ffmpeg.StdoutPipe()
	if err != nil {
		log.Println("Error creating FFmpeg stdout pipe:", err)
		mb.BotStatus = Err
		return
	}

	ytdlp.Stderr = log.Writer()
	ffmpeg.Stderr = log.Writer()

	err = ytdlp.Start()
	if err != nil {
		log.Println("Error starting yt-dlp:", err)
		mb.BotStatus = Err
		return
	}

	err = ffmpeg.Start()
	if err != nil {
		log.Println("Error starting FFmpeg:", err)
		ytdlp.Process.Kill()
		mb.BotStatus = Err
		return
	}

	opusEncoder, err := gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		log.Println("Error creating Opus encoder:", err)
		ffmpeg.Process.Kill()
		ytdlp.Process.Kill()
		mb.BotStatus = Err
		return
	}

	opusEncoder.SetBitrate(96000)

	opusEncoder.SetApplication(gopus.Voip)

	// Buffer for reading PCM data
	ffmpegbuf := make([]int16, frameSize*channels)

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
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "NekoSan",
			IconURL: session.State.User.AvatarURL(""),
		},
	}
	session.ChannelMessageSendEmbed(event.ChannelID, embed)

	playing := true
	for playing {
		// Read PCM data from ffmpeg
		err = binary.Read(ffmpegout, binary.LittleEndian, &ffmpegbuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			log.Println("End of audio stream")
			break
		}
		if err != nil {
			log.Println("Error reading from FFmpeg:", err)
			break
		}

		// Encode PCM data to Opus
		opus, err := opusEncoder.Encode(ffmpegbuf, frameSize, maxBytes)
		if err != nil {
			log.Println("Error encoding to Opus:", err)
			break
		}

		// Send Opus data to Discord
		select {
		case mb.VoiceConn.OpusSend <- opus:
			// Data sent successfully
		case <-mb.SkipInterrupt:
			log.Println("Skipping current song")
			playing = false
		case <-mb.StopPlayback:
			log.Println("Stopping playback")
			playing = false
			ytdlp.Process.Kill()
			ffmpeg.Process.Kill()
			mb.BotStatus = Resting
			mb.VoiceConn.Speaking(false)
			return
		}
	}

	// Clean up
	ytdlp.Process.Kill()
	ffmpeg.Process.Kill()
	mb.VoiceConn.Speaking(false)

	// Play next song if available
	if len(mb.QueueList) > 0 {
		time.Sleep(250 * time.Millisecond)
		log.Println("Playing next song in queue")
		go mb.PlaySong(session, event)
	} else {
		log.Println("Queue empty, stopping playback")
		mb.BotStatus = Resting
	}
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
	if mb.VoiceConn != nil {
		mb.StopPlayback <- true
		mb.VoiceConn.Disconnect()
		mb.VoiceConn = nil
	}
	mb.BotStatus = Resting
}

func (mb *GuildAudioManager) Skip(session *discordgo.Session, event *discordgo.MessageCreate) {
	if len(mb.QueueList) == 0 {
		mb.Stop()
	} else {
		if len(mb.SkipInterrupt) == 0 {
			mb.SkipInterrupt <- true
			session.ChannelMessageSend(event.ChannelID, "Skipping current song.")
		}
	}
}
