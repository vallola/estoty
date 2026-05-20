package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func UpdateAccountMetadata(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	if !ok {
		return "", runtime.NewError("Invalid user ID", UNAUTHENTICATED)
	}

	account, err := nk.AccountGetId(ctx, userID)
	
	if err != nil {
		logger.WithField("err", err).Error("Failed to find an account by ID")
		return "", runtime.NewError("Account not found", NOT_FOUND)
	}

	metadata := make(map[string]interface{})

	err = json.Unmarshal([]byte(payload), &metadata)
	
	if err != nil {
		logger.WithField("err", err).Error("Failed to unmarshal the payload")
		return "", runtime.NewError("Invalid payload", INVALID_ARGUMENT)
	}

	if err := nk.AccountUpdateId(
		ctx,
		userID,
		account.User.Username,
		metadata,
		account.User.DisplayName,
		account.User.Timezone,
		account.User.Location,
		account.User.LangTag,
		account.User.AvatarUrl,
	); err != nil {
		logger.WithField("err", err).Error("Failed to update account metadata")
		return "", runtime.NewError("Internal error", INTERNAL)
	}
	
	return "{}", nil
}
