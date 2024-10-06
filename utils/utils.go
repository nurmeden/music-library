package utils

import (
	"database/sql"
	"fmt"
	"io/ioutil"
)

const migrationPath = "./migrations"

func RunMigrations(db *sql.DB) error {
	files, err := ioutil.ReadDir(migrationPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		content, err := ioutil.ReadFile(fmt.Sprintf(migrationPath+"/%s", f.Name()))
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return err
		}
	}

	return nil
}
