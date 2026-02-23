package handler

import (
	"encoding/json"
	"net/http"
	"setlist/api/apierror"
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
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	createdSong, err := h.SongService.Create(r.Context(), payload, bandID)
	if err != nil {
		writeAppError(w, apierror.InternalError("création de chanson"))
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
		writeAppError(w, apierror.InternalError("récupération des chansons"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func (h SongHandler) GetSong(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeAppError(w, apierror.InvalidRequest("Identifiant de chanson invalide."))
		return
	}

	song, err := h.SongService.GetByID(r.Context(), id, bandID)
	if err != nil {
		writeAppError(w, apierror.NotFound("Chanson"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

func (h SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeAppError(w, apierror.InvalidRequest("Identifiant de chanson invalide."))
		return
	}

	var payload service.UpdateSongPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeAppError(w, apierror.InvalidRequest("Corps de la requête invalide."))
		return
	}

	updatedSong, err := h.SongService.Update(r.Context(), id, bandID, payload)
	if err != nil {
		writeAppError(w, apierror.InternalError("mise à jour de chanson"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSong)
}

func (h SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	bandID, _ := r.Context().Value(middleware.BandIDKey).(int)
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeAppError(w, apierror.InvalidRequest("Identifiant de chanson invalide."))
		return
	}

	err = h.SongService.SoftDelete(r.Context(), id, bandID)
	if err != nil {
		writeAppError(w, apierror.NotFound("Chanson"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
