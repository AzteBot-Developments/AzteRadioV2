package main

import (
	"fmt"

	"github.com/AzteBot-Developments/AzteMusic/internal/data/models/dax"
	"github.com/bwmarrin/discordgo"
)

// ONLY RUN THESE IF THE RUNTIME BENFITS OF A DB

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

			// Saved configuration, update player
			// TODO

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

		// Saved configuration, update player
		// TODO

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
