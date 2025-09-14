package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"setlist/api/model"
	"setlist/api/repository"
)

type SetlistService struct {
	SetlistRepo repository.SetlistRepository
}

type CreateSetlistPayload struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type UpdateSetlistPayload struct {
	Name  string `json:"name"`
	Color string `json:"color"`
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

var hexColorRegex = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$`)

func (s SetlistService) Create(ctx context.Context, payload CreateSetlistPayload, bandID int) (model.Setlist, error) {
	if payload.Name == "" {
		return model.Setlist{}, errors.New("setlist name cannot be empty")
	}

	if payload.Color == "" {
		return model.Setlist{}, errors.New("color cannot be empty")
	}

	if !hexColorRegex.MatchString(payload.Color) {
		return model.Setlist{}, fmt.Errorf("invalid color format: %s (must be hex like #RRGGBB or #RGB)", payload.Color)
	}

	return s.SetlistRepo.CreateSetlist(ctx, payload.Name, payload.Color, bandID)
}

func (s SetlistService) Update(ctx context.Context, id int, bandID int, payload UpdateSetlistPayload) (model.Setlist, error) {
	if payload.Name == "" {
		return model.Setlist{}, errors.New("setlist name cannot be empty")
	}

	if payload.Color == "" {
		return model.Setlist{}, errors.New("color cannot be empty")
	}

	if !hexColorRegex.MatchString(payload.Color) {
		return model.Setlist{}, fmt.Errorf("invalid color format: %s", payload.Color)
	}

	return s.SetlistRepo.UpdateSetlist(ctx, id, bandID, payload.Name, payload.Color)
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

	return SetlistDetails{
		Setlist: setlist,
		Items:   items,
	}, nil
}

func (s SetlistService) AddItem(ctx context.Context, setlistID int, payload AddItemPayload) (model.SetlistItem, error) {
	item := model.SetlistItem{
		SetlistID: setlistID,
		ItemType:  payload.ItemType,
		Notes:     sql.NullString{String: payload.Notes, Valid: payload.Notes != ""},
	}

	if payload.ItemType == "song" {
		item.SongID = sql.NullInt32{Int32: int32(payload.ItemID), Valid: true}
	} else if payload.ItemType == "interlude" {
		item.InterludeID = sql.NullInt32{Int32: int32(payload.ItemID), Valid: true}
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
	notes := sql.NullString{String: payload.Notes, Valid: payload.Notes != ""}
	return s.SetlistRepo.UpdateSetlistItem(ctx, itemID, bandID, notes)
}

func (s SetlistService) DeleteItem(ctx context.Context, itemID int, bandID int) error {
	return s.SetlistRepo.DeleteSetlistItem(ctx, itemID, bandID)
}
