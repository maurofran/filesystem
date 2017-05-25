package plugins

import "github.com/maurofran/filesystem"

type plugin struct {
	fs filesystem.Interface
}

func (p *plugin) SetFileSystem(fs filesystem.Interface) {
	p.fs = fs
}
