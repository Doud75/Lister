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
    tables := []string{"setlist_items", "setlists", "songs", "interludes", "band_users", "users", "bands"}
    for _, table := range tables {
        _, err := db.Exec(context.Background(), "TRUNCATE "+table+" RESTART IDENTITY CASCADE")
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
    err = db.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
        "testuser", hashedPassword).Scan(&userID)
    if err != nil {
        log.Fatalf("Failed to seed user: %v", err)
    }

    _, err = db.Exec(ctx, "INSERT INTO band_users (user_id, band_id, role) VALUES ($1, $2, $3)", userID, bandID, "admin")
    if err != nil {
        log.Fatalf("Failed to link user to band: %v", err)
    }

    // --- Songs ---
    var song1ID, song2ID, song3ID, songToDeleteID int
    db.QueryRow(ctx, "INSERT INTO songs (band_id, title, album_name, duration_seconds, tempo) VALUES ($1, 'Song Title 1', 'Test Album', 185, 120) RETURNING id", bandID).Scan(&song1ID)
    db.QueryRow(ctx, "INSERT INTO songs (band_id, title, duration_seconds) VALUES ($1, 'Song Title 2', 210) RETURNING id", bandID).Scan(&song2ID)
    db.QueryRow(ctx, "INSERT INTO songs (band_id, title, album_name, duration_seconds) VALUES ($1, 'Another Song To Add', 'Test Album', 150) RETURNING id", bandID).Scan(&song3ID)
    db.QueryRow(ctx, "INSERT INTO songs (band_id, title, album_name, duration_seconds) VALUES ($1, 'Song To Delete', 'Test Album', 100) RETURNING id", bandID).Scan(&songToDeleteID)

    // --- Interludes ---
    var interlude1ID, interlude2ID int
    db.QueryRow(ctx, "INSERT INTO interludes (band_id, title, speaker, script, duration_seconds) VALUES ($1, 'Test Interlude', 'Le Chanteur', 'Quelques mots pour le public ici.', 60) RETURNING id", bandID).Scan(&interlude1ID)
    db.QueryRow(ctx, "INSERT INTO interludes (band_id, title, speaker, script, duration_seconds) VALUES ($1, 'Interlude To Add', 'Le Bassiste', 'Solo de basse.', 45) RETURNING id", bandID).Scan(&interlude2ID)

    // --- Setlist 1 (pour les tests de détail et de suppression d'item) ---
    var setlist1ID int
    err = db.QueryRow(ctx, "INSERT INTO setlists (band_id, name, color) VALUES ($1, 'Test Setlist', '#ff0000') RETURNING id", bandID).Scan(&setlist1ID)
    if err != nil {
        log.Fatalf("Failed to seed setlist 1: %v", err)
    }
    db.Exec(ctx, "INSERT INTO setlist_items (setlist_id, position, item_type, song_id, notes) VALUES ($1, 0, 'song', $2, 'Commencer avec le riff de guitare.')", setlist1ID, song1ID)
    db.Exec(ctx, "INSERT INTO setlist_items (setlist_id, position, item_type, interlude_id) VALUES ($1, 1, 'interlude', $2)", setlist1ID, interlude1ID)
    db.Exec(ctx, "INSERT INTO setlist_items (setlist_id, position, item_type, song_id) VALUES ($1, 2, 'song', $2)", setlist1ID, song2ID)
}
