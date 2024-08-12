package bot

import (
	"fmt"

	"github.com/AzteBot-Developments/AzteMusic/pkg"
	"github.com/AzteBot-Developments/AzteMusic/src/libs/data/models/dax"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleSlashSetRadioConfig(i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {

	if b.AzteradioConfigurationRepository != nil {

		designatedRadioChannel := data.Options[0].ChannelValue(b.Session)

		if designatedRadioChannel.Type == discordgo.ChannelTypeGuildVoice {
			cfg, _ := b.AzteradioConfigurationRepository.GetConfiguration(i.GuildID)
			if cfg == nil {
				err := b.AzteradioConfigurationRepository.SaveConfiguration(dax.AzteradioConfiguration{
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

			err := b.AzteradioConfigurationRepository.UpdateConfiguration(dax.AzteradioConfiguration{
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
		} else {
			return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Can't set the AzteRadio to play its tracklist on a channel other than a voice channel.",
				},
			})
		}
	}

	return b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "No configuration available to set.",
		},
	})

}

func (b *Bot) handleSlashRemoveRadioConfig(i *discordgo.InteractionCreate, _ discordgo.ApplicationCommandInteractionData) error {

	if b.AzteradioConfigurationRepository != nil {
		err := b.AzteradioConfigurationRepository.RemoveConfiguration(i.GuildID)
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

	if b.AzteradioConfigurationRepository != nil {
		config, err := b.AzteradioConfigurationRepository.GetConfiguration(i.GuildID)
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

		var prettyTargetChannelId = "*none selected.*"
		if config.DefaultRadioChannelId != "" {
			prettyTargetChannelId = config.DefaultRadioChannelId
		}

		// Build configuration output for displaying
		embed := pkg.NewEmbed().
			SetTitle(fmt.Sprintf("🤖🎵   `%s` Configuration", b.Environment.BotName)).
			SetDescription(prettyGuildIdentifier).
			SetColor(000000).
			AddField("Guild ID", fmt.Sprintf("`%s`", config.GuildId), false).
			AddField("Channel ID *(playing radio automatically on this channel)*", prettyTargetChannelId, false)
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
