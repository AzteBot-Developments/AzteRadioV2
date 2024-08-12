package commands

import "github.com/bwmarrin/discordgo"

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
