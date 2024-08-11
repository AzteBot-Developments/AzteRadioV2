package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgolink/v3/lavalink"

	"github.com/AzteBot-Developments/AzteMusic/pkg/shared"
)

func (b *Bot) queue(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := b.Queues.Get(event.GuildID)
	if queue == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	if len(queue.Tracks) == 0 {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There are no songs on this queue.",
			},
		})
	}

	// Calculate the total length in time of the playlist
	var totalDurationSec int64
	for _, track := range queue.Tracks {
		totalDurationSec += track.Info.Length.Seconds()
	}

	// Get current track playing and add to embed
	currentTrack, player := b.GetCurrentTrack(event.GuildID)

	// Build embed response for the queue response
	embed := shared.NewEmbed().
		SetTitle(fmt.Sprintf("🎵  Queue - %s", BotName)).
		SetDescription(
			fmt.Sprintf(
				"Currently playing `%s` (%s) at %s / %s.\n\nQueue Duration: %s\nThere are %d other songs in this queue.\nThe first %d tracks in the queue can be seen below.", currentTrack.Info.Title, *currentTrack.Info.URI, formatPosition(player.Position()), formatPosition(currentTrack.Info.Length), shared.FormatDuration(totalDurationSec), len(queue.Tracks), 10)).
		SetThumbnail(*currentTrack.Info.ArtworkURL).
		SetColor(000000)

	// Build a list of discordgo embed fields out of the songs on the queue
	for index, track := range queue.Tracks {
		title := fmt.Sprintf("%d. `%s` (%s)", index+1, track.Info.Title, *track.Info.URI)
		text := ""
		embed.AddField(title, text, false)
	}

	// Truncate & paginate
	embed.Truncate()

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func (b *Bot) nowPlaying(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := b.Lavalink.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No player found.",
			},
		})
	}

	track := player.Track()
	if track == nil {
		return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No track found.",
			},
		})
	}

	embed := shared.NewEmbed().
		SetTitle("🎵  Now Playing").
		SetDescription(
			fmt.Sprintf("`%s` (%s).\n%s / %s", track.Info.Title, *track.Info.URI, formatPosition(player.Position()), formatPosition(track.Info.Length))).
		SetThumbnail(*track.Info.ArtworkURL).
		SetColor(000000)

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func formatPosition(position lavalink.Duration) string {
	if position == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
}

func (b *Bot) help(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {

	embed := shared.NewEmbed().
		SetTitle(fmt.Sprintf("🎵  `%s` Slash Commands Guide", BotName)).
		SetDescription(fmt.Sprintf("See below the available slash commands for `%s`.", BotName)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Build a list of discordgo embed fields out of the available slash commands
	for _, command := range Commands {

		text := command.Description
		title := fmt.Sprintf("`/%s`", command.Name)

		if len(command.Options) > 0 {
			for _, param := range command.Options {
				var required string
				if param.Required {
					required = "required"
				} else {
					required = "optional"
				}
				title += fmt.Sprintf(" `[%s (%s) - %s]`", param.Name, required, param.Description)
			}
		}

		embed.AddField(title, text, false)
	}

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}
