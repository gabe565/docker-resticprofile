package clix

import (
	"path/filepath"
	"strconv"

	"github.com/urfave/cli/v3"
)

func SecretFile(path *string, name string) cli.ValueSource {
	return &secretFileSource{path: path, name: name}
}

type secretFileSource struct {
	path *string
	name string
}

func (s *secretFileSource) getPath() string {
	if s.path == nil {
		return ""
	}
	return *s.path
}

func (s *secretFileSource) Lookup() (string, bool) {
	return cli.File(filepath.Join(s.getPath(), s.name)).Lookup()
}

func (s *secretFileSource) String() string {
	return "file " + filepath.Join(s.getPath(), s.name)
}

func (s *secretFileSource) GoString() string {
	return "&secretFileSource{path: " + strconv.Quote(s.getPath()) + ", name: " + strconv.Quote(s.name) + "}"
}
