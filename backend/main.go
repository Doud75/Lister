package main

import (
	"fmt"
	"log"
	"net/http"
	"setlist/api/handler"
	"setlist/api/middleware"
	"setlist/api/repository"
	"setlist/api/service"
	"setlist/config"
	"setlist/db"
)

func main() {
	cfg := config.Load()
	dbPool := db.NewConnection(cfg.DatabaseURL)
	defer dbPool.Close()

	userRepo := repository.UserRepository{DB: dbPool}
	refreshTokenRepo := repository.RefreshTokenRepository{DB: dbPool}
	userService := service.UserService{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshTokenRepo,
		JWTSecret:        cfg.JWTSecret,
	}
	authService := service.AuthService{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshTokenRepo,
		JWTSecret:        cfg.JWTSecret,
	}
	userHandler := handler.UserHandler{UserService: userService}
	authHandler := handler.AuthHandler{AuthService: authService}
	bandHandler := handler.BandHandler{UserService: userService}

	interludeRepo := repository.InterludeRepository{DB: dbPool}
	interludeService := service.InterludeService{InterludeRepo: interludeRepo}
	interludeHandler := handler.InterludeHandler{InterludeService: interludeService}

	infoRepo := repository.InfoRepository{DB: dbPool}
	infoHandler := handler.InfoHandler{InfoRepo: infoRepo, UserRepo: userRepo}

	setlistRepo := repository.SetlistRepository{DB: dbPool}
	setlistService := service.SetlistService{SetlistRepo: setlistRepo, InterludeRepo: interludeRepo}
	setlistHandler := handler.SetlistHandler{SetlistService: setlistService}

	songRepo := repository.SongRepository{DB: dbPool}
	songService := service.SongService{SongRepo: songRepo}
	songHandler := handler.SongHandler{SongService: songService}

	authMiddleware := middleware.JWTAuth(cfg.JWTSecret, userRepo)
	adminMiddleware := middleware.AdminOnly(userRepo)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitEnabled)

	mux := http.NewServeMux()

	mux.Handle("/api/auth/login", rateLimiter.LimitMiddleware(http.HandlerFunc(userHandler.Login)))
	mux.Handle("/api/auth/signup", rateLimiter.LimitMiddleware(http.HandlerFunc(userHandler.Signup)))
	mux.HandleFunc("/api/auth/refresh", authHandler.RefreshToken)
	mux.HandleFunc("/api/auth/logout", authHandler.Logout)
	mux.Handle("PUT /api/user/password", authMiddleware(http.HandlerFunc(userHandler.UpdatePassword)))
	mux.Handle("GET /api/user/info", authMiddleware(http.HandlerFunc(infoHandler.GetCurrentUserInfo)))
	mux.Handle("GET /api/user/search", authMiddleware(http.HandlerFunc(userHandler.SearchUsers)))

	mux.Handle("GET /api/bands/{bandId}/members", authMiddleware(http.HandlerFunc(bandHandler.GetMembers)))
	mux.Handle("POST /api/bands/{bandId}/members", authMiddleware(adminMiddleware(http.HandlerFunc(bandHandler.InviteMember))))
	mux.Handle("DELETE /api/bands/{bandId}/members/{userId}", authMiddleware(adminMiddleware(http.HandlerFunc(bandHandler.RemoveMember))))

	mux.Handle("POST /api/setlist", authMiddleware(http.HandlerFunc(setlistHandler.CreateSetlist)))
	mux.Handle("GET /api/setlist", authMiddleware(http.HandlerFunc(setlistHandler.GetSetlists)))
	mux.Handle("GET /api/setlist/{id}", authMiddleware(http.HandlerFunc(setlistHandler.GetSetlistDetails)))
	mux.Handle("PUT /api/setlist/{id}", authMiddleware(adminMiddleware(http.HandlerFunc(setlistHandler.UpdateSetlist))))
	mux.Handle("DELETE /api/setlist/{id}", authMiddleware(adminMiddleware(http.HandlerFunc(setlistHandler.DeleteSetlist))))

	mux.Handle("POST /api/setlist/{id}/duplicate", authMiddleware(adminMiddleware(http.HandlerFunc(setlistHandler.DuplicateSetlist))))
	mux.Handle("POST /api/setlist/{id}/items", authMiddleware(http.HandlerFunc(setlistHandler.AddItem)))
	mux.Handle("PUT /api/setlist/{id}/items/order", authMiddleware(http.HandlerFunc(setlistHandler.UpdateItemOrder)))
	mux.Handle("PUT /api/setlist/item/{itemId}", authMiddleware(http.HandlerFunc(setlistHandler.UpdateItem)))
	mux.Handle("DELETE /api/setlist/item/{itemId}", authMiddleware(http.HandlerFunc(setlistHandler.DeleteItem)))

	mux.Handle("POST /api/song", authMiddleware(http.HandlerFunc(songHandler.CreateSong)))
	mux.Handle("GET /api/song", authMiddleware(http.HandlerFunc(songHandler.GetSongs)))
	mux.Handle("GET /api/song/{id}", authMiddleware(http.HandlerFunc(songHandler.GetSong)))
	mux.Handle("PUT /api/song/{id}", authMiddleware(http.HandlerFunc(songHandler.UpdateSong)))
	mux.Handle("DELETE /api/song/{id}", authMiddleware(http.HandlerFunc(songHandler.DeleteSong)))

	mux.Handle("POST /api/interlude", authMiddleware(http.HandlerFunc(interludeHandler.CreateInterlude)))
	mux.Handle("GET /api/interlude", authMiddleware(http.HandlerFunc(interludeHandler.GetInterludes)))
	mux.Handle("PUT /api/interlude/{id}", authMiddleware(http.HandlerFunc(interludeHandler.UpdateInterlude)))

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})

	port := "8089"
	address := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Backend server starting on %s\n", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
