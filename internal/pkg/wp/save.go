package wp

import (
	"errors"
	"image"
	"io"
	"log"
	"os"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/mkuznets/wp/internal/pkg/utils"
)

type SaveCommand struct {
	Command
}

func (cmd *SaveCommand) Execute(args []string) error {

	if cmd.Config.Fs.Gallery == "" {
		return errors.New("gallery path is not configured (fs.gallery)")
	}

	f, err := os.Open(cmd.CurrentPath())
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return errors.New("could not find most recently used wallpaper")
		}
		return err
	}

	_, format, err := image.DecodeConfig(f)
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}
	mtime := fi.ModTime().Format("20060102_150405")
	path := filepath.Join(cmd.Config.Fs.Gallery, mtime+"."+format)

	if _, err := os.Stat(path); err == nil {
		log.Printf("the wallpaper is already in the gallery: %s", filepath.Base(path))
		return nil
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := utils.WriteFile(path, f); err != nil {
		return err
	}

	log.Printf("Saved at %s\n", path)

	return nil
}
