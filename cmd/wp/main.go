package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mkuznets/wp/internal/pkg/wp"
)

type Commander interface {
	Init(opts interface{}) error
	Execute(args []string) error
}

type Options struct {
	File     *wp.FileCommand     `command:"file" description:"set wallpaper from the given local file"`
	Unsplash *wp.UnsplashCommand `command:"unsplash" description:"set random wallpaper from Unsplash"`
	Tick     *wp.TickCommand     `command:"tick" description:"update wallpaper if the current one is old enough"`
	Init     *wp.InitCommand     `command:"init" description:"restore the most recent wallpaper"`
	Save     *wp.SaveCommand     `command:"save" description:"save the current wallpaper to the gallery"`
	Config   *wp.ConfigCommand   `command:"config" description:"print config"`
	Common   *wp.Options         `group:"Common Options"`
}

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	var opts Options
	parser := flags.NewParser(&opts, flags.Default)

	parser.CommandHandler = func(command flags.Commander, args []string) error {
		c := command.(Commander)
		if err := c.Init(opts.Common); err != nil {
			return err
		}
		if err := c.Execute(args); err != nil {
			return err
		}
		return nil
	}

	if _, err := parser.Parse(); err != nil {
		switch e := err.(type) {
		case *flags.Error:
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		os.Exit(1)
	}
}
