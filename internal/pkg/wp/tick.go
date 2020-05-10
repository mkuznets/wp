package wp

import (
	"log"
	"os"
	"time"
)

type TickCommand struct {
	Force bool `short:"f" long:"force" description:"ignore TTL of the most recent wallpaper"`
	Command
}

func (cmd *TickCommand) Execute(args []string) error {
	if cmd.Force || cmd.needUpdate() {
		return cmd.SetUnsplash()
	}
	return nil
}

func (cmd *TickCommand) needUpdate() bool {
	fi, err := os.Stat(cmd.CurrentPath())
	if err == nil {
		elapsed := time.Since(fi.ModTime())
		if elapsed < cmd.Config.Tick.TTL {
			log.Printf("current wallpaper is recent enough (updated %v ago), do nothing", elapsed.Round(time.Minute))
			return false
		}
	} else {
		log.Print("could not find the most recent wallpaper, TTL is not checked")
	}
	return true
}
