package main

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func GetGameConfig(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	_, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	if !ok {
		return "", runtime.NewError("Invalid user ID", UNAUTHENTICATED)
	}

	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{{
		Collection: Collection,
		Key:        Key,
		UserID:     OwnerID,
	}})
	
	if err != nil {
		logger.WithField("err", err).Error("Failed to to find game config records")
		return "", runtime.NewError("Game config lookup failure", INTERNAL)
	}

	if len(objects) == 0 {
		logger.Error("Game config records are empty")
		return "", runtime.NewError("Game config not found", NOT_FOUND)
	}
	
	return objects[0].Value, nil
}
