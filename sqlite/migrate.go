package sqlite

import "database/sql"

func Migrate(db *sql.DB, schema string) error {
	_, err := db.Exec(schema)
	return err
}
