package plugins

import (
	"errors"

	"github.com/maurofran/filesystem"
)

// ForceRename will rename a file, overwriting if does not exists.
type ForceRename struct {
	plugin
}

// Method is the name of the method to be used to invoke the plugin.
func (p *ForceRename) Method() string {
	return "ForceRename"
}

// Handle the invocation of plugin.
func (p *ForceRename) Handle(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return false, errors.New("path and newpath arguments are required")
	}
	path, ok := args[0].(filesystem.Path)
	if !ok {
		return false, errors.New("path must be an instance of filesystem.Path")
	}
	newPath, ok := args[1].(filesystem.Path)
	if !ok {
		return false, errors.New("newPath must be an instance of filesystem.Path")
	}
	deleted, err := p.fs.Delete(newPath)
	if err != nil {
		if filesystem.IsFileNotFound(err) {
			deleted = true
		} else {
			return false, err
		}
	}
	if deleted {
		err := p.fs.Move(path, newPath)
		return true, err
	}
	return false, nil
}
