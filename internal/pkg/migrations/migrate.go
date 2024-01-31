package migrations

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migration(migrationFilePath string, dbURl string) error {
	m, err := migrate.New(migrationFilePath, dbURl)
	if err != nil {
		return err
	}

	return m.Up()
}
