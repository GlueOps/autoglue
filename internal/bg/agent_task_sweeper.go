package bg

import (
	"context"
	"fmt"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/riverqueue/river"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// AgentTaskSweeperArgs drives the only path that can settle a task whose agent
// stopped answering.
//
// Without it the task plane deadlocks. A task only leaves assigned or started
// when the agent reports a terminal outcome, and an agent that never comes back
// never reports one — so the task holds the agent's single in-flight slot
// forever and its ClusterRun sits neither succeeded nor failed, with nothing in
// the system able to move either.
//
// It dead-letters rather than retries, and that is not timidity. The make
// targets are supplied rather than authored here, so "the container may have
// completed and we never heard" and "the container half-applied a terraform
// plan" are indistinguishable from the control plane, and re-running the second
// one is how you corrupt state. Somebody has to look.
type AgentTaskSweeperArgs struct{}

func (AgentTaskSweeperArgs) Kind() string { return "agent_task_sweeper" }

func (AgentTaskSweeperArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{Queue: QueueMaintenance, MaxAttempts: 2}
}

type AgentTaskSweeperResult struct {
	Status       string `json:"status"`
	DeadLettered int    `json:"dead_lettered"`
}

type AgentTaskSweeperWorker struct {
	river.WorkerDefaults[AgentTaskSweeperArgs]
	db *gorm.DB
}

func (w *AgentTaskSweeperWorker) Timeout(*river.Job[AgentTaskSweeperArgs]) time.Duration {
	return 2 * time.Minute
}

func (w *AgentTaskSweeperWorker) Work(ctx context.Context, _ *river.Job[AgentTaskSweeperArgs]) error {
	cutoff := time.Now().UTC().Add(-agentLostAfter())

	var stranded []models.AgentTask
	if err := w.db.WithContext(ctx).
		Joins("JOIN agents ON agents.id = agent_tasks.agent_id").
		Where("agent_tasks.state IN ?", []string{
			models.AgentTaskStateAssigned, models.AgentTaskStateStarted,
		}).
		// A NULL last_seen_at means the agent enrolled and never called home at
		// all, so it is judged from enrolment instead. Leaving NULL out would
		// make the never-started case the one thing the sweeper cannot reach.
		Where("COALESCE(agents.last_seen_at, agents.enrolled_at) < ?", cutoff).
		Find(&stranded).Error; err != nil {
		return fmt.Errorf("find stranded agent tasks: %w", err)
	}

	swept := 0
	for _, task := range stranded {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := w.deadLetter(ctx, task); err != nil {
			// One stranded task failing must not cost the rest their sweep;
			// they are already stuck and the next tick retries this one.
			log.Warn().Err(err).
				Str("task_id", task.ID.String()).
				Str("run_id", task.RunID.String()).
				Msg("[agent_task_sweeper] could not dead-letter stranded task")
			continue
		}
		swept++
		log.Warn().
			Str("task_id", task.ID.String()).
			Str("agent_id", task.AgentID.String()).
			Str("run_id", task.RunID.String()).
			Msg("[agent_task_sweeper] agent stopped calling home; task dead-lettered for review")
	}

	if err := river.RecordOutput(ctx, AgentTaskSweeperResult{
		Status:       "ok",
		DeadLettered: swept,
	}); err != nil {
		log.Warn().Err(err).Msg("[agent_task_sweeper] could not record output")
	}
	return nil
}

// deadLetter settles the task and its run together. One transaction, because a
// task marked terminal without its run settled swaps one kind of limbo for
// another: the agent's slot frees, the next task is assigned, and the run it
// came from is still sitting there claiming to be in progress.
func (w *AgentTaskSweeperWorker) deadLetter(ctx context.Context, task models.AgentTask) error {
	now := time.Now().UTC()

	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// The state predicate is the race guard: an agent that reports a
		// terminal outcome between the SELECT above and this UPDATE wins, and
		// the sweep silently does nothing rather than overwriting a real result
		// with a guess.
		res := tx.Model(&models.AgentTask{}).
			Where("id = ? AND state IN ?", task.ID, []string{
				models.AgentTaskStateAssigned, models.AgentTaskStateStarted,
			}).
			Updates(map[string]any{
				"state":              models.AgentTaskStateDeadLettered,
				"dead_letter_reason": models.AgentTaskDeadLetterAgentLost,
				"error":              "agent stopped calling home; outcome unknown",
				"ended_at":           now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return nil
		}

		// Only settle a run that is still open. A redelivered or late result
		// must not resurrect one that already reached a terminal state.
		return tx.Model(&models.ClusterRun{}).
			Where("id = ? AND status IN ?", task.RunID, []string{
				models.ClusterRunStatusQueued, models.ClusterRunStatusRunning,
			}).
			Updates(map[string]any{
				"status":      models.ClusterRunStatusFailed,
				"error":       "dead-lettered (" + models.AgentTaskDeadLetterAgentLost + "): agent stopped calling home",
				"finished_at": now,
				"updated_at":  now,
			}).Error
	})
}

// agentLostAfter is how long an agent may go quiet before its work is written
// off. It has to exceed the agent's own poll interval by a wide margin: a
// bastion mid-way through a multi-hour make target still calls home on the
// config loop, so silence really does mean the process is gone rather than busy.
func agentLostAfter() time.Duration {
	if m := viper.GetInt("agent.lost_after_minutes"); m > 0 {
		return time.Duration(m) * time.Minute
	}
	return 15 * time.Minute
}
