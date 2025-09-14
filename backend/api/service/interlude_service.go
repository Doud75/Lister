package service

import (
	"context"
	"database/sql"
	"errors"
	"setlist/api/model"
	"setlist/api/repository"
)

type CreateInterludePayload struct {
	Title           string  `json:"title"`
	Speaker         *string `json:"speaker"`
	Script          *string `json:"script"`
	DurationSeconds *int    `json:"duration_seconds"`
}

type InterludeService struct {
	InterludeRepo repository.InterludeRepository
}

func (s InterludeService) Create(ctx context.Context, payload CreateInterludePayload, bandID int) (model.Interlude, error) {
	if payload.Title == "" {
		return model.Interlude{}, errors.New("interlude title cannot be empty")
	}

	interlude := model.Interlude{
		BandID: bandID,
		Title:  payload.Title,
	}

	if payload.Speaker != nil {
		interlude.Speaker = sql.NullString{String: *payload.Speaker, Valid: true}
	}
	if payload.Script != nil {
		interlude.Script = sql.NullString{String: *payload.Script, Valid: true}
	}
	if payload.DurationSeconds != nil {
		interlude.DurationSeconds = sql.NullInt32{Int32: int32(*payload.DurationSeconds), Valid: true}
	}

	return s.InterludeRepo.CreateInterlude(ctx, interlude)
}

func (s InterludeService) GetAllForBand(ctx context.Context, bandID int) ([]model.Interlude, error) {
	return s.InterludeRepo.GetAllInterludesByBandID(ctx, bandID)
}

func (s InterludeService) GetByID(ctx context.Context, id int, bandID int) (model.Interlude, error) {
	return s.InterludeRepo.GetInterludeByID(ctx, id, bandID)
}

func (s InterludeService) Update(ctx context.Context, id int, bandID int, payload CreateInterludePayload) error {
	if payload.Title == "" {
		return errors.New("interlude title cannot be empty")
	}

	interlude := model.Interlude{
		ID:     id,
		BandID: bandID,
		Title:  payload.Title,
	}

	if payload.Speaker != nil {
		interlude.Speaker = sql.NullString{String: *payload.Speaker, Valid: true}
	}
	if payload.Script != nil {
		interlude.Script = sql.NullString{String: *payload.Script, Valid: true}
	}
	if payload.DurationSeconds != nil {
		interlude.DurationSeconds = sql.NullInt32{Int32: int32(*payload.DurationSeconds), Valid: true}
	}

	return s.InterludeRepo.UpdateInterlude(ctx, interlude)
}

func (s InterludeService) Delete(ctx context.Context, id int, bandID int) error {
	return s.InterludeRepo.DeleteInterlude(ctx, id, bandID)
}
