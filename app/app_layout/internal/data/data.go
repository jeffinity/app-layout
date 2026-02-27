package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/jeffinity/singularity/friendly"
	"github.com/jeffinity/singularity/migratex"
	"github.com/jeffinity/singularity/pgx"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/jeffinity/app-layout/app/app_layout/internal/conf"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewPostgres,
	NewRedis,
	NewHelloRepo,
	NewAllMigrator,
)

// Data .
type Data struct {
	pg  *gorm.DB
	rdb *redis.ClusterClient
}

func NewAllMigrator(data *Data) *migratex.Migrator {
	return migratex.NewAllMigrator(data.pg, []any{
		&Hello{},
	})
}

// NewData .
func NewData(c *conf.Bootstrap, pg *gorm.DB, rdb *redis.ClusterClient, logger log.Logger) (*Data, func(), error) {

	mLog := log.NewHelper(log.With(logger, "module", "app_layout/data"))
	cleanup := func() {
		mLog.Info("closing the data resources")
	}

	return &Data{
		pg:  pg,
		rdb: rdb,
	}, cleanup, nil
}

func NewPostgres(c *conf.Bootstrap, logger log.Logger) (*gorm.DB, error) {
	return pgx.NewPostgres(c.GetLog().GetLevel(), c.GetData().GetPostgres().GetDsn(), logger)
}

func NewRedis(rootCtx context.Context, c *conf.Bootstrap, mLogger log.Logger) (*redis.ClusterClient, func(), error) {
	rc := c.GetData().GetRedisCluster()
	return friendly.NewRedisCluster(rootCtx, mLogger, rc.GetSeeds(), rc.GetPassword(), rc.GetReadOnly())
}
