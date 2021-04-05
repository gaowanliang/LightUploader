// +build nacl

/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"os"
)

// TrapSignal is a helper for simplify your infinite loop codes.
func TrapSignals(onTrapped func(s os.Signal)) (waiter func()) {
	done := make(chan struct{}, 1)
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	//
	// go func() {
	// 	s := <-sigs
	// 	cmdr.Logger.Debugf("receive signal '%v'", s)
	//
	// 	onTrapped(s)
	//
	// 	// for _, s := range servers {
	// 	// 	s.Stop()
	// 	// }
	// 	// cmdr.Logger.Infof("done")
	// 	done <- struct{}{}
	// }()

	waiter = func() {
		for {
			select {
			case <-done:
				// os.Exit(1)
				// cmdr.Logger.Infof("done got.")
				return
			}
		}
	}

	return
}
