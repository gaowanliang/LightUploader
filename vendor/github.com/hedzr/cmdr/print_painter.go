/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bytes"
)

type (
	// Painter to support the genManual, genMarkdown, printHelpScreen.
	Painter interface {
		Printf(fmtStr string, args ...interface{})
		Print(fmtStr string, args ...interface{})

		FpPrintHeader(command *Command)
		FpPrintHelpTailLine(command *Command)

		FpUsagesTitle(command *Command, title string)
		FpUsagesLine(command *Command, fmt, appName, cmdList, cmdsTitle, tailPlaceHolder string)
		FpDescTitle(command *Command, title string)
		FpDescLine(command *Command)
		FpExamplesTitle(command *Command, title string)
		FpExamplesLine(command *Command)

		FpCommandsTitle(command *Command)
		FpCommandsGroupTitle(group string)
		FpCommandsLine(command *Command) (bufL, bufR bytes.Buffer)
		FpFlagsTitle(command *Command, flag *Flag, title string)
		FpFlagsGroupTitle(group string)
		FpFlagsLine(command *Command, flag *Flag, maxShort int, defValStr string) (bufL, bufR bytes.Buffer)

		Flush()

		Results() []byte

		// clear any internal states and reset itself
		Reset()
	}
)
