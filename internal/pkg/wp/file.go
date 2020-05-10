package wp

import (
	"fmt"
)

type FileCommand struct {
	Command
	Args struct {
		Path string `positional-arg-name:"FILE"`
	} `positional-args:"1" required:"1"`
}

func (cmd *FileCommand) Execute(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown arguments: %v", args)
	}
	return cmd.SetFile(cmd.Args.Path)
}
