package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(redisURL string) *redis.Client {
	if redisURL == "" {
		log.Println("[cache] REDIS_URL non définie — cache Redis désactivé (mode fail-safe).")
		return nil
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("[cache] URL Redis invalide (%s) — cache désactivé : %v\n", redisURL, err)
		return nil
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("[cache] Impossible de joindre Redis — cache désactivé : %v\n", err)
		client.Close()
		return nil
	}

	log.Println("[cache] Connexion Redis établie ✓")
	return client
}

func SongKey(bandID int) string {
	return fmt.Sprintf("band:%d:songs", bandID)
}

func ProfileKey(userID int, bandID int) string {
	return fmt.Sprintf("user:%d:band:%d:profile", userID, bandID)
}

func SetlistKey(bandID int) string {
	return fmt.Sprintf("band:%d:setlists", bandID)
}

func Get(ctx context.Context, client *redis.Client, key string) (string, bool) {
	if client == nil {
		return "", false
	}
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func Set(ctx context.Context, client *redis.Client, key string, value string, ttl time.Duration) {
	if client == nil {
		return
	}
	if err := client.Set(ctx, key, value, ttl).Err(); err != nil {
		log.Printf("[cache] Erreur lors de l'écriture de la clé %s : %v\n", key, err)
	}
}

func Delete(ctx context.Context, client *redis.Client, key string) {
	if client == nil {
		return
	}
	if err := client.Del(ctx, key).Err(); err != nil {
		log.Printf("[cache] Erreur lors de la suppression de la clé %s : %v\n", key, err)
	}
}
