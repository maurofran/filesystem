package filesystem

import (
	"io"
	"regexp"
	"strings"
	"time"
)

func splitPath(path Path) (string, Path, error) {
	match, err := regexp.MatchString("^[a-z]+://", string(path))
	if err != nil {
		return "", "", err
	}
	if !match {
		return "", "", invalidPathError(path)
	}
	idx := strings.Index(string(path), "://")
	prefix := string(path[:idx])
	subPath := path[idx+3:]
	return prefix, subPath, nil
}

// MountManager is the interface exposed by objects that allows to mount more file systems.
type MountManager interface {
	Interface
	// Mount the provided manager with provided prefix.
	Mount(prefix string, mgr Interface) error
	// Unmount the provided prefix.
	Unmount(prefix string) error
}

type mountManager struct {
	managers map[string]Interface
}

// EmptyMountManager will create a new empty mount manager.
func EmptyMountManager() MountManager {
	return &mountManager{}
}

func (mm *mountManager) Mount(prefix string, mgr Interface) error {
	if _, ok := mm.managers[prefix]; ok {
		return mountExistsError(prefix)
	}
	mm.managers[prefix] = mgr
	return nil
}

func (mm *mountManager) Unmount(prefix string) error {
	if _, ok := mm.managers[prefix]; !ok {
		return mountNotFoundError(prefix)
	}
	delete(mm.managers, prefix)
	return nil
}

func (mm *mountManager) managerFor(path Path) (Interface, Path, error) {
	prefix, subPath, err := splitPath(path)
	if err != nil {
		return nil, "", err
	}
	mgr, ok := mm.managers[prefix]
	if !ok {
		return nil, "", mountNotFoundError(prefix)
	}
	return mgr, subPath, nil
}

// Has will check if a file exists.
func (mm *mountManager) Has(path Path) (bool, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return false, err
	}
	return mgr.Has(subPath)
}

// Read the file at provided path.
func (mm *mountManager) Read(path Path) (string, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return "", err
	}
	return mgr.Read(subPath)
}

// ReadStream will read the file at provided path as a stream.
func (mm *mountManager) ReadStream(path Path) (io.ReadCloser, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return nil, err
	}
	return mgr.ReadStream(subPath)
}

// Write the supplied content at supplied path, creating the file.
func (mm *mountManager) Write(path Path, content string) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.Write(subPath, content)
}

// WriteStream will write the content of provided reader at supplied path, creating the file.
func (mm *mountManager) WriteStream(path Path, r io.Reader) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.WriteStream(subPath, r)
}

// Update the supplied content at supplied path, returning an error if file does not exists.
func (mm *mountManager) Update(path Path, content string) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.Update(subPath, content)
}

// Update with the content of supplied reader at supplied path, returning an error if file does not exists
func (mm *mountManager) UpdateStream(path Path, r io.Reader) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.UpdateStream(subPath, r)
}

// Put the supplied content at supplied path, creating the file if does not exists.
func (mm *mountManager) Put(path Path, content string) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.Put(subPath, content)
}

// Puth the content of supplied reader at supplied path, creating the file if does not exists.
func (mm *mountManager) PutStream(path Path, r io.Reader) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.PutStream(subPath, r)
}

// Deletes a file at provided path.
func (mm *mountManager) Delete(path Path) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.Delete(subPath)
}

// ReadAndDelete will read the file at provided path and delete after read.
func (mm *mountManager) ReadAndDelete(path Path) (string, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return "", err
	}
	return mgr.ReadAndDelete(subPath)
}

// Move the file at supplied path to new path.
func (mm *mountManager) Move(path, newpath Path) error {
	mgr1, subPath1, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	mgr2, subPath2, err := mm.managerFor(newpath)
	if err != nil {
		return err
	}
	if &mgr1 == &mgr2 {
		// The source and target managers are the same
		return mgr1.Move(subPath1, subPath2)
	}
	source, err := mgr1.ReadStream(subPath1)
	defer source.Close()
	if err != nil {
		return err
	}
	err = mgr2.WriteStream(subPath2, source)
	if err != nil {
		return err
	}
	return mgr1.Delete(subPath1)
}

// Copy the file at supplied path to new path.
func (mm *mountManager) Copy(path, newpath Path) error {
	mgr1, subPath1, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	mgr2, subPath2, err := mm.managerFor(newpath)
	if err != nil {
		return err
	}
	if &mgr1 == &mgr2 {
		return mgr1.Copy(subPath1, subPath2)
	}
	source, err := mgr1.ReadStream(subPath1)
	defer source.Close()
	if err != nil {
		return err
	}
	return mgr2.WriteStream(subPath2, source)
}

// GetMimeType will retrieve the mime type of file at supplied path.
func (mm *mountManager) GetMimeType(path Path) (string, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return "", err
	}
	return mgr.GetMimeType(subPath)
}

// GetTimestamp will retrieve the timestamp of file at supplied path.
func (mm *mountManager) GetTimestamp(path Path) (time.Time, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return time.Now(), err
	}
	return mgr.GetTimestamp(subPath)
}

// GetFileSize will retrieve the size of file at supplied path.
func (mm *mountManager) GetFileSize(path Path) (int64, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return 0, err
	}
	return mgr.GetFileSize(subPath)
}

// GetMetadata will retrieve the metadata of file at supplied path.
func (mm *mountManager) GetMetadata(path Path) (Metadata, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return nil, err
	}
	return mgr.GetMetadata(subPath)
}

// CreateDir will create a new directory at provided path.
func (mm *mountManager) CreateDir(path Path) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.CreateDir(subPath)
}

// DeleteDir will delete the directory at provided path.
func (mm *mountManager) DeleteDir(path Path) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.DeleteDir(subPath)
}

// Get the visibility of file at supplied path.
func (mm *mountManager) GetVisibility(path Path) (Visibility, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return 0, err
	}
	return mgr.GetVisibility(subPath)
}

// Set the visibility of file at supplied path.
func (mm *mountManager) SetVisibility(path Path, v Visibility) error {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return err
	}
	return mgr.SetVisibility(subPath, v)
}

// List the contents of given path.
func (mm *mountManager) ListContents(path Path, recursive bool) ([]Metadata, error) {
	mgr, subPath, err := mm.managerFor(path)
	if err != nil {
		return nil, err
	}
	return mgr.ListContents(subPath, recursive)
}
