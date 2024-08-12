package bot

import (
	"context"
	"log"
	"time"

	"github.com/AzteBot-Developments/AzteMusic/src/libs/jobs"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func (b *Bot) OnReady(s *discordgo.Session, event *discordgo.Ready) {

	// Params
	const repeatPlaylistCount int = 3
	const syncRadioStatesFrequency int = 10
	const shufflePlaylistForEachPlayer bool = false

	// Initial lavalink setup unless it was setup already
	if !b.HasLavaLinkClient {
		b.SetupLavalink()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		b.AddLavalinkNode(ctx,
			b.Environment.NodeName,
			b.Environment.NodeAddress,
			b.Environment.NodePassword,
			b.Environment.NodeSecure,
		)
	}

	// Any initial setup for the music service !
	// i.e join designated server, play designated playlist, etc.

	// Set the playing status
	if b.Environment.StatusText != "" {
		s.UpdateGameStatus(0, b.Environment.StatusText)
	}

	// DEFAULT STRATEGY FOR MAIN GUILD (OTA)
	if b.Environment.DefaultDesignatedChannelId != "" {
		if err := s.ChannelVoiceJoinManual(b.Environment.DefaultGuildId, b.Environment.DefaultDesignatedChannelId, false, false); err != nil {
			log.Fatalf("Could not join designated voice channel (%s) for guild %s (onReady): %v", b.Environment.DefaultDesignatedChannelId, b.Environment.DefaultGuildId, err)
		}

		// Play designated playlist on loop, FOREVER :')
		if b.Environment.DefaultDesignatedPlaylistUrl != "" {
			go b.PlayOnStartupFromSource(b.Environment.DefaultGuildId, b.Environment.DefaultDesignatedChannelId, event, b.Environment.DefaultDesignatedPlaylistUrl, repeatPlaylistCount, shufflePlaylistForEachPlayer)
		}

		// Also run a cron to check whether there is anything playing - if there isn't, shuffle and play the designated playlist
		var numSec int = 300
		ticker := time.NewTicker(time.Duration(numSec) * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					serverQueue := b.Queues.Get(b.Environment.DefaultGuildId)
					if len(serverQueue.Tracks) == 0 || !ServiceIsPlayingTrackForGuild(b, b.Environment.DefaultGuildId) {
						go b.PlayOnStartupFromSource(b.Environment.DefaultGuildId, b.Environment.DefaultDesignatedChannelId, event, b.Environment.DefaultDesignatedPlaylistUrl, repeatPlaylistCount, shufflePlaylistForEachPlayer)
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}

	// BACKGROUND JOBS
	go jobs.ProcessSyncRadioStates(b.AzteradioConfigurationRepository, s, b.Lavalink, b.Queues, syncRadioStatesFrequency, repeatPlaylistCount, b.Environment.DefaultDesignatedPlaylistUrl)
}
