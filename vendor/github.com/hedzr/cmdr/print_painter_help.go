/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bytes"
	"fmt"
	"github.com/hedzr/cmdr/tool"
	"strconv"
	"strings"
)

type (
	helpPainter struct {
	}
)

func (s *helpPainter) Reset() {
	if w := internalGetWorker().rootCommand.ow; w != nil {
		_ = w.Flush()
	}
}

func (s *helpPainter) Flush() {
	if w := internalGetWorker().rootCommand.ow; w != nil {
		_ = w.Flush()
	}
}

func (s *helpPainter) Results() (res []byte) {
	return
}

func (s *helpPainter) bufPrintf(sb *bytes.Buffer, fmtStr string, args ...interface{}) {
	_, _ = sb.WriteString(fmt.Sprintf(fmtStr, args...))
}

func (s *helpPainter) Printf(fmtStr string, args ...interface{}) {
	fp(fmtStr, args...)
}

func (s *helpPainter) Print(fmtStr string, args ...interface{}) {
	fp0(fmtStr, args...)
}

func (s *helpPainter) FpPrintHeader(command *Command) {
	if len(command.root.Header) == 0 {
		s.Printf("%v by %v - v%v", command.root.Copyright, command.root.Author, command.root.Version)
	} else {
		s.Printf("%v", command.root.Header)
	}
}

func (s *helpPainter) FpPrintHelpTailLine(command *Command) {
	if internalGetWorker().enableHelpCommands {
		if GetNoColorMode() {
			s.Printf(fmtTailLineNC, internalGetWorker().helpTailLine)
		} else {
			s.Printf(fmtTailLine, CurrentGroupTitleColor, internalGetWorker().helpTailLine)
		}
	}
}

func (s *helpPainter) FpUsagesTitle(command *Command, title string) {
	s.Printf("\n%s:", title)
	// s.Printf("\n\x1b[%dm\x1b[%dm%s\x1b[0m", BgNormal, DarkColor, title)
	// fp("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", BgDim, DarkColor, StripOrderPrefix(group))
}

func (s *helpPainter) FpUsagesLine(command *Command, fmt, appName, cmdList, cmdsTitle, tailPlaceHolder string) {
	if strings.HasPrefix(cmdList, appName) {
		appName = ""
	} else {
		cmdList = " " + cmdList
	}
	if len(tailPlaceHolder) > 0 {
		tailPlaceHolder = command.TailPlaceHolder
	} else {
		tailPlaceHolder = "[tail args...]"
	}
	s.Printf("    %s%v%s%s [Options] [Parent/Global Options]"+fmt, appName, cmdList, cmdsTitle, tailPlaceHolder)
}

func (s *helpPainter) FpDescTitle(command *Command, title string) {
	s.Printf("\n%s:", title)
}

func (s *helpPainter) FpDescLine(command *Command) {
	s.Printf("    %v", command.Description)
}

func (s *helpPainter) FpExamplesTitle(command *Command, title string) {
	s.Printf("\n%s:", title)
}

func (s *helpPainter) FpExamplesLine(command *Command) {
	str := tplApply(command.Examples, command.root)
	for _, line := range strings.Split(str, "\n") {
		s.Printf("    %v", line)
	}
}

func (s *helpPainter) FpCommandsTitle(command *Command) {
	var title string
	if command.owner == nil {
		title = "Commands"
	} else {
		title = "Sub-Commands"
	}
	s.Printf("\n%s:", title)
}

func (s *helpPainter) FpCommandsGroupTitle(group string) {
	if group != UnsortedGroup {
		if GetNoColorMode() {
			s.Printf(fmtCmdGroupTitleNC, tool.StripOrderPrefix(group))
		} else {
			s.Printf(fmtCmdGroupTitle, CurrentGroupTitleColor, tool.StripOrderPrefix(group))
		}
	}
}

func (s *helpPainter) FpCommandsLine(command *Command) (bufL, bufR bytes.Buffer) {
	if !command.Hidden {
		if len(command.Deprecated) > 0 {
			if GetNoColorMode() {
				s.bufPrintf(&bufL, fmtCmdlineDepNCL, command.GetTitleNames())
				s.bufPrintf(&bufR, fmtCmdlineDepNCR, command.Description, command.Deprecated)
			} else {
				s.bufPrintf(&bufL, fmtCmdlineDepL, BgNormal, CurrentDescColor, command.GetTitleNames())
				s.bufPrintf(&bufR, fmtCmdlineDepR, command.Description, command.Deprecated)
			}
		} else {
			if GetNoColorMode() {
				s.bufPrintf(&bufL, fmtCmdlineNCL, command.GetTitleNames())
				s.bufPrintf(&bufR, fmtCmdlineNCR, command.Description)
			} else {
				// s.Printf("  %-48s%v", command.GetTitleNames(), command.Description)
				// s.Printf("\n\x1b[%dm\x1b[%dm%s\x1b[0m", BgNormal, DarkColor, title)
				// s.Printf("  [\x1b[%dm\x1b[%dm%s\x1b[0m]", BgDim, DarkColor, StripOrderPrefix(group))
				s.bufPrintf(&bufL, fmtCmdlineL, command.GetTitleNames())
				s.bufPrintf(&bufR, fmtCmdlineR, BgNormal, CurrentDescColor, command.Description)
			}
		}
	}
	return
}

// func (s *helpPainter) FpFlagsSssTitle(flag *Flag) {
// 	var title string
// 	if flag.owner == nil {
// 		title = "Commands"
// 	} else {
// 		title = "Sub-Commands"
// 	}
// 	s.Printf("\n%s:", title)
// }

func (s *helpPainter) FpFlagsTitle(command *Command, flag *Flag, title string) {
	s.Printf("\n%s:", title)
}

func (s *helpPainter) FpFlagsGroupTitle(group string) {
	if group != UnsortedGroup {
		if GetNoColorMode() {
			s.Printf(fmtGroupTitleNC, tool.StripOrderPrefix(group))
		} else {
			// fp("  [%s]:", StripOrderPrefix(group))
			// // echo -e "Normal \e[2mDim"
			// _, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m\x1b[2m\x1b[%dm[%04d]\x1b[0m%-48s \x1b[2m\x1b[%dm%s\x1b[0m ",
			// 	levelColor, levelText, DarkColor, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message, DarkColor, caller)
			s.Printf(fmtGroupTitle, CurrentGroupTitleColor, tool.StripOrderPrefix(group))
		}
	}
}

func (s *helpPainter) FpFlagsLine(command *Command, flg *Flag, maxShort int, defValStr string) (bufL, bufR bytes.Buffer) {
	if len(flg.ValidArgs) > 0 {
		defValStr = fmt.Sprintf("%v, in %v", defValStr, flg.ValidArgs)
	}
	if flg.Min >= 0 && flg.Max > 0 {
		defValStr = fmt.Sprintf("%v, in [%v..%v]", defValStr, flg.Min, flg.Max)
	}

	var envKeys string
	if len(flg.EnvVars) > 0 {
		var sb strings.Builder
		for _, k := range flg.EnvVars {
			if len(strings.TrimSpace(k)) > 0 {
				sb.WriteString(strings.TrimSpace(k))
				sb.WriteRune(',')
			}
		}
		if sb.Len() > 0 {
			envKeys = fmt.Sprintf(" [env: %v]", strings.TrimRight(sb.String(), ","))
		}
	}

	if len(flg.Deprecated) > 0 {
		if GetNoColorMode() {
			s.bufPrintf(&bufL, fmtFlagsDepNCL, // "  %-48s%s%s [deprecated since %v]",
				flg.GetTitleFlagNamesByMax(",", maxShort))
			s.bufPrintf(&bufR, fmtFlagsDepNCR, // "  %-48s%s%s [deprecated since %v]",
				flg.Description, envKeys, defValStr, flg.Deprecated)
		} else {
			s.bufPrintf(&bufL, fmtFlagsDepL, // "  \x1b[%dm\x1b[%dm%-48s%s\x1b[%dm\x1b[%dm%s\x1b[0m [deprecated since %v]",
				BgNormal, CurrentDescColor, flg.GetTitleFlagNamesByMax(",", maxShort))
			s.bufPrintf(&bufR, fmtFlagsDepR, // "  \x1b[%dm\x1b[%dm%-48s%s\x1b[%dm\x1b[%dm%s\x1b[0m [deprecated since %v]",
				flg.Description, BgItalic, CurrentDefaultValueColor, envKeys, defValStr, flg.Deprecated)
		}
	} else {
		if GetNoColorMode() {
			s.bufPrintf(&bufL, fmtFlagsNCL, flg.GetTitleFlagNamesByMax(",", maxShort))
			s.bufPrintf(&bufR, fmtFlagsNCR, flg.Description, envKeys, defValStr)
		} else {
			s.bufPrintf(&bufL, fmtFlagsL, // "  %-48s\x1b[%dm\x1b[%dm%s\x1b[%dm\x1b[%dm%s\x1b[0m",
				flg.GetTitleFlagNamesByMax(",", maxShort))
			s.bufPrintf(&bufR, fmtFlagsR, BgNormal, CurrentDescColor, flg.Description,
				BgItalic, CurrentDefaultValueColor, envKeys, defValStr)
		}
	}
	return
}

func initTabStop(ts int) {
	// defaultTabStop = ts
	defaultTabStop = ts

	var s = strconv.Itoa(defaultTabStop)

	fmtCmdGroupTitle = "  [\x1b[2m\x1b[%dm%s\x1b[0m]"
	fmtCmdGroupTitleNC = "  [%s]"

	fmtCmdlineL = "  %-" + s + "s"
	fmtCmdlineR = "\x1b[%dm\x1b[%dm%s\x1b[0m"
	fmtCmdlineDepL = "  \x1b[%dm\x1b[%dm%-" + s + "s"
	fmtCmdlineDepR = "%s\x1b[0m [deprecated since %v]"
	fmtCmdlineNCL = "  %-" + s + "s"
	fmtCmdlineNCR = "%s"
	fmtCmdlineDepNCL = "  %-" + s + "s"
	fmtCmdlineDepNCR = "%s [deprecated since %v]"

	fmtGroupTitle = "  [\x1b[2m\x1b[%dm%s\x1b[0m]"
	fmtGroupTitleNC = "  [%s]"

	fmtFlagsDepL = "  \x1b[%dm\x1b[%dm%-" + s + "s"
	fmtFlagsDepR = "%s\x1b[%dm\x1b[%dm%v%s\x1b[0m [deprecated since %v]"
	fmtFlagsL = "  %-" + s + "s"
	fmtFlagsR = "\x1b[%dm\x1b[%dm%s\x1b[%dm\x1b[%dm%v%s\x1b[0m"
	fmtFlagsNCL = "  %-" + s + "s"
	fmtFlagsNCR = "%s%v%s"
	fmtFlagsDepNCL = "  %-" + s + "s"
	fmtFlagsDepNCR = "%s%v%s [deprecated since %v]"

	fmtTailLine = "\x1b[2m\x1b[%dm%s\x1b[0m"
	fmtTailLineNC = "%s"
}

var (
	defaultTabStop                       = 33
	fmtCmdGroupTitle, fmtCmdGroupTitleNC string
	fmtCmdlineL, fmtCmdlineR             string
	fmtCmdlineDepL, fmtCmdlineDepR       string
	fmtCmdlineNCL, fmtCmdlineNCR         string
	fmtCmdlineDepNCL, fmtCmdlineDepNCR   string
	fmtGroupTitle, fmtGroupTitleNC       string
	fmtFlagsL, fmtFlagsR                 string
	fmtFlagsDepL, fmtFlagsDepR           string
	fmtFlagsNCL, fmtFlagsNCR             string
	fmtFlagsDepNCL, fmtFlagsDepNCR       string
	fmtTailLine, fmtTailLineNC           string
)

const defaultTailLine = `
Type '-h'/'-?' or '--help' to get command help screen. 
More: '-D'/'--debug'['--env'|'--raw'|'--more'], '-V'/'--version', '-#'/'--build-info', '--no-color', '--strict-mode', '--no-env-overrides'...`
