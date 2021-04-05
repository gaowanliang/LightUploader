// +build darwin dragonfly freebsd linux netbsd openbsd windows aix arm_linux solaris
// +build !nacl

/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func fsWatcherRoutine(s *Options, configDir string, filesWatching []string, initWG *sync.WaitGroup) {
	// effw.Lock()
	// if cmdrExitingForFsWatcher != nil {
	// 	effw.Unlock()
	// 	return
	// }
	//
	// cmdrExitingForFsWatcher = make(chan struct{}, 16)
	// effw.Unlock()
	// // initExitingChannelForFsWatcher()

	watcher, err := fsnotify.NewWatcher()
	if err == nil {
		defer watcher.Close()

		eventsWG := &sync.WaitGroup{}
		eventsWG.Add(1)
		go fsWatchRunner(s, configDir, filesWatching, watcher, eventsWG)
		_ = watcher.Add(configDir)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	} else {
		stopExitingChannelForFsWatcher()
	}
}

func fsWatchRunner(s *Options, configDir string, filesWatching []string, watcher *fsnotify.Watcher, eventsWG *sync.WaitGroup) {
	defer func() {
		// effw.Lock()
		// defer effw.Unlock()
		// if cmdrExitingForFsWatcher != nil {
		// 	close(cmdrExitingForFsWatcher)
		// 	cmdrExitingForFsWatcher = nil
		// }
		eventsWG.Done()
	}()
	for {
		select {
		case event, ok := <-watcher.Events:
			// ok == false: 'Events' channel is closed
			if ok {
				// log.Debugf("ooo | fsnotify.watcher |%v", event.String())
				// currentConfigFile, _ := filepath.EvalSymlinks(filename)
				// we only care about the config file with the following cases:
				// 1 - if the config file was modified or created
				// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
				const writeOrCreateMask = fsnotify.Write | fsnotify.Create

				if event.Op&writeOrCreateMask != 0 {
					suffixIsValid := testCfgSuffix(event.Name)
					if suffixIsValid {
						inside := strings.HasPrefix(filepath.Clean(event.Name), configDir)
						include := testArrayContains(event.Name, filesWatching)
						if inside || include {
							file, err := os.Open(event.Name)
							if err != nil {
								log.Printf("ERROR: os.Open() returned %v\n", err)
							} else {
								err = s.mergeConfigFile(bufio.NewReader(file), event.Name, path.Ext(event.Name))
								if err != nil {
									log.Printf("ERROR: os.Open() returned %v\n", err)
								}
								s.reloadConfig()
								_ = file.Close()

								if !include {
									filesWatching = append(filesWatching, event.Name)
								}
							}
						}
					}
				}
			}

		case err, ok := <-watcher.Errors:
			if ok { // 'Errors' channel is not closed
				// log.Printf("watcher error: %v\n", err)
				log.Printf("Watcher error: %v\n", err)
			}
			return

		case <-cmdrExitingForFsWatcher:
			return
		}
	}
}

// stopExitingChannelForFsWatcher stop fs watcher explicitly
func stopExitingChannelForFsWatcher() {
	effw.Lock()
	defer effw.Unlock()
	if cmdrExitingForFsWatcher != nil && GetBool("watching") {
		cmdrExitingForFsWatcher <- struct{}{}
	}
}

// func initExitingChannelForFsWatcher() {
// 	effw.Lock()
// 	defer effw.Unlock()
// 	if cmdrExitingForFsWatcher == nil {
// 		cmdrExitingForFsWatcher = make(chan struct{}, 16)
// 	}
// }

var cmdrExitingForFsWatcher = make(chan struct{}, 16)
var effw sync.RWMutex
