package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/heroiclabs/nakama-common/runtime"
)

const (
	OwnerID	   		= "00000000-0000-0000-0000-000000000000"
	Collection 		= "config"
	Key		   		= "game"
	CurrentVersion	= 1
)

type GameConfig struct {
	Version			int			`json:"version"`
	WelcomeMessage	string		`json:"welcome_message"`
	XpRate			float64		`json:"xp_rate"`
	RarityOptions	[]string	`json:"rarity_options"`
}

func seedGameConfig(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule) error {
	logger.Info("Seeding the game config!")

	env, _ := ctx.Value(runtime.RUNTIME_CTX_ENV).(map[string]string)

	path, ok := env["GAME_CONFIG_PATH"]
	if !ok || path == "" {
		path = "/nakama/data/modules/game_config.json"
		logger.Warn("GAME_CONFIG_PATH env variable is not set, using default: %s", path)
	}

	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{{
		Collection: Collection,
		Key:        Key,
		UserID:     OwnerID,
	}})
	
	if err != nil {
		logger.WithField("err", err).Error("Failed to check for an existing game config")
		return err 
	}
	
	data, err := os.ReadFile(path)
	
	if err != nil {
		logger.WithField("err", err).Error("Failed to read %s game config file", path)
		return err 
	}

	var config GameConfig

	err = json.Unmarshal(data, &config)
	
	if err != nil {
		logger.WithField("err", err).Error("Failed to parse %s game config file", path)
		return err 
	}

	if len(objects) == 0 {
		return writeGameConfig(ctx, nk, &config,  "", logger)
	}

	var stored GameConfig
	if err := json.Unmarshal([]byte(objects[0].Value), &stored); err != nil {
		logger.WithField("err", err).Warn("Failed to parse existing game config, overwriting")
		return writeGameConfig(ctx, nk, &config, objects[0].Version, logger)
	}

	if stored.Version == config.Version {
		logger.Info("Game config version unchanged (%d), skipping", stored.Version)
		return nil
	} else if stored.Version > config.Version {
		logger.WithField("err", err).Error("Stored game config has a never version (stored is %d, loaded is %d), aborting",
			stored.Version,
			config.Version,
		)
		return err
	}

	logger.Info("Loading new game config, from version %d to %d", stored.Version, config.Version)
	return writeGameConfig(ctx, nk, &config, objects[0].Version, logger)
}

func writeGameConfig(ctx context.Context, nk runtime.NakamaModule, config *GameConfig, 
	prevVersion string, logger runtime.Logger) error {

	data, err := json.Marshal(config)

	if err != nil {
		logger.WithField("err", err).Error("Failed to marshal game config")
		return err 
	}

	write := &runtime.StorageWrite{
		Collection:      Collection,
		Key:             Key,
		UserID:          OwnerID,
		Value:           string(data),
		PermissionRead:  2,
		PermissionWrite: 0,
		Version:         prevVersion,
	}

	if _, err := nk.StorageWrite(ctx, []*runtime.StorageWrite{write}); err != nil {
		logger.WithField("err", err).Error("Failed to write game config")
		return err 
	}

	logger.Info("New game config written with version %d", config.Version)

	return nil
}
