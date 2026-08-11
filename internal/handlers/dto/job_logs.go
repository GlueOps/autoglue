package dto

import "time"

// JobLogChunk is one batched piece of output from a background job.
// swagger:model JobLogChunk
type JobLogChunk struct {
	ID        int64     `json:"id" example:"1042"`
	Stream    string    `json:"stream" example:"stdout" enums:"stdout,system"`
	Chunk     string    `json:"chunk"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
}

// JobLogPage is a cursor page of job output.
//
// Poll by passing the previous `next_cursor` back as `after`. When `done` is
// true the underlying work has finished and no further output will appear, so
// the client can stop polling.
//
// swagger:model JobLogPage
type JobLogPage struct {
	Items      []JobLogChunk `json:"items"`
	NextCursor int64         `json:"next_cursor" example:"1042"`
	Done       bool          `json:"done" example:"false"`
}
