/*
 * Copyright © 2019 Hedzr Yeh.
 */

package logex

// import "github.com/coreos/go-systemd/journal"
//
// func ij(target string, foreground bool) (can_use_log_file, journal_mode bool) {
// 	// Only log the warning severity or above.
// 	can_use_log_file = true
// 	journal_mode = journal.Enabled() && target == "journal"
// 	// daemon mode 才会发送给 journal
// 	if !journal_mode || foreground == true {
// 		can_use_log_file = false
// 		// sink = NewJournalSink()
// 	} else {
// 		// sink = NewDefaultSink()
// 	}
// 	return
// }
