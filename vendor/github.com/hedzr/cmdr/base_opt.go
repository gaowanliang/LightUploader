/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"strings"
)

// HasParent detects whether owner is available or not
func (s *BaseOpt) HasParent() bool {
	return s.owner != nil
}

// GetTitleName returns name/full/short string
func (s *BaseOpt) GetTitleName() string {
	if len(s.Name) != 0 {
		return s.Name
	}
	if len(s.Full) > 0 {
		return s.Full
	}
	if len(s.Short) > 0 {
		return s.Short
	}
	// for _, ss := range s.Aliases {
	// 	return ss
	// }
	return ""
}

// GetTitleNamesArray returns short,full,aliases names
func (s *BaseOpt) GetTitleNamesArray() []string {
	var a []string
	if len(s.Short) != 0 {
		a = uniAddStr(a, s.Short)
	}
	if len(s.Full) > 0 {
		a = uniAddStr(a, s.Full)
	}
	a = uniAddStrs(a, s.Aliases...)
	return a
}

// GetShortTitleNamesArray returns short name as an array
func (s *BaseOpt) GetShortTitleNamesArray() []string {
	var a []string
	if len(s.Short) != 0 {
		a = uniAddStr(a, s.Short)
	}
	return a
}

// GetLongTitleNamesArray returns long name and aliases as an array
func (s *BaseOpt) GetLongTitleNamesArray() []string {
	var a []string
	if len(s.Full) > 0 {
		a = uniAddStr(a, s.Full)
	}
	a = uniAddStrs(a, s.Aliases...)
	return a
}

// GetTitleNames return the joint string of short,full,aliases names
func (s *BaseOpt) GetTitleNames() string {
	return s.GetTitleNamesBy(", ")
}

// GetTitleNamesBy returns the joint string of short,full,aliases names
func (s *BaseOpt) GetTitleNamesBy(delimChar string) string {
	var a = s.GetTitleNamesArray()
	str := strings.Join(a, delimChar)
	return str
}
