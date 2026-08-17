package store

import (
	"context"
	"fmt"
)

// DesiredResource is one entry of a desired-state snapshot.
type DesiredResource struct {
	ResourceType string
	ResourceID   string
	ResourceName string
	Phase        int
	DependsOn    string // JSON array of resource ids
	Spec         string // JSON object
	SpecHash     string
	Required     bool
}

// ReplaceDesiredSnapshot stores a whole generation atomically and moves the
// current-generation pointer to it.
//
// Whole snapshot, not a diff: a diff protocol requires the agent to have every
// prior generation to reconstruct the present, and an agent that was switched
// off for a week has none of them. A snapshot is correct from any starting
// state, including a wiped disk.
//
// Older generations are left in place; PruneGenerationsBefore removes them once
// the new one is applied, so a failed apply can still read what it came from.
func (s *Store) ReplaceDesiredSnapshot(ctx context.Context, generation int64, resources []DesiredResource) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin snapshot tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	// Re-delivery of a generation must not double-insert. Clearing first makes
	// the whole operation idempotent, which matters because the sync fetch is a
	// safe method and may be retried by anything.
	if _, err := tx.ExecContext(ctx,
		`DELETE FROM desired_resources WHERE generation = ?`, generation); err != nil {
		return fmt.Errorf("clear generation %d: %w", generation, err)
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO desired_resources
			(generation, resource_type, resource_id, resource_name, phase,
			 depends_on_json, spec_json, spec_hash, required)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, r := range resources {
		dependsOn := r.DependsOn
		if dependsOn == "" {
			dependsOn = "[]"
		}
		spec := r.Spec
		if spec == "" {
			spec = "{}"
		}
		if _, err := stmt.ExecContext(ctx,
			generation, r.ResourceType, r.ResourceID, r.ResourceName, r.Phase,
			dependsOn, spec, r.SpecHash, boolToInt(r.Required),
		); err != nil {
			return fmt.Errorf("insert %s/%s: %w", r.ResourceType, r.ResourceID, err)
		}
	}

	// The pointer moves last and inside the same transaction, so a crash
	// mid-write leaves the agent on the previous generation rather than
	// pointing at a snapshot that is only half stored.
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO agent_meta(key, value) VALUES(?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		metaCurrentGeneration, fmt.Sprint(generation)); err != nil {
		return fmt.Errorf("set current generation: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit snapshot: %w", err)
	}
	return nil
}

// ListDesiredResources returns a generation's resources in dependency order:
// phase first, then a stable tiebreak so two runs plan the same sequence.
func (s *Store) ListDesiredResources(ctx context.Context, generation int64) ([]DesiredResource, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT resource_type, resource_id, resource_name, phase,
		       depends_on_json, spec_json, spec_hash, required
		FROM desired_resources
		WHERE generation = ?
		ORDER BY phase, resource_type, resource_id`, generation)
	if err != nil {
		return nil, fmt.Errorf("list generation %d: %w", generation, err)
	}
	defer func() { _ = rows.Close() }()

	var out []DesiredResource
	for rows.Next() {
		var r DesiredResource
		var required int
		if err := rows.Scan(&r.ResourceType, &r.ResourceID, &r.ResourceName, &r.Phase,
			&r.DependsOn, &r.Spec, &r.SpecHash, &required); err != nil {
			return nil, fmt.Errorf("scan desired resource: %w", err)
		}
		r.Required = required != 0
		out = append(out, r)
	}
	return out, rows.Err()
}

// PruneGenerationsBefore drops superseded snapshots. Call it only after a
// generation has been applied, so a failed reconcile can still read the state
// it was moving away from.
func (s *Store) PruneGenerationsBefore(ctx context.Context, generation int64) error {
	if _, err := s.db.ExecContext(ctx,
		`DELETE FROM desired_resources WHERE generation < ?`, generation); err != nil {
		return fmt.Errorf("prune generations before %d: %w", generation, err)
	}
	return nil
}

// ReconcileStatus is the agent's own record of one resource's outcome.
type ReconcileStatus struct {
	ResourceType      string
	ResourceID        string
	DesiredGeneration int64
	State             string
	LastError         string
}

// UpsertReconcileStatus records the outcome of reconciling one resource.
//
// An unrecognised resource type belongs here as ReconcileFailed with a clear
// LastError, never as a silent skip: agents are deployed software and will run
// behind the API, so an unknown type is an expected event. Skipping it would
// let the agent report convergence it has not achieved.
func (s *Store) UpsertReconcileStatus(ctx context.Context, st ReconcileStatus) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO reconcile_status
			(resource_type, resource_id, desired_generation, state, last_error, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(resource_type, resource_id) DO UPDATE SET
			desired_generation = excluded.desired_generation,
			state              = excluded.state,
			last_error         = excluded.last_error,
			updated_at         = excluded.updated_at`,
		st.ResourceType, st.ResourceID, st.DesiredGeneration, st.State, st.LastError, nowRFC3339())
	if err != nil {
		return fmt.Errorf("upsert reconcile status %s/%s: %w", st.ResourceType, st.ResourceID, err)
	}
	return nil
}

// ListReconcileStatus returns every recorded resource outcome, for the report
// the agent posts back to the control plane.
func (s *Store) ListReconcileStatus(ctx context.Context) ([]ReconcileStatus, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT resource_type, resource_id, desired_generation, state, last_error
		FROM reconcile_status
		ORDER BY resource_type, resource_id`)
	if err != nil {
		return nil, fmt.Errorf("list reconcile status: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var out []ReconcileStatus
	for rows.Next() {
		var st ReconcileStatus
		if err := rows.Scan(&st.ResourceType, &st.ResourceID, &st.DesiredGeneration,
			&st.State, &st.LastError); err != nil {
			return nil, fmt.Errorf("scan reconcile status: %w", err)
		}
		out = append(out, st)
	}
	return out, rows.Err()
}

// GenerationConverged reports whether every resource in the generation has been
// applied. This is what may advance AppliedGeneration — never the mere fact
// that a reconcile pass finished, since a pass that failed every resource also
// finishes.
func (s *Store) GenerationConverged(ctx context.Context, generation int64) (bool, error) {
	var unconverged int
	err := s.db.QueryRowContext(ctx, `
		SELECT count(*)
		FROM desired_resources d
		LEFT JOIN reconcile_status r
		  ON r.resource_type = d.resource_type AND r.resource_id = d.resource_id
		WHERE d.generation = ?
		  AND d.required = 1
		  AND (r.state IS NULL OR r.state != ? OR r.desired_generation < ?)`,
		generation, ReconcileApplied, generation).Scan(&unconverged)
	if err != nil {
		return false, fmt.Errorf("check convergence of %d: %w", generation, err)
	}
	return unconverged == 0, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
