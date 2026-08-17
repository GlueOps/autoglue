package store

import "testing"

func snapshot(t *testing.T, s *Store, gen int64, res ...DesiredResource) {
	t.Helper()
	if err := s.ReplaceDesiredSnapshot(t.Context(), gen, res); err != nil {
		t.Fatalf("snapshot gen %d: %v", gen, err)
	}
}

func res(id string, phase int, required bool) DesiredResource {
	return DesiredResource{
		ResourceType: "cluster_shape",
		ResourceID:   id,
		Phase:        phase,
		Spec:         `{"k":"v"}`,
		SpecHash:     "h-" + id,
		Required:     required,
	}
}

func TestSnapshotMovesCurrentGenerationAndIsIdempotent(t *testing.T) {
	s, ctx := open(t)

	if gen, _ := s.CurrentGeneration(ctx); gen != 0 {
		t.Fatalf("fresh store current generation = %d, want 0", gen)
	}

	snapshot(t, s, 7, res("a", 0, true), res("b", 1, true))
	if gen, _ := s.CurrentGeneration(ctx); gen != 7 {
		t.Fatalf("current generation = %d, want 7", gen)
	}

	// The sync fetch is a safe method, so the same generation may be delivered
	// more than once. That must not duplicate rows.
	snapshot(t, s, 7, res("a", 0, true), res("b", 1, true))
	got, err := s.ListDesiredResources(ctx, 7)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("redelivery produced %d rows, want 2", len(got))
	}
}

// Ordering is declarative: phase decides sequence, not arrival order.
func TestDesiredResourcesComeBackInPhaseOrder(t *testing.T) {
	s, ctx := open(t)
	snapshot(t, s, 1, res("z", 0, true), res("a", 2, true), res("m", 1, true))

	got, err := s.ListDesiredResources(ctx, 1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	want := []string{"z", "m", "a"}
	for i, w := range want {
		if got[i].ResourceID != w {
			t.Fatalf("order = %v..., want %v", got[i].ResourceID, want)
		}
	}
}

// Convergence must mean "every required resource applied *at this generation*".
// A pass that merely finished is not convergence: a pass that failed everything
// also finishes.
func TestGenerationConvergence(t *testing.T) {
	s, ctx := open(t)
	snapshot(t, s, 5, res("a", 0, true), res("b", 0, true))

	if ok, _ := s.GenerationConverged(ctx, 5); ok {
		t.Fatal("converged with nothing applied")
	}

	apply := func(id string, gen int64, state string) {
		t.Helper()
		if err := s.UpsertReconcileStatus(ctx, ReconcileStatus{
			ResourceType: "cluster_shape", ResourceID: id,
			DesiredGeneration: gen, State: state,
		}); err != nil {
			t.Fatalf("upsert %s: %v", id, err)
		}
	}

	apply("a", 5, ReconcileApplied)
	if ok, _ := s.GenerationConverged(ctx, 5); ok {
		t.Fatal("converged with one of two applied")
	}

	apply("b", 5, ReconcileFailed)
	if ok, _ := s.GenerationConverged(ctx, 5); ok {
		t.Fatal("a failed resource counted as converged")
	}

	apply("b", 5, ReconcileApplied)
	if ok, _ := s.GenerationConverged(ctx, 5); !ok {
		t.Fatal("all applied but not converged")
	}
}

// The subtle one: a resource applied at an older generation is stale, not done.
// Without the generation comparison an agent would report convergence to a
// snapshot it never actually applied.
func TestStaleStatusDoesNotCountAsConverged(t *testing.T) {
	s, ctx := open(t)
	snapshot(t, s, 1, res("a", 0, true))
	if err := s.UpsertReconcileStatus(ctx, ReconcileStatus{
		ResourceType: "cluster_shape", ResourceID: "a",
		DesiredGeneration: 1, State: ReconcileApplied,
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if ok, _ := s.GenerationConverged(ctx, 1); !ok {
		t.Fatal("generation 1 should be converged")
	}

	// New snapshot, same resource, nothing re-applied yet.
	snapshot(t, s, 2, res("a", 0, true))
	if ok, _ := s.GenerationConverged(ctx, 2); ok {
		t.Fatal("status from generation 1 counted as convergence to generation 2")
	}
}

func TestOptionalResourcesDoNotBlockConvergence(t *testing.T) {
	s, ctx := open(t)
	snapshot(t, s, 3, res("required", 0, true), res("optional", 0, false))

	if err := s.UpsertReconcileStatus(ctx, ReconcileStatus{
		ResourceType: "cluster_shape", ResourceID: "required",
		DesiredGeneration: 3, State: ReconcileApplied,
	}); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if ok, _ := s.GenerationConverged(ctx, 3); !ok {
		t.Fatal("an unapplied optional resource blocked convergence")
	}
}

// Prune only after applying: a failed reconcile still needs to read what it was
// moving away from.
func TestPruneKeepsTheNamedGeneration(t *testing.T) {
	s, ctx := open(t)
	snapshot(t, s, 1, res("a", 0, true))
	snapshot(t, s, 2, res("a", 0, true))

	if err := s.PruneGenerationsBefore(ctx, 2); err != nil {
		t.Fatalf("prune: %v", err)
	}
	if got, _ := s.ListDesiredResources(ctx, 1); len(got) != 0 {
		t.Fatalf("generation 1 survived prune: %+v", got)
	}
	if got, _ := s.ListDesiredResources(ctx, 2); len(got) != 1 {
		t.Fatal("prune removed the generation it was told to keep")
	}
}

func TestAppliedGenerationRoundTrips(t *testing.T) {
	s, ctx := open(t)
	if err := s.SetAppliedGeneration(ctx, 42); err != nil {
		t.Fatalf("set: %v", err)
	}
	if got, _ := s.AppliedGeneration(ctx); got != 42 {
		t.Fatalf("applied generation = %d, want 42", got)
	}
}
