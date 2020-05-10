package utils

import (
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func ExpandHome(path string) string {
	if len(path) > 1 && path[:2] == "~/" {
		u, err := user.Current()
		if err != nil {
			log.Fatalf("could not get current user: %v", err)
		}
		return filepath.Join(u.HomeDir, path[1:])
	}
	return path
}

func WriteFile(path string, r io.Reader) (err error) {
	fCurr, err := os.Create(path)
	defer func() {
		if e := fCurr.Close(); e != nil {
			err = e
		}
	}()

	if err != nil {
		return err
	}
	if _, err := io.Copy(fCurr, r); err != nil {
		return err
	}
	if err := fCurr.Sync(); err != nil {
		return err
	}

	return
}
