/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"github.com/hedzr/cmdr/tool"
	"reflect"
	"time"
)

type optCommandImpl struct {
	working *Command
	parent  OptCmd
}

func (s *optCommandImpl) ToCommand() *Command {
	return s.working
}

func (s *optCommandImpl) AddOptFlag(flag OptFlag) {
	if flag != nil && s != nil && s.working != nil {
		s.working.Flags = uniAddFlg(s.working.Flags, flag.ToFlag())
	}
}

func (s *optCommandImpl) AddFlag(flag *Flag) {
	if flag != nil {
		s.working.Flags = uniAddFlg(s.working.Flags, flag)
	}
}

func (s *optCommandImpl) AddOptCmd(opt OptCmd) {
	if opt != nil {
		cmd := opt.ToCommand()

		// optCtx.current = cmd

		s.working.SubCommands = uniAddCmd(s.working.SubCommands, cmd)

		// opt = &subCmdOpt{optCommandImpl: optCommandImpl{working: cmd, parent: s}}
	}
}

func (s *optCommandImpl) AddCommand(cmd *Command) {
	if cmd != nil {
		s.working.SubCommands = uniAddCmd(s.working.SubCommands, cmd)
	}
}

func (s *optCommandImpl) AttachTo(opt OptCmd) {
	if opt != nil {
		opt.AddOptCmd(s)
	}
}

func (s *optCommandImpl) AttachToCommand(cmd *Command) {
	if cmd != nil {
		cmd.SubCommands = uniAddCmd(cmd.SubCommands, s.working)
	}
}

func (s *optCommandImpl) AttachToRoot(root *RootCommand) {
	if root != nil {
		root.SubCommands = uniAddCmd(root.SubCommands, s.working)
	}
}

func (s *optCommandImpl) Titles(long, short string, aliases ...string) (opt OptCmd) {
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

func (s *optCommandImpl) Short(short string) (opt OptCmd) {
	s.working.Short = short
	opt = s
	return
}

func (s *optCommandImpl) Long(long string) (opt OptCmd) {
	s.working.Full = long
	if tool.HasOrderPrefix(long) {
		s.working.Full = tool.StripOrderPrefix(long)
		s.working.Name = long
	}
	opt = s
	return
}

func (s *optCommandImpl) Name(name string) (opt OptCmd) {
	s.working.Name = name
	opt = s
	return
}

func (s *optCommandImpl) Aliases(aliases ...string) (opt OptCmd) {
	s.working.Aliases = uniAddStrs(s.working.Aliases, aliases...)
	opt = s
	return
}

func (s *optCommandImpl) Description(oneLine string, long ...string) (opt OptCmd) {
	s.working.Description = oneLine
	for _, l := range long {
		s.working.LongDescription = l
	}
	opt = s
	return
}

func (s *optCommandImpl) Examples(examples string) (opt OptCmd) {
	s.working.Examples = examples
	opt = s
	return
}

func (s *optCommandImpl) Group(group string) (opt OptCmd) {
	s.working.Group = group
	opt = s
	return
}

func (s *optCommandImpl) Hidden(hidden bool) (opt OptCmd) {
	s.working.Hidden = hidden
	opt = s
	return
}

func (s *optCommandImpl) Deprecated(deprecation string) (opt OptCmd) {
	s.working.Deprecated = deprecation
	opt = s
	return
}

func (s *optCommandImpl) Action(action Handler) (opt OptCmd) {
	s.working.Action = action
	opt = s
	return
}

func (s *optCommandImpl) PreAction(pre Handler) (opt OptCmd) {
	// s.workingFlag.ExternalTool = envKeyName
	s.working.PreAction = pre
	opt = s
	return
}

func (s *optCommandImpl) PostAction(post Invoker) (opt OptCmd) {
	// s.workingFlag.ExternalTool = envKeyName
	s.working.PostAction = post
	opt = s
	return
}

func (s *optCommandImpl) TailPlaceholder(placeholder string) (opt OptCmd) {
	// s.workingFlag.ExternalTool = envKeyName
	s.working.TailPlaceHolder = placeholder
	opt = s
	return
}

func (s *optCommandImpl) Bool() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &boolOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) String() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &stringOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) StringSlice() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &stringSliceOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) IntSlice() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &intSliceOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Int() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &intOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Uint() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &uintOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Int64() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &int64Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Uint64() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &uint64Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Float32() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &float32Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Float64() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &float64Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Complex64() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &complex64Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Complex128() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &complex128Opt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) Duration() (opt OptFlag) {
	flg := &Flag{}
	s.working.Flags = uniAddFlg(s.working.Flags, flg)
	return &durationOpt{optFlagImpl: optFlagImpl{working: flg, parent: s}}
}

func (s *optCommandImpl) NewFlag(typ OptFlagType) (opt OptFlag) {
	var flg OptFlag

	switch typ {
	case OptFlagTypeInt:
		flg = s.Int()
	case OptFlagTypeUint:
		flg = s.Uint()
	case OptFlagTypeInt64:
		flg = s.Int64()
	case OptFlagTypeUint64:
		flg = s.Uint64()
	case OptFlagTypeString:
		flg = s.String()
	case OptFlagTypeStringSlice:
		flg = s.StringSlice()
	case OptFlagTypeIntSlice:
		flg = s.IntSlice()
	case OptFlagTypeFloat32:
		flg = s.Float32()
	case OptFlagTypeFloat64:
		flg = s.Float64()
	case OptFlagTypeComplex64:
		flg = s.Complex64()
	case OptFlagTypeComplex128:
		flg = s.Complex128()
	case OptFlagTypeDuration:
		flg = s.Duration()
	default:
		flg = s.Bool()
	}

	flg.SetOwner(s)

	opt = flg
	return
}

func (s *optCommandImpl) newFlagVC(vv reflect.Type, defaultValue interface{}) (flg OptFlag) {
	switch vv.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32:
		flg = s.Int()
	case reflect.Uint, reflect.Uint16, reflect.Uint32:
		flg = s.Uint()
	case reflect.Int64:
		if _, ok := defaultValue.(time.Duration); ok {
			flg = s.Duration()
		} else {
			flg = s.Int64()
		}
	case reflect.Uint64:
		flg = s.Uint64()
	case reflect.String:
		flg = s.String()
	case reflect.Slice:
		flg = s.newFlagVCSlice(vv.Elem(), defaultValue)
	case reflect.Float32:
		flg = s.Float32()
	case reflect.Float64:
		flg = s.Float64()
	case reflect.Complex64:
		flg = s.Complex64()
	case reflect.Complex128:
		flg = s.Complex128()
	default:
		flg = s.Bool()
	}
	return
}

func (s *optCommandImpl) newFlagVCSlice(elt reflect.Type, defaultValue interface{}) (flg OptFlag) {
	switch elt.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		flg = s.IntSlice()
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// flg = s.UintSlice()
		flg = s.IntSlice()
	case reflect.String:
		flg = s.StringSlice()
	}
	return
}

func (s *optCommandImpl) NewFlagV(defaultValue interface{}, titles ...string) (opt OptFlag) {
	var flg OptFlag
	var vv = reflect.TypeOf(defaultValue)
	flg = s.newFlagVC(vv, defaultValue)
	if flg != nil {
		flg.DefaultValue(defaultValue, "")
		flg.SetOwner(s)
	}
	opt = flg

	if opt != nil && len(titles) > 0 {
		opt.Long(titles[0])
		if len(titles) > 1 {
			opt.Short(titles[1])
			if len(titles) > 2 {
				opt.Aliases(titles[2:]...)
			}
		}
	}
	return
}

func (s *optCommandImpl) NewSubCommand(titles ...string) (opt OptCmd) {
	cmd := &Command{root: internalGetWorker().rootCommand}

	optCtx.current = cmd

	s.working.SubCommands = uniAddCmd(s.working.SubCommands, cmd)

	opt = &subCmdOpt{optCommandImpl: optCommandImpl{working: cmd, parent: s}}

	if len(titles) > 0 {
		opt.Long(titles[0])
		if len(titles) > 1 {
			opt.Short(titles[1])
			if len(titles) > 2 {
				opt.Aliases(titles[2:]...)
			}
		}
	}
	return
}

func (s *optCommandImpl) OwnerCommand() (opt OptCmd) {
	opt = s.parent
	return
}

func (s *optCommandImpl) SetOwner(opt OptCmd) {
	s.parent = opt
	return
}

func (s *optCommandImpl) RootCommand() (root *RootCommand) {
	root = optCtx.root
	return
}
