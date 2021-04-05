// +build darwin dragonfly freebsd linux netbsd openbsd windows aix arm_linux solaris
// +build !nacl

/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"os"
	"os/signal"
	"syscall"
)

// TrapSignals is a helper for simplify your infinite loop codes.
//
// Usage
//
//  func enteringLoop() {
// 	  waiter := cmdr.TrapSignals(func(s os.Signal) {
// 	    cmdr.Logger.Debugf("receive signal '%v' in onTrapped()", s)
// 	  })
// 	  go waiter()
//  }
//
//
//
func TrapSignals(onTrapped func(s os.Signal), signals ...os.Signal) (waiter func()) {
	done := make(chan bool, 1)
	waiter = TrapSignalsEnh(done, onTrapped)
	return
}

// TrapSignalsEnh is a helper for simplify your infinite loop codes.
//
// Usage
//
//  func enteringLoop() {
//    done := make(chan bool, 1)
//    go func(done chan bool){
//       // your main serve loop
//       done <- true   // to end the TrapSignalsEnh waiter by manually, instead of os signals caught.
//    }(done)
// 	  waiter := cmdr.TrapSignalsEnh(done, func(s os.Signal) {
// 	    cmdr.Logger.Debugf("receive signal '%v' in onTrapped()", s)
// 	  })
// 	  go waiter()
//  }
//
//
//
func TrapSignalsEnh(done chan bool, onTrapped func(s os.Signal), signals ...os.Signal) (waiter func()) {
	sigs := make(chan os.Signal, 1)
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT}
	}
	signal.Notify(sigs, signals...)

	go func() {
		defer close(sigs)
		defer close(done)
		for {
			select {
			case s := <-sigs:
				if !silent() {
					Logger.Printf("receive signal '%v'", s)
				}

				onTrapped(s)

				// for _, s := range servers {
				// 	s.Stop()
				// }
				// Logger.Infof("done")
				done <- false
				return
			case <-done:
				if !silent() {
					Logger.Printf("receive done sig and return for-select go-routine")
				}
				return
			}
		}
	}()

	waiter = func() {
		for {
			select {
			case byManual := <-done:
				if byManual {
					done <- true // stop os signals for-select looper
				}
				return // os.Exit(1) // log.Infof("done got.")
			}
		}
	}

	return
}

// SignalToSelf trigger the sig signal to the current process
func SignalToSelf(sig os.Signal) (err error) {
	var p *os.Process
	if p, err = os.FindProcess(os.Getpid()); err != nil {
		Logger.Printf("error: can't find process with pid=%v: %v", os.Getpid(), err)
	}
	err = p.Signal(sig)
	return
}

// SignalQuitSignal post a SIGQUIT signal to the current process
func SignalQuitSignal() {
	_ = SignalToSelf(syscall.SIGQUIT)
}

// SignalTermSignal post a SIGTERM signal to the current process
func SignalTermSignal() {
	_ = SignalToSelf(syscall.SIGTERM)
}

func silent() bool {
	return GetQuietMode()
}
