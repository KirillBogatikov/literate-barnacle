package service

import (
	"context"
	"time"
)

type ContextProvider func() context.Context
type TimeoutContextProvider func(duration time.Duration) (context.Context, context.CancelFunc)
