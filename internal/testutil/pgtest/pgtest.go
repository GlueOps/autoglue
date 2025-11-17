package pgtest

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	const port uint32 = 55432

	cfg := embeddedpostgres.
		DefaultConfig().
		Database("autoglue_test").
		Username("autoglue").
		Password("autoglue").
		Port(port).
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

	// Use the same model list as app.NewRuntime so schema matches prod
	if err := db.Run(
		dbConn,
		&models.Job{},
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
		&models.Cluster{},
		&models.Credential{},
		&models.Domain{},
		&models.RecordSet{},
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
