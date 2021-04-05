/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"github.com/hedzr/cmdr/tool"
	"regexp"
	"strings"
)

type optFlagImpl struct {
	working *Flag
	parent  OptCmd
}

func (s *optFlagImpl) ToFlag() *Flag {
	return s.working
}

func (s *optFlagImpl) AttachTo(opt OptCmd) {
	if s != nil && s.working != nil && opt != nil {
		opt.AddOptFlag(s)
	}
}

func (s *optFlagImpl) AttachToCommand(cmd *Command) {
	if cmd != nil {
		cmd.Flags = uniAddFlg(cmd.Flags, s.ToFlag())
	}
}

func (s *optFlagImpl) AttachToRoot(root *RootCommand) {
	if root != nil {
		root.Command.Flags = uniAddFlg(root.Command.Flags, s.ToFlag())
	}
}

func (s *optFlagImpl) Titles(long, short string, aliases ...string) (opt OptFlag) {
	s.working.Short = short
	s.working.Full = long
	if tool.HasOrderPrefix(long) {
		s.working.Full = tool.StripOrderPrefix(long)
		s.working.Name = long
	}
	s.working.Aliases = uniAddStrs(s.working.Aliases, aliases...)
	opt = s
	return
}

func (s *optFlagImpl) Short(short string) (opt OptFlag) {
	s.working.Short = short
	opt = s
	return
}

func (s *optFlagImpl) Long(long string) (opt OptFlag) {
	s.working.Full = long
	if tool.HasOrderPrefix(long) {
		s.working.Full = tool.StripOrderPrefix(long)
		s.working.Name = long
	}
	opt = s
	return
}

func (s *optFlagImpl) Name(name string) (opt OptFlag) {
	s.working.Name = name
	opt = s
	return
}

func (s *optFlagImpl) Aliases(aliases ...string) (opt OptFlag) {
	s.working.Aliases = uniAddStrs(s.working.Aliases, aliases...)
	opt = s
	return
}

func (s *optFlagImpl) Description(oneLineDesc string, longDesc ...string) (opt OptFlag) {
	s.working.Description = oneLineDesc

	for _, long := range longDesc {
		s.working.LongDescription = long

		if len(s.working.Description) == 0 {
			s.working.Description = long
		}
	}

	if b := regexp.MustCompile("`(.+)`").Find([]byte(s.working.Description)); len(b) > 2 {
		ph := strings.ToUpper(strings.Trim(string(b), "`"))
		s.Placeholder(ph)
	}

	opt = s
	return
}

func (s *optFlagImpl) Examples(examples string) (opt OptFlag) {
	s.working.Examples = examples
	opt = s
	return
}

func (s *optFlagImpl) Group(group string) (opt OptFlag) {
	s.working.Group = group
	opt = s
	return
}

func (s *optFlagImpl) Hidden(hidden bool) (opt OptFlag) {
	s.working.Hidden = hidden
	opt = s
	return
}

func (s *optFlagImpl) Deprecated(deprecation string) (opt OptFlag) {
	s.working.Deprecated = deprecation
	opt = s
	return
}

func (s *optFlagImpl) Action(action Handler) (opt OptFlag) {
	s.working.Action = action
	opt = s
	return
}

func (s *optFlagImpl) ToggleGroup(group string) (opt OptFlag) {
	s.working.ToggleGroup = group
	opt = s
	return
}

func (s *optFlagImpl) DefaultValue(val interface{}, placeholder string) (opt OptFlag) {
	s.working.DefaultValue = val
	s.working.DefaultValuePlaceholder = placeholder
	opt = s
	return
}

// Placeholder to specify the text string that will be appended
// to the end of a flag expr; it is used into help screen.
// For example, `Placeholder("PASSWORD")` will take the form like:
//   -p, --password=PASSWORD, --pwd, --passwd    to input password
func (s *optFlagImpl) Placeholder(placeholder string) (opt OptFlag) {
	s.working.DefaultValuePlaceholder = placeholder
	opt = s
	return
}

func (s *optFlagImpl) ExternalTool(envKeyName string) (opt OptFlag) {
	s.working.ExternalTool = envKeyName
	opt = s
	return
}

func (s *optFlagImpl) ValidArgs(list ...string) (opt OptFlag) {
	s.working.ValidArgs = list
	opt = s
	return
}

func (s *optFlagImpl) HeadLike(enable bool, min, max int64) (opt OptFlag) {
	s.working.HeadLike = enable
	s.working.Min, s.working.Max = min, max
	opt = s
	return
}

func (s *optFlagImpl) EnvKeys(keys ...string) (opt OptFlag) {
	s.working.EnvVars = uniAddStrs(s.working.EnvVars, keys...)
	opt = s
	return
}

func (s *optFlagImpl) Required(required ...bool) (opt OptFlag) {
	var b bool = true
	for _, bb := range required {
		b = bb
	}
	s.working.Required = b
	opt = s
	return
}

func (s *optFlagImpl) OnSet(f func(keyPath string, value interface{})) (opt OptFlag) {
	s.working.onSet = f
	opt = s
	return
}

func (s *optFlagImpl) SetOwner(opt OptCmd) {
	s.parent = opt
	return
}

func (s *optFlagImpl) OwnerCommand() (opt OptCmd) {
	opt = s.parent
	return
}

func (s *optFlagImpl) RootCommand() (root *RootCommand) {
	root = optCtx.root
	return
}
