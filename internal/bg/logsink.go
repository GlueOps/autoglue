package bg

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/glueops/autoglue/internal/models"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

const (
	// logFlushBytes is the buffered size that forces a write. A `make` target
	// on a bastion is chatty, and one row per line would be thousands of
	// INSERTs per run.
	logFlushBytes = 4 << 10

	// logFlushInterval bounds how stale a tail can be, so a job that goes quiet
	// mid-step still shows its last output to anyone watching.
	logFlushInterval = time.Second

	// logMaxTailBytes bounds what a sink keeps in memory to hand back to the
	// caller as the "output" string. A multi-hour run can emit far more than we
	// want to hold or put in an error message.
	logMaxTailBytes = 64 << 10
)

// LogSink is an io.Writer that appends job output to job_logs.
//
// Writes are batched: a flush happens once the buffer reaches logFlushBytes or
// logFlushInterval elapses, whichever comes first. Close drains whatever is
// left. It is safe for concurrent use, because stdout and stderr are pumped by
// separate goroutines into the same sink.
type LogSink struct {
	db        *gorm.DB
	jobID     int64
	orgID     uuid.UUID
	subject   string
	subjectID uuid.UUID
	stream    string

	mu     sync.Mutex
	buf    []byte
	tail   []byte
	closed bool

	cancel context.CancelFunc
	done   chan struct{}
}

// NewLogSink starts a sink writing against one subject. Always Close it, or the
// final partial chunk is lost and the ticker goroutine leaks.
func NewLogSink(db *gorm.DB, jobID int64, orgID uuid.UUID, subjectType string, subjectID uuid.UUID) *LogSink {
	ctx, cancel := context.WithCancel(context.Background())

	s := &LogSink{
		db:        db,
		jobID:     jobID,
		orgID:     orgID,
		subject:   subjectType,
		subjectID: subjectID,
		stream:    models.JobLogStreamStdout,
		cancel:    cancel,
		done:      make(chan struct{}),
	}

	go s.flushLoop(ctx)
	return s
}

func (s *LogSink) flushLoop(ctx context.Context) {
	defer close(s.done)

	t := time.NewTicker(logFlushInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			s.flush()
		case <-ctx.Done():
			return
		}
	}
}

// Write buffers p, flushing if it crosses the size threshold. It never returns
// an error: losing a log line must not fail the job that produced it.
func (s *LogSink) Write(p []byte) (int, error) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return len(p), nil
	}

	s.buf = append(s.buf, p...)

	s.tail = append(s.tail, p...)
	if len(s.tail) > logMaxTailBytes {
		s.tail = s.tail[len(s.tail)-logMaxTailBytes:]
	}

	over := len(s.buf) >= logFlushBytes
	s.mu.Unlock()

	if over {
		s.flush()
	}
	return len(p), nil
}

// System records a line from autoglue itself rather than the remote host, so a
// reader can tell "starting make bootstrap" from the command's own output.
func (s *LogSink) System(msg string) {
	s.append(models.JobLogStreamSystem, []byte(msg+"\n"))
}

func (s *LogSink) flush() {
	s.mu.Lock()
	if len(s.buf) == 0 {
		s.mu.Unlock()
		return
	}
	chunk := string(s.buf)
	s.buf = s.buf[:0]
	s.mu.Unlock()

	s.append(s.stream, []byte(chunk))
}

func (s *LogSink) append(stream string, chunk []byte) {
	if len(chunk) == 0 {
		return
	}

	row := models.JobLog{
		JobID:          s.jobID,
		OrganizationID: s.orgID,
		SubjectType:    s.subject,
		SubjectID:      s.subjectID,
		Stream:         stream,
		Chunk:          string(chunk),
	}

	if err := s.db.Create(&row).Error; err != nil {
		// Deliberately non-fatal: the job's own work matters more than its
		// narration, and this is already on an error path often enough.
		log.Warn().Err(err).
			Str("subject_type", s.subject).
			Str("subject_id", s.subjectID.String()).
			Msg("[joblog] failed to persist chunk")
	}
}

// Tail returns up to logMaxTailBytes of the most recent output. Used for error
// messages and for the value the old code returned from CombinedOutput.
func (s *LogSink) Tail() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return string(s.tail)
}

// Close drains the buffer and stops the flush loop. Safe to call twice.
func (s *LogSink) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	s.mu.Unlock()

	s.cancel()
	<-s.done

	// Final drain after the loop has stopped, so nothing races us.
	s.mu.Lock()
	chunk := string(s.buf)
	s.buf = nil
	s.mu.Unlock()

	if chunk != "" {
		s.append(s.stream, []byte(chunk))
	}
	return nil
}

var _ io.WriteCloser = (*LogSink)(nil)
