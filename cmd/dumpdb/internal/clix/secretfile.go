package clix

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func SecretFile(path *string, name string) cli.ValueSource {
	return &secretFileSource{path: path, name: name}
}

type secretFileSource struct {
	path *string
	name string
}

func (s *secretFileSource) Lookup() (string, bool) {
	return cli.File(filepath.Join(*s.path, s.name)).Lookup()
}

func (s *secretFileSource) String() string {
	return fmt.Sprintf("file %s", filepath.Join(*s.path, s.name))
}

func (s *secretFileSource) GoString() string {
	return fmt.Sprintf("&secretFileSource{path: %q, name: %q}", *s.path, s.name)
}
