package main

import (
	"fmt"

	"github.com/AzteBot-Developments/AzteMusic/internal/data/models/dax"
	"github.com/AzteBot-Developments/AzteMusic/pkg/shared"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleSlashSetRadioConfig(i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {

	if AzteradioConfigurationRepository != nil {

		designatedRadioChannel := data.Options[0].ChannelValue(b.Session)

		cfg, _ := AzteradioConfigurationRepository.GetConfiguration(i.GuildID)
		if cfg == nil {
			err := AzteradioConfigurationRepository.SaveConfiguration(dax.AzteradioConfiguration{
				GuildId:               i.GuildID,
				DefaultRadioChannelId: designatedRadioChannel.ID,
			})
			if err != nil {
				fmt.Printf("An error ocurred while saving the initial configuration for guild %s: %v\n", i.GuildID, err)
				return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to save AzteRadio configuration.",
					},
				})
			}

			return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Saved AzteRadio configuration.",
				},
			})
		}

		err := AzteradioConfigurationRepository.UpdateConfiguration(dax.AzteradioConfiguration{
			GuildId:               i.GuildID,
			DefaultRadioChannelId: designatedRadioChannel.ID,
		})
		if err != nil {
			fmt.Printf("An error ocurred while updating the configuration for guild %s: %v\n", i.GuildID, err)
			return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to update AzteRadio configuration.",
				},
			})
		}

		return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Saved AzteRadio configuration.",
			},
		})
	}

	return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "No configuration available to set.",
		},
	})

}

func (b *Bot) handleSlashRemoveRadioConfig(i *discordgo.InteractionCreate, _ discordgo.ApplicationCommandInteractionData) error {

	if AzteradioConfigurationRepository != nil {
		err := AzteradioConfigurationRepository.RemoveConfiguration(i.GuildID)
		if err != nil {
			fmt.Printf("An error ocurred while removing the configuration for guild %s: %v\n", i.GuildID, err)
			return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to remove AzteRadio configuration.",
				},
			})
		}

		return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Removed AzteRadio configuration.",
			},
		})
	}

	return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "No AzteRadio configuration available to remove.",
		},
	})

}

func (b *Bot) handleSlashSeeRadioConfig(i *discordgo.InteractionCreate, _ discordgo.ApplicationCommandInteractionData) error {

	if AzteradioConfigurationRepository != nil {
		config, err := AzteradioConfigurationRepository.GetConfiguration(i.GuildID)
		if err != nil {
			fmt.Printf("An error ocurred while retrieving the configuration for guild %s: %v\n", i.GuildID, err)
			return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to retrieve AzteRadio configuration.",
				},
			})
		}

		var prettyGuildIdentifier string = ""
		guild, err := b.Session.Guild(i.GuildID)
		if err != nil {
			fmt.Println("Error fetching guild:", err)
		} else {
			prettyGuildIdentifier = fmt.Sprintf("These are the settings currently stored in the AzteRadio database for guild `%s`.", guild.Name)
		}

		// Build configuration output for displaying
		embed := shared.NewEmbed().
			SetTitle(fmt.Sprintf("🤖🎵   `%s` Configuration", BotName)).
			SetDescription(prettyGuildIdentifier).
			SetColor(000000).
			AddField("Guild ID", fmt.Sprintf("`%s`", config.GuildId), false).
			AddField("Channel ID (*playing radio automatically on this channel*)", fmt.Sprintf("`%s`", config.DefaultRadioChannelId), false)
		return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
			},
		})
	}

	return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "No AzteRadio configuration available to retrieve.",
		},
	})

}
