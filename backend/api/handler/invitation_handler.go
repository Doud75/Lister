package handler

import (
	"net/http"
	"setlist/api/apierror"
	"setlist/api/repository"
	"setlist/api/service"
	"strings"
)

type InvitationHandler struct {
	InvitationService service.InvitationService
}

type CreateInvitationPayload struct {
	Role string `json:"role"`
}

type CreateInvitationResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

func (h InvitationHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) error {
	bandID, err := GetBandID(r)
	if err != nil {
		return err
	}

	payload, err := DecodeJSON[CreateInvitationPayload](r)
	if err != nil {
		payload = CreateInvitationPayload{Role: "member"}
	}

	if payload.Role == "" {
		payload.Role = "member"
	}

	token, expiresAt, err := h.InvitationService.CreateInvitation(r.Context(), bandID, payload.Role)
	if err != nil {
		return apierror.InternalError("création de l'invitation")
	}

	RespondCreated(w, CreateInvitationResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
	})
	return nil
}

type GetInvitationResponse struct {
	Valid    bool   `json:"valid"`
	BandName string `json:"band_name"`
	Role     string `json:"role,omitempty"`
}

func (h InvitationHandler) GetInvitation(w http.ResponseWriter, r *http.Request) error {
	token := r.PathValue("token")
	if token == "" {
		return apierror.InvalidRequest("Token manquant")
	}

	detail, err := h.InvitationService.GetInvitationDetails(r.Context(), token)
	if err != nil {
		if err == repository.ErrInvitationNotFound {
			return apierror.NewUserError(apierror.ErrNotFound, "Invitation introuvable", http.StatusNotFound)
		}
		if err == repository.ErrInvitationExpired {
			return apierror.NewUserError(apierror.ErrInvalidRequest, "Invitation expirée", http.StatusGone)
		}
		return apierror.InternalError("récupération de l'invitation")
	}

	RespondOK(w, GetInvitationResponse{
		Valid:    true,
		BandName: detail.BandName,
		Role:     detail.Role,
	})
	return nil
}

type AcceptInvitationResponse struct {
	Success  bool   `json:"success"`
	BandID   int    `json:"band_id"`
	BandName string `json:"band_name"`
}

func (h InvitationHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) error {
	token := r.PathValue("token")
	if token == "" {
		return apierror.InvalidRequest("Token manquant")
	}

	userID, err := GetUserID(r)
	if err != nil {
		return err
	}

	detail, err := h.InvitationService.AcceptInvitation(r.Context(), token, userID)
	if err != nil {
		if strings.Contains(err.Error(), "conflict: user is already a member") {
			return apierror.NewUserError(apierror.ErrInvalidRequest, "Déjà membre", http.StatusConflict)
		}
		if err == repository.ErrInvitationNotFound {
			return apierror.NewUserError(apierror.ErrNotFound, "Invitation introuvable", http.StatusNotFound)
		}
		if err == repository.ErrInvitationExpired {
			return apierror.NewUserError(apierror.ErrInvalidRequest, "Invitation expirée", http.StatusGone)
		}
		return apierror.InternalError("acceptation de l'invitation")
	}

	RespondOK(w, AcceptInvitationResponse{
		Success:  true,
		BandID:   detail.BandID,
		BandName: detail.BandName,
	})
	return nil
}
