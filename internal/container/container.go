package container

import (
	"blogo/internal/config"
	"blogo/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Container struct {
	db      *pgxpool.Pool
	config  *config.Config
	queries sqlc.Querier
}

func NewContainer(db *pgxpool.Pool, config *config.Config, queries sqlc.Querier) *Container {
	return &Container{
		db:      db,
		config:  config,
		queries: queries,
	}
}

func (c *Container) GetDB() *pgxpool.Pool {
	return c.db
}

func (c *Container) GetConfig() *config.Config {
	return c.config
}

func (c *Container) GetQueries() sqlc.Querier {
	return c.queries
}
