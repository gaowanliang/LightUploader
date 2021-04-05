/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "github.com/hedzr/cmdr/tool"

type (
	// UnknownOptionHandler for WithSimilarThreshold/SetUnknownOptionHandler
	UnknownOptionHandler func(isFlag bool, title string, cmd *Command, args []string) (fallbackToDefaultDetector bool)
)

var (
	unknownOptionHandler UnknownOptionHandler
)

// // SetUnknownOptionHandler enables your customized wrong command/flag processor.
// // internal processor supports smart suggestions for those wrong commands and flags.
// //
// // Deprecated: from v1.5.5, replaced with WithUnknownOptionHandler
// func SetUnknownOptionHandler(handler UnknownOptionHandler) {
// 	unknownOptionHandler = handler
// }

func unknownCommand(pkg *ptpkg, cmd *Command, args []string) {
	if internalGetWorker().noUnknownCmdTip {
		return
	}

	ferr("\n\x1b[%dmUnknown command:\x1b[0m %v", BgBoldOrBright, pkg.a)
	if unknownOptionHandler != nil {
		if !unknownOptionHandler(false, pkg.a, cmd, args) {
			return
		}
	}
	unknownCommandDetector(pkg, cmd, args)
}

func unknownFlag(pkg *ptpkg, cmd *Command, args []string) {
	if internalGetWorker().noUnknownCmdTip {
		return
	}

	ferr("\n\x1b[%dmUnknown flag:\x1b[0m %v", BgBoldOrBright, pkg.a)
	if unknownOptionHandler != nil && !pkg.short {
		if !unknownOptionHandler(true, pkg.a, cmd, args) {
			return
		}
	}
	unknownFlagDetector(pkg, cmd, args)
}

func unknownCommandDetector(pkg *ptpkg, cmd *Command, args []string) {
	ever := false
	for k := range cmd.plainCmds {
		distance := float64(defaultStringMetric.Calc(pkg.a, k)) / tool.StringMetricFactor
		if distance >= internalGetWorker().similarThreshold {
			ferr("  - do you mean: %v", k)
			ever = true
		}
	}

	// sndSrc := soundex(pkg.a)
	// ever := false
	// for k := range cmd.plainCmds {
	// 	snd := soundex(k)
	// 	if sndSrc == snd {
	// 		ferr("  - do you mean: %v", k)
	// 		ever = true
	// 		// } else {
	// 		// 	ferr("  . %v -> %v: --%v -> %v", pkg.a, sndSrc, k, snd)
	// 	}
	// }

	if !ever && cmd.HasParent() {
		unknownCommandDetector(pkg, cmd.GetOwner(), args)
	}
}

func unknownFlagDetector(pkg *ptpkg, cmd *Command, args []string) {
	if !pkg.short {
		ever := false
		str := tool.StripPrefix(pkg.a, "--")
		for k := range cmd.plainLongFlags {
			distance := float64(defaultStringMetric.Calc(str, k)) / tool.StringMetricFactor
			if distance >= internalGetWorker().similarThreshold {
				ferr("  - do you mean: --%v", k)
				ever = true
				// } else {
				// 	ferr("  ? '%v' - '%v': %v", pkg.a, k, distance)
			}
		}
		if !ever && cmd.HasParent() {
			unknownFlagDetector(pkg, cmd.GetOwner(), args)
		}
	}

	// sndSrc := soundex(pkg.a)
	// if !pkg.short {
	// 	ever := false
	// 	for k := range cmd.plainLongFlags {
	// 		snd := soundex(k)
	// 		if sndSrc == snd {
	// 			ferr("  - do you mean: --%v", k)
	// 			ever = true
	// 			// } else {
	// 			// 	ferr("  . %v -> %v: --%v -> %v", pkg.a, sndSrc, k, snd)
	// 		}
	// 	}
	// 	if !ever && cmd.HasParent() {
	// 		unknownFlagDetector(pkg, cmd.GetOwner(), args)
	// 	}
	// }
}
