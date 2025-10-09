package user

import (
	"context"
	"embed"
)

//go:embed migrations/*.sql
var FS embed.FS

func InitializeRepos(ctx context.Context) error {
	return nil
}
