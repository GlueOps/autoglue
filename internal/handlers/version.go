package handlers

import (
	"net/http"
	"runtime"
	"runtime/debug"
	"strconv"

	"github.com/glueops/autoglue/internal/utils"
	"github.com/glueops/autoglue/internal/version"
)

type VersionResponse struct {
	Version    string `json:"version"  example:"1.4.2"`
	Commit     string `json:"commit"   example:"a1b2c3d"`
	Built      string `json:"built"    example:"2025-11-08T12:34:56Z"`
	BuiltBy    string `json:"builtBy"  example:"ci"`
	Go         string `json:"go"     example:"go1.23.3"`
	GOOS       string `json:"goOS"   example:"linux"`
	GOARCH     string `json:"goArch" example:"amd64"`
	VCS        string `json:"vcs,omitempty"        example:"git"`
	Revision   string `json:"revision,omitempty"   example:"a1b2c3d4e5f6abcdef"`
	CommitTime string `json:"commitTime,omitempty" example:"2025-11-08T12:31:00Z"`
	Modified   *bool  `json:"modified,omitempty"   example:"false"`
}

// Version godoc
//
//	@Summary		Service version information
//	@Description	Returns build/runtime metadata for the running service.
//	@Tags			Meta
//	@ID				Version                 // operationId
//	@Produce		json
//	@Success		200	{object}	VersionResponse
//	@Router			/version [get]
func Version(w http.ResponseWriter, r *http.Request) {
	resp := VersionResponse{
		Version: version.Version,
		Commit:  version.Commit,
		Built:   version.Date,
		BuiltBy: version.BuiltBy,
		Go:      runtime.Version(),
		GOOS:    runtime.GOOS,
		GOARCH:  runtime.GOARCH,
	}

	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs":
				resp.VCS = s.Value
			case "vcs.revision":
				resp.Revision = s.Value
			case "vcs.time":
				resp.CommitTime = s.Value
			case "vcs.modified":
				if b, err := strconv.ParseBool(s.Value); err == nil {
					resp.Modified = &b
				}
			}
		}
	}
	utils.WriteJSON(w, http.StatusOK, resp)
}
