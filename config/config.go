package config

import (
	"github.com/BurntSushi/toml"
)

var configs = map[string]interface{}{}

// Require registers a request for package configuration.
func Require(key string, dest interface{}) {
	configs[key] = dest
}

// Load reads the config file and distributes provided configuration to
// requested destinations.
func Load(path string) error {
	var conf map[string]toml.Primitive
	meta, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return err
	}

	for key, prim := range conf {
		dest, ok := configs[key]
		if !ok {
			continue
		}
		err := meta.PrimitiveDecode(prim, dest)
		if err != nil {
			return err
		}
	}
	return nil
}
