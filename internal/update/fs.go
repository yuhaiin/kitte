package update

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func writeFileAtomic(path string, data []byte, mode fs.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	f, err := os.CreateTemp(dir, "."+filepath.Base(path)+"-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)

	if err := f.Chmod(mode); err != nil {
		f.Close()
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Rename(tmp, path); err == nil {
		return nil
	}

	// Windows does not replace an existing destination with Rename.
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return os.Rename(tmp, path)
}

func replaceDir(path string, files map[string][]byte) error {
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}

	tmp, err := os.MkdirTemp(parent, "."+filepath.Base(path)+"-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		if filepath.Base(name) != name {
			return fmt.Errorf("invalid generated filename %q", name)
		}
		if err := os.WriteFile(filepath.Join(tmp, name), files[name], 0o644); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
