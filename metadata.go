package filesystem

import (
	"io"
	"time"
)

// Metadata is the interface used to represent file metadata
type Metadata interface {
	MimeType() string
	Timestamp() time.Time
	Visibility() Visibility
	Size() int64
}

// Content is the interface used to provide content
type Content interface {
	Content() string
}

// ContentStream is the interface used to provide content stream
type ContentStream interface {
	Content() io.ReadCloser
}

// Metadata is the type used to provide metadata about files.
type metadata map[string]interface{}
