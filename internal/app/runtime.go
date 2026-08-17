package app

import (
	"context"
	"log"

	"github.com/glueops/autoglue/internal/config"
	"github.com/glueops/autoglue/internal/db"
	"github.com/glueops/autoglue/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	"gorm.io/gorm"
)

type Runtime struct {
	Cfg  config.Config
	DB   *gorm.DB
	Pool *pgxpool.Pool
}

func NewRuntime() *Runtime {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	d := db.Open(cfg.DbURL)

	err = db.Run(d,
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
	)

	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	ctx := context.Background()

	pool, err := db.OpenPool(ctx, cfg.DbURL)
	if err != nil {
		log.Fatalf("Error opening pgx pool: %v", err)
	}

	if err := migrateRiver(ctx, pool); err != nil {
		log.Fatalf("Error migrating river schema: %v", err)
	}

	dropLegacyJobsTable(d)

	return &Runtime{
		Cfg:  cfg,
		DB:   d,
		Pool: pool,
	}
}

// migrateRiver brings the river_* tables up to date. It is idempotent and a
// no-op once the schema matches the linked River version.
func migrateRiver(ctx context.Context, pool *pgxpool.Pool) error {
	migrator, err := rivermigrate.New(riverpgxv5.New(pool), nil)
	if err != nil {
		return err
	}
	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	return err
}

// dropLegacyJobsTable removes the archer-era `jobs` table. River owns
// river_job now, and nothing reads the old table: serve used to delete every
// scheduled/queued row on boot anyway, so there is no history worth keeping.
//
// This is deliberately best-effort and non-fatal. Once every environment has
// booted at least once on River, this call and the whole function can go.
func dropLegacyJobsTable(d *gorm.DB) {
	if !d.Migrator().HasTable("jobs") {
		return
	}
	if err := d.Exec(`DROP TABLE IF EXISTS jobs`).Error; err != nil {
		log.Printf("warning: could not drop legacy jobs table: %v", err)
		return
	}
	log.Printf("dropped legacy archer jobs table")
}

// Close releases the pgx pool. The GORM handle is left alone: it is process
// scoped and torn down on exit.
func (r *Runtime) Close() {
	if r.Pool != nil {
		r.Pool.Close()
	}
}
