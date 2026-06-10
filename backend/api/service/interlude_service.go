package service

import (
	"context"
	"errors"
	"setlist/api/model"
	"setlist/api/repository"
)

var (
	ErrInterludeNotFound      = errors.New("interlude not found or does not belong to the user's band")
	ErrInterludeTitleRequired = errors.New("interlude title cannot be empty")
)

type CreateInterludePayload struct {
	Title           string  `json:"title"`
	Speaker         *string `json:"speaker"`
	Script          *string `json:"script"`
	DurationSeconds *int    `json:"duration_seconds"`
}

type UpdateInterludePayload struct {
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
		return model.Interlude{}, ErrInterludeTitleRequired
	}

	interlude := model.Interlude{
		BandID:          bandID,
		Title:           payload.Title,
		Speaker:         payload.Speaker,
		Script:          payload.Script,
		DurationSeconds: ptrInt32(payload.DurationSeconds),
	}

	return s.InterludeRepo.CreateInterlude(ctx, interlude)
}

func (s InterludeService) GetAllForBand(ctx context.Context, bandID int) ([]model.Interlude, error) {
	return s.InterludeRepo.GetAllInterludesByBandID(ctx, bandID)
}

func (s InterludeService) Update(ctx context.Context, id int, bandID int, payload UpdateInterludePayload) (model.Interlude, error) {
	interlude, err := s.InterludeRepo.GetInterludeByID(ctx, id, bandID)
	if err != nil {
		return model.Interlude{}, mapNotFound(err, ErrInterludeNotFound)
	}

	if payload.Title != "" {
		interlude.Title = payload.Title
	} else {
		return model.Interlude{}, ErrInterludeTitleRequired
	}

	if payload.Speaker != nil {
		interlude.Speaker = payload.Speaker
	}
	if payload.DurationSeconds != nil {
		interlude.DurationSeconds = ptrInt32(payload.DurationSeconds)
	}

	return s.InterludeRepo.UpdateInterlude(ctx, interlude)
}
