package main

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kwalter26/scoreit-api-go/api"
	"github.com/kwalter26/scoreit-api-go/api/middleware"
	db "github.com/kwalter26/scoreit-api-go/db/sqlc"
	"github.com/kwalter26/scoreit-api-go/util"
	_ "github.com/lib/pq"
	_ "github.com/nats-io/nkeys"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "google.golang.org/grpc"
	"os"
)

func main() {
	config, err := util.LoadConfig(".", false)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config:")
	}

	if config.Environment == util.Development || config.Environment == util.Testing {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db:")
	}

	// run db migrations
	runDBMigrations(config.MigrationUrl, config.DBSource)

	store := db.NewStore(conn)

	runGINServer(config, store)

}

// runDBMigrations runs db migrations
func runDBMigrations(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msgf("cannot migrate db:%s", migrationURL)
	}
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("failed to apply migration:")
	}
	log.Info().Msg("db migration completed")
}

// runGINServer runs gin server but is not used anymore
func runGINServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store, middleware.NewAuth0JwtValidator(config))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gin server")
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
