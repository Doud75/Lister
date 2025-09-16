package main

import (
	"context"
	"log"
	"setlist/auth"
	"setlist/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	log.Println("Cleaning database...")
	cleanDatabase(db)

	log.Println("Seeding database...")
	seed(db)

	log.Println("Seeding complete!")
}

func cleanDatabase(db *pgxpool.Pool) {
	tables := []string{"setlist_items", "setlists", "songs", "interludes", "users", "bands"}
	for _, table := range tables {
		_, err := db.Exec(context.Background(), "DELETE FROM "+table)
		if err != nil {
			log.Printf("Could not clean table %s: %v\n", table, err)
		}
	}
}

func seed(db *pgxpool.Pool) {
	ctx := context.Background()

	var bandID int
	err := db.QueryRow(ctx, "INSERT INTO bands (name) VALUES ('The Testers') RETURNING id").Scan(&bandID)
	if err != nil {
		log.Fatalf("Failed to seed band: %v", err)
	}

	hashedPassword, _ := auth.HashPassword("password123")
	var userID int
	err = db.QueryRow(ctx, "INSERT INTO users (band_id, username, password_hash) VALUES ($1, $2, $3) RETURNING id",
		bandID, "testuser", hashedPassword).Scan(&userID)
	if err != nil {
		log.Fatalf("Failed to seed user: %v", err)
	}

	var song1ID, song2ID, interludeID, setlistID int

	db.QueryRow(ctx, "INSERT INTO songs (band_id, title) VALUES ($1, 'Song Title 1') RETURNING id", bandID).Scan(&song1ID)
	db.QueryRow(ctx, "INSERT INTO songs (band_id, title) VALUES ($1, 'Song Title 2') RETURNING id", bandID).Scan(&song2ID)

	db.QueryRow(ctx, "INSERT INTO interludes (band_id, title, speaker, script) VALUES ($1, 'Test Interlude', 'Le Chanteur', 'Quelques mots pour le public ici.') RETURNING id", bandID).Scan(&interludeID)

	err = db.QueryRow(ctx, "INSERT INTO setlists (band_id, name, color) VALUES ($1, 'Test Setlist', '#ff0000') RETURNING id", bandID).Scan(&setlistID)
	if err != nil {
		log.Fatalf("Failed to seed setlist: %v", err)
	}

	db.Exec(ctx, "INSERT INTO setlist_items (setlist_id, position, item_type, song_id, notes) VALUES ($1, 0, 'song', $2, 'Commencer avec le riff de guitare.')", setlistID, song1ID)
	db.Exec(ctx, "INSERT INTO setlist_items (setlist_id, position, item_type, interlude_id) VALUES ($1, 1, 'interlude', $2)", setlistID, interludeID)
	db.Exec(ctx, "INSERT INTO setlist_items (setlist_id, position, item_type, song_id) VALUES ($1, 2, 'song', $2)", setlistID, song2ID)
}
