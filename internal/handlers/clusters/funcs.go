package clusters

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ensureNodePoolsBelongToOrg(orgID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return errors.New("empty ids")
	}
	var count int64
	if err := db.DB.Model(&models.NodePool{}).
		Where("organization_id = ? AND id IN ?", orgID, ids).
		Count(&count).Error; err != nil {
		return err
	}
	if count != int64(len(ids)) {
		return errors.New("some node pools do not belong to org")
	}
	return nil
}

func ensureServerBelongsToOrgWithRole(orgID uuid.UUID, id uuid.UUID, role string) error {
	var count int64
	if err := db.DB.Model(&models.Server{}).
		Where("organization_id = ? AND id = ? AND role = ?", orgID, id, role).
		Count(&count).Error; err != nil {
		return err
	}
	if count != 1 {
		return errors.New("server not found in org or role mismatch")
	}
	return nil
}

func toResp(c models.Cluster, includePools, includeBastion bool) clusterResponse {
	out := clusterResponse{
		ID:                  c.ID,
		Name:                c.Name,
		Provider:            c.Provider,
		Region:              c.Region,
		Status:              c.Status,
		ClusterLoadBalancer: c.ClusterLoadBalancer,
		ControlLoadBalancer: c.ControlLoadBalancer,
	}
	if includePools {
		out.NodePools = make([]nodePoolBrief, 0, len(c.NodePools))
		for _, p := range c.NodePools {
			nb := nodePoolBrief{ID: p.ID, Name: p.Name}
			fmt.Printf("node pool %s\n", p.Name)
			fmt.Printf("node pool labels %d\n", len(p.Labels))
			if len(p.Labels) > 0 {
				nb.Labels = make([]labelBrief, 0, len(p.Labels))
				for _, l := range p.Labels {
					nb.Labels = append(nb.Labels, labelBrief{ID: l.ID, Key: l.Key, Value: l.Value})
				}
			}

			fmt.Printf("node pool annotations %d\n", len(p.Annotations))

			if len(p.Annotations) > 0 {
				nb.Annotations = make([]annotationBrief, 0, len(p.Annotations))
				for _, a := range p.Annotations {
					nb.Annotations = append(nb.Annotations, annotationBrief{ID: a.ID, Key: a.Key, Value: a.Value})
				}
			}

			fmt.Printf("node pool taints %d\n", len(p.Taints))

			if len(p.Taints) > 0 {
				nb.Taints = make([]taintBrief, 0, len(p.Taints))
				for _, t := range p.Taints {
					nb.Taints = append(nb.Taints, taintBrief{ID: t.ID, Key: t.Key, Value: t.Value, Effect: t.Effect})
				}
			}

			if len(p.Servers) > 0 {
				nb.Servers = make([]serverBrief, 0, len(p.Servers))
				for _, s := range p.Servers {
					nb.Servers = append(nb.Servers, serverBrief{ID: s.ID, Hostname: s.Hostname, Role: s.Role, Status: s.Status, IP: s.IPAddress})
				}
			}

			out.NodePools = append(out.NodePools, nb)
		}
	}
	if includeBastion && c.BastionServer != nil {
		out.BastionServer = &serverBrief{
			ID:       c.BastionServer.ID,
			Hostname: c.BastionServer.Hostname,
			IP:       c.BastionServer.IPAddress,
			Role:     c.BastionServer.Role,
			Status:   c.BastionServer.Status,
		}
	}
	return out
}

func contains(xs []string, want string) bool {
	for _, x := range xs {
		if strings.TrimSpace(x) == want {
			return true
		}
	}
	return false
}

func errorsIsNotFound(err error) bool { return err == gorm.ErrRecordNotFound }

func parseUUIDs(ids []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(ids))
	for _, s := range ids {
		u, err := uuid.Parse(strings.TrimSpace(s))
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

const maxJSONBytes int64 = 1 << 20

func readJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if ct := r.Header.Get("Content-Type"); ct != "" {
		mt, _, err := mime.ParseMediaType(ct)
		if err != nil || mt != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return false
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBytes)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		var maxErr *http.MaxBytesError

		switch {
		case errors.As(err, &maxErr):
			http.Error(w, fmt.Sprintf("request body too large (max %d bytes)", maxJSONBytes), http.StatusRequestEntityTooLarge)
		case errors.Is(err, io.EOF):
			http.Error(w, "request body must not be empty", http.StatusBadRequest)
		case errors.As(err, &syntaxErr):
			http.Error(w, fmt.Sprintf("malformed JSON at character %d", syntaxErr.Offset), http.StatusBadRequest)
		case errors.Is(err, io.ErrUnexpectedEOF):
			http.Error(w, "malformed JSON", http.StatusBadRequest)
		case errors.As(err, &typeErr):
			// Example: expected string but got number for field "name"
			field := typeErr.Field
			if field == "" && len(typeErr.Struct) > 0 {
				field = typeErr.Struct
			}
			http.Error(w, fmt.Sprintf("invalid value for %q (expected %s)", field, typeErr.Type.String()), http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			// Extract the field name from the error message.
			field := strings.Trim(strings.TrimPrefix(err.Error(), "json: unknown field "), "\"")
			http.Error(w, fmt.Sprintf("unknown field %q", field), http.StatusBadRequest)
		default:
			http.Error(w, "invalid json", http.StatusBadRequest)
		}
		return false
	}

	if dec.More() {
		// Try to read one more token/value; if not EOF, there was extra content.
		var extra any
		if err := dec.Decode(&extra); err != io.EOF {
			http.Error(w, "body must contain only a single JSON object", http.StatusBadRequest)
			return false
		}
	}

	return true
}
