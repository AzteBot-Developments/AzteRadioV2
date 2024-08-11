package runtime

import (
	"os"

	"github.com/AzteBot-Developments/AzteMusic/internal/data/repositories"
)

// Connection strings
var MySqlAztebotRootConnectionString = os.Getenv("DB_AZTEBOT_ROOT_CONNSTRING")

// Repos
var AzteradioConfigurationRepository = repositories.NewAzteradioConfigurationRepository(MySqlAztebotRootConnectionString)
