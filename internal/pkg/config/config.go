package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/mkuznets/wp/internal/pkg/utils"
	"gopkg.in/yaml.v2"
)

const defaultCollection = 1065976

type Unsplash struct {
	Token       string
	Collections []int
}

type Tick struct {
	TTL time.Duration
}

type Fs struct {
	Gallery string
}

type Config struct {
	Fs       Fs
	Tick     Tick
	Unsplash Unsplash
	origin   string
}

func (c *Config) Origin() string {
	if c.origin != "" {
		return c.origin
	}
	return "defaults"
}

func New(path string) (*Config, error) {
	if path != "" {
		conf, err := FromFile(utils.ExpandHome(path))
		if err != nil {
			return nil, err
		}
		return conf, nil
	}

	var confPath string

	if configHome, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		confPath = filepath.Join(configHome, "wp.yaml")
	} else {
		confPath = utils.ExpandHome("~/.config/wp.yaml")
	}

	return FromFile(confPath)
}

func FromFile(path string) (*Config, error) {
	config := Default()

	f, err := os.Open(path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			// config file is missing, quietly loading defaults
			return config, nil
		}
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.Size() == 0 {
		return config, nil
	}

	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)

	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	if absPath, err := filepath.Abs(path); err == nil {
		config.origin = absPath
	}

	config.Fs.Gallery = utils.ExpandHome(config.Fs.Gallery)

	return config, nil
}

func Default() *Config {
	config := &Config{
		Fs: Fs{},
		Tick: Tick{
			TTL: 3 * time.Hour,
		},
		Unsplash: Unsplash{
			Collections: []int{defaultCollection},
		},
	}
	return config
}
