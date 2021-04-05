/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "github.com/hedzr/cmdr/tool"

func (w *ExecWorker) helpOptMatching(pkg *ptpkg, goCommand **Command, args []string) (matched, stop bool, err error) {
	// pkg.needHelp = true
	// pkg.needFlagsHelp = true
	ra := (args)[pkg.i:]
	if len(ra) > 0 {
		ra = ra[1:]
	}

	stop = true
	pkg.lastCommandHeld = false
	pkg.needHelp = false
	pkg.needFlagsHelp = false
	if w.onSwitchCharHit != nil {
		err = w.onSwitchCharHit(*goCommand, pkg.a, ra)
	} else {
		err = defaultOnSwitchCharHit(*goCommand, pkg.a, ra)
	}
	return
}

func (w *ExecWorker) cmdMatching(pkg *ptpkg, goCommand **Command, args []string) (matched, stop bool, err error) {
	// command, files
	if cmd, ok := (*goCommand).plainCmds[pkg.a]; ok {
		cmd.strHit = pkg.a
		*goCommand = cmd
		matched = true
		flog("    -> command %q hit (a=%q, idx=%v)...", cmd.GetTitleName(), pkg.a, pkg.i)
		stop, err = w.cmdMatched(pkg, *goCommand, args)
		return
	}

	if len((*goCommand).SubCommands) == 0 { // (*goCommand).Action != nil &&
		// the args remained are files, not sub-commands.
		pkg.i--
		pkg.lastCommandHeld = true
		pkg.iLastCommand = pkg.i
		return
	}

	if w.treatUnknownCommandAsArgs {
		stop = true
		return
	}

	flog("    . adding unknown command %q", pkg.a)
	pkg.unknownCmds = append(pkg.unknownCmds, pkg.a)
	unknownCommand(pkg, *goCommand, args)
	return
}

func (w *ExecWorker) cmdMatched(pkg *ptpkg, goCommand *Command, args []string) (stop bool, err error) {
	pkg.iLastCommand = pkg.i

	if len((*goCommand).SubCommands) == 0 { // (*goCommand).Action != nil &&
		// the args remained are files, not sub-commands.
		pkg.lastCommandHeld = true
		stop = true
	}

	return
}

func (w *ExecWorker) flagsPrepare(pkg *ptpkg, goCommand **Command, args []string) (stop bool, err error) {
	if len(pkg.a) > 1 && (pkg.a[1] == '-' || pkg.a[1] == '~') {
		if len(pkg.a) == 2 {
			// disableParser = true // '--': ignore the following args // PassThrough hit!
			stop = true
			pkg.lastCommandHeld = false
			pkg.needHelp = false
			pkg.needFlagsHelp = false
			ra := args[pkg.i:]
			if len(ra) > 0 {
				ra = ra[1:]
			}
			if w.onPassThruCharHit != nil {
				err = w.onPassThruCharHit(*goCommand, pkg.a, ra)
			} else {
				err = defaultOnPassThruCharHit(*goCommand, pkg.a, ra)
			}
			return
		}

		// long flag
		pkg.fn = pkg.a[2:]
		pkg.findValueAttached(&pkg.fn)

	} else {

		// short flag
		if (*goCommand).headLikeFlag != nil && tool.IsDigitHeavy(pkg.a[1:]) {
			// println("head-like")
			pkg.short = true
			pkg.flg = (*goCommand).headLikeFlag
			pkg.val = pkg.a[1:]
			pkg.fn = pkg.flg.Short
			pkg.found = true
			err = pkg.processTypeIntCore(args)
			return
		}

		pkg.suffix = pkg.a[len(pkg.a)-1]
		if pkg.suffix == '+' || pkg.suffix == '-' {
			pkg.a = pkg.a[0 : len(pkg.a)-1]
		} else {
			pkg.suffix = 0
		}

		if i := pkg.matchShortFlag(*goCommand, pkg.a); i >= 0 {
			pkg.fn = pkg.a[1:i]
			pkg.savedFn = pkg.a[i:]
		} else {
			pkg.fn = pkg.a[1:2]     // from one char
			pkg.savedFn = pkg.a[2:] // save others
		}
		pkg.short = true
		pkg.findValueAttached(&pkg.savedFn)
	}
	return
}

func (w *ExecWorker) flagsMatching(pkg *ptpkg, cc *Command, goCommand **Command, args []string) (matched, stop bool, err error) {
	var upLevel bool
goUp:
	pkg.found = false
	if pkg.short {
		a := "-" + pkg.fn + pkg.savedFn
		flog("    .  . matching short flag for %q", a)
		if i := pkg.matchShortFlag(cc, a); i >= 0 {
			pkg.fn = a[1:i]
			pkg.savedFn = a[i:]
			pkg.flg, matched = cc.plainShortFlags[pkg.fn]
		}
	} else {
		flog("    .  . matching long flag for --%v", pkg.fn)
		matched = w.matchForLongFlags(pkg.fn, cc, pkg)
	}

	if matched {
		if upLevel, stop, err = w.flagsMatched(pkg, *goCommand, args); stop || err != nil {
			return
		}
		if upLevel {
			goto goUp
		}
	} else {
		if cc.owner != nil {
			// match the flag within parent's flags set.
			cc = cc.owner
			goto goUp
		}
		if !pkg.assigned && pkg.short {
			// try matching 2-chars short opt
			if len(pkg.savedFn) > 0 {
				fnf := pkg.fn + pkg.savedFn
				pkg.fn = fnf[0:2]
				pkg.savedFn = fnf[2:]
				*goCommand = pkg.savedGoCommand
				if (*goCommand).owner != nil {
					goto goUp
				}
			}
		}

		pkg.unknownFlags = append(pkg.unknownFlags, pkg.a)
		unknownFlag(pkg, *goCommand, args)
	}
	return
}

func (w *ExecWorker) flagsMatched(pkg *ptpkg, goCommand *Command, args []string) (upLevel, stop bool, err error) {
	pkg.flg.times++

	if err = pkg.tryExtractingValue(args); err != nil {
		stop = true
		return
	}

	if pkg.found {
		// if !GetBoolP(getPrefix(), "quiet") {
		// 	logrus.Debugf("-- flag '%v' hit, go ahead...", pkg.flg.GetTitleName())
		// }
		if pkg.flg.Action != nil {
			if err = pkg.flg.Action(goCommand, w.getRemainArgs(pkg, args)); err == ErrShouldBeStopException {
				stop = true
				err = nil
				return
			}
		}
		if isBool(pkg.flg.DefaultValue) || isNil1(pkg.flg.DefaultValue) {
			flog("    .  . [tryToggleGroup] %q = %v", pkg.fn, pkg.val)
			pkg.tryToggleGroup()
		}

		if !pkg.assigned {
			if len(pkg.savedFn) > 0 && len(pkg.savedVal) == 0 {
				pkg.fn = pkg.savedFn[0:1]
				pkg.savedFn = pkg.savedFn[1:]
				// goto GO_UP
				upLevel = true
			}
		} else {
			flog("    .  . [value assigned] %q = %v", pkg.fn, pkg.val)
		}
	}
	return
}

func (w *ExecWorker) matchForLongFlags(fn string, cc *Command, pkg *ptpkg) (ok bool) {
	var ln = len(fn)
	for ; ln > 1; ln-- {
		fn = pkg.fn[0:ln]
		pkg.flg, ok = cc.plainLongFlags[fn]
		if ok {
			if ln < len(pkg.fn) {
				pkg.val = pkg.fn[ln:]
				pkg.fn = fn
				pkg.assigned = true
			}
			break
		}
	}
	return
}
