package plugins

import "errors"
import "github.com/maurofran/filesystem"

// EmptyDir is the plugin that will remove directory content.
type EmptyDir struct {
	plugin
}

// Method is the name of method to be used to invoke the plugin.
func (p *EmptyDir) Method() string {
	return "EmptyDir"
}

// Handle the invocation of empty dirs
func (p *EmptyDir) Handle(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("Missing dirname argument")
	}
	dirname, ok := args[0].(filesystem.Path)
	if !ok {
		return nil, errors.New("Invalid dirname argument: filesystem.Path required")
	}
	listing, err := p.fs.ListContents(dirname, false)
	if err != nil {
		return nil, err
	}
	for _, item := range listing {
		itemPath := item["path"].(filesystem.Path)
		if item["type"] == "dir" {
			if err := p.fs.DeleteDir(itemPath); err != nil {
				return nil, err
			}
		} else {
			_, err := p.fs.Delete(itemPath)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}
