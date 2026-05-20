package main

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func PrivatePing(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	if ok && userID != "" {
		return "", runtime.NewError("You're not authorized", PERMISSION_DENIED)
	}

	return "", nil
}
