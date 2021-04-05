/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"fmt"
	"strings"
)

// AppendPostActions adds the global post-action to cmdr system
func (c *RootCommand) AppendPostActions(fns ...func(cmd *Command, args []string)) {
	for _, fn := range fns {
		c.PostActions = append(c.PostActions, fn)
	}
}

// PrintHelp prints help screen
func (c *Command) PrintHelp(justFlags bool) {
	internalGetWorker().printHelp(c, justFlags)
}

// PrintVersion prints versions information
func (c *Command) PrintVersion() {
	internalGetWorker().showVersion()
}

// PrintBuildInfo print building information
func (c *Command) PrintBuildInfo() {
	internalGetWorker().showBuildInfo()
}

// GetRoot returns the `RootCommand`
func (c *Command) GetRoot() *RootCommand {
	return c.root
}

// GetOwner returns the parent command object
func (c *Command) GetOwner() *Command {
	return c.owner
}

// IsRoot returns true if this command is a RootCommand
func (c *Command) IsRoot() bool {
	return c == &c.root.Command
}

// GetHitStr returns the matched command string
func (c *Command) GetHitStr() string {
	return c.strHit
}

// FindSubCommand find sub-command with `longName` from `cmd`
func (c *Command) FindSubCommand(longName string) (res *Command) {
	// return FindSubCommand(longName, c)
	for _, cx := range c.SubCommands {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	return
}

// FindSubCommandRecursive find sub-command with `longName` from `cmd` recursively
func (c *Command) FindSubCommandRecursive(longName string) (res *Command) {
	// return FindSubCommandRecursive(longName, c)
	for _, cx := range c.SubCommands {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	for _, cx := range c.SubCommands {
		if len(cx.SubCommands) > 0 {
			if res = cx.FindSubCommandRecursive(longName); res != nil {
				return
			}
		}
	}
	return
}

// FindFlag find flag with `longName` from `cmd`
func (c *Command) FindFlag(longName string) (res *Flag) {
	// return FindFlag(longName, c)
	for _, cx := range c.Flags {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	return
}

// FindFlagRecursive find flag with `longName` from `cmd` recursively
func (c *Command) FindFlagRecursive(longName string) (res *Flag) {
	// return FindFlagRecursive(longName, c)
	for _, cx := range c.Flags {
		if longName == cx.Full {
			res = cx
			return
		}
	}
	for _, cx := range c.SubCommands {
		// if len(cx.SubCommands) > 0 {
		if res = cx.FindFlagRecursive(longName); res != nil {
			return
		}
		// }
	}
	return
}

// // HasParent detects whether owner is available or not
// func (c *BaseOpt) HasParent() bool {
// 	return c.owner != nil
// }

// GetName returns the name of a `Command`.
func (c *Command) GetName() string {
	if len(c.Name) > 0 {
		return c.Name
	}
	if len(c.Full) > 0 {
		return c.Full
	}
	panic("The `Full` or `Name` must be non-empty for a command or flag")
}

// GetDottedNamePath return the dotted key path of this command
// in the options store.
// For example, the returned string just like: 'server.start'.
// NOTE that there is no OptiontPrefixes in this key path. For
// more information about Option Prefix, refer
// to [WithOptionsPrefix]
func (c *Command) GetDottedNamePath() string {
	return internalGetWorker().backtraceCmdNames(c, false)
}

// GetQuotedGroupName returns the group name quoted string.
func (c *Command) GetQuotedGroupName() string {
	if len(strings.TrimSpace(c.Group)) == 0 {
		return ""
	}
	i := strings.Index(c.Group, ".")
	if i >= 0 {
		return fmt.Sprintf("[%v]", c.Group[i+1:])
	}
	return fmt.Sprintf("[%v]", c.Group)
}

// GetExpandableNamesArray returns the names array of command, includes short name and long name.
func (c *Command) GetExpandableNamesArray() []string {
	var a []string
	if len(c.Full) > 0 {
		a = append(a, c.Full)
	}
	if len(c.Short) > 0 {
		a = append(a, c.Short)
	}
	return a
}

// GetExpandableNames returns the names comma splitted string.
func (c *Command) GetExpandableNames() string {
	a := c.GetExpandableNamesArray()
	if len(a) == 1 {
		return a[0]
	} else if len(a) > 1 {
		return fmt.Sprintf("{%v}", strings.Join(a, ","))
	}
	return c.Name
}

// GetParentName returns the owner command name
func (c *Command) GetParentName() string {
	if c.owner != nil {
		//return c.owner.GetName()
		if len(c.owner.Name) > 0 {
			return c.owner.Name
		}
		if len(c.owner.Full) > 0 {
			return c.owner.Full
		}
		// panic("The `Full` or `Name` must be non-empty for a command or flag")
	}
	return c.GetRoot().AppName
}

// GetSubCommandNamesBy returns the joint string of subcommands
func (c *Command) GetSubCommandNamesBy(delimChar string) string {
	var a []string
	for _, sc := range c.SubCommands {
		if !sc.Hidden {
			a = append(a, sc.GetTitleNamesBy(delimChar))
		}
	}
	return strings.Join(a, delimChar)
}
