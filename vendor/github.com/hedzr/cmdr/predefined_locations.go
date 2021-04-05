// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/conf"
	"github.com/hedzr/cmdr/tool"
	"os"
	"path"
	"strings"
)

func (w *ExecWorker) parsePredefinedLocation() (err error) {
	// pre-detects for `--config xxx`, `--config=xxx`, `--configxxx`
	if ix, str, yes := partialContains(os.Args, "--config"); yes {
		var location string
		if i := strings.Index(str, "="); i > 0 {
			location = str[i+1:]
		} else if len(str) > 8 {
			location = str[8:]
		} else if ix+1 < len(os.Args) {
			location = os.Args[ix+1]
		}

		location = tool.StripQuotes(location)
		flog("--> preprocess / buildXref / parsePredefinedLocation: %q", location)

		if len(location) > 0 && FileExists(location) {
			if yes, err = IsDirectory(location); yes {
				if FileExists(path.Join(location, w.confDFolderName)) {
					setPredefinedLocations(location + "/%s.yml")
				} else {
					setPredefinedLocations(location + "/%s/%s.yml")
				}
			} else if yes, err = IsRegularFile(location); yes {
				setPredefinedLocations(location)
			}
		}
	}
	return
}

func (w *ExecWorker) checkMoreLocations(rootCmd *RootCommand) (err error) {
	if w.watchChildConfigFiles {
		a1, a2 := ".$APPNAME.yml", ".$APPNAME/*.yml"
		a3, a4 := os.ExpandEnv(a1), os.ExpandEnv(a2)
		b := FileExists(a3)
		if b {
			w.predefinedLocations = append(w.predefinedLocations, a3)
		}
		b = FileExists(a4)
		if b {
			//
		}
	}
	return
}

func (w *ExecWorker) loadFromPredefinedLocations(rootCmd *RootCommand) (err error) {
	err = w.checkMoreLocations(w.rootCommand)
	var mainFile, subDir string
	mainFile, subDir, err = w.loadFromLocations(rootCmd, w.getExpandedPredefinedLocations(), true)
	if err == nil {
		conf.CfgFile = mainFile
		flog("--> preprocess / buildXref / loadFromPredefinedLocations: %q loaded (CFG_DIR=%v)", mainFile, subDir)
		//flog("--> loadFromPredefinedLocations(): %q loaded", fn)
	}
	return
}

func (w *ExecWorker) loadFromAlterLocations(rootCmd *RootCommand) (err error) {
	err = w.checkMoreLocations(w.rootCommand)
	var mainFile, subDir string
	mainFile, subDir, err = w.loadFromLocations(rootCmd, w.getExpandedAlterLocations(), false)
	if err == nil {
		flog("--> preprocess / buildXref / loadFromAlterLocations: %q loaded (CFG_DIR=%v)", mainFile, subDir)
	}
	return
}

func (w *ExecWorker) loadFromLocations(rootCmd *RootCommand, locations []string, main bool) (mainFile, subDir string, err error) {
	// and now, loading the external configuration files
	for _, s := range locations {
		fn := s
		switch strings.Count(fn, "%s") {
		case 2:
			fn = fmt.Sprintf(s, rootCmd.AppName, rootCmd.AppName)
		case 1:
			fn = fmt.Sprintf(s, rootCmd.AppName)
		}

		b := FileExists(fn)
		if !b {
			fn = replaceAll(fn, ".yml", ".yaml")
			b = FileExists(fn)
		}
		if b {
			mainFile, subDir, err = w.rxxtOptions.LoadConfigFile(fn, main)
			break
		}
	}
	return
}

// getExpandedAlterLocations for internal using
func (w *ExecWorker) getExpandedAlterLocations() (locations []string) {
	for _, d := range internalGetWorker().alterLocations {
		locations = uniAddStr(locations, normalizeDir(d))
	}
	return
}

func setAlterLocations(locations ...string) {
	internalGetWorker().alterLocations = locations
}

// getExpandedPredefinedLocations for internal using
func (w *ExecWorker) getExpandedPredefinedLocations() (locations []string) {
	for _, d := range internalGetWorker().predefinedLocations {
		locations = uniAddStr(locations, normalizeDir(d))
	}
	return
}

// GetPredefinedLocations return the searching locations for loading config files.
func GetPredefinedLocations() []string {
	return internalGetWorker().predefinedLocations
}

// // SetPredefinedLocations to customize the searching locations for loading config files.
// //
// // It MUST be invoked before `cmdr.Exec`. Such as:
// // ```go
// //     SetPredefinedLocations([]string{"./config", "~/.config/cmdr/", "$GOPATH/running-configs/cmdr"})
// // ```
// func SetPredefinedLocations(locations []string) {
// 	uniqueWorker.predefinedLocations = locations
// }

func setPredefinedLocations(locations ...string) {
	internalGetWorker().predefinedLocations = locations
}
