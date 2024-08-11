package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

func (b *Bot) GetCurrentTrack(guildId string) (*lavalink.Track, disgolink.Player) {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(guildId))
	if player == nil {
		return nil, nil
	}

	track := player.Track()
	if track == nil {
		return nil, nil
	}

	return track, player
}

func (b *Bot) AddToQueueFromSource(guildId string, url string, repeatCount int) {
	playlistUrl := url

	if !urlPattern.MatchString(playlistUrl) && !searchPattern.MatchString(playlistUrl) {
		playlistUrl = lavalink.SearchTypeYouTube.Apply(playlistUrl)
	}

	queue := b.Queues.Get(guildId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	b.Lavalink.BestNode().LoadTracksHandler(ctx, playlistUrl, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			queue.Add(track)
		},
		func(playlist lavalink.Playlist) {
			// Repeat the queue `repeatCount` times
			for i := 0; i < repeatCount; i++ {
				queue.Add(playlist.Tracks[0:]...)
			}
		},
		func(tracks []lavalink.Track) {
			queue.Add(tracks[0])
		},
		nil,
		nil,
	))
}

// Plays a YT track or playlist from the given source URL.
func (b *Bot) PlayOnStartupFromSource(guildId string, channelId string, event *discordgo.Ready, url string, repeatCount int) error {

	playlistUrl := url

	if !urlPattern.MatchString(playlistUrl) && !searchPattern.MatchString(playlistUrl) {
		playlistUrl = lavalink.SearchTypeYouTube.Apply(playlistUrl)
	}

	player := b.Lavalink.Player(snowflake.MustParse(guildId))
	queue := b.Queues.Get(guildId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	b.Lavalink.BestNode().LoadTracksHandler(ctx, playlistUrl, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			if player.Track() == nil {
				toPlay = &track
			} else {
				queue.Add(track)
			}
		},
		func(playlist lavalink.Playlist) {
			if player.Track() == nil {
				toPlay = &playlist.Tracks[0]
				queue.Add(playlist.Tracks[1:]...)
				// Repeat the queue `repeatCount` times
				for i := 0; i < repeatCount; i++ {
					queue.Add(playlist.Tracks[0:]...)
				}
				queue.Shuffle()
			} else {
				queue.Add(playlist.Tracks...)
				queue.Shuffle()
			}
		},
		func(tracks []lavalink.Track) {
			if player.Track() == nil {
				toPlay = &tracks[0]
			} else {
				queue.Add(tracks[0])
			}
		},
		nil,
		nil,
	))
	if toPlay == nil {
		return nil
	}

	if err := b.Session.ChannelVoiceJoinManual(guildId, channelId, false, false); err != nil {
		log.Fatalf("Could not join channel (2) at startup: %v", err)
		return err
	}

	return player.Update(context.TODO(), lavalink.WithTrack(*toPlay))
}

// Plays a YT track or playlist from the given source URL for a specific guild.
func (b *Bot) PlayOnStartupFromSourceForGuild(guildId string, event *discordgo.Ready, designatedChannelId string, url string, repeatCount int) error {

	playlistUrl := url

	if !urlPattern.MatchString(playlistUrl) && !searchPattern.MatchString(playlistUrl) {
		playlistUrl = lavalink.SearchTypeYouTube.Apply(playlistUrl)
	}

	player := b.Lavalink.Player(snowflake.MustParse(guildId))
	queue := b.Queues.Get(guildId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	b.Lavalink.BestNode().LoadTracksHandler(ctx, playlistUrl, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			if player.Track() == nil {
				toPlay = &track
			} else {
				queue.Add(track)
			}
		},
		func(playlist lavalink.Playlist) {
			if player.Track() == nil {
				toPlay = &playlist.Tracks[0]
				queue.Add(playlist.Tracks[1:]...)
				// Repeat the queue `repeatCount` times
				for i := 0; i < repeatCount; i++ {
					queue.Add(playlist.Tracks[0:]...)
				}
				queue.Shuffle()
			} else {
				queue.Add(playlist.Tracks...)
				queue.Shuffle()
			}
		},
		func(tracks []lavalink.Track) {
			if player.Track() == nil {
				toPlay = &tracks[0]
			} else {
				queue.Add(tracks[0])
			}
		},
		nil,
		nil,
	))
	if toPlay == nil {
		return nil
	}

	if err := b.Session.ChannelVoiceJoinManual(guildId, designatedChannelId, false, false); err != nil {
		log.Fatalf("Could not join channel (2) at startup: %v", err)
		return err
	}

	return player.Update(context.TODO(), lavalink.WithTrack(*toPlay))
}

func ServiceIsPlayingTrack(b *Bot, guildId string) bool {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(guildId))
	if player == nil {
		return false
	}

	track := player.Track()

	return track != nil
}

func ClientPlayerIsPlayingTrack(l disgolink.Client, guildId string) bool {
	player := l.ExistingPlayer(snowflake.MustParse(guildId))
	if player == nil {
		return false
	}

	track := player.Track()

	return track != nil
}

func MemberIsAdmin(guildId string, s *discordgo.Session, i discordgo.Interaction, m discordgo.Member) bool {

	hasAdminPermissions := false

	for _, roleID := range m.Roles {
		role, err := s.State.Role(guildId, roleID)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving role from Discord: %v\n", err)
			continue
		}

		if role.Permissions&discordgo.PermissionAdministrator != 0 {
			hasAdminPermissions = true
			break
		}
	}

	// Check if the member has the "Administrator" permission directly (e.g., server owner or other)
	if i.Member.Permissions&discordgo.PermissionAdministrator != 0 {
		hasAdminPermissions = true
	}

	if hasAdminPermissions {
		return true
	} else {
		return false
	}
}
