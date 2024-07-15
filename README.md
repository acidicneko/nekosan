# NekoSan

A music/mod Discord bot written in Golang!
<img src="https://github.com/acidicneko/nekosan/blob/main/assets/cat-rounded.png?raw=true" align="right" alt="logo">

> Who would believe such pleasure from a wee ball O' fur?

NekoSan is a descendent of [ArruBot](https://github.com/acidicneko/ArruBot).

ArruBot was written in Java 17 and supported moderation and music operations.
It died due to major changes in YT APIs and lavaplayer package, however its spirit continued to live.

NekoSan continues its legacy with a much faster and compiled language Go!

## Compiling
The steps below assume that you have already installed the Go(v1.22.4) toolchain on your system.

Clone the repository
```sh
git clone https://github.com/acidicneko/nekosan.git
```

Chdir into the cloned directory
```sh
cd nekosan
```

Compile the executable
```sh
go build -v
```

## Running the bot
The compiled executable expects a few environment variables and `yt-dlp` binary to be present on the host system.

- `BOT_TOKEN`: Your discord bot token here.
- `BOT_PREFIX`: Bot will listen to all the messages prefixed with it.
- `YT_API_KEY`: Provide a Youtube API Key if you plan on using YT API to fetch search queries from YT(Paid, faster).
- `yt-dlp`: Install `yt-dlp` if you plan on using it to fetch search queries from YT(Free, slower).

Launch the executable
```sh
./nekosan -login
```

<hr>

*look it's a cat!*
