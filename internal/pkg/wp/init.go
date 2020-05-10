package wp

import (
	"errors"
	"os"
	"time"
)

type InitCommand struct {
	Command
	Delay time.Duration `short:"d" long:"delay" description:"wait before the operation (useful for xinit/xsession scripts)"`
}

func (cmd *InitCommand) Execute(args []string) error {
	if cmd.Delay.Seconds() > 0 {
		time.Sleep(cmd.Delay)
	}

	if err := cmd.SetFile(cmd.CurrentPath()); err != nil {
		if _, ok := err.(*os.PathError); ok {
			return errors.New("could not find most recently used wallpaper")
		}
		return err
	}

	return nil
}
