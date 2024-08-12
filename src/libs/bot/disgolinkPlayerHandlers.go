package bot

import (
	"context"
	"fmt"
	"log"

	"github.com/AzteBot-Developments/AzteMusic/pkg"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

func (b *Bot) onPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
	b.Session.UpdateGameStatus(0, b.Environment.StatusText)
}

func (b *Bot) onPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
	b.Session.UpdateGameStatus(0, player.Track().Info.Title)
}

func (b *Bot) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	guildId := event.GuildID().String()
	if guildId == b.Environment.DefaultGuildId {
		// Only update the status when the main guild player changes state
		// in order to simulate some kind of "broadcast"
		// such that the status matches what the current song on the main station is
		b.Session.UpdateGameStatus(0, event.Track.Info.Title)
	}
}

func (b *Bot) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {

	if !event.Reason.MayStartNext() {
		return
	}

	guildId := event.GuildID().String()

	if guildId == b.Environment.DefaultGuildId {
		// Only update the status when the main guild player changes state
		// in order to simulate some kind of "broadcast"
		// such that the status matches what the current song on the main station is
		b.Session.UpdateGameStatus(0, b.Environment.StatusText)
	}

	queue := b.Queues.Get(guildId)

	// in the case of the radio service, we can check here whether the queue is empty
	// if it is, play form url again
	if len(queue.Tracks) < 2 && b.Environment.DefaultDesignatedPlaylistUrl != "" {
		b.AddToQueueFromSource(guildId, b.Environment.DefaultDesignatedPlaylistUrl, 3)
	}

	var (
		nextTrack lavalink.Track
		ok        bool
	)
	switch queue.Type {
	case pkg.QueueTypeNormal:
		nextTrack, ok = queue.Next()

	case pkg.QueueTypeRepeatTrack:
		nextTrack = event.Track

	case pkg.QueueTypeRepeatQueue:
		queue.Add(event.Track)
		nextTrack, ok = queue.Next()
	}

	if !ok {
		// retry to play designated playlist
		if b.Environment.DefaultDesignatedPlaylistUrl != "" {
			b.AddToQueueFromSource(guildId, b.Environment.DefaultDesignatedPlaylistUrl, 3)
		} else {
			// No tracks on the queue, or could not play next, so can safely disconnect from the VC to save resources.
			if err := b.Session.ChannelVoiceJoinManual(guildId, "", false, false); err != nil {
				fmt.Printf("[onTrackEnd] Error ocurred when disconnecting from VC: %v", err)
			}
			return
		}
	}

	if err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack)); err != nil {
		log.Fatal("Failed to play next track: ", err)
	}
}

func (b *Bot) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	fmt.Printf("onTrackException: %v\n", event)
	b.Session.UpdateGameStatus(0, b.Environment.StatusText)
}

func (b *Bot) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	fmt.Printf("onTrackStuck: %v\n", event)
	b.Session.UpdateGameStatus(0, b.Environment.StatusText)
}

func (b *Bot) onWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	fmt.Printf("onWebSocketClosed: %v\n", event)
	b.Session.UpdateGameStatus(0, b.Environment.StatusText)
}

func (b *Bot) onUnknownEvent(player disgolink.Player, event lavalink.UnknownEvent) {
	fmt.Printf("onWebSocketClosed: %v\n", event)
	b.Session.UpdateGameStatus(0, b.Environment.StatusText)
}
