package clumad

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// SaltOps object
type SaltOps struct {
	confpath string
}

// NewSaltOps is a SaltOps constructor
func NewSaltOps() *SaltOps {
	sop := new(SaltOps)
	return sop
}

// SetSaltConfigDir sets Salt configuration directory
func (sop *SaltOps) SetSaltConfigDir(path string) *SaltOps {
	sop.confpath = path
	return sop
}

// GetConfDOption looks for a key of the additional Salt configuration.
// if it was in any of config files of a ".d" directory.
// If none is found, an empty string is returned (and it indicates
// a default value is used or main config has it).
func (sop *SaltOps) GetConfDOption(conf string, key string) interface{} {
	var value interface{}
	var files []string
	filepath.Walk(path.Join(sop.confpath, conf+".d"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	for _, confFile := range files {
		value = sop.GetConfOpion(confFile, key)
		if value != nil {
			return value
		}
	}

	return nil
}

// GetConfOpion looks for a key of the main Salt configuration.
// If none is found, an empty string is returned (and it indicates
// a default value is used)
func (sop *SaltOps) GetConfOpion(conf string, key string) interface{} {
	var confContent map[string]interface{}
	var confFile string
	if strings.Contains(conf, "/") { // Assumed relative or absolute path
		confFile = conf
	} else {
		confFile = path.Join(sop.confpath, conf)
	}
	fh, err := os.Open(confFile)
	if err != nil {
		panic("Unable to open configuration file: " + confFile)
	}
	defer fh.Close()
	confmap, err := ioutil.ReadAll(fh)
	if err != nil {
		panic("Error reading config: " + err.Error())
	}
	if err := yaml.Unmarshal(confmap, &confContent); err != nil {
		panic("Error parsing config: " + err.Error())
	}

	return confContent[key]
}

/////////////////////////////////
//
// SysOp object
type SysOp struct {
	saltOps *SaltOps
}

// NewSysOp is a SysOp constructor
func NewSysOp() *SysOp {
	sysop := new(SysOp)
	sysop.saltOps = NewSaltOps()
	return sysop
}

// GetSaltOps returns all methods to Salt Operation configuration on the client
func (sysop *SysOp) GetSaltOps() *SaltOps {
	return sysop.saltOps
}
