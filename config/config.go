package config

import (
	"github.com/BurntSushi/toml"
)

var configs = map[string]interface{}{}

// Require registers a request for package configuration.
func Require(key string, dest interface{}) {
	configs[key] = dest
}

// Load ...
func Load(src string) error {
	return load(src, toml.Decode)
}

// LoadFile reads the config file and distributes provided configuration to
// requested destinations.
func LoadFile(path string) error {
	return load(path, toml.DecodeFile)
}

func load(from string, withFn func(string, interface{}) (toml.MetaData, error)) error {
	var sections map[string]toml.Primitive
	meta, err := withFn(from, &sections)
	if err != nil {
		return err
	}
	return decode(meta, sections)
}

func decode(meta toml.MetaData, sections map[string]toml.Primitive) error {
	for section, conf := range sections {
		dest, ok := configs[section]
		if !ok {
			continue
		}

		err := meta.PrimitiveDecode(conf, dest)
		if err != nil {
			return err
		}
	}
	return nil
}
