package filesystem

import (
	"errors"
	"io"
	"time"
)

// Interface is exposed for file system management.
type Interface interface {
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
	// Write the supplied content at supplied path, creating the file.
	Write(path Path, content string, config map[string]interface{}) error
	// WriteStream will write the content of provided reader at supplied path, creating the file.
	WriteStream(path Path, r io.Reader, config map[string]interface{}) error
	// Deletes a file at provided path.
	Delete(path Path) (bool, error)
	// ReadAndDelete will read the file at provided path and delete after read.
	ReadAndDelete(path Path) (string, error)
	// Move the file at supplied path to new path.
	Move(path, newpath Path) error
	// Copy the file at supplied path to new path.
	Copy(path, newpath Path) error
	// CreateDir will create a new directory at provided path.
	CreateDir(path Path, config map[string]interface{}) error
	// DeleteDir will delete the directory at provided path.
	DeleteDir(path Path) error
	// Set the visibility of file at supplied path.
	SetVisibility(path Path, v Visibility) error
	// Update the supplied content at supplied path, returning an error if file does not exists.
	Update(path Path, content string, config map[string]interface{}) error
	// Update with the content of supplied reader at supplied path, returning an error if file does not exists
	UpdateStream(path Path, r io.Reader, config map[string]interface{}) error
	// Put the supplied content at supplied path, creating the file if does not exists.
	Put(path Path, content string, config map[string]interface{}) error
	// Puth the content of supplied reader at supplied path, creating the file if does not exists.
	PutStream(path Path, r io.Reader, config map[string]interface{}) error
}

// New will create a new file system.
func New(adapter Adapter, config *Config) (Interface, error) {
	if adapter == nil {
		return nil, errors.New("adapter is required")
	}
	fs := new(filesystem)
	fs.adapter = adapter
	fs.SetConfig(config)
	return fs, nil
}

type filesystem struct {
	Configurable
	Pluggable
	adapter Adapter
}

// Adapter will retrieve the file system adapter.
func (fs *filesystem) Adapter() Adapter {
	return fs.adapter
}

// Has will check if a file exists.
func (fs *filesystem) Has(path Path) (bool, error) {
	path, err := normalizePath(path)
	if err != nil {
		return false, err
	}
	if path == emptyPath {
		return false, nil
	}
	// Delegeate to adapter.
	return fs.Adapter().Has(path)
}

// Read the file at provided path.
func (fs *filesystem) Read(path Path) (string, error) {
	path, err := normalizePath(path)
	if err != nil {
		return "", err
	}
	meta, err := fs.Adapter().Read(path)
	if err != nil {
		return "", err
	}
	return meta.Content(), nil
}

// ReadStream will read the file at provided path as a stream.
func (fs *filesystem) ReadStream(path Path) (io.ReadCloser, error) {
	path, err := normalizePath(path)
	if err != nil {
		return nil, err
	}
	meta, err := fs.Adapter().ReadStream(path)
	if err != nil {
		return nil, err
	}
	return meta.Content(), nil
}

// GetMimeType will retrieve the mime type of file at supplied path.
func (fs *filesystem) GetMimeType(path Path) (string, error) {
	path, err := normalizePath(path)
	if err != nil {
		return "", err
	}
	meta, err := fs.Adapter().GetMetadata(path)
	if err != nil {
		return "", err
	}
	return meta.MimeType(), nil
}

// GetTimestamp will retrieve the timestamp of file at supplied path.
func (fs *filesystem) GetTimestamp(path Path) (time.Time, error) {
	path, err := normalizePath(path)
	if err != nil {
		return time.Now(), err
	}
	meta, err := fs.Adapter().GetMetadata(path)
	if err != nil {
		return time.Now(), err
	}
	return meta.Timestamp(), nil
}

// GetFileSize will retrieve the size of file at supplied path.
func (fs *filesystem) GetFileSize(path Path) (int64, error) {
	path, err := normalizePath(path)
	if err != nil {
		return 0, err
	}
	meta, err := fs.Adapter().GetMetadata(path)
	if err != nil {
		return 0, err
	}
	return meta.Size(), nil
}

// GetMetadata will retrieve the metadata of file at supplied path.
func (fs *filesystem) GetMetadata(path Path) (Metadata, error) {
	path, err := normalizePath(path)
	if err != nil {
		return nil, err
	}
	return fs.Adapter().GetMetadata(path)
}

// Get the visibility of file at supplied path.
func (fs *filesystem) GetVisibility(path Path) (Visibility, error) {
	path, err := normalizePath(path)
	if err != nil {
		return Visibility(0), err
	}
	meta, err := fs.Adapter().GetMetadata(path)
	if err != nil {
		return Visibility(0), err
	}
	return meta.Visibility(), nil
}

// List the contents of given path.
func (fs *filesystem) ListContents(path Path, recursive bool) ([]Metadata, error) {
	path, err := normalizePath(path)
	if err != nil {
		return nil, err
	}
	// TODO ContentListFormatter????
	return fs.Adapter().ListContents(path, recursive)
}

// Write the supplied content at supplied path, creating the file.
func (fs *filesystem) Write(path Path, content string, config map[string]interface{}) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	if err := fs.assertAbsent(path); err != nil {
		return err
	}
	return fs.Adapter().Write(path, content, fs.PrepareConfig(config))
}

// WriteStream will write the content of provided reader at supplied path, creating the file.
func (fs *filesystem) WriteStream(path Path, r io.Reader, config map[string]interface{}) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	if err := fs.assertAbsent(path); err != nil {
		return err
	}
	return fs.Adapter().WriteStream(path, r, fs.PrepareConfig(config))
}

// Deletes a file at provided path.
func (fs *filesystem) Delete(path Path) (bool, error) {
	path, err := normalizePath(path)
	if err != nil {
		return false, err
	}
	if err := fs.assertPresent(path); err != nil {
		return false, err
	}
	return fs.Adapter().Delete(path)
}

// ReadAndDelete will read the file at provided path and delete after read.
func (fs *filesystem) ReadAndDelete(path Path) (string, error) {
	path, err := normalizePath(path)
	if err != nil {
		return "", err
	}
	if err := fs.assertPresent(path); err != nil {
		return "", err
	}
	content, err := fs.Read(path)
	if err != nil {
		return "", err
	}
	if _, err = fs.Delete(path); err != nil {
		return "", err
	}
	return content, nil
}

// Move the file at supplied path to new path.
func (fs *filesystem) Move(path, newpath Path) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	newpath, err = normalizePath(newpath)
	if err != nil {
		return err
	}
	if err := fs.assertPresent(path); err != nil {
		return err
	}
	if err := fs.assertAbsent(newpath); err != nil {
		return err
	}
	return fs.Adapter().Move(path, newpath)
}

// Copy the file at supplied path to new path.
func (fs *filesystem) Copy(path, newpath Path) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	newpath, err = normalizePath(newpath)
	if err != nil {
		return err
	}
	if err := fs.assertPresent(path); err != nil {
		return err
	}
	if err := fs.assertAbsent(path); err != nil {
		return err
	}
	return fs.Adapter().Copy(path, newpath)
}

// CreateDir will create a new directory at provided path.
func (fs *filesystem) CreateDir(path Path, config map[string]interface{}) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	return fs.Adapter().CreateDir(path, fs.PrepareConfig(config))
}

// DeleteDir will delete the directory at provided path.
func (fs *filesystem) DeleteDir(path Path) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	if path == emptyPath {
		return errors.New("Root directories cannot be deleted")
	}
	return fs.Adapter().DeleteDir(path)
}

// Set the visibility of file at supplied path.
func (fs *filesystem) SetVisibility(path Path, v Visibility) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	if err := fs.assertPresent(path); err != nil {
		return err
	}
	return fs.Adapter().SetVisibility(path, v)
}

// Update the supplied content at supplied path, returning an error if file does not exists.
func (fs *filesystem) Update(path Path, content string, config map[string]interface{}) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	if err := fs.assertPresent(path); err != nil {
		return err
	}
	return fs.Adapter().Update(path, content, fs.PrepareConfig(config))
}

// Update with the content of supplied reader at supplied path, returning an error if file does not exists
func (fs *filesystem) UpdateStream(path Path, r io.Reader, config map[string]interface{}) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	if err := fs.assertPresent(path); err != nil {
		return err
	}
	return fs.Adapter().UpdateStream(path, r, fs.PrepareConfig(config))
}

// Put the supplied content at supplied path, creating the file if does not exists.
func (fs *filesystem) Put(path Path, content string, config map[string]interface{}) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	cfg := fs.PrepareConfig(config)

	hasPath, err := fs.Has(path)
	if err != nil {
		return err
	}
	if hasPath && fs.canOverwrite() {
		return fs.Adapter().Update(path, content, cfg)
	}
	return fs.Adapter().Write(path, content, cfg)
}

// Puth the content of supplied reader at supplied path, creating the file if does not exists.
func (fs *filesystem) PutStream(path Path, r io.Reader, config map[string]interface{}) error {
	path, err := normalizePath(path)
	if err != nil {
		return err
	}
	cfg := fs.PrepareConfig(config)
	hasPath, err := fs.Has(path)
	if err != nil {
		return err
	}
	if hasPath && fs.canOverwrite() {
		return fs.Adapter().UpdateStream(path, r, cfg)
	}
	return fs.Adapter().WriteStream(path, r, cfg)
}

func (fs *filesystem) canOverwrite() bool {
	return fs.Adapter().CanOverwrite()
}

func (fs *filesystem) assertPresent(path Path) error {
	v, ok := fs.Config().Get("disableAsserts", false).(bool)
	if !ok || v {
		pathPresent, err := fs.Has(path)
		if err != nil {
			return err
		}
		if !pathPresent {
			return fileNotFoundError{missingPath: path}
		}
	}
	return nil
}

func (fs *filesystem) assertAbsent(path Path) error {
	v, ok := fs.Config().Get("disableAsserts", false).(bool)
	if !ok || v {
		pathPresent, err := fs.Has(path)
		if err != nil {
			return err
		}
		if pathPresent {
			return fileExistsError{existingPath: path}
		}
	}
	return nil
}
