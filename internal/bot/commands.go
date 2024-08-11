package main

import (
	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "help",
		Description: "Returns a guide for the slash commands of the AzteRadio",
	},
	{
		Name:        "now-playing",
		Description: "Shows the current playing song",
	},
	{
		Name:        "queue",
		Description: "Shows the current queue of the AzteRadio",
	},
	// CONFIGURATION COMMANDS
	{
		Name:        "radio-set-cfg",
		Description: "Configures the settings of the AzteRadio application for this guild",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionChannel,
				Name:        "default-radio-channel",
				Description: "The bot will automatically join this channel whenever possible to play the default playlist",
				Required:    true,
			},
		},
	},
	{
		Name:        "radio-rm-cfg",
		Description: "Clear the AzteRadio configurations for this guild",
	},
	{
		Name:        "radio-config",
		Description: "Displays the current configurations for the AzteRadio",
	},
}

func (b *Bot) RegisterCommands() {
	b.Handlers = map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error{
		"now-playing": b.nowPlaying,
		"queue":       b.queue,
		"help":        b.help,
		// CONFIGURATION COMMANDS
		"radio-config":  b.handleSlashSeeRadioConfig,
		"radio-set-cfg": b.handleSlashSetRadioConfig,
		"radio-rm-cfg":  b.handleSlashRemoveRadioConfig,
	}

	if AzteradioConfigurationRepository != nil {
		configs, _ := AzteradioConfigurationRepository.GetAll()
		if len(configs) != 0 {
			for _, config := range configs {
				// if _, err := b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, config.GuildId, Commands); err != nil {
				// 	fmt.Printf("could not bulk overwrite app commands for guild with ID %s: %v\n", config.GuildId, err)
				// }
				go b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, config.GuildId, Commands)
			}
		} else {
			go b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, DefaultGuildId, Commands)
		}
	} else {
		go b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, DefaultGuildId, Commands)
	}
}

func (b *Bot) RegisterCommandsForGuild(guildId string) {
	go b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, guildId, Commands)
}
