package wp

import (
	"fmt"
)

type UnsplashCommand struct {
	Command
}

func (cmd *UnsplashCommand) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown arguments: %v", args)
	}
	return cmd.SetUnsplash()
}
