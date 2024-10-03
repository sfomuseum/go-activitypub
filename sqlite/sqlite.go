package sqlite

import (
	"context"
	"database/sql"
	"fmt"
)

// SetupConnection executes a series of performance-enhacing SQLite PRAGMA statements.
func SetupConnection(ctx context.Context, conn *sql.DB) error {

	// conn.SetMaxOpenConns(1)

	pragma := []string{
		"PRAGMA JOURNAL_MODE=OFF",
		"PRAGMA SYNCHRONOUS=OFF",
		"PRAGMA LOCKING_MODE=EXCLUSIVE",
		// https://www.gaia-gis.it/gaia-sins/spatialite-cookbook/html/system.html
		"PRAGMA PAGE_SIZE=4096",
		"PRAGMA CACHE_SIZE=1000000",
	}

	for _, p := range pragma {

		_, err := conn.Exec(p)

		if err != nil {
			return fmt.Errorf("Failed to set pragma '%s', %w", p, err)
		}
	}

	return nil
}
