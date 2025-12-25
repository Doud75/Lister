// backend/db/seed/main.go

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

    checkErr := func(err error, message string) {
        if err != nil {
            log.Fatalf("%s: %v", message, err)
        }
    }

    // --- Scénario de base pour les tests ---
    var bandID int
    err := db.QueryRow(ctx, "INSERT INTO bands (name) VALUES ('The Testers') RETURNING id").Scan(&bandID)
    checkErr(err, "Failed to seed band 'The Testers'")

    // Admin user "testuser"
    hashedPasswordAdmin, _ := auth.HashPassword("password123")
    var adminUserID int
    err = db.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
        "testuser", hashedPasswordAdmin).Scan(&adminUserID)
    checkErr(err, "Failed to seed admin user 'testuser'")
    _, err = db.Exec(ctx, "INSERT INTO band_users (user_id, band_id, role, color) VALUES ($1, $2, $3, $4)", adminUserID, bandID, "admin", "#FDFFB6")
    checkErr(err, "Failed to link admin user to band")

    // Member user "memberuser"
    hashedPasswordMember, _ := auth.HashPassword("password123")
    var memberUserID int
    err = db.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
        "memberuser", hashedPasswordMember).Scan(&memberUserID)
    checkErr(err, "Failed to seed member user 'memberuser'")
    _, err = db.Exec(ctx, "INSERT INTO band_users (user_id, band_id, role, color) VALUES ($1, $2, $3, $4)", memberUserID, bandID, "member", "#CAFFBF")
    checkErr(err, "Failed to link member user to band")

    // Utilisateur dédié pour le test de changement de mot de passe
    passwordChangeUserHash, _ := auth.HashPassword("password123")
    var passwordChangeUserID int
    err = db.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
        "passwordChangeUser", passwordChangeUserHash).Scan(&passwordChangeUserID)
    checkErr(err, "Failed to seed passwordChangeUser")
    _, err = db.Exec(ctx, "INSERT INTO band_users (user_id, band_id, role, color) VALUES ($1, $2, $3, $4)", passwordChangeUserID, bandID, "member", "#9BF6FF")
    checkErr(err, "Failed to link passwordChangeUser to band")

    // Songs
    var song1ID, song2ID, song3ID, songToDeleteID int
    err = db.QueryRow(ctx, "INSERT INTO songs (band_id, title, album_name, duration_seconds, tempo) VALUES ($1, 'Song Title 1', 'Test Album', 185, 120) RETURNING id", bandID).Scan(&song1ID)
    checkErr(err, "Failed to seed song 1")
    err = db.QueryRow(ctx, "INSERT INTO songs (band_id, title, duration_seconds) VALUES ($1, 'Song Title 2', 210) RETURNING id", bandID).Scan(&song2ID)
    checkErr(err, "Failed to seed song 2")
    err = db.QueryRow(ctx, "INSERT INTO songs (band_id, title, album_name, duration_seconds) VALUES ($1, 'Another Song To Add', 'Test Album', 150) RETURNING id", bandID).Scan(&song3ID)
    checkErr(err, "Failed to seed song 3")
    err = db.QueryRow(ctx, "INSERT INTO songs (band_id, title, album_name, duration_seconds) VALUES ($1, 'Song To Delete', 'Test Album', 100) RETURNING id", bandID).Scan(&songToDeleteID)
    checkErr(err, "Failed to seed song to delete")

    // Interludes
    var interlude1ID, interlude2ID int
    err = db.QueryRow(ctx, "INSERT INTO interludes (band_id, title, speaker, script, duration_seconds) VALUES ($1, 'Test Interlude', 'Le Chanteur', 'Quelques mots pour le public ici.', 60) RETURNING id", bandID).Scan(&interlude1ID)
    checkErr(err, "Failed to seed interlude 1")
    err = db.QueryRow(ctx, "INSERT INTO interludes (band_id, title, speaker, script, duration_seconds) VALUES ($1, 'Interlude To Add', 'Le Bassiste', 'Solo de basse.', 45) RETURNING id", bandID).Scan(&interlude2ID)
    checkErr(err, "Failed to seed interlude 2")

    // Setlist
    var setlist1ID int
    err = db.QueryRow(ctx, "INSERT INTO setlists (band_id, name, color) VALUES ($1, 'Test Setlist', '#ff0000') RETURNING id", bandID).Scan(&setlist1ID)
    checkErr(err, "Failed to seed setlist 1")

    _, err = db.Exec(ctx, "INSERT INTO setlist_items (setlist_id, position, item_type, song_id, notes) VALUES ($1, 0, 'song', $2, 'Commencer avec le riff de guitare.')", setlist1ID, song1ID)
    checkErr(err, "Failed to seed setlist item 1")
    _, err = db.Exec(ctx, "INSERT INTO setlist_items (setlist_id, position, item_type, interlude_id) VALUES ($1, 1, 'interlude', $2)", setlist1ID, interlude1ID)
    checkErr(err, "Failed to seed setlist item 2")
    _, err = db.Exec(ctx, "INSERT INTO setlist_items (setlist_id, position, item_type, song_id) VALUES ($1, 2, 'song', $2)", setlist1ID, song2ID)
    checkErr(err, "Failed to seed setlist item 3")

    // --- Scénario Multi-Groupes ---
    log.Println("Seeding data for multi-group test...")

    multiGroupHashedPassword, _ := auth.HashPassword("password123")
    var multiGroupUserID int
    err = db.QueryRow(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id", "multiGroupUser", multiGroupHashedPassword).Scan(&multiGroupUserID)
    checkErr(err, "Failed to seed multiGroupUser")

    var bandAID, bandBID int
    err = db.QueryRow(ctx, "INSERT INTO bands (name) VALUES ('Band A') RETURNING id").Scan(&bandAID)
    checkErr(err, "Failed to seed Band A")
    err = db.QueryRow(ctx, "INSERT INTO bands (name) VALUES ('Band B') RETURNING id").Scan(&bandBID)
    checkErr(err, "Failed to seed Band B")

    _, err = db.Exec(ctx, "INSERT INTO band_users (user_id, band_id, role, color) VALUES ($1, $2, 'admin', '#A0C4FF'), ($1, $3, 'member', '#BDB2FF')", multiGroupUserID, bandAID, bandBID)
    checkErr(err, "Failed to link multiGroupUser to bands")

    _, err = db.Exec(ctx, "INSERT INTO songs (band_id, title) VALUES ($1, 'Chanson A1')", bandAID)
    checkErr(err, "Failed to seed song for Band A")
    _, err = db.Exec(ctx, "INSERT INTO setlists (band_id, name, color) VALUES ($1, 'Setlist A', '#ff0000')", bandAID)
    checkErr(err, "Failed to seed setlist for Band A")

    _, err = db.Exec(ctx, "INSERT INTO songs (band_id, title) VALUES ($1, 'Chanson B1')", bandBID)
    checkErr(err, "Failed to seed song for Band B")
    _, err = db.Exec(ctx, "INSERT INTO setlists (band_id, name, color) VALUES ($1, 'Setlist B', '#0000ff')", bandBID)
    checkErr(err, "Failed to seed setlist for Band B")
}
