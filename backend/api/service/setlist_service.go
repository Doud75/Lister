package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"setlist/api/middleware"
	"setlist/api/model"
	"setlist/api/repository"
	"setlist/api/validator"
)

type SetlistService struct {
	SetlistRepo   repository.SetlistRepository
	InterludeRepo repository.InterludeRepository
}

type CreateSetlistPayload struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type UpdateSetlistPayload struct {
	Name       *string `json:"name"`
	Color      *string `json:"color"`
	IsArchived *bool   `json:"is_archived"`
}

type SetlistDetails struct {
	model.Setlist
	Items []model.SetlistItem `json:"items"`
}

type AddItemPayload struct {
	ItemType string `json:"item_type"`
	ItemID   int    `json:"item_id"`
	Notes    string `json:"notes"`
}

type UpdateOrderPayload struct {
	ItemIDs []int `json:"item_ids"`
}

type UpdateItemPayload struct {
	Notes string `json:"notes"`
}

type DuplicateSetlistPayload struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

var hexColorRegex = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)

func (s SetlistService) Create(ctx context.Context, payload CreateSetlistPayload, bandID int) (model.Setlist, error) {
	if payload.Name == "" {
		return model.Setlist{}, errors.New("setlist name cannot be empty")
	}
	if payload.Color == "" || !hexColorRegex.MatchString(payload.Color) {
		return model.Setlist{}, fmt.Errorf("invalid color format: %s", payload.Color)
	}

	return s.SetlistRepo.CreateSetlist(ctx, s.SetlistRepo.DB, validator.Sanitize(payload.Name), payload.Color, bandID)
}

func (s SetlistService) Update(ctx context.Context, id int, bandID int, payload UpdateSetlistPayload) (model.Setlist, error) {
	setlist, err := s.SetlistRepo.GetSetlistByID(ctx, id, bandID)
	if err != nil {
		return model.Setlist{}, errors.New("setlist not found or does not belong to user's band")
	}

	if payload.Name != nil {
		if *payload.Name == "" {
			return model.Setlist{}, errors.New("setlist name cannot be empty")
		}
		setlist.Name = validator.Sanitize(*payload.Name)
	}
	if payload.Color != nil {
		if *payload.Color == "" || !hexColorRegex.MatchString(*payload.Color) {
			return model.Setlist{}, fmt.Errorf("invalid color format")
		}
		setlist.Color = *payload.Color
	}
	if payload.IsArchived != nil {
		setlist.IsArchived = *payload.IsArchived
	}

	return s.SetlistRepo.UpdateSetlist(ctx, setlist)
}

func (s SetlistService) Delete(ctx context.Context, setlistID int, bandID int) error {
	_, err := s.SetlistRepo.GetSetlistByID(ctx, setlistID, bandID)
	if err != nil {
		return errors.New("setlist not found or does not belong to the user's band")
	}
	return s.SetlistRepo.DeleteSetlist(ctx, setlistID, bandID)
}

func (s SetlistService) GetAllForBand(ctx context.Context, bandID int) ([]model.Setlist, error) {
	return s.SetlistRepo.GetSetlistsByBandID(ctx, bandID)
}

func (s SetlistService) GetDetails(ctx context.Context, id int, bandID int) (SetlistDetails, error) {
	setlist, err := s.SetlistRepo.GetSetlistByID(ctx, id, bandID)
	if err != nil {
		return SetlistDetails{}, err
	}
	items, err := s.SetlistRepo.GetSetlistItemsBySetlistID(ctx, id)
	if err != nil {
		return SetlistDetails{}, err
	}
	return SetlistDetails{Setlist: setlist, Items: items}, nil
}

func (s SetlistService) AddItem(ctx context.Context, setlistID int, payload AddItemPayload) (model.SetlistItem, error) {
	item := model.SetlistItem{
		SetlistID: setlistID,
		ItemType:  payload.ItemType,
		Notes:     sql.NullString{String: validator.Sanitize(payload.Notes), Valid: payload.Notes != ""},
	}

	if payload.ItemType == "song" {
		item.SongID = sql.NullInt32{Int32: int32(payload.ItemID), Valid: true}
	} else if payload.ItemType == "interlude" {
		item.InterludeID = sql.NullInt32{Int32: int32(payload.ItemID), Valid: true}
		bandID, ok := ctx.Value(middleware.BandIDKey).(int)
		if !ok {
			return model.SetlistItem{}, errors.New("band ID not found in context")
		}
		interlude, err := s.InterludeRepo.GetInterludeByID(ctx, payload.ItemID, bandID)
		if err != nil {
			return model.SetlistItem{}, fmt.Errorf("could not retrieve interlude: %w", err)
		}
		item.Notes = interlude.Script
	} else {
		return model.SetlistItem{}, errors.New("invalid item type")
	}
	return s.SetlistRepo.AddItemToSetlist(ctx, item)
}

func (s SetlistService) UpdateOrder(ctx context.Context, setlistID int, payload UpdateOrderPayload) error {
	if len(payload.ItemIDs) == 0 {
		return nil
	}
	return s.SetlistRepo.UpdateItemOrder(ctx, setlistID, payload.ItemIDs)
}

func (s SetlistService) UpdateItem(ctx context.Context, itemID int, bandID int, payload UpdateItemPayload) (model.SetlistItem, error) {
	notes := sql.NullString{String: validator.Sanitize(payload.Notes), Valid: payload.Notes != ""}
	return s.SetlistRepo.UpdateSetlistItem(ctx, itemID, bandID, notes)
}

func (s SetlistService) DeleteItem(ctx context.Context, itemID int, bandID int) error {
	return s.SetlistRepo.DeleteSetlistItem(ctx, itemID, bandID)
}

func (s SetlistService) Duplicate(ctx context.Context, originalSetlistID int, bandID int, newName, newColor string) (model.Setlist, error) {
	tx, err := s.SetlistRepo.DB.Begin(ctx)
	if err != nil {
		return model.Setlist{}, err
	}
	defer tx.Rollback(ctx)

	_, err = s.SetlistRepo.GetSetlistByID(ctx, originalSetlistID, bandID)
	if err != nil {
		return model.Setlist{}, errors.New("original setlist not found or does not belong to the user's band")
	}

	originalItems, err := s.SetlistRepo.GetSetlistItemsBySetlistID(ctx, originalSetlistID)
	if err != nil {
		return model.Setlist{}, err
	}

	newSetlist, err := s.SetlistRepo.CreateSetlist(ctx, tx, newName, newColor, bandID)
	if err != nil {
		return model.Setlist{}, err
	}

	if err := s.SetlistRepo.CopyItemsToNewSetlist(ctx, tx, newSetlist.ID, originalItems); err != nil {
		return model.Setlist{}, err
	}

	return newSetlist, tx.Commit(ctx)
}
