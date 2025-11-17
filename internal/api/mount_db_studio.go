package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	pgapi "github.com/sosedoff/pgweb/pkg/api"
	pgclient "github.com/sosedoff/pgweb/pkg/client"
	pgcmd "github.com/sosedoff/pgweb/pkg/command"
)

func MountDbStudio(dbURL, prefix string, readonly bool) (http.Handler, error) {
	// Normalize prefix for pgweb:
	//  - no leading slash
	//  - always trailing slash if not empty
	prefix = strings.Trim(prefix, "/")
	if prefix != "" {
		prefix = prefix + "/"
	}

	pgcmd.Opts = pgcmd.Options{
		URL:         dbURL,
		Prefix:      prefix, // e.g. "db-studio/"
		ReadOnly:    readonly,
		Sessions:    false,
		LockSession: true,
		SkipOpen:    true,
	}

	cli, err := pgclient.NewFromUrl(dbURL, nil)
	if err != nil {
		return nil, err
	}
	if readonly {
		_ = cli.SetReadOnlyMode()
	}

	if err := cli.Test(); err != nil {
		return nil, err
	}

	pgapi.DbClient = cli

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Recovery())

	pgapi.SetupRoutes(g)
	pgapi.SetupMetrics(g)

	return g, nil
}
