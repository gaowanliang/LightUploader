/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"github.com/hedzr/cmdr/tool"
	"strings"
)

// GetTriggeredTimes returns the matched times
func (s *Flag) GetTriggeredTimes() int {
	return s.times
}

// GetTitleFlagNames temp
func (s *Flag) GetTitleFlagNames() string {
	return s.GetTitleFlagNamesBy(",")
}

// GetDescZsh temp
func (s *Flag) GetDescZsh() (desc string) {
	desc = s.Description
	if len(desc) == 0 {
		desc = tool.EraseAnyWSs(s.GetTitleZshFlagName())
	}
	// desc = replaceAll(desc, " ", "\\ ")
	return
}

// GetTitleZshFlagName temp
func (s *Flag) GetTitleZshFlagName() (str string) {
	if len(s.Full) > 0 {
		str += "--" + s.Full
	} else if len(s.Short) == 1 {
		str += "-" + s.Short
	}
	return
}

// GetTitleZshFlagNames temp
func (s *Flag) GetTitleZshFlagNames(delimChar string) (str string) {
	if len(s.Short) == 1 {
		str += "-" + s.Short + delimChar
	}
	if len(s.Full) > 0 {
		str += "--" + s.Full
	}
	return
}

// GetTitleZshFlagNamesArray temp
func (s *Flag) GetTitleZshFlagNamesArray() (ary []string) {
	if len(s.Short) == 1 || len(s.Short) == 2 {
		if len(s.DefaultValuePlaceholder) > 0 {
			ary = append(ary, "-"+s.Short+"=") // +s.DefaultValuePlaceholder)
		} else {
			ary = append(ary, "-"+s.Short)
		}
	}
	if len(s.Full) > 0 {
		if len(s.DefaultValuePlaceholder) > 0 {
			ary = append(ary, "--"+s.Full+"=") // +s.DefaultValuePlaceholder)
		} else {
			ary = append(ary, "--"+s.Full)
		}
	}
	return
}

// GetTitleFlagNamesBy temp
func (s *Flag) GetTitleFlagNamesBy(delimChar string) string {
	return s.GetTitleFlagNamesByMax(delimChar, len(s.Short))
}

// GetTitleFlagNamesByMax temp
func (s *Flag) GetTitleFlagNamesByMax(delimChar string, maxShort int) string {
	var sb strings.Builder

	if len(s.Short) == 0 {
		// if no flag.Short,
		sb.WriteString(strings.Repeat(" ", maxShort))
	} else {
		sb.WriteRune('-')
		sb.WriteString(s.Short)
		sb.WriteString(delimChar)
		if len(s.Short) < maxShort {
			sb.WriteString(strings.Repeat(" ", maxShort-len(s.Short)))
		}
	}

	if len(s.Short) == 0 {
		sb.WriteRune(' ')
		sb.WriteRune(' ')
	}
	sb.WriteRune(' ')
	sb.WriteString("--")
	sb.WriteString(s.Full)
	if len(s.DefaultValuePlaceholder) > 0 {
		// str += fmt.Sprintf("=\x1b[2m\x1b[%dm%s\x1b[0m", DarkColor, s.DefaultValuePlaceholder)
		sb.WriteString(fmt.Sprintf("=%s", s.DefaultValuePlaceholder))
	}

	for _, sz := range s.Aliases {
		sb.WriteString(delimChar)
		sb.WriteString("--")
		sb.WriteString(sz)
	}
	return sb.String()
}
