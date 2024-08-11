package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AzteBot-Developments/AzteMusic/internal/data/models/dax"
	"github.com/AzteBot-Developments/AzteMusic/internal/runtime"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection.
func (b *Bot) onReady(s *discordgo.Session, event *discordgo.Ready) {

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

	// Join designated channel for all configured servers
	var configs []dax.AzteradioConfiguration
	if runtime.AzteradioConfigurationRepository != nil {
		var err error
		configs, err = runtime.AzteradioConfigurationRepository.GetAll()
		if err != nil {
			log.Fatalf("Could not retrieve radio configs for guilds: %v", err)
		}
	}

	repeatPlaylistCount := 3

	if len(configs) == 0 {
		fmt.Println("STRAT 1")
		if DefaultDesignatedChannelId != "" {
			if err := s.ChannelVoiceJoinManual(DefaultGuildId, DefaultDesignatedChannelId, false, false); err != nil {
				log.Fatalf("Could not join designated voice channel (%s) for guild %s (onReady): %v", DefaultDesignatedChannelId, DefaultGuildId, err)
			}

			// Play designated playlist on loop, FOREVER :')
			if DefaultDesignatedPlaylistUrl != "" {
				if err := b.PlayOnStartupFromSource(DefaultGuildId, DefaultDesignatedChannelId, event, DefaultDesignatedChannelId, repeatPlaylistCount); err != nil {
					log.Fatalf("Could not play default radio playlist on channel (%s) for guild %s (onReady): %v", DefaultDesignatedChannelId, DefaultGuildId, err)
				}
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
							if err := b.PlayOnStartupFromSource(DefaultGuildId, DefaultDesignatedChannelId, event, DefaultDesignatedPlaylistUrl, repeatPlaylistCount); err != nil {
								log.Fatalf("Could not play default radio playlist on channel (%s) for guild %s (onReady CRON): %v", DefaultDesignatedChannelId, DefaultGuildId, err)
							}
						}
					case <-quit:
						ticker.Stop()
						return
					}
				}
			}()
		}
	} else {
		fmt.Println("STRAT 2")
		for _, config := range configs {
			if config.DefaultRadioChannelId != "" {
				if err := s.ChannelVoiceJoinManual(config.GuildId, config.DefaultRadioChannelId, false, false); err != nil {
					log.Fatalf("Could not join designated voice channel (%s) for guild %s (onReady): %v", config.DefaultRadioChannelId, config.GuildId, err)
				}

				// Play designated playlist on loop, FOREVER :')
				if DefaultDesignatedPlaylistUrl != "" {
					if err := b.PlayOnStartupFromSourceForGuild(config.GuildId, event, config.DefaultRadioChannelId, DefaultDesignatedPlaylistUrl, repeatPlaylistCount); err != nil {
						log.Fatalf("Could not play default radio playlist on channel (%s) for guild %s (onReady): %v", config.DefaultRadioChannelId, config.GuildId, err)
					}
				}

				// Also run a cron to check whether there is anything playing - if there isn't, shuffle and play the designated playlist
				var numSec int = 60 * 5
				ticker := time.NewTicker(time.Duration(numSec) * time.Second)
				quit := make(chan struct{})
				go func() {
					for {
						select {
						case <-ticker.C:
							serverQueue := b.Queues.Get(config.GuildId)
							if len(serverQueue.Tracks) == 0 || !ServiceIsPlayingTrack(b, config.GuildId) {
								if err := b.PlayOnStartupFromSourceForGuild(config.GuildId, event, config.DefaultRadioChannelId, DefaultDesignatedPlaylistUrl, repeatPlaylistCount); err != nil {
									log.Fatalf("Could not play default radio playlist on channel (%s) for guild %s (onReady CRON): %v", config.DefaultRadioChannelId, config.GuildId, err)
								}
							}
						case <-quit:
							ticker.Stop()
							return
						}
					}
				}()
			}
		}
	}
}

func (b *Bot) onGuildCreate(_ *discordgo.Session, event *discordgo.GuildCreate) {
	if runtime.AzteradioConfigurationRepository != nil {
		// HACK Aug 2024: It seems that GuildCreate fires everytime onReady too ???? a bit strange...
		cfg, err := runtime.AzteradioConfigurationRepository.GetConfiguration(event.ID)
		if cfg == nil && err != nil {
			err := runtime.AzteradioConfigurationRepository.SaveConfiguration(dax.AzteradioConfiguration{
				GuildId:               event.ID,
				DefaultRadioChannelId: "",
			})
			if err != nil {
				fmt.Printf("An error ocurred while saving the initial configuration for guild %s: %v\n", event.ID, err)
			}
		}
	}
}

func (b *Bot) onGuildDelete(_ *discordgo.Session, event *discordgo.GuildDelete) {

	if runtime.AzteradioConfigurationRepository != nil {

		err := runtime.AzteradioConfigurationRepository.RemoveConfiguration(event.ID)
		if err != nil {
			fmt.Printf("An error ocurred while removing the configuration for guild %s: %v\n", event.ID, err)
		}
	}

}
