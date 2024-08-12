package bot

import (
	"log"

	"github.com/AzteBot-Developments/AzteMusic/pkg"
	"github.com/AzteBot-Developments/AzteMusic/src/libs/commands"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) RegisterCommands() {

	b.Handlers = map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error{
		// GUIDE COMMANDS
		"help": b.help,

		// PLAYER COMMANDS
		"now-playing": b.nowPlaying,
		"queue":       b.queue,

		// CONFIGURATION COMMANDS
		"radio-config":  b.handleSlashSeeRadioConfig,
		"radio-set-cfg": b.handleSlashSetRadioConfig,
		"radio-rm-cfg":  b.handleSlashRemoveRadioConfig,
	}

	if b.AzteradioConfigurationRepository != nil {
		configs, _ := b.AzteradioConfigurationRepository.GetAll()
		if len(configs) != 0 {
			for _, config := range configs {
				// if _, err := b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, config.GuildId, Commands); err != nil {
				// 	fmt.Printf("could not bulk overwrite app commands for guild with ID %s: %v\n", config.GuildId, err)
				// }
				go b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, config.GuildId, commands.Commands)
			}
		} else {
			go b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, b.Environment.DefaultGuildId, commands.Commands)
		}
	} else {
		go b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, b.Environment.DefaultGuildId, commands.Commands)
	}
}

func (b *Bot) RegisterCommandsForGuild(guildId string) {
	go b.Session.ApplicationCommandBulkOverwrite(b.Session.State.User.ID, guildId, commands.Commands)
}

// GENERIC InteractionCreate HANDLER
func (b *Bot) onApplicationCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {

	data := event.ApplicationCommandData()

	// If allowed roles are configured, only allow a user with one of these roles to execute an app command
	// The app commands which require role permissions are defined here
	if pkg.StringInSlice(data.Name, b.Environment.RestrictedCommands) {
		if event.Type == discordgo.InteractionApplicationCommand {
			// Check if the user has the allowed role
			hasAllowedRole := pkg.MemberIsAdmin(event.GuildID, session, *event.Interaction, *event.Member)

			if !hasAllowedRole {
				// If the user doesn't have the allowed role, send a response
				session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You do not have the required role to use this command.",
					},
				})
				return
			}
		}
	}

	handler, ok := b.Handlers[data.Name]
	if !ok {
		log.Println("unknown command: ", data.Name)
		return
	}
	if err := handler(event, data); err != nil {
		log.Printf("error ocurred for %s: %v\n", data.Name, err)
		return
	}
}
