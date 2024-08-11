package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/AzteBot-Developments/AzteMusic/internal/data/models/dax"
	"github.com/AzteBot-Developments/AzteMusic/internal/jobs"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {

	// Params
	const repeatPlaylistCount int = 3
	const syncRadioStatesFrequency int = 10

	// Initial lavalink setup unless it was setup already
	if !b.HasLavaLinkClient {
		b.SetupLavalink()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		b.AddLavalinkNode(ctx)
	}

	// Any initial setup for the music service !
	// i.e join designated server, play designated playlist, etc.

	// Set the playing status
	if StatusText != "" {
		s.UpdateGameStatus(0, StatusText)
	}

	// DEFAULT STRATEGY FOR MAIN GUILD (OTA)
	if DefaultDesignatedChannelId != "" {
		if err := s.ChannelVoiceJoinManual(DefaultGuildId, DefaultDesignatedChannelId, false, false); err != nil {
			log.Fatalf("Could not join designated voice channel (%s) for guild %s (onReady): %v", DefaultDesignatedChannelId, DefaultGuildId, err)
		}

		// Play designated playlist on loop, FOREVER :')
		if DefaultDesignatedPlaylistUrl != "" {
			go b.PlayOnStartupFromSource(DefaultGuildId, DefaultDesignatedChannelId, event, DefaultDesignatedChannelId, repeatPlaylistCount)
		}

		// Also run a cron to check whether there is anything playing - if there isn't, shuffle and play the designated playlist
		var numSec int = 300
		ticker := time.NewTicker(time.Duration(numSec) * time.Second)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					serverQueue := b.Queues.Get(DefaultGuildId)
					if len(serverQueue.Tracks) == 0 || !ServiceIsPlayingTrack(b, DefaultGuildId) {
						go b.PlayOnStartupFromSource(DefaultGuildId, DefaultDesignatedChannelId, event, DefaultDesignatedPlaylistUrl, repeatPlaylistCount)
					}
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
	}

	// BACKGROUND JOBS
	go jobs.ProcessSyncRadioStates(AzteradioConfigurationRepository, s, b.Lavalink, b.Queues, syncRadioStatesFrequency, repeatPlaylistCount, DefaultDesignatedPlaylistUrl)
}

func (b *Bot) onGuildCreate(_ *discordgo.Session, event *discordgo.GuildCreate) {

	fmt.Printf("Registering radio app for guild %s (%s)\n", event.ID, event.Name)

	go b.RegisterCommandsForGuild(event.ID)

	if AzteradioConfigurationRepository != nil {
		// HACK Aug 2024: It seems that GuildCreate fires everytime onReady too ???? a bit strange...
		cfg, err := AzteradioConfigurationRepository.GetConfiguration(event.ID)
		if cfg == nil {
			if err == sql.ErrNoRows {
				err := AzteradioConfigurationRepository.SaveConfiguration(dax.AzteradioConfiguration{
					GuildId:               event.ID,
					DefaultRadioChannelId: "",
				})
				if err != nil {
					fmt.Printf("An error ocurred while saving the initial configuration for guild %s: %v\n", event.ID, err)
				}
			} else {
				fmt.Printf("An error ocurred while retrieving the radio configuration for guild %s: %v\n", event.ID, err)
			}
		}
	}
}

func (b *Bot) onGuildDelete(_ *discordgo.Session, event *discordgo.GuildDelete) {
	if AzteradioConfigurationRepository != nil {
		err := AzteradioConfigurationRepository.RemoveConfiguration(event.ID)
		if err != nil {
			fmt.Printf("An error ocurred while removing the configuration for guild %s: %v\n", event.ID, err)
		}
	}
}
