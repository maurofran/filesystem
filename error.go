package filesystem

import (
	"fmt"
)

/*
 * This package will try to implement the best practice provided in
 * https://dave.cheney.net/2014/12/24/inspecting-errors page.
 */

type fileSystemNotFoundError struct {
	name string
}

func (e fileSystemNotFoundError) Error() string {
	return fmt.Sprintf("File system not found: %s", e.name)
}

func (e fileSystemNotFoundError) FileSystemName() string {
	return e.name
}

// IsFileSystemNotFound will check if provided error is a file system not found and return the file system name
func IsFileSystemNotFound(err error) (bool, string) {
	type fileSystemNotFound interface {
		FileSystemName() string
	}

	if fsnf, ok := err.(fileSystemNotFound); ok {
		return ok, fsnf.FileSystemName()
	}
	return false, ""
}

type fileNotFoundError struct {
	missingPath Path
}

func (e fileNotFoundError) Error() string {
	return fmt.Sprintf("File not found at path: %s", e.missingPath)
}

func (e fileNotFoundError) MissingPath() Path {
	return e.missingPath
}

// IsFileNotFound will check if the provided error is a file not found error and return the missing path if true
func IsFileNotFound(err error) (bool, Path) {
	type fileNotFound interface {
		MissingPath() Path
	}

	if fnf, ok := err.(fileNotFound); ok {
		return ok, fnf.MissingPath()
	}
	return false, ""
}

type fileExistsError struct {
	existingPath Path
}

func (e fileExistsError) Error() string {
	return fmt.Sprintf("File already exists at path: %s", e.existingPath)
}

func (e fileExistsError) ExistingPath() Path {
	return e.existingPath
}

// IsFileExists will check if the provided error is a file exists error and return the existing path if true
func IsFileExists(err error) (bool, Path) {
	type fileExists interface {
		ExistingPath() Path
	}

	if fe, ok := err.(fileExists); ok {
		return ok, fe.ExistingPath()
	}
	return false, ""
}

type linkNotSupportedError struct {
	linkPath Path
}

func (e linkNotSupportedError) Error() string {
	return fmt.Sprintf("Links are not supported, encoundered link at %s", e.linkPath)
}

func (e linkNotSupportedError) LinkPath() Path {
	return e.linkPath
}

// IsLinkNotSupported will check if provided error is a link not supported error, returning the link path if true
func IsLinkNotSupported(err error) (bool, Path) {
	type linkNotSupported interface {
		LinkPath() Path
	}

	if lns, ok := err.(linkNotSupported); ok {
		return ok, lns.LinkPath()
	}
	return false, ""
}

type unreadableFileError struct {
	unreadablePath Path
}

func (e unreadableFileError) Error() string {
	return fmt.Sprintf("Unreadable file encountered at path: %s", e.unreadablePath)
}

func (e unreadableFileError) UnreadablePath() Path {
	return e.unreadablePath
}

// IsUnreadableFile will check if provided error is a unreadable file error, returning the path that it's unreadable
func IsUnreadableFile(err error) (bool, Path) {
	type unreadableFileError interface {
		UnreadablePath() Path
	}

	if uf, ok := err.(unreadableFileError); ok {
		return ok, uf.UnreadablePath()
	}
	return false, ""
}
