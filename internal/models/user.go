package models

import (
	"context"

	"github.com/cameo-engineering/tonconnect"
)

type UserState struct {
	Ctx       context.Context
	S         *tonconnect.Session
	Connected bool
}
