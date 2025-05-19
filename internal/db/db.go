package db

import (
	"context"
	"livescore/internal/models"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB(connString string) error {
	var err error
	Pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		return err
	}

	// Run migrations
	return RunMigrations()
}

// RunMigrations executes the SQL migrations to set up the database schema
func RunMigrations() error {
	// Get the directory of the current file
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// Read migrations file
	migrationPath := filepath.Join(dir, "migrations.sql")
	migrationSQL, err := os.ReadFile(migrationPath)
	if err != nil {
		return err
	}

	// Execute migrations
	_, err = Pool.Exec(context.Background(), string(migrationSQL))
	return err
}

// InsertLeagues inserts or updates a list of leagues in the database
func InsertLeagues(ctx context.Context, leagues []models.League) error {
	for _, league := range leagues {
		_, err := Pool.Exec(ctx, `
			INSERT INTO leagues (id, sport_id, country_id, name, active, short_code, image_path, type, sub_type)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
			ON CONFLICT (id) DO UPDATE SET
				sport_id=EXCLUDED.sport_id,
				country_id=EXCLUDED.country_id,
				name=EXCLUDED.name,
				active=EXCLUDED.active,
				short_code=EXCLUDED.short_code,
				image_path=EXCLUDED.image_path,
				type=EXCLUDED.type,
				sub_type=EXCLUDED.sub_type
		`,
			league.ID, league.SportID, league.CountryID, league.Name, league.Active, league.ShortCode, league.ImagePath, league.Type, league.SubType,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetLeagues retrieves all leagues from the database
func GetLeagues(ctx context.Context) ([]models.League, error) {
	rows, err := Pool.Query(ctx, `SELECT id, sport_id, country_id, name, active, short_code, image_path, type, sub_type FROM leagues`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leagues []models.League
	for rows.Next() {
		var l models.League
		if err := rows.Scan(&l.ID, &l.SportID, &l.CountryID, &l.Name, &l.Active, &l.ShortCode, &l.ImagePath, &l.Type, &l.SubType); err != nil {
			return nil, err
		}
		leagues = append(leagues, l)
	}
	return leagues, nil
}
