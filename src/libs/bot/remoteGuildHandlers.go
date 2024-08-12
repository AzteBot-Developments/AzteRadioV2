package bot

import (
	"database/sql"
	"fmt"

	"github.com/AzteBot-Developments/AzteMusic/src/libs/data/models/dax"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) OnGuildCreate(_ *discordgo.Session, event *discordgo.GuildCreate) {

	fmt.Printf("Registering radio app for guild %s (%s)\n", event.ID, event.Name)

	go b.RegisterCommandsForGuild(event.ID)

	if b.AzteradioConfigurationRepository != nil {
		// HACK Aug 2024: It seems that GuildCreate fires everytime onReady too ???? a bit strange...
		cfg, err := b.AzteradioConfigurationRepository.GetConfiguration(event.ID)
		if cfg == nil {
			if err == sql.ErrNoRows {
				err := b.AzteradioConfigurationRepository.SaveConfiguration(dax.AzteradioConfiguration{
					GuildId:               event.ID,
					DefaultRadioChannelId: "",
				})
				if err != nil {
					fmt.Printf("An error ocurred while saving the initial configuration for guild %s: %v\n", event.ID, err)
				}
			} else {
				fmt.Printf("An error ocurred while retrieving the radio configuration for guild %s: %v\n", event.ID, err)
			}
		}
	}
}

func (b *Bot) OnGuildDelete(_ *discordgo.Session, event *discordgo.GuildDelete) {
	if b.AzteradioConfigurationRepository != nil {
		err := b.AzteradioConfigurationRepository.RemoveConfiguration(event.ID)
		if err != nil {
			fmt.Printf("An error ocurred while removing the configuration for guild %s: %v\n", event.ID, err)
		}
	}
}
