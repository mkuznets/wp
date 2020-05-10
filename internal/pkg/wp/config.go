package wp

import (
	"log"

	"github.com/mkuznets/wp/internal/pkg/config"
	"gopkg.in/yaml.v2"
)

type ConfigCommand struct {
	Defaults bool `long:"defaults" description:"show default settings"`
	Command
}

func (cmd *ConfigCommand) Execute(args []string) error {

	var conf *config.Config
	if cmd.Defaults {
		conf = config.Default()
	} else {
		conf = cmd.Config
	}

	out, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	log.Printf("# Origin: %s\n", conf.Origin())
	log.Println(string(out))
	return nil
}
