package config

type Environment struct {
	// Bot App Configuration
	Token                            string
	BotAppId                         string
	MySqlAztebotRootConnectionString string
	DefaultGuildId                   string
	DefaultDesignatedChannelId       string
	BotName                          string
	DefaultDesignatedPlaylistUrl     string
	StatusText                       string
	RestrictedCommands               []string

	// Lavalink Node Configuration
	NodeName     string
	NodeAddress  string
	NodePassword string
	NodeSecure   bool
}
