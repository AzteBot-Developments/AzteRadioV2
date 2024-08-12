package pkg

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

var urlPattern = regexp.MustCompile("^https?://[-a-zA-Z0-9+&@#/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#/%=~_|]?")
var searchPattern = regexp.MustCompile(`^(.{2})search:(.+)`)

func ClientPlayerIsPlayingTrack(l disgolink.Client, guildId string) bool {
	player := l.ExistingPlayer(snowflake.MustParse(guildId))
	if player == nil {
		return false
	}

	track := player.Track()

	return track != nil
}

func PlayerCurrentChannelId(l disgolink.Client, guildId string) string {
	player := l.ExistingPlayer(snowflake.MustParse(guildId))
	if player == nil {
		return ""
	}

	channelId := player.ChannelID().String()

	return channelId
}

func PlayFromUrlForGuildChannelById(guildId string, session *discordgo.Session, client disgolink.Client, queues *QueueManager, designatedChannelId string, url string, repeatCount int, shuffle bool) error {

	playlistUrl := url

	if !urlPattern.MatchString(playlistUrl) && !searchPattern.MatchString(playlistUrl) {
		playlistUrl = lavalink.SearchTypeYouTube.Apply(playlistUrl)
	}

	player := client.Player(snowflake.MustParse(guildId))
	queue := queues.Get(guildId)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	client.BestNode().LoadTracksHandler(ctx, playlistUrl, disgolink.NewResultHandler(
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

				if shuffle {
					queue.Shuffle()
				}
			} else {
				queue.Add(playlist.Tracks...)

				if shuffle {
					queue.Shuffle()
				}
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

	if err := session.ChannelVoiceJoinManual(guildId, designatedChannelId, false, false); err != nil {
		log.Fatalf("Could not join channel (2) at startup: %v", err)
		return err
	}

	return player.Update(context.TODO(), lavalink.WithTrack(*toPlay))
}
