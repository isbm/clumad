package uccd

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	data map[string]interface{}
}

func NewConfig(cfgpath string) *Config {
	cfg := new(Config)

	fh, err := os.Open(cfgpath)
	if err != nil {
		cfg.data = make(map[string]interface{})
	} else {
		defer fh.Close()
		confmap, err := ioutil.ReadAll(fh)
		if err != nil {
			panic("Error reading config: " + err.Error())
		}
		if err := yaml.Unmarshal(confmap, &cfg.data); err != nil {
			panic("Error parsing config: " + err.Error())
		}
	}

	return cfg
}

// GetContent returns a content of the configuration
func (cfg *Config) GetContent() map[string]interface{} {
	return cfg.data
}

// String returns a string type of a config value.
func (cfg *Config) String(key string, overlay string) string {
	if overlay != "" {
		return overlay
	}
	return fmt.Sprintf("%s", cfg.data[key])
}

// Int returns an integer type of a config value.
func (cfg *Config) Int(key string, overlay string) int {
	var v string
	if overlay != "" {
		v = overlay
	} else {
		v = fmt.Sprintf("%d", cfg.data[key])
	}
	data, err := strconv.Atoi(v)
	if err != nil {
		panic(err)
	}
	return data
}
