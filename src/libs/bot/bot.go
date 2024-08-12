package bot

import (
	"context"
	"log"

	"github.com/AzteBot-Developments/AzteMusic/pkg"
	"github.com/AzteBot-Developments/AzteMusic/src/libs/config"
	"github.com/AzteBot-Developments/AzteMusic/src/libs/data/repositories"
	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

type Bot struct {
	// App configuration
	Environment config.Environment
	// Session details
	Session *discordgo.Session
	// Player details
	Lavalink          disgolink.Client
	Handlers          map[string]func(event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error
	Queues            *pkg.QueueManager
	HasLavaLinkClient bool
	// Other dependencies (repos, services, etc.)
	AzteradioConfigurationRepository repositories.AzteradioConfigurationsDataRepository
}

func NewBot(env config.Environment) *Bot {
	newBot := &Bot{
		Environment: env,
		Queues: &pkg.QueueManager{
			Queues: make(map[string]*pkg.Queue),
		},
		AzteradioConfigurationRepository: repositories.NewAzteradioConfigurationRepository(env.MySqlAztebotRootConnectionString),
	}
	newBot.HasLavaLinkClient = false
	return newBot
}

func GetAuthenticatedBotSession(token string) *discordgo.Session {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	return session
}

func (b *Bot) SetIntents() {
	b.Session.State.TrackVoice = true
	b.Session.Identify.Intents = discordgo.IntentGuilds | discordgo.IntentsGuildVoiceStates
}

func (b *Bot) SetupLavalink() {
	b.Lavalink = disgolink.New(snowflake.MustParse(b.Session.State.User.ID),
		disgolink.WithListenerFunc(b.onPlayerPause),
		disgolink.WithListenerFunc(b.onPlayerResume),
		disgolink.WithListenerFunc(b.onTrackStart),
		disgolink.WithListenerFunc(b.onTrackEnd),
		disgolink.WithListenerFunc(b.onTrackException),
		disgolink.WithListenerFunc(b.onTrackStuck),
		disgolink.WithListenerFunc(b.onWebSocketClosed),
		disgolink.WithListenerFunc(b.onUnknownEvent),
	)
	b.HasLavaLinkClient = true
}

func (b *Bot) AddVoiceHandlers() {
	b.Session.AddHandler(b.onApplicationCommand)
	b.Session.AddHandler(b.onVoiceStateUpdate)
	b.Session.AddHandler(b.onVoiceServerUpdate)
}

func (b *Bot) AddLavalinkNode(ctx context.Context, nodeName string, nodeAddr string, nodePass string, secure bool) {
	node, err := b.Lavalink.AddNode(ctx, disgolink.NodeConfig{
		Name:     nodeName,
		Address:  nodeAddr,
		Password: nodePass,
		Secure:   secure,
	})
	if err != nil {
		panic(err)
	}
	version, err := node.Version(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("lavalink node version: %s", version)
}
