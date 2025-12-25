package repository

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "setlist/api/model"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
    "github.com/jackc/pgx/v5/pgxpool"
)

var (
    ErrDuplicateBand     = errors.New("band with this name already exists")
    ErrDuplicateUsername = errors.New("username already exists")
)

var PreferredColors = []string{
    "#FDFFB6",
    "#CAFFBF",
    "#9BF6FF",
    "#A0C4FF",
    "#BDB2FF",
    "#FFC6FF",
    "#FFADAD",
    "#FFD6A5",
}

type UserRepository struct {
    DB *pgxpool.Pool
}

func generateRandomPastelColor() string {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    red := 200 + r.Intn(56)
    green := 200 + r.Intn(56)
    blue := 200 + r.Intn(56)

    return fmt.Sprintf("#%02X%02X%02X", red, green, blue)
}

func (r *UserRepository) getNextAvailableColor(ctx context.Context, bandID int) (string, error) {
    rows, err := r.DB.Query(ctx, "SELECT color FROM band_users WHERE band_id = $1", bandID)
    if err != nil {
        return PreferredColors[0], nil
    }
    defer rows.Close()

    usedColors := make(map[string]bool)
    for rows.Next() {
        var c string
        if err := rows.Scan(&c); err == nil {
            usedColors[c] = true
        }
    }

    for _, color := range PreferredColors {
        if !usedColors[color] {
            return color, nil
        }
    }

    for i := 0; i < 5; i++ {
        randomColor := generateRandomPastelColor()
        if !usedColors[randomColor] {
            return randomColor, nil
        }
    }

    return generateRandomPastelColor(), nil
}

func (r *UserRepository) CreateBandAndUser(ctx context.Context, bandName, username, passwordHash string) (model.User, model.Band, error) {
    tx, err := r.DB.Begin(ctx)
    if err != nil {
        return model.User{}, model.Band{}, err
    }
    defer tx.Rollback(ctx)

    var band model.Band
    bandQuery := `INSERT INTO bands (name) VALUES ($1) RETURNING id, name`
    err = tx.QueryRow(ctx, bandQuery, bandName).Scan(&band.ID, &band.Name)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return model.User{}, model.Band{}, ErrDuplicateBand
        }
        return model.User{}, model.Band{}, err
    }

    var user model.User
    userQuery := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, created_at`
    err = tx.QueryRow(ctx, userQuery, username, passwordHash).Scan(&user.ID, &user.Username, &user.CreatedAt)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return model.User{}, model.Band{}, ErrDuplicateUsername
        }
        return model.User{}, model.Band{}, err
    }

    firstColor := PreferredColors[0]
    linkQuery := `INSERT INTO band_users (user_id, band_id, role, color) VALUES ($1, $2, $3, $4)`
    _, err = tx.Exec(ctx, linkQuery, user.ID, band.ID, "admin", firstColor)
    if err != nil {
        return model.User{}, model.Band{}, err
    }

    return user, band, tx.Commit(ctx)
}

func (r *UserRepository) CreateUserAndAddToBand(ctx context.Context, bandID int, username, passwordHash, role string) (model.User, error) {
    color, err := r.getNextAvailableColor(ctx, bandID)
    if err != nil {
        return model.User{}, err
    }

    tx, err := r.DB.Begin(ctx)
    if err != nil {
        return model.User{}, err
    }
    defer tx.Rollback(ctx)

    var user model.User
    userQuery := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, username, created_at`
    err = tx.QueryRow(ctx, userQuery, username, passwordHash).Scan(&user.ID, &user.Username, &user.CreatedAt)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return model.User{}, ErrDuplicateUsername
        }
        return model.User{}, err
    }

    linkQuery := `INSERT INTO band_users (user_id, band_id, role, color) VALUES ($1, $2, $3, $4)`
    _, err = tx.Exec(ctx, linkQuery, user.ID, bandID, role, color)
    if err != nil {
        return model.User{}, err
    }

    return user, tx.Commit(ctx)
}

func (r *UserRepository) GetMembersByBandID(ctx context.Context, bandID int) ([]model.BandMember, error) {
    members := make([]model.BandMember, 0)
    query := `
		SELECT u.id, u.username, bu.role, bu.color
		FROM users u
		JOIN band_users bu ON u.id = bu.user_id
		WHERE bu.band_id = $1
		ORDER BY u.username
	`
    rows, err := r.DB.Query(ctx, query, bandID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var member model.BandMember
        if err := rows.Scan(&member.ID, &member.Username, &member.Role, &member.Color); err != nil {
            return nil, err
        }
        members = append(members, member)
    }
    return members, rows.Err()
}

func (r *UserRepository) RemoveUserFromBand(ctx context.Context, bandID int, userID int) error {
    var adminCount int
    countQuery := `SELECT COUNT(*) FROM band_users WHERE band_id = $1 AND role = 'admin'`
    err := r.DB.QueryRow(ctx, countQuery, bandID).Scan(&adminCount)
    if err != nil {
        return err
    }

    if adminCount <= 1 {
        var userRole string
        roleQuery := `SELECT role FROM band_users WHERE user_id = $1 AND band_id = $2`
        err := r.DB.QueryRow(ctx, roleQuery, userID, bandID).Scan(&userRole)
        if err != nil {
            return err
        }
        if userRole == "admin" {
            return errors.New("cannot remove the last admin of the band")
        }
    }

    query := `DELETE FROM band_users WHERE user_id = $1 AND band_id = $2`
    cmdTag, err := r.DB.Exec(ctx, query, userID, bandID)
    if err != nil {
        return err
    }
    if cmdTag.RowsAffected() == 0 {
        return pgx.ErrNoRows
    }
    return nil
}

func (r *UserRepository) GetUserRoleInBand(ctx context.Context, userID int, bandID int) (string, error) {
    var role string
    query := `SELECT role FROM band_users WHERE user_id = $1 AND band_id = $2`
    err := r.DB.QueryRow(ctx, query, userID, bandID).Scan(&role)
    return role, err
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (model.User, error) {
    var user model.User
    query := `SELECT id, password_hash, username FROM users WHERE username = $1`
    err := r.DB.QueryRow(ctx, query, username).Scan(&user.ID, &user.PasswordHash, &user.Username)
    if err != nil {
        return model.User{}, err
    }
    return user, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, id int) (model.User, error) {
    var user model.User
    query := `SELECT id, password_hash, username FROM users WHERE id = $1`
    err := r.DB.QueryRow(ctx, query, id).Scan(&user.ID, &user.PasswordHash, &user.Username)
    if err != nil {
        return model.User{}, err
    }
    return user, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, newHash string) error {
    query := `UPDATE users SET password_hash = $1 WHERE id = $2`
    cmdTag, err := r.DB.Exec(ctx, query, newHash, userID)
    if err != nil {
        return err
    }
    if cmdTag.RowsAffected() == 0 {
        return errors.New("user not found or no update was needed")
    }
    return nil
}

func (r *UserRepository) FindBandsByUserID(ctx context.Context, userID int) ([]model.Band, error) {
    var bands []model.Band
    query := `
		SELECT b.id, b.name FROM bands b
		JOIN band_users bu ON b.id = bu.band_id
		WHERE bu.user_id = $1
		ORDER BY b.name
	`
    rows, err := r.DB.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var band model.Band
        if err := rows.Scan(&band.ID, &band.Name); err != nil {
            return nil, err
        }
        bands = append(bands, band)
    }

    return bands, rows.Err()
}

func (r *UserRepository) IsUserInBand(ctx context.Context, userID int, bandID int) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM band_users WHERE user_id = $1 AND band_id = $2)`
    err := r.DB.QueryRow(ctx, query, userID, bandID).Scan(&exists)
    return exists, err
}

func (r *UserRepository) AddUserToBand(ctx context.Context, userID, bandID int, role string) error {
    color, err := r.getNextAvailableColor(ctx, bandID)
    if err != nil {
        return err
    }

    query := `INSERT INTO band_users (user_id, band_id, role, color) VALUES ($1, $2, $3, $4)`
    _, err = r.DB.Exec(ctx, query, userID, bandID, role, color)
    return err
}

func (r *UserRepository) SearchUsersByUsername(ctx context.Context, usernameQuery string) ([]model.User, error) {
    users := make([]model.User, 0)
    query := `SELECT id, username FROM users WHERE username ILIKE $1 LIMIT 10`

    rows, err := r.DB.Query(ctx, query, usernameQuery+"%")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var user model.User
        if err := rows.Scan(&user.ID, &user.Username); err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, rows.Err()
}
