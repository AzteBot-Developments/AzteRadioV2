package repositories

import (
	"fmt"

	"github.com/AzteBot-Developments/AzteMusic/internal/data/models/dax"
)

type AzteradioConfigurationsDataRepository interface {
	GetAll() ([]dax.AzteradioConfiguration, error)
	GetConfiguration(guildId string) (*dax.AzteradioConfiguration, error)
	SaveConfiguration(config dax.AzteradioConfiguration) error
	RemoveConfiguration(guildId string) error
}

type AzteradioConfigurationRepository struct {
	DbContext AztebotDbContext
}

func NewAzteradioConfigurationRepository(connString string) *AzteradioConfigurationRepository {

	if connString == "" {
		return nil
	}

	repo := AzteradioConfigurationRepository{AztebotDbContext{
		ConnectionString: connString,
	}}
	repo.DbContext.Connect()
	return &repo
}

func (r AzteradioConfigurationRepository) GetAll() ([]dax.AzteradioConfiguration, error) {

	var configs []dax.AzteradioConfiguration

	rows, err := r.DbContext.SqlDb.Query("SELECT * FROM AzteradioConfigurations")
	if err != nil {
		return nil, fmt.Errorf("an error ocurred while retrieving all radio configs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var config dax.AzteradioConfiguration
		if err := rows.Scan(&config.GuildId, &config.DefaultRadioChannelId); err != nil {
			return nil, fmt.Errorf("error in AzteradioConfiguration GetAll: %v", err)
		}
		configs = append(configs, config)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in AzteradioConfiguration GetAll: %v", err)
	}

	return configs, nil
}

func (r AzteradioConfigurationRepository) GetConfiguration(guildId string) (*dax.AzteradioConfiguration, error) {

	query := "SELECT * FROM AzteradioConfigurations WHERE guildId = ?"
	row := r.DbContext.SqlDb.QueryRow(query, guildId)

	var item dax.AzteradioConfiguration
	err := row.Scan(&item.GuildId,
		&item.DefaultRadioChannelId)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r AzteradioConfigurationRepository) SaveConfiguration(config dax.AzteradioConfiguration) error {

	stmt, err := r.DbContext.SqlDb.Prepare(`
		INSERT INTO 
			AzteradioConfigurations(
				guildId, 
				defaultRadioChannelId
			)
		VALUES(?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(config.GuildId, config.DefaultRadioChannelId)
	if err != nil {
		return err
	}

	return nil
}

func (r AzteradioConfigurationRepository) UpdateConfiguration(config dax.AzteradioConfiguration) error {

	stmt, err := r.DbContext.SqlDb.Prepare(`
	UPDATE AzteradioConfigurations SET 
		defaultRadioChannelId = ?
	WHERE guildId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(config.DefaultRadioChannelId, config.GuildId)
	if err != nil {
		return err
	}

	return nil
}

func (r AzteradioConfigurationRepository) RemoveConfiguration(guildId string) error {

	stmt, err := r.DbContext.SqlDb.Prepare(`
	DELETE FROM AzteradioConfigurations WHERE guildId = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(guildId)
	if err != nil {
		return err
	}

	return nil
}
