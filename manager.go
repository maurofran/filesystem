package filesystem

import (
	"io"
	"time"
)

// Read is the interface exposed for file system reading
type Read interface {
	// Has will check if a file exists.
	Has(path Path) (bool, error)
	// Read the file at provided path.
	Read(path Path) (string, error)
	// ReadStream will read the file at provided path as a stream.
	ReadStream(path Path) (io.ReadCloser, error)
	// GetMimeType will retrieve the mime type of file at supplied path.
	GetMimeType(path Path) (string, error)
	// GetTimestamp will retrieve the timestamp of file at supplied path.
	GetTimestamp(path Path) (time.Time, error)
	// GetFileSize will retrieve the size of file at supplied path.
	GetFileSize(path Path) (int64, error)
	// GetMetadata will retrieve the metadata of file at supplied path.
	GetMetadata(path Path) (Metadata, error)
	// Get the visibility of file at supplied path.
	GetVisibility(path Path) (Visibility, error)
	// List the contents of given path.
	ListContents(path Path, recursive bool) ([]Metadata, error)
}

// Write is the interface exposed for file system writing.
type Write interface {
	// Write the supplied content at supplied path, creating the file.
	Write(path Path, content string) error
	// WriteStream will write the content of provided reader at supplied path, creating the file.
	WriteStream(path Path, r io.Reader) error
	// Deletes a file at provided path.
	Delete(path Path) error
	// ReadAndDelete will read the file at provided path and delete after read.
	ReadAndDelete(path Path) (string, error)
	// Move the file at supplied path to new path.
	Move(path, newpath Path) error
	// Copy the file at supplied path to new path.
	Copy(path, newpath Path) error
	// CreateDir will create a new directory at provided path.
	CreateDir(path Path) error
	// DeleteDir will delete the directory at provided path.
	DeleteDir(path Path) error
	// Set the visibility of file at supplied path.
	SetVisibility(path Path, v Visibility) error
}

// Update is the interface exposed for file system update.
type Update interface {
	// Update the supplied content at supplied path, returning an error if file does not exists.
	Update(path Path, content string) error
	// Update with the content of supplied reader at supplied path, returning an error if file does not exists
	UpdateStream(path Path, r io.Reader) error
	// Put the supplied content at supplied path, creating the file if does not exists.
	Put(path Path, content string) error
	// Puth the content of supplied reader at supplied path, creating the file if does not exists.
	PutStream(path Path, r io.Reader) error
}

// Interface is interface exposed by file system objects.
type Interface interface {
	Read
	Write
	Update
}

type filesystem struct {
	Configurable
	Pluggable
	adapter Adapter
}
