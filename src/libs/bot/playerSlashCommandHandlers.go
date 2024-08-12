package bot

import (
	"fmt"

	"github.com/AzteBot-Developments/AzteMusic/pkg"
	"github.com/AzteBot-Developments/AzteMusic/src/libs/commands"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"
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
	currentTrack, player := b.GetCurrentTrackForGuild(event.GuildID)

	// Build embed response for the queue response
	embed := pkg.NewEmbed().
		SetTitle(fmt.Sprintf("🎵  Queue - %s", b.Environment.BotName)).
		SetDescription(
			fmt.Sprintf(
				"Currently playing `%s` (%s) at %s / %s.\n\nQueue Duration: %s\nThere are %d other songs in this queue.\nThe first %d tracks in the queue can be seen below.", currentTrack.Info.Title, *currentTrack.Info.URI, pkg.FormatPosition(player.Position()), pkg.FormatPosition(currentTrack.Info.Length), pkg.FormatDuration(totalDurationSec), len(queue.Tracks), 10)).
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

	embed := pkg.NewEmbed().
		SetTitle("🎵  Now Playing").
		SetDescription(
			fmt.Sprintf("`%s` (%s).\n%s / %s", track.Info.Title, *track.Info.URI, pkg.FormatPosition(player.Position()), pkg.FormatPosition(track.Info.Length))).
		SetThumbnail(*track.Info.ArtworkURL).
		SetColor(000000)

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func (b *Bot) help(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {

	embed := pkg.NewEmbed().
		SetTitle(fmt.Sprintf("🎵  `%s` Guide", b.Environment.BotName)).
		SetDescription(fmt.Sprintf("The AzteRadio is a music app which automatically plays the Azteca Essentials playlist on your Discord server, and all it needs is a voice channel that the bot can join and then it'll sort out the rest!\n\nThe available slash commands for `%s` can be seen below.", b.Environment.BotName)).
		SetThumbnail("https://i.postimg.cc/262tK7VW/148c9120-e0f0-4ed5-8965-eaa7c59cc9f2-2.jpg").
		SetColor(000000)

	// Build a list of discordgo embed fields out of the available slash commands
	for _, command := range commands.Commands {

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
				title += fmt.Sprintf(" `[%s (%s)]`", param.Name, required)
			}
		}

		embed.AddField(title, text, false)
	}

	embed.AddLineBreakField()
	embed.AddField(fmt.Sprintf("Configuring the `%s`", b.Environment.BotName), "\n\nIn order to configure your bot to automatically play Azteca's Essentials on one of your voice channels, you'll have to run the `/radio-set-cfg` slash command.", false)

	return b.Session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}
