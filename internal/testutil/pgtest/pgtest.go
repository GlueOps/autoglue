package pgtest

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// EnvPort lets a caller pin the embedded Postgres port. When unset, a free
// port is chosen at runtime: a fixed port collides with anything else already
// listening (another checkout, a local Postgres, a stray container) and turns
// an unrelated process into an unexplained test failure.
const EnvPort = "AUTOGLUE_TEST_PG_PORT"

// listenPort returns the port embedded Postgres should bind.
func listenPort() uint32 {
	if v := os.Getenv(EnvPort); v != "" {
		p, err := strconv.ParseUint(v, 10, 16)
		if err != nil || p == 0 {
			log.Printf("pgtest: ignoring invalid %s=%q, selecting a free port", EnvPort, v)
		} else {
			return uint32(p)
		}
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		// Nothing better to fall back to; Start will report the real failure.
		log.Printf("pgtest: could not reserve a free port (%v), falling back to 55432", err)
		return 55432
	}
	port := uint32(l.Addr().(*net.TCPAddr).Port)
	_ = l.Close()
	return port
}

// runtimePath isolates the extracted Postgres runtime per test binary.
//
// `go test ./...` runs packages in parallel, and embedded-postgres does
// os.RemoveAll on the runtime path at the start of every Start. Two packages
// sharing the default (~/.embedded-postgres-go/extracted) therefore delete the
// directory out from under each other mid-extraction, and the loser fails with
// a rename into a directory that no longer exists. The library's own error text
// names this fix. It only bites on a cold cache, so it looks fine locally and
// fails in CI.
//
// Only the extracted runtime needs isolating: the downloaded archive is written
// to a temp file and renamed into the cache, which the library does explicitly
// to survive concurrent extraction. So this costs one extra unpack, not an
// extra download.
//
// Keyed on the test binary name (bg.test, handlers.test) rather than the PID:
// unique per package, but stable across runs, so the unpacked copy is still
// reused instead of being rebuilt every time.
func runtimePath() string {
	return filepath.Join(os.TempDir(), "autoglue-pgtest", filepath.Base(os.Args[0]))
}

var (
	once    sync.Once
	epg     *embeddedpostgres.EmbeddedPostgres
	gdb     *gorm.DB
	initErr error
	dsn     string
)

// initDB is called once via sync.Once. It starts embedded Postgres,
// opens a GORM connection and runs the same migrations as NewRuntime.
func initDB() {
	port := listenPort()

	cfg := embeddedpostgres.
		DefaultConfig().
		Database("autoglue_test").
		Username("autoglue").
		Password("autoglue").
		Port(port).
		RuntimePath(runtimePath()).
		StartTimeout(30 * time.Second)

	epg = embeddedpostgres.NewDatabase(cfg)
	if err := epg.Start(); err != nil {
		initErr = fmt.Errorf("start embedded postgres: %w", err)
		return
	}

	dsn = fmt.Sprintf(
		"host=127.0.0.1 port=%d user=%s password=%s dbname=%s sslmode=disable",
		port,
		"autoglue",
		"autoglue",
		"autoglue_test",
	)

	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		initErr = fmt.Errorf("open gorm: %w", err)
		return
	}

	// Keep this list identical to app.NewRuntime, in the same order, so tests
	// migrate the same schema production does. Anything missing here surfaces
	// as an unrelated-looking failure in whichever test first Preloads it.
	if err := db.Run(
		dbConn,
		&models.MasterKey{},
		&models.SigningKey{},
		&models.User{},
		&models.Organization{},
		&models.Account{},
		&models.Membership{},
		&models.APIKey{},
		&models.UserEmail{},
		&models.RefreshToken{},
		&models.OrganizationKey{},
		&models.SshKey{},
		&models.Server{},
		&models.Taint{},
		&models.Label{},
		&models.Annotation{},
		&models.NodePool{},
		&models.Credential{},
		&models.Domain{},
		&models.RecordSet{},
		&models.LoadBalancer{},
		&models.Cluster{},
		&models.Action{},
		&models.ClusterRun{},
		&models.ClusterMetadata{},
		&models.JobLog{},
		&models.ClusterDesiredState{},
		&models.DesiredResource{},
		&models.Agent{},
		&models.AgentEnrollmentTicket{},
		&models.AgentTask{},
		&models.AgentReconcileStatus{},
	); err != nil {
		initErr = fmt.Errorf("migrate: %w", err)
		return
	}

	gdb = dbConn
}

// DB returns a lazily-initialized *gorm.DB backed by embedded Postgres.
//
// Call this from any test that needs a real DB. If init fails, the test
// will fail immediately with a clear message.
func DB(t *testing.T) *gorm.DB {
	t.Helper()
	once.Do(initDB)
	if initErr != nil {
		t.Fatalf("failed to init embedded postgres: %v", initErr)
	}
	return gdb
}

// URL returns the DSN for the embedded Postgres instance, useful for code
// that expects a DB URL (e.g. bg.NewJobs).
func URL(t *testing.T) string {
	t.Helper()
	DB(t) // ensure initialized
	return dsn
}

// Stop stops the embedded Postgres process. Call from TestMain in at
// least one package, or let the OS clean it up on process exit.
func Stop() {
	if epg != nil {
		if err := epg.Stop(); err != nil {
			log.Printf("stop embedded postgres: %v", err)
		}
	}
}
