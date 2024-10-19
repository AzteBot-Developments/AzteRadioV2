package bot

import (
	"context"
	"time"

	"github.com/AzteBot-Developments/AzteMusic/src/libs/jobs"
	"github.com/bwmarrin/discordgo"
)

// Called once the Discord servers confirm a succesful connection
func (b *Bot) OnReady(s *discordgo.Session, event *discordgo.Ready) {

	// Params
	const repeatPlaylistCount int = 2
	const syncRadioStatesFrequency int = 100

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

	// BACKGROUND JOBS
	go jobs.ProcessSyncRadioStates(b.AzteradioConfigurationRepository, s, b.Lavalink, b.Queues, syncRadioStatesFrequency, repeatPlaylistCount, b.Environment.DefaultDesignatedPlaylistUrl)
}
