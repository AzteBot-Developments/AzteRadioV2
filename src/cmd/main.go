package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/AzteBot-Developments/AzteMusic/src/libs/bot"
	"github.com/AzteBot-Developments/AzteMusic/src/libs/config"
)

var LavalinkUseSecure, _ = strconv.ParseBool(os.Getenv("LAVALINK_NODE_SECURE"))
var Bot = bot.NewBot(
	config.Environment{
		Token:                            os.Getenv("BOT_TOKEN"),
		BotAppId:                         os.Getenv("BOT_APP_ID"),
		MySqlAztebotRootConnectionString: os.Getenv("DB_AZTEBOT_ROOT_CONNSTRING"),
		DefaultGuildId:                   os.Getenv("GUILD_ID"),
		DefaultDesignatedChannelId:       os.Getenv("DESIGNATED_VOICE_CHANNEL_ID"),
		BotName:                          os.Getenv("BOT_NAME"),
		DefaultDesignatedPlaylistUrl:     os.Getenv("DESIGNATED_PLAYLIST_URL"),
		StatusText:                       os.Getenv("STATUS_TEXT"),
		RestrictedCommands:               strings.Split(os.Getenv("RESTRICTED_COMMANDS"), ","),
		NodeName:                         os.Getenv("LAVALINK_NODE_NAME"),
		NodeAddress:                      os.Getenv("LAVALINK_NODE_ADDRESS"),
		NodePassword:                     os.Getenv("LAVALINK_NODE_PASSWORD"),
		NodeSecure:                       LavalinkUseSecure,
	})

func main() {

	// Retrieve an authenticated Discord bot session through the token provided as an env variable
	Bot.Session = bot.GetAuthenticatedBotSession(Bot.Environment.Token)

	// Set the required intents for the bot's operation and what states it tracks
	Bot.SetIntents()

	// Register the handlers for the Discord session (onReady, onVoiceUpdate, etc.)
	Bot.AddVoiceHandlers()
	Bot.Session.AddHandler(Bot.OnReady)
	Bot.Session.AddHandler(Bot.OnGuildCreate)
	Bot.Session.AddHandler(Bot.OnGuildDelete)

	// Connect the authenticated bot session to the Discord servers
	if err := Bot.Session.Open(); err != nil {
		panic(err)
	}
	defer Bot.Session.Close()

	// Register the bot's slash commands (play, shuffle, skip, etc.)
	Bot.RegisterCommands()

	// Setup the Lavalink client for the bot if it hasn't been setup already
	if !Bot.HasLavaLinkClient {
		Bot.SetupLavalink()
		// Connect to the associated LavaLink server
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		Bot.AddLavalinkNode(ctx,
			Bot.Environment.NodeName,
			Bot.Environment.NodeAddress,
			Bot.Environment.NodePassword,
			Bot.Environment.NodeSecure,
		)
	}

	log.Printf("AzteRadio bot application is now running.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
