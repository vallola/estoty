package main

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	var err error

	err = seedGameConfig(ctx, logger, nk)
	
	if err != nil {
		return err
	}

    logger.Info("Registering RPC methods!")

	if err = initializer.RegisterRpc("update_account_metadata", UpdateAccountMetadata); err != nil {
		logger.WithField("err", err).Error("Unable to register RPC method update_account_metadata")
		return err
	}
	
	if err = initializer.RegisterRpc("get_game_config", GetGameConfig); err != nil {
		logger.WithField("err", err).Error("Unable to register RPC method get_game_config")
		return err
	}
	
	if err = initializer.RegisterRpc("private_ping", PrivatePing); err != nil {
		logger.WithField("err", err).Error("Unable to register RPC method private_ping")
		return err
	}
    
	logger.Info("Go RPC method registered!")

    return nil
}
