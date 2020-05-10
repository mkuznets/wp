package wp

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/mkuznets/wp/internal/pkg/config"
	"github.com/mkuznets/wp/internal/pkg/unsplash"
	"github.com/mkuznets/wp/internal/pkg/utils"
	"github.com/mkuznets/wp/internal/pkg/xwallpaper"
)

// Options is a group of common options for all wp subcommands.
type Options struct {
	Config string `short:"c" long:"config" description:"custom config path"`
}

// Command is a common part of all wp subcommands.
type Command struct {
	Config  *config.Config
	workdir string
}

func (cmd *Command) Init(opts interface{}) error {
	options, ok := opts.(*Options)
	if !ok {
		panic("type mismatch")
	}

	conf, err := config.New(options.Config)
	if err != nil {
		return fmt.Errorf("could not read config: %v", err)
	}
	cmd.Config = conf

	// Setup working directory
	if stateHome, ok := os.LookupEnv("XDG_STATE_HOME"); ok {
		cmd.workdir = filepath.Join(stateHome, "wp")
	} else if dataHome, ok := os.LookupEnv("XDG_DATA_HOME"); ok {
		cmd.workdir = filepath.Join(dataHome, "wp")
	} else {
		cmd.workdir = utils.ExpandHome("~/.local/wp")
	}
	if err := os.MkdirAll(cmd.workdir, 0755); err != nil {
		return fmt.Errorf("could not create working directory: %v", err)
	}

	// Setup gallery directory
	if cmd.Config.Fs.Gallery != "" {
		if err := os.MkdirAll(cmd.Config.Fs.Gallery, 0755); err != nil {
			return fmt.Errorf("could not create gallery directory: %v", err)
		}
	}

	return nil
}

func (cmd *Command) CurrentPath() string {
	return filepath.Join(cmd.workdir, "current")
}

func (cmd *Command) SetUnsplash() error {
	if cmd.Config.Unsplash.Token == "" {
		return errors.New("unsplash token is not configured (unsplash.token)")
	}

	client := unsplash.New(cmd.Config.Unsplash.Token, cmd.Config.Unsplash.Collections)
	content, err := client.Random()
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewBuffer(content))
	if err != nil {
		return err
	}

	if err := xwallpaper.New().Set(img); err != nil {
		return err
	}

	fCurr, err := os.Create(cmd.CurrentPath())
	defer func() {
		if e := fCurr.Close(); e != nil {
			err = e
		}
	}()

	if err := utils.WriteFile(cmd.CurrentPath(), bytes.NewBuffer(content)); err != nil {
		return fmt.Errorf("could save wallpaper: %v", err)
	}

	return nil
}

func (cmd *Command) SetFile(path string) error {
	path = utils.ExpandHome(path)

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return errors.New("could not open file: regular file expected")
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	//noinspection GoUnhandledErrorResult
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("could not decode image: %v", err)
	}

	if err := xwallpaper.New().Set(img); err != nil {
		return fmt.Errorf("could not set wallpaper: %v", err)
	}

	if path == cmd.CurrentPath() {
		return nil
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := utils.WriteFile(cmd.CurrentPath(), f); err != nil {
		return fmt.Errorf("could not save wallpaper: %v", err)
	}

	return nil
}
