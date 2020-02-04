package clumad

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

// SaltOps object
type SaltOps struct {
	confpath   string
	saltminion *exec.Cmd
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
func (sop *SaltOps) GetConfDOption(conf string, key string) (interface{}, string) {
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
			return value, confFile
		}
	}

	return nil, ""
}

// GetConfOpion looks for a key of the main Salt configuration.
// If none is found, an empty string is returned (and it indicates
// a default value is used)
func (sop *SaltOps) GetConfOpion(conf string, key string) interface{} {
	return sop.getConfStruct(conf)[key]
}

// Get YAML content
func (sop *SaltOps) getConfStruct(conf string) map[string]interface{} {
	var confContent map[string]interface{}
	var confFile string
	if strings.Contains(conf, "/") { // Assumed relative or absolute path
		confFile = conf
	} else {
		confFile = path.Join(sop.confpath, conf)
	}
	fh, err := os.Open(confFile)
	if err != nil {
		confContent = make(map[string]interface{})
		return confContent
	} else {
		defer fh.Close()
	}
	confmap, err := ioutil.ReadAll(fh)
	if err != nil {
		panic("Error reading config: " + err.Error())
	}
	if err := yaml.Unmarshal(confmap, &confContent); err != nil {
		panic("Error parsing config: " + err.Error())
	}
	return confContent
}

// Backup main configuration file, if it wasnt.
func (sop *SaltOps) backupConfigFile(confPath string) {
	nfo, err := os.Stat(confPath + "~")
	if nfo != nil && nfo.IsDir() {
		panic("This is a directory")
	}
	if !os.IsNotExist(err) {
		src, err := os.Open(confPath)
		if err != nil {
			panic("Unable to open config file: " + err.Error())
		}
		defer src.Close()
		dst, err := os.OpenFile(confPath+"~", os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			panic("Unable to open backup config file: " + err.Error())
		}
		defer dst.Close()
		_, err = io.Copy(dst, src)
		if err != nil {
			panic("Unable to create backup config file: " + err.Error())
		}
	}
}

// Write a configuration option to the main config file
func (sop *SaltOps) SetConfOption(conf string, key string, value interface{}) {
	var confPath string
	if strings.Contains(conf, "/") {
		confPath = conf
	} else {
		confPath = path.Join(sop.confpath, conf)
	}
	sop.backupConfigFile(confPath)
	config := sop.getConfStruct(confPath)
	config[key] = value
	data, err := yaml.Marshal(&config)
	if err != nil {
		panic("Unable to render config data: " + err.Error())
	}
	err = ioutil.WriteFile(confPath, data, 0600)
	if err != nil {
		panic("Unable to write config to the file: " + err.Error())
	}
}

// SetConfOption writes a configuration option to a .d "drop-in" config file.
// If parameter conf is already a full path, it will be used.
// Otherwise new config file "/path/to/foo.d/clumad.conf" will be used.
func (sop *SaltOps) SetConfDOption(conf string, key string, value interface{}) {
	var confPath string
	if !strings.Contains(conf, "/") { // Assume new drop-in should be created
		confPath = path.Join(sop.confpath, conf+".d", "clumad.conf")
	} else {
		confPath = conf
	}
	sop.SetConfOption(confPath, key, value)
}

// StartSaltMinion starts Salt Minion in background and watches it.
func (sop *SaltOps) StartSaltMinion() {
	if sop.saltminion == nil {
		sop.saltminion = exec.Command("salt-minion", "-l", "debug")
		sop.saltminion.Env = append(os.Environ(), "FOO=BAR")
		sop.saltminion.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		if err := sop.saltminion.Start(); err != nil {
			panic("Unable to run Salt Minion: " + err.Error())
		}
	}
}

// StopSaltMinion terminates currently running Salt Minion.
func (sop *SaltOps) StopSaltMinion() {
	fmt.Println("Stopping salt minion")
	if sop.saltminion != nil {
		syscall.Kill(-sop.saltminion.Process.Pid, syscall.SIGKILL)
		sop.saltminion = nil
	}
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
