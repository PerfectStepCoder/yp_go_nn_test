package servers

import (
	"context"
)

type Server interface {
	Start(addr string) error
	Stop(ctx context.Context) error
}
