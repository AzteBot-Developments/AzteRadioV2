package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/AzteBot-Developments/AzteMusic/internal/data/models/dax"
	"github.com/AzteBot-Developments/AzteMusic/internal/data/repositories"
	"github.com/AzteBot-Developments/AzteMusic/pkg/shared"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

func ProcessSyncRadioStates(repo repositories.AzteradioConfigurationsDataRepository, s *discordgo.Session, client disgolink.Client, queues *shared.QueueManager, sFrequency int, repeatPlaylistCount int, defaultDesignatedPlaylistUrl string) {

	// Local configs for the cron
	const shufflePlaylistForEachPlayer bool = false

	fmt.Println("[CRON] Starting Task ProcessSyncRadioStates() at", time.Now(), "running every", sFrequency, "seconds")

	ticker := time.NewTicker(time.Duration(sFrequency) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				go syncRadioStates(repo, s, client, queues, repeatPlaylistCount, defaultDesignatedPlaylistUrl, shufflePlaylistForEachPlayer)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func syncRadioStates(repo repositories.AzteradioConfigurationsDataRepository, s *discordgo.Session, client disgolink.Client, queues *shared.QueueManager, repeatPlaylistCount int, defaultDesignatedPlaylistUrl string, shuffle bool) {
	if repo != nil {
		var configs []dax.AzteradioConfiguration
		var err error
		configs, err = repo.GetAll()
		if err != nil {
			log.Fatalf("Could not retrieve radio configs for guilds: %v", err)
		}
		for _, config := range configs {
			if config.DefaultRadioChannelId != "" {
				// If the guild has an assigned DefaultRadioChannelId
				if !shared.ClientPlayerIsPlayingTrack(client, config.GuildId) || (shared.PlayerCurrentChannelId(client, config.GuildId) != config.DefaultRadioChannelId) {
					// If the radio is not playing on the guild
					// or if it's playing on a channel different to the assigned one
					if err := shared.PlayOnStartupFromSourceForGuild(config.GuildId, s, client, queues, config.DefaultRadioChannelId, defaultDesignatedPlaylistUrl, repeatPlaylistCount, shuffle); err != nil {
						log.Fatalf("Could not play default radio playlist on channel (%s) for guild %s (onReady CRON): %v", config.DefaultRadioChannelId, config.GuildId, err)
					}
				}
			}
		}
	}
}
