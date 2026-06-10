package handler

import (
	"errors"
	"net/http"
	"setlist/api/apierror"
	"setlist/api/service"
)

type SongHandler struct {
	SongService service.SongService
}

// mapSongError translates the song service's sentinel errors into typed API
// errors; anything else is reported as an internal error on the operation.
func mapSongError(err error, operation string) error {
	switch {
	case errors.Is(err, service.ErrSongNotFound):
		return apierror.NotFound("Chanson")
	case errors.Is(err, service.ErrSongTitleRequired):
		return apierror.ValidationFailed("Le titre de la chanson est requis.")
	default:
		return apierror.InternalError(operation)
	}
}

func (h SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	payload, err := DecodeJSON[service.CreateSongPayload](r)
	if err != nil {
		return err
	}

	createdSong, err := h.SongService.Create(r.Context(), payload, bandID)
	if err != nil {
		return mapSongError(err, "création de chanson")
	}

	RespondCreated(w, createdSong)
	return nil
}

func (h SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	songs, err := h.SongService.GetAllForBand(r.Context(), bandID)
	if err != nil {
		return apierror.InternalError("récupération des chansons")
	}

	RespondOK(w, songs)
	return nil
}

func (h SongHandler) GetSong(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	id, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant de chanson invalide.")
	}

	song, err := h.SongService.GetByID(r.Context(), id, bandID)
	if err != nil {
		return mapSongError(err, "récupération de la chanson")
	}

	RespondOK(w, song)
	return nil
}

func (h SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	id, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant de chanson invalide.")
	}

	payload, err := DecodeJSON[service.UpdateSongPayload](r)
	if err != nil {
		return err
	}

	updatedSong, err := h.SongService.Update(r.Context(), id, bandID, payload)
	if err != nil {
		return mapSongError(err, "mise à jour de chanson")
	}

	RespondOK(w, updatedSong)
	return nil
}

func (h SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	id, err := GetIntParam(r, "id")
	if err != nil {
		return apierror.InvalidRequest("Identifiant de chanson invalide.")
	}

	if err := h.SongService.SoftDelete(r.Context(), id, bandID); err != nil {
		return mapSongError(err, "suppression de chanson")
	}

	RespondNoContent(w)
	return nil
}
