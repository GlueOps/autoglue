package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/glueops/autoglue/internal/api/httpmiddleware"
	"github.com/glueops/autoglue/internal/handlers/dto"
	"github.com/glueops/autoglue/internal/models"
	"github.com/glueops/autoglue/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ---------- Helpers ----------

func normLowerNoDot(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	return strings.TrimSuffix(s, ".")
}

func fqdn(domain string, rel string) string {
	d := normLowerNoDot(domain)
	r := normLowerNoDot(rel)
	if r == "" || r == "@" {
		return d
	}
	return r + "." + d
}

func canonicalJSONAny(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var anyv any
	if err := json.Unmarshal(b, &anyv); err != nil {
		return nil, err
	}
	return marshalSortedDNS(anyv)
}

func marshalSortedDNS(v any) ([]byte, error) {
	switch vv := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(vv))
		for k := range vv {
			keys = append(keys, k)
		}
		sortStrings(keys)
		var buf bytes.Buffer
		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			kb, _ := json.Marshal(k)
			buf.Write(kb)
			buf.WriteByte(':')
			b, err := marshalSortedDNS(vv[k])
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}
		buf.WriteByte('}')
		return buf.Bytes(), nil
	case []any:
		var buf bytes.Buffer
		buf.WriteByte('[')
		for i, e := range vv {
			if i > 0 {
				buf.WriteByte(',')
			}
			b, err := marshalSortedDNS(e)
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

func sortStrings(a []string) {
	for i := 0; i < len(a); i++ {
		for j := i + 1; j < len(a); j++ {
			if a[j] < a[i] {
				a[i], a[j] = a[j], a[i]
			}
		}
	}
}

func sha256HexBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

/* Fingerprint (provider-agnostic) */
type desiredRecord struct {
	ZoneID string   `json:"zone_id"`
	FQDN   string   `json:"fqdn"`
	Type   string   `json:"type"`
	TTL    *int     `json:"ttl,omitempty"`
	Values []string `json:"values,omitempty"`
}

func computeFingerprint(zoneID, fqdn, typ string, ttl *int, values datatypes.JSON) (string, error) {
	var vals []string
	if len(values) > 0 && string(values) != "null" {
		if err := json.Unmarshal(values, &vals); err != nil {
			return "", err
		}
		sortStrings(vals)
	}
	payload := &desiredRecord{
		ZoneID: zoneID, FQDN: fqdn, Type: strings.ToUpper(typ), TTL: ttl, Values: vals,
	}
	can, err := canonicalJSONAny(payload)
	if err != nil {
		return "", err
	}
	return sha256HexBytes(can), nil
}

func mustSameOrgDomainWithCredential(db *gorm.DB, orgID uuid.UUID, credID uuid.UUID) error {
	var cred models.Credential
	if err := db.Where("id = ? AND organization_id = ?", credID, orgID).First(&cred).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("credential not found or belongs to different org")
		}
		return err
	}
	if cred.Provider != "aws" || cred.ScopeKind != "service" {
		return fmt.Errorf("credential must be AWS Route 53 service scoped")
	}
	var scope map[string]any
	if err := json.Unmarshal(cred.Scope, &scope); err != nil {
		return fmt.Errorf("credential scope invalid json: %w", err)
	}
	if strings.ToLower(fmt.Sprint(scope["service"])) != "route53" {
		return fmt.Errorf("credential scope.service must be route53")
	}
	return nil
}

// ---------- Domain Handlers ----------

// ListDomains godoc
//
//	@ID				ListDomains
//	@Summary		List domains (org scoped)
//	@Description	Returns domains for X-Org-ID. Filters: `domain_name`, `status`, `q` (contains).
//	@Tags			DNS
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			domain_name	query		string	false	"Exact domain name (lowercase, no trailing dot)"
//	@Param			status		query		string	false	"pending|provisioning|ready|failed"
//	@Param			q			query		string	false	"Domain contains (case-insensitive)"
//	@Success		200			{array}		dto.DomainResponse
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"db error"
//	@Router			/dns/domains [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListDomains(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		q := db.Model(&models.Domain{}).Where("organization_id = ?", orgID)
		if v := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("domain_name"))); v != "" {
			q = q.Where("LOWER(domain_name) = ?", v)
		}
		if v := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("status"))); v != "" {
			q = q.Where("status = ?", v)
		}
		if needle := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("q"))); needle != "" {
			q = q.Where("LOWER(domain_name) LIKE ?", "%"+needle+"%")
		}

		var rows []models.Domain
		if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		out := make([]dto.DomainResponse, 0, len(rows))
		for i := range rows {
			out = append(out, domainOut(&rows[i]))
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetDomain godoc
//
//	@ID			GetDomain
//	@Summary	Get a domain (org scoped)
//	@Tags		DNS
//	@Produce	json
//	@Param		X-Org-ID	header		string	false	"Organization UUID"
//	@Param		id			path		string	true	"Domain ID (UUID)"
//	@Success	200			{object}	dto.DomainResponse
//	@Failure	401			{string}	string	"Unauthorized"
//	@Failure	403			{string}	string	"organization required"
//	@Failure	404			{string}	string	"not found"
//	@Router		/dns/domains/{id} [get]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func GetDomain(db *gorm.DB) http.HandlerFunc {
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
		var row models.Domain
		if err := db.Where("organization_id = ? AND id = ?", orgID, id).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "domain not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		utils.WriteJSON(w, http.StatusOK, domainOut(&row))
	}
}

// CreateDomain godoc
//
//	@ID				CreateDomain
//	@Summary		Create a domain (org scoped)
//	@Description	Creates a domain bound to a Route 53 scoped credential. Archer will backfill ZoneID if omitted.
//	@Tags			DNS
//	@Accept			json
//	@Produce		json
//	@Param			X-Org-ID	header		string					false	"Organization UUID"
//	@Param			body		body		dto.CreateDomainRequest	true	"Domain payload"
//	@Success		201			{object}	dto.DomainResponse
//	@Failure		400			{string}	string	"validation error"
//	@Failure		401			{string}	string	"Unauthorized"
//	@Failure		403			{string}	string	"organization required"
//	@Failure		500			{string}	string	"db error"
//	@Router			/dns/domains [post]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func CreateDomain(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		var in dto.CreateDomainRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}
		if err := dto.DNSValidate(in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		credID, _ := uuid.Parse(in.CredentialID)
		if err := mustSameOrgDomainWithCredential(db, orgID, credID); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "invalid_credential", err.Error())
			return
		}

		row := &models.Domain{
			OrganizationID: orgID,
			DomainName:     normLowerNoDot(in.DomainName),
			ZoneID:         strings.TrimSpace(in.ZoneID),
			Status:         "pending",
			LastError:      "",
			CredentialID:   credID,
		}
		if err := db.Create(row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusCreated, domainOut(row))
	}
}

// UpdateDomain godoc
//
//	@ID			UpdateDomain
//	@Summary	Update a domain (org scoped)
//	@Tags		DNS
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header		string					false	"Organization UUID"
//	@Param		id			path		string					true	"Domain ID (UUID)"
//	@Param		body		body		dto.UpdateDomainRequest	true	"Fields to update"
//	@Success	200			{object}	dto.DomainResponse
//	@Failure	400			{string}	string	"validation error"
//	@Failure	403			{string}	string	"organization required"
//	@Failure	404			{string}	string	"not found"
//	@Router		/dns/domains/{id} [patch]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func UpdateDomain(db *gorm.DB) http.HandlerFunc {
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
		var row models.Domain
		if err := db.Where("organization_id = ? AND id = ?", orgID, id).First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "domain not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		var in dto.UpdateDomainRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}
		if err := dto.DNSValidate(in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		if in.DomainName != nil {
			row.DomainName = normLowerNoDot(*in.DomainName)
		}
		if in.CredentialID != nil {
			credID, _ := uuid.Parse(*in.CredentialID)
			if err := mustSameOrgDomainWithCredential(db, orgID, credID); err != nil {
				utils.WriteError(w, http.StatusBadRequest, "invalid_credential", err.Error())
				return
			}
			row.CredentialID = credID
			row.Status = "pending"
			row.LastError = ""
		}
		if in.ZoneID != nil {
			row.ZoneID = strings.TrimSpace(*in.ZoneID)
		}
		if in.Status != nil {
			row.Status = *in.Status
			if row.Status == "pending" {
				row.LastError = ""
			}
		}
		if err := db.Save(&row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, domainOut(&row))
	}
}

// DeleteDomain godoc
//
//	@ID			DeleteDomain
//	@Summary	Delete a domain
//	@Tags		DNS
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Domain ID (UUID)"
//	@Success	204
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Router		/dns/domains/{id} [delete]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func DeleteDomain(db *gorm.DB) http.HandlerFunc {
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
		res := db.Where("organization_id = ? AND id = ?", orgID, id).Delete(&models.Domain{})
		if res.Error != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", res.Error.Error())
			return
		}
		if res.RowsAffected == 0 {
			utils.WriteError(w, http.StatusNotFound, "not_found", "domain not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// ---------- Record Set Handlers ----------

// ListRecordSets godoc
//
//	@ID				ListRecordSets
//	@Summary		List record sets for a domain
//	@Description	Filters: `name`, `type`, `status`.
//	@Tags			DNS
//	@Produce		json
//	@Param			X-Org-ID	header		string	false	"Organization UUID"
//	@Param			domain_id	path		string	true	"Domain ID (UUID)"
//	@Param			name		query		string	false	"Exact relative name or FQDN (server normalizes)"
//	@Param			type		query		string	false	"RR type (A, AAAA, CNAME, TXT, MX, NS, SRV, CAA)"
//	@Param			status		query		string	false	"pending|provisioning|ready|failed"
//	@Success		200			{array}		dto.RecordSetResponse
//	@Failure		403			{string}	string	"organization required"
//	@Failure		404			{string}	string	"domain not found"
//	@Router			/dns/domains/{domain_id}/records [get]
//	@Security		BearerAuth
//	@Security		OrgKeyAuth
//	@Security		OrgSecretAuth
func ListRecordSets(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		did, err := uuid.Parse(chi.URLParam(r, "domain_id"))
		if err != nil {
			log.Info().Msg(err.Error())
			utils.WriteError(w, http.StatusBadRequest, "bad_id", "invalid domain UUID:")
			return
		}
		var domain models.Domain
		if err := db.Where("organization_id = ? AND id = ?", orgID, did).First(&domain).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "domain not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}

		q := db.Model(&models.RecordSet{}).Where("domain_id = ?", did)
		if v := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("name"))); v != "" {
			dn := strings.ToLower(domain.DomainName)
			rel := v
			// normalize apex or FQDN into relative
			if v == dn || v == dn+"." {
				rel = ""
			} else {
				rel = strings.TrimSuffix(v, "."+dn)
				rel = normLowerNoDot(rel)
			}
			q = q.Where("LOWER(name) = ?", rel)
		}
		if v := strings.TrimSpace(strings.ToUpper(r.URL.Query().Get("type"))); v != "" {
			q = q.Where("type = ?", v)
		}
		if v := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("status"))); v != "" {
			q = q.Where("status = ?", v)
		}

		var rows []models.RecordSet
		if err := q.Order("created_at DESC").Find(&rows).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", "db error")
			return
		}
		out := make([]dto.RecordSetResponse, 0, len(rows))
		for i := range rows {
			out = append(out, recordOut(&rows[i]))
		}
		utils.WriteJSON(w, http.StatusOK, out)
	}
}

// GetRecordSet godoc
//
//	@ID			GetRecordSet
//	@Summary	Get a record set (org scoped)
//	@Tags		DNS
//	@Produce	json
//	@Param		X-Org-ID	header		string	false	"Organization UUID"
//	@Param		id			path		string	true	"Record Set ID (UUID)"
//	@Success	200			{object}	dto.RecordSetResponse
//	@Failure	403			{string}	string	"organization required"
//	@Failure	404			{string}	string	"not found"
//	@Router		/dns/records/{id} [get]
func GetRecordSet(db *gorm.DB) http.HandlerFunc {
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

		var row models.RecordSet
		if err := db.
			Joins("Domain").
			Where(`record_sets.id = ? AND "Domain"."organization_id" = ?`, id, orgID).
			First(&row).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "record set not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		utils.WriteJSON(w, http.StatusOK, recordOut(&row))
	}
}

// CreateRecordSet godoc
//
//	@ID			CreateRecordSet
//	@Summary	Create a record set (pending; Archer will UPSERT to Route 53)
//	@Tags		DNS
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header		string						false	"Organization UUID"
//	@Param		domain_id	path		string						true	"Domain ID (UUID)"
//	@Param		body		body		dto.CreateRecordSetRequest	true	"Record set payload"
//	@Success	201			{object}	dto.RecordSetResponse
//	@Failure	400			{string}	string	"validation error"
//	@Failure	403			{string}	string	"organization required"
//	@Failure	404			{string}	string	"domain not found"
//	@Router		/dns/domains/{domain_id}/records [post]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func CreateRecordSet(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID, ok := httpmiddleware.OrgIDFrom(r.Context())
		if !ok {
			utils.WriteError(w, http.StatusForbidden, "org_required", "specify X-Org-ID")
			return
		}
		did, err := uuid.Parse(chi.URLParam(r, "domain_id"))
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_id", "invalid domain UUID")
			return
		}
		var domain models.Domain
		if err := db.Where("organization_id = ? AND id = ?", orgID, did).First(&domain).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "domain not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		var in dto.CreateRecordSetRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}
		if err := dto.DNSValidate(in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		t := strings.ToUpper(in.Type)
		if t == "CNAME" && len(in.Values) != 1 {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", "CNAME requires exactly one value")
			return
		}

		rel := normLowerNoDot(in.Name)
		fq := fqdn(domain.DomainName, rel)

		// Pre-flight: block duplicate tuple and protect from non-autoglue rows
		var existing models.RecordSet
		if err := db.Where("domain_id = ? AND LOWER(name) = ? AND type = ?",
			domain.ID, strings.ToLower(rel), t).First(&existing).Error; err == nil {
			if existing.Owner != "" && existing.Owner != "autoglue" {
				utils.WriteError(w, http.StatusConflict, "ownership_conflict",
					"record with the same (name,type) exists but is not owned by autoglue")
				return
			}
			utils.WriteError(w, http.StatusConflict, "already_exists",
				"a record with the same (name,type) already exists; use PATCH to modify")
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		valuesJSON, _ := json.Marshal(in.Values)
		fp, err := computeFingerprint(domain.ZoneID, fq, t, in.TTL, datatypes.JSON(valuesJSON))
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "fingerprint_error", err.Error())
			return
		}

		row := &models.RecordSet{
			DomainID:    domain.ID,
			Name:        rel,
			Type:        t,
			TTL:         in.TTL,
			Values:      datatypes.JSON(valuesJSON),
			Fingerprint: fp,
			Status:      "pending",
			LastError:   "",
			Owner:       "autoglue",
		}
		if err := db.Create(row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusCreated, recordOut(row))
	}
}

// UpdateRecordSet godoc
//
//	@ID			UpdateRecordSet
//	@Summary	Update a record set (flips to pending for reconciliation)
//	@Tags		DNS
//	@Accept		json
//	@Produce	json
//	@Param		X-Org-ID	header		string						false	"Organization UUID"
//	@Param		id			path		string						true	"Record Set ID (UUID)"
//	@Param		body		body		dto.UpdateRecordSetRequest	true	"Fields to update"
//	@Success	200			{object}	dto.RecordSetResponse
//	@Failure	400			{string}	string	"validation error"
//	@Failure	403			{string}	string	"organization required"
//	@Failure	404			{string}	string	"not found"
//	@Router		/dns/records/{id} [patch]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func UpdateRecordSet(db *gorm.DB) http.HandlerFunc {
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

		var row models.RecordSet
		if err := db.
			Joins("Domain").
			Where(`record_sets.id = ? AND "Domain"."organization_id" = ?`, id, orgID).
			First(&row).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.WriteError(w, http.StatusNotFound, "not_found", "record set not found")
				return
			}
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		var domain models.Domain
		if err := db.Where("id = ?", row.DomainID).First(&domain).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}

		var in dto.UpdateRecordSetRequest
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "bad_json", err.Error())
			return
		}
		if err := dto.DNSValidate(in); err != nil {
			utils.WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		if row.Owner != "" && row.Owner != "autoglue" {
			utils.WriteError(w, http.StatusConflict, "ownership_conflict",
				"record is not owned by autoglue; refuse to modify")
			return
		}

		// Mutations
		if in.Name != nil {
			row.Name = normLowerNoDot(*in.Name)
		}
		if in.Type != nil {
			row.Type = strings.ToUpper(*in.Type)
		}
		if in.TTL != nil {
			row.TTL = in.TTL
		}
		if in.Values != nil {
			t := row.Type
			if in.Type != nil {
				t = strings.ToUpper(*in.Type)
			}
			if t == "CNAME" && len(*in.Values) != 1 {
				utils.WriteError(w, http.StatusBadRequest, "validation_error", "CNAME requires exactly one value")
				return
			}
			b, _ := json.Marshal(*in.Values)
			row.Values = datatypes.JSON(b)
		}

		if in.Status != nil {
			row.Status = *in.Status
		} else {
			row.Status = "pending"
			row.LastError = ""
		}

		fq := fqdn(domain.DomainName, row.Name)
		fp, err := computeFingerprint(domain.ZoneID, fq, row.Type, row.TTL, row.Values)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "fingerprint_error", err.Error())
			return
		}
		row.Fingerprint = fp

		if err := db.Save(&row).Error; err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", err.Error())
			return
		}
		utils.WriteJSON(w, http.StatusOK, recordOut(&row))
	}
}

// DeleteRecordSet godoc
//
//	@ID			DeleteRecordSet
//	@Summary	Delete a record set (API removes row; worker can optionally handle external deletion policy)
//	@Tags		DNS
//	@Produce	json
//	@Param		X-Org-ID	header	string	false	"Organization UUID"
//	@Param		id			path	string	true	"Record Set ID (UUID)"
//	@Success	204
//	@Failure	403	{string}	string	"organization required"
//	@Failure	404	{string}	string	"not found"
//	@Router		/dns/records/{id} [delete]
//	@Security	BearerAuth
//	@Security	OrgKeyAuth
//	@Security	OrgSecretAuth
func DeleteRecordSet(db *gorm.DB) http.HandlerFunc {
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
		sub := db.Model(&models.RecordSet{}).
			Select("record_sets.id").
			Joins("JOIN domains ON domains.id = record_sets.domain_id").
			Where("record_sets.id = ? AND domains.organization_id = ?", id, orgID)

		res := db.Where("id IN (?)", sub).Delete(&models.RecordSet{})
		if res.Error != nil {
			utils.WriteError(w, http.StatusInternalServerError, "db_error", res.Error.Error())
			return
		}
		if res.RowsAffected == 0 {
			utils.WriteError(w, http.StatusNotFound, "not_found", "record set not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// ---------- Out mappers ----------

func domainOut(m *models.Domain) dto.DomainResponse {
	return dto.DomainResponse{
		ID:             m.ID.String(),
		OrganizationID: m.OrganizationID.String(),
		DomainName:     m.DomainName,
		ZoneID:         m.ZoneID,
		Status:         m.Status,
		LastError:      m.LastError,
		CredentialID:   m.CredentialID.String(),
		CreatedAt:      m.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:      m.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func recordOut(r *models.RecordSet) dto.RecordSetResponse {
	vals := r.Values
	if len(vals) == 0 {
		vals = datatypes.JSON("[]")
	}
	return dto.RecordSetResponse{
		ID:          r.ID.String(),
		DomainID:    r.DomainID.String(),
		Name:        r.Name,
		Type:        r.Type,
		TTL:         r.TTL,
		Values:      []byte(vals),
		Fingerprint: r.Fingerprint,
		Status:      r.Status,
		LastError:   r.LastError,
		Owner:       r.Owner,
		CreatedAt:   r.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   r.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
