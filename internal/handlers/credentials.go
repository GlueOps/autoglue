package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ListCredentials godoc
//
//	@ID				ListCredentials
//	@Summary		List credentials (metadata only)
//	@Description	Returns credential metadata for the current org. Secrets are never returned.
//	@Tags			Credentials
//	@Produce		json
//	@Param			X-Org-ID			header		string	false	"Organization ID (UUID)"
//	@Param			credential_provider	query		string	false	"Filter by provider (e.g., aws)"
//	@Param			kind				query		string	false	"Filter by kind (e.g., aws_access_key)"
//	@Param			scope_kind			query		string	false	"Filter by scope kind (credential_provider/service/resource)"
//	@Success		200					{array}		dto.CredentialOut
//	@Failure		401					{string}	string	"Unauthorized"
//	@Failure		403					{string}	string	"organization required"
//	@Failure		500					{string}	string	"internal server error"
//	@Router			/credentials [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListCredentials(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		q := db.Where("organization_id = ?", orgID)
		if v := r.URL.Query().Get("credential_provider"); v != "" {
			q = q.Where("provider = ?", v)
		}
		if v := r.URL.Query().Get("kind"); v != "" {
			q = q.Where("kind = ?", v)
		}
		if v := r.URL.Query().Get("scope_kind"); v != "" {
			q = q.Where("scope_kind = ?", v)
		}

		var rows []models.Credential
		if err := q.Order("updated_at DESC").Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		out := make([]dto.CredentialOut, 0, len(rows))
		for i := range rows {
			out = append(out, credOut(&rows[i]))
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetCredential godoc
//
//	@ID			GetCredential
//	@Summary	Get credential by ID (metadata only)
//	@Tags		Credentials
//	@Produce	json
//	@Param		X-Org-ID	header		string	false	"Organization ID (UUID)"
//	@Param		id			path		string	true	"Credential ID (UUID)"
//	@Success	200			{object}	dto.CredentialOut
//	@Failure	401			{string}	string	"Unauthorized"
//	@Failure	403			{string}	string	"organization required"
//	@Failure	500			{string}	string	"internal server error"
//	@Router		/credentials/{id} [get]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func GetCredential(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_id", "invalid UUID")
			return
		}

		var row models.Credential
		if err := db.Where("organization_id = ? AND id = ?", orgID, id).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "credential not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, credOut(&row))
	}
}

// CreateCredential godoc
//
//	@ID			CreateCredential
//	@Summary	Create a credential (encrypts secret)
//	@Tags		Credentials
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header		string						false	"Organization ID (UUID)"
//	@Param		body		body		dto.CreateCredentialRequest	true	"Credential payload"
//	@Success	201			{object}	dto.CredentialOut
//	@Failure	401			{string}	string	"Unauthorized"
//	@Failure	403			{string}	string	"organization required"
//	@Failure	500			{string}	string	"internal server error"
//	@Router		/credentials [post]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func CreateCredential(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		var in dto.CreateCredentialRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		if err := dto.Validate.Struct(in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}

		cred, err := SaveCredentialWithScope(
			r.Context(), db, orgID,
			in.CredentialProvider, in.Kind, in.SchemaVersion,
			in.ScopeKind, in.ScopeVersion, json.RawMessage(in.Scope), json.RawMessage(in.Secret),
			in.Name, in.AccountID, in.Region,
		)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "save_failed", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusCreated, credOut(cred))
	}
}

// UpdateCredential godoc
//
//	@ID			UpdateCredential
//	@Summary	Update credential metadata and/or rotate secret
//	@Tags		Credentials
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header		string						false	"Organization ID (UUID)"
//	@Param		id			path		string						true	"Credential ID (UUID)"
//	@Param		body		body		dto.UpdateCredentialRequest	true	"Fields to update"
//	@Success	200			{object}	dto.CredentialOut
//	@Failure	403			{string}	string	"X-Org-ID required"
//	@Failure	404			{string}	string	"not found"
//	@Router		/credentials/{id} [patch]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func UpdateCredential(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_id", "invalid UUID")
			return
		}

		var row models.Credential
		if err := db.Where("organization_id = ? AND id = ?", orgID, id).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "credential not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		var in dto.UpdateCredentialRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}

		// Update metadata
		if in.Name != nil {
			row.Name = *in.Name
		}
		if in.AccountID != nil {
			row.AccountID = *in.AccountID
		}
		if in.Region != nil {
			row.Region = *in.Region
		}

		// Update scope (re-validate + fingerprint)
		if in.ScopeKind != nil || in.Scope != nil || in.ScopeVersion != nil {
			newKind := row.ScopeKind
			if in.ScopeKind != nil {
				newKind = *in.ScopeKind
			}
			newVersion := row.ScopeVersion
			if in.ScopeVersion != nil {
				newVersion = *in.ScopeVersion
			}
			if in.Scope == nil {
				utils.WriteError(w, http.StatusBadRequest, "validation_error", "scope must be provided when changing scope kind/version")
				return
			}
			prScopes := dto.ScopeRegistry[row.Provider]
			kScopes := prScopes[newKind]
			sdef := kScopes[newVersion]
			dst := sdef.New()
			if err := json.Unmarshal(*in.Scope, dst); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid_scope_json", err.Error())
				return
			}
			if err := sdef.Validate(dst); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid_scope", err.Error())
				return
			}
			canonScope, err := canonicalJSON(dst)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "canon_error", err.Error())
				return
			}
			row.Scope = canonScope
			row.ScopeKind = newKind
			row.ScopeVersion = newVersion
			row.ScopeFingerprint = sha256Hex(canonScope)
		}

		// Rotate secret
		if in.Secret != nil {
			// validate against current Provider/Kind/SchemaVersion
			def := dto.CredentialRegistry[row.Provider][row.Kind][row.SchemaVersion]
			dst := def.New()
			if err := json.Unmarshal(*in.Secret, dst); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid_secret_json", err.Error())
				return
			}
			if err := def.Validate(dst); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid_secret", err.Error())
				return
			}
			canonSecret, err := canonicalJSON(dst)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "canon_error", err.Error())
				return
			}
			cipher, iv, tag, err := utils.EncryptForOrg(orgID, canonSecret, db)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, "encrypt_error", err.Error())
				return
			}
			row.EncryptedData = cipher
			row.IV = iv
			row.Tag = tag
		}

		if err := db.Save(&row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, credOut(&row))
	}
}

// DeleteCredential godoc
//
//	@ID			DeleteCredential
//	@Summary	Delete credential
//	@Tags		Credentials
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization ID (UUID)"
//	@Param		id			path	string	true	"Credential ID (UUID)"
//	@Success	204
//	@Failure	404	{string}	string	"not found"
//	@Router		/credentials/{id} [delete]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func DeleteCredential(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_id", "invalid UUID")
			return
		}
		res := db.Where("organization_id = ? AND id = ?", orgID, id).Delete(&models.Credential{})
		if res.Error != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", res.Error.Error())
			return
		}
		if res.RowsAffected == 0 {
			utils.WriteError(w, http.StatusNotFound, "not_found", "credential not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// RevealCredential godoc
//
//	@ID			RevealCredential
//	@Summary	Reveal decrypted secret (one-time read)
//	@Tags		Credentials
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header		string	false	"Organization ID (UUID)"
//	@Param		id			path		string	true	"Credential ID (UUID)"
//	@Success	200			{object}	map[string]any
//	@Failure	403			{string}	string	"organization required"
//	@Failure	404			{string}	string	"not found"
//	@Router		/credentials/{id}/reveal [post]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func RevealCredential(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_id", "invalid UUID")
			return
		}

		var row models.Credential
		if err := db.Where("organization_id = ? AND id = ?", orgID, id).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "credential not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		plain, err := utils.DecryptForOrg(orgID, row.EncryptedData, row.IV, row.Tag, db)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "decrypt_error", err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, plain)
	}
}

// -- Helpers

func canonicalJSON(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return marshalSorted(m)
}

func marshalSorted(v any) ([]byte, error) {
	switch vv := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(vv))
		for k := range vv {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		buf := bytes.NewBufferString("{")
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			kb, _ := json.Marshal(k)
			buf.Write(kb)
			buf.WriteByte(':')
			b, err := marshalSorted(vv[k])
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}
		buf.WriteByte('}')
		return buf.Bytes(), nil
	case []any:
		buf := bytes.NewBufferString("[")
		for i, e := range vv {
			if i > 0 {
				buf.WriteByte(',')
			}
			b, err := marshalSorted(e)
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}
		buf.WriteByte(']')
		return buf.Bytes(), nil
	default:
		return json.Marshal(v)
	}
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// SaveCredentialWithScope validates secret+scope, encrypts, fingerprints, and stores.
func SaveCredentialWithScope(
	ctx context.Context,
	db *gorm.DB,
	orgID uuid.UUID,
	provider, kind string,
	schemaVersion int,
	scopeKind string,
	scopeVersion int,
	rawScope json.RawMessage,
	rawSecret json.RawMessage,
	name, accountID, region string,
) (*models.Credential, error) {
	// 1) secret shape
	pv, ok := dto.CredentialRegistry[provider]
	if !ok {
		return nil, fmt.Errorf("unknown provider %q", provider)
	}
	kv, ok := pv[kind]
	if !ok {
		return nil, fmt.Errorf("unknown kind %q for provider %q", kind, provider)
	}
	def, ok := kv[schemaVersion]
	if !ok {
		return nil, fmt.Errorf("unsupported schema version %d for %s/%s", schemaVersion, provider, kind)
	}

	secretDst := def.New()
	if err := json.Unmarshal(rawSecret, secretDst); err != nil {
		return nil, fmt.Errorf("payload is not valid JSON for %s/%s: %w", provider, kind, err)
	}
	if err := def.Validate(secretDst); err != nil {
		return nil, fmt.Errorf("invalid %s/%s: %w", provider, kind, err)
	}

	// 2) scope shape
	prScopes, ok := dto.ScopeRegistry[provider]
	if !ok {
		return nil, fmt.Errorf("no scopes registered for provider %q", provider)
	}
	kScopes, ok := prScopes[scopeKind]
	if !ok {
		return nil, fmt.Errorf("invalid scope_kind %q for provider %q", scopeKind, provider)
	}
	sdef, ok := kScopes[scopeVersion]
	if !ok {
		return nil, fmt.Errorf("unsupported scope version %d for %s/%s", scopeVersion, provider, scopeKind)
	}

	scopeDst := sdef.New()
	if err := json.Unmarshal(rawScope, scopeDst); err != nil {
		return nil, fmt.Errorf("invalid scope JSON: %w", err)
	}
	if err := sdef.Validate(scopeDst); err != nil {
		return nil, fmt.Errorf("invalid scope: %w", err)
	}

	// 3) canonicalize scope (also what we persist in plaintext)
	canonScope, err := canonicalJSON(scopeDst)
	if err != nil {
		return nil, err
	}
	fp := sha256Hex(canonScope) // or HMAC if you have a server-side key

	// 4) canonicalize + encrypt secret
	canonSecret, err := canonicalJSON(secretDst)
	if err != nil {
		return nil, err
	}
	cipher, iv, tag, err := utils.EncryptForOrg(orgID, canonSecret, db)
	if err != nil {
		return nil, fmt.Errorf("encrypt: %w", err)
	}

	cred := &models.Credential{
		OrganizationID:   orgID,
		Provider:         provider,
		Kind:             kind,
		SchemaVersion:    schemaVersion,
		Name:             name,
		ScopeKind:        scopeKind,
		Scope:            datatypes.JSON(canonScope),
		ScopeVersion:     scopeVersion,
		AccountID:        accountID,
		Region:           region,
		ScopeFingerprint: fp,
		EncryptedData:    cipher,
		IV:               iv,
		Tag:              tag,
	}

	if err := db.WithContext(ctx).Create(cred).Error; err != nil {
		return nil, err
	}
	return cred, nil
}

// credOut converts model → response DTO
func credOut(c *models.Credential) dto.CredentialOut {
	return dto.CredentialOut{
		ID:                 c.ID.String(),
		CredentialProvider: c.Provider,
		Kind:               c.Kind,
		SchemaVersion:      c.SchemaVersion,
		Name:               c.Name,
		ScopeKind:          c.ScopeKind,
		ScopeVersion:       c.ScopeVersion,
		Scope:              dto.RawJSON(c.Scope),
		AccountID:          c.AccountID,
		Region:             c.Region,
		CreatedAt:          c.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:          c.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
