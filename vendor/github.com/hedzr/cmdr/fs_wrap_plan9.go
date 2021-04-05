// +build plan9

// Copyright © 2019 Hedzr Yeh.

package cmdr

import (
	"sync"
)

func fsWatcherRoutine(s *Options, configDir string, filesWatching []string, initWG *sync.WaitGroup) {
	initWG.Done() // done initializing the watch in this go routine, so the parent routine can move on...
}

func fsWatchRunner(s *Options, configDir string, watcher *fsnotify.Watcher, eventsWG *sync.WaitGroup) {
	// eventsWG.Done()
}

// stopExitingChannelForFsWatcher stop fs watcher explicitly
func stopExitingChannelForFsWatcher() {
	effw.Lock()
	defer effw.Unlock()
	if cmdrExitingForFsWatcher != nil && GetBool("watching") {
		cmdrExitingForFsWatcher <- struct{}{}
	}
}

var cmdrExitingForFsWatcher = make(chan struct{}, 16)
var effw sync.RWMutex
