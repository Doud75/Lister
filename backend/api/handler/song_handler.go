package handler

import (
	"encoding/json"
	"net/http"
	"setlist/api/middleware"
	"setlist/api/service"
	"strconv"
)

type SongHandler struct {
	SongService service.SongService
}

func (h SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)

	var payload service.CreateSongPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdSong, err := h.SongService.Create(r.Context(), payload, bandID)
	if err != nil {
		writeError(w, "Failed to create song", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSong)
}

func (h SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	songs, err := h.SongService.GetAllForBand(r.Context(), bandID)
	if err != nil {
		writeError(w, "Failed to retrieve songs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h SongHandler) GetSong(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	id, _ := strconv.Atoi(r.PathValue("id"))

	song, err := h.SongService.GetByID(r.Context(), id, bandID)
	if err != nil {
		writeError(w, "Song not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

func (h SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	id, _ := strconv.Atoi(r.PathValue("id"))

	var payload service.CreateSongPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.SongService.Update(r.Context(), id, bandID, payload); err != nil {
		writeError(w, "Failed to update song", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	id, _ := strconv.Atoi(r.PathValue("id"))

	if err := h.SongService.Delete(r.Context(), id, bandID); err != nil {
		writeError(w, "Failed to delete song", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
