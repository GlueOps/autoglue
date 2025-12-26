package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ListActions godoc
//
//	@ID				ListActions
//	@Summary		List available actions
//	@Description	Returns all admin-configured actions.
//	@Tags			Actions
//	@Produce		json
//	@Success		200	{array}		dto.ActionResponse
//	@Failure		401	{string}	string	"Unauthorized"
//	@Failure		500	{string}	string	"db error"
//	@Router			/admin/actions [get]
//	@Security		BearerAuth
func ListActions(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rows []models.Action
		if err := db.Order("label ASC").Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		out := make([]dto.ActionResponse, 0, len(rows))
		for _, a := range rows {
			out = append(out, actionToDTO(a))
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetAction godoc
//
//	@ID				GetAction
//	@Summary		Get a single action by ID
//	@Description	Returns a single action.
//	@Tags			Actions
//	@Produce		json
//	@Param			actionID	path		string	true	"Action ID"
//	@Success		200			{object}	dto.ActionResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/admin/actions/{actionID} [get]
//	@Security		BearerAuth
func GetAction(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actionID, err := uuid.Parse(chi.URLParam(r, "actionID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_action_id", "invalid action id")
			return
		}

		var row models.Action
		if err := db.Where("id = ?", actionID).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "action not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, actionToDTO(row))
	}
}

// CreateAction godoc
//
//	@ID				CreateAction
//	@Summary		Create an action
//	@Description	Creates a new admin-configured action.
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.CreateActionRequest	true	"payload"
//	@Success		201		{object}	dto.ActionResponse
//	@Failure		400		{string}	string	"bad request"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{string}	string	"db error"
//	@Router			/admin/actions [post]
//	@Security		BearerAuth
func CreateAction(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in dto.CreateActionRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		label := strings.TrimSpace(in.Label)
		desc := strings.TrimSpace(in.Description)
		target := strings.TrimSpace(in.MakeTarget)

		if label == "" {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", "label is required")
			return
		}
		if desc == "" {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", "description is required")
			return
		}
		if target == "" {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", "make_target is required")
			return
		}

		row := models.Action{
			Label:       label,
			Description: desc,
			MakeTarget:  target,
		}

		if err := db.Create(&row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		utils.WriteJSON(w, http.StatusCreated, actionToDTO(row))
	}
}

// UpdateAction godoc
//
//	@ID				UpdateAction
//	@Summary		Update an action
//	@Description	Updates an action. Only provided fields are modified.
//	@Tags			Actions
//	@Accept			json
//	@Produce		json
//	@Param			actionID	path		string					true	"Action ID"
//	@Param			body		body		dto.UpdateActionRequest	true	"payload"
//	@Success		200			{object}	dto.ActionResponse
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/admin/actions/{actionID} [patch]
//	@Security		BearerAuth
func UpdateAction(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actionID, err := uuid.Parse(chi.URLParam(r, "actionID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_action_id", "invalid action id")
			return
		}

		var in dto.UpdateActionRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		var row models.Action
		if err := db.Where("id = ?", actionID).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "action not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		if in.Label != nil {
			v := strings.TrimSpace(*in.Label)
			if v == "" {
				utils.WriteError(w, http.StatusBadRequest, "validation_error", "label cannot be empty")
				return
			}
			row.Label = v
		}
		if in.Description != nil {
			v := strings.TrimSpace(*in.Description)
			if v == "" {
				utils.WriteError(w, http.StatusBadRequest, "validation_error", "description cannot be empty")
				return
			}
			row.Description = v
		}
		if in.MakeTarget != nil {
			v := strings.TrimSpace(*in.MakeTarget)
			if v == "" {
				utils.WriteError(w, http.StatusBadRequest, "validation_error", "make_target cannot be empty")
				return
			}
			row.MakeTarget = v
		}

		if err := db.Save(&row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		utils.WriteJSON(w, http.StatusOK, actionToDTO(row))
	}
}

// DeleteAction godoc
//
//	@ID				DeleteAction
//	@Summary		Delete an action
//	@Description	Deletes an action.
//	@Tags			Actions
//	@Produce		json
//	@Param			actionID	path		string	true	"Action ID"
//	@Success		204			{string}	string	"deleted"
//	@Failure		400			{string}	string	"bad request"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		404			{string}	string	"not found"
//	@Failure		500			{string}	string	"db error"
//	@Router			/admin/actions/{actionID} [delete]
//	@Security		BearerAuth
func DeleteAction(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actionID, err := uuid.Parse(chi.URLParam(r, "actionID"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_action_id", "invalid action id")
			return
		}

		tx := db.Where("id = ?", actionID).Delete(&models.Action{})
		if tx.Error != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		if tx.RowsAffected == 0 {
			utils.WriteError(w, http.StatusNotFound, "not_found", "action not found")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func actionToDTO(a models.Action) dto.ActionResponse {
	return dto.ActionResponse{
		ID:          a.ID,
		Label:       a.Label,
		Description: a.Description,
		MakeTarget:  a.MakeTarget,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}
