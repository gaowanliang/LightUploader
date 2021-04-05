/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"gopkg.in/hedzr/errors.v2"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// GetOptions returns the global options instance (rxxtOptions),
// ie. cmdr Options Store
func GetOptions() *Options {
	return internalGetWorker().rxxtOptions
}

// GetUsedConfigFile returns the main config filename (generally
// it's `<appname>.yml`)
func GetUsedConfigFile() string {
	return internalGetWorker().rxxtOptions.usedConfigFile
}

// GetUsedConfigSubDir returns the sub-directory `conf.d` of config files.
// Note that it be always normalized now.
// Sometimes it might be empty string ("") if `conf.d` have not been found.
func GetUsedConfigSubDir() string {
	return internalGetWorker().rxxtOptions.usedConfigSubDir
}

// GetUsingConfigFiles returns all loaded config files, includes
// the main config file and children in sub-directory `conf.d`.
func GetUsingConfigFiles() []string {
	return internalGetWorker().rxxtOptions.configFiles
}

// GetWatchingConfigFiles returns all config files being watched.
func GetWatchingConfigFiles() []string {
	return internalGetWorker().rxxtOptions.filesWatching
}

// var rwlCfgReload = new(sync.RWMutex)

// AddOnConfigLoadedListener adds an functor on config loaded
// and merged
func AddOnConfigLoadedListener(c ConfigReloaded) {
	defer internalGetWorker().rxxtOptions.rwlCfgReload.Unlock()
	internalGetWorker().rxxtOptions.rwlCfgReload.Lock()

	// rwlCfgReload.RLock()
	if _, ok := internalGetWorker().rxxtOptions.onConfigReloadedFunctions[c]; ok {
		// rwlCfgReload.RUnlock()
		return
	}

	// rwlCfgReload.RUnlock()
	// rwlCfgReload.Lock()

	// defer rwlCfgReload.Unlock()

	internalGetWorker().rxxtOptions.onConfigReloadedFunctions[c] = true
}

// RemoveOnConfigLoadedListener remove an functor on config
// loaded and merged
func RemoveOnConfigLoadedListener(c ConfigReloaded) {
	w := internalGetWorker()
	defer w.rxxtOptions.rwlCfgReload.Unlock()
	w.rxxtOptions.rwlCfgReload.Lock()
	delete(w.rxxtOptions.onConfigReloadedFunctions, c)
}

// SetOnConfigLoadedListener enable/disable an functor on config
// loaded and merged
func SetOnConfigLoadedListener(c ConfigReloaded, enabled bool) {
	w := internalGetWorker()
	defer w.rxxtOptions.rwlCfgReload.Unlock()
	w.rxxtOptions.rwlCfgReload.Lock()
	w.rxxtOptions.onConfigReloadedFunctions[c] = enabled
}

// LoadConfigFile loads a yaml config file and merge the settings
// into `rxxtOptions`
// and load files in the `conf.d` child directory too.
func LoadConfigFile(file string) (mainFile, subDir string, err error) {
	return internalGetWorker().rxxtOptions.LoadConfigFile(file, true)
}

// LoadConfigFile loads a yaml config file and merge the settings
// into `rxxtOptions`
// and load files in the `conf.d` child directory too.
func (s *Options) LoadConfigFile(file string, main bool) (mainFile, subDir string, err error) {

	defer func() {
		s.batchMerging = false
		s.mapOrphans()
	}()

	s.batchMerging = true

	if !FileExists(file) {
		// log.Warnf("%v NOT EXISTS. PWD=%v", file, GetCurrentDir())
		return // not error, just ignore loading
	}

	if err = s.loadConfigFile(file); err != nil {
		return
	}

	mainFile = file
	dirWatch := path.Dir(mainFile)
	enableWatching := internalGetWorker().watchMainConfigFileToo
	confDFolderName := internalGetWorker().confDFolderName
	var filesWatching []string
	if main {
		dir := dirWatch
		s.usedConfigFile = mainFile
		_ = os.Setenv("CFG_DIR", dir)

		if internalGetWorker().watchMainConfigFileToo {
			filesWatching = append(filesWatching, s.usedConfigFile)
		}

		subDir = path.Join(dir, confDFolderName)
		if !FileExists(subDir) {
			subDir = ""
			if len(filesWatching) == 0 {
				return
			}
		}

		subDir, err = filepath.Abs(subDir)
		if err == nil {
			err = filepath.Walk(subDir, s.visit)
			if err == nil {
				if !internalGetWorker().watchMainConfigFileToo {
					dirWatch = subDir
				}
				filesWatching = append(filesWatching, s.configFiles...)
				enableWatching = true
			}
			// don't bring the minor error for sub-dir walking back to main caller
			err = nil
			// log.Fatalf("ERROR: filepath.Walk() returned %v\n", err)
		}

		s.usedConfigSubDir = subDir
	} else {
		s.usedAlterConfigFile = mainFile
	}
	err = s.doWatchConfigFile(enableWatching, confDFolderName, dirWatch, filesWatching)
	return
}

func (s *Options) doWatchConfigFile(enableWatching bool, confDFolderName, dirWatch string, filesWatching []string) (err error) {
	if internalGetWorker().watchChildConfigFiles {
		var dir string
		confDFolderName = os.ExpandEnv(".$APPNAME")
		dir, err = filepath.Abs(confDFolderName)
		if err == nil && FileExists(dir) {
			err = filepath.Walk(dir, s.visit)
			if err == nil {
				filesWatching = append(filesWatching, s.configFiles...)
				enableWatching = true
			}
			// don't bring the minor error for sub-dir walking back to main caller
			err = nil
		}
	}

	if enableWatching {
		s.watchConfigDir(dirWatch, filesWatching)
	}
	flog("the watching config files: %v", s.filesWatching)
	flog("the loaded config files: %v", s.configFiles)
	return
}

// Load a yaml config file and merge the settings into `Options`
func (s *Options) loadConfigFile(file string) (err error) {
	var m map[string]interface{}
	m, err = s.loadConfigFileAsMap(file)
	if err == nil {
		err = s.loopMap("", m)
	}
	//if err != nil {
	//	return
	//}
	return
}

func (s *Options) loadConfigFileAsMap(file string) (m map[string]interface{}, err error) {
	var (
		b  []byte
		mm map[string]map[string]interface{}
	)

	b, _ = ioutil.ReadFile(file)

	m = make(map[string]interface{})
	switch path.Ext(file) {
	case ".toml", ".ini", ".conf", "toml":
		mm = make(map[string]map[string]interface{})
		err = toml.Unmarshal(b, &mm)
		if err == nil {
			err = s.loopMapMap("", mm)
		}
		if err != nil {
			return
		}
		return

	case ".json", "json":
		err = json.Unmarshal(b, &m)
	default:
		err = yaml.Unmarshal(b, &m)
	}
	return
}

func (s *Options) mergeConfigFile(fr io.Reader, src string, ext string) (err error) {
	var (
		m   map[string]interface{}
		buf *bytes.Buffer
	)

	buf = new(bytes.Buffer)
	_, err = buf.ReadFrom(fr)

	m = make(map[string]interface{})
	switch ext {
	case ".toml", ".ini", ".conf", "toml":
		err = toml.Unmarshal(buf.Bytes(), &m)
	case ".json", "json":
		err = json.Unmarshal(buf.Bytes(), &m)
	default:
		err = yaml.Unmarshal(buf.Bytes(), &m)
	}

	if err == nil {
		err = s.loopMap("", m)
	}
	if err != nil {
		ferr("unsatisfied config file `%s` while importing: %v", src, err)
		return
	}

	return
}

func (s *Options) visit(path string, f os.FileInfo, e error) (err error) {
	// fmt.Printf("Visited: %s, e: %v\n", path, e)
	flog("    visiting: %v, e: %v", path, e)
	err = e
	if f != nil && !f.IsDir() && e == nil {
		// log.Infof("    path: %v, ext: %v", path, filepath.Ext(path))
		ext := filepath.Ext(path)
		switch ext {
		case ".yml", ".yaml", ".json", ".toml", ".ini", ".conf": // , "yml", "yaml":
			var file *os.File
			file, err = os.Open(path)
			// if err != nil {
			// log.Warnf("ERROR: os.Open() returned %v", err)
			// } else {
			if err == nil {
				defer file.Close()
				flog("    visited and merging: %v", file.Name())
				if err = s.mergeConfigFile(bufio.NewReader(file), file.Name(), ext); err != nil {
					err = errors.New("error in merging config file '%s': %v", path, err)
					return
				}
				s.configFiles = uniAddStr(s.configFiles, path)
			} else {
				err = errors.New("error in merging config file '%s': %v", path, err)
			}
		}
	}
	return
}

func (s *Options) reloadConfig() {
	// log.Debugf("\n\nConfig file changed: %s\n", e.String())

	defer s.rwlCfgReload.RUnlock()
	s.rwlCfgReload.RLock()

	for x, ok := range s.onConfigReloadedFunctions {
		if ok {
			x.OnConfigReloaded()
		}
	}
}

func (s *Options) watchConfigDir(configDir string, filesWatching []string) {
	if internalGetWorker().doNotWatchingConfigFiles || GetBoolR("no-watch-conf-dir") {
		return
	}

	if len(configDir) == 0 || len(filesWatching) == 0 {
		return
	}

	initWG := &sync.WaitGroup{}
	initWG.Add(1)
	// initExitingChannelForFsWatcher()
	s.filesWatching = filesWatching
	go fsWatcherRoutine(s, configDir, filesWatching, initWG)
	initWG.Wait() // make sure that the go routine above fully ended before returning
	s.SetNx("watching", true)
}

func testCfgSuffix(name string) bool {
	for _, suf := range []string{".yaml", ".yml", ".json", ".toml", ".ini", ".conf"} {
		if strings.HasSuffix(name, suf) {
			return true
		}
	}
	return false
}

func testArrayContains(s string, container []string) (contained bool) {
	for _, ss := range container {
		if ss == s {
			contained = true
			break
		}
	}
	return
}
