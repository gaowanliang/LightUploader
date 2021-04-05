/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import "time"

type (
	// RootCmdOpt for fluent api
	RootCmdOpt struct {
		optCommandImpl
	}

	// cmdOpt for fluent api
	cmdOpt struct {
		optCommandImpl
	}

	// subCmdOpt for fluent api
	subCmdOpt struct {
		optCommandImpl
	}

	// stringOpt for fluent api
	stringOpt struct {
		optFlagImpl
	}

	// stringSliceOpt for fluent api
	stringSliceOpt struct {
		optFlagImpl
	}

	// intSliceOpt for fluent api
	intSliceOpt struct {
		optFlagImpl
	}

	// uintSliceOpt for fluent api
	uintSliceOpt struct {
		optFlagImpl
	}

	// boolOpt for fluent api
	boolOpt struct {
		optFlagImpl
	}

	// intOpt for fluent api
	intOpt struct {
		optFlagImpl
	}

	// uintOpt for fluent api
	uintOpt struct {
		optFlagImpl
	}

	// int64Opt for fluent api
	int64Opt struct {
		optFlagImpl
	}

	// uint64Opt for fluent api
	uint64Opt struct {
		optFlagImpl
	}

	// float32Opt for fluent api
	float32Opt struct {
		optFlagImpl
	}

	// float64Opt for fluent api
	float64Opt struct {
		optFlagImpl
	}

	// complex64Opt for fluent api
	complex64Opt struct {
		optFlagImpl
	}

	// complex128Opt for fluent api
	complex128Opt struct {
		optFlagImpl
	}

	// durationOpt for fluent api
	durationOpt struct {
		optFlagImpl
	}
)

// Header for fluent api
func (s *RootCmdOpt) Header(header string) *RootCmdOpt {
	optCtx.root.Header = header
	return s
}

// Copyright for fluent api
func (s *RootCmdOpt) Copyright(copyright, author string) *RootCmdOpt {
	optCtx.root.Copyright = copyright
	optCtx.root.Author = author
	return s
}

// AddGlobalPreAction attaches your pre handler to the global Pre-Actions list
func (s *RootCmdOpt) AddGlobalPreAction(pre Handler) *RootCmdOpt {
	optCtx.root.PreActions = append(optCtx.root.PreActions, pre)
	return s
}

// AddGlobalPostAction attaches your post handler to the global Post-Actions list
func (s *RootCmdOpt) AddGlobalPostAction(post Invoker) *RootCmdOpt {
	optCtx.root.PostActions = append(optCtx.root.PostActions, post)
	return s
}

// func (s *RootCmdOpt) Command(cmdOpt *cmdOpt) *RootCmdOpt {
// 	optCtx.root.Command = *cmdOpt.workingFlag
// 	return s
// }

// func (s *RootCmdOpt) SubCmd() (opt OptCmd) {
// 	cmd := &Command{}
// 	optCtx.root.SubCommands = append(optCtx.root.SubCommands, cmd)
// 	optCtx.current = cmd
// 	return &subCmdOpt{optCommandImpl: optCommandImpl{workingFlag: cmd},}
// }

// NewCmdFrom creates a wrapped Command object as OptCmd, and make it as the current working item.
func NewCmdFrom(cmd *Command) (opt OptCmd) {
	optCtx.current = cmd
	return &cmdOpt{optCommandImpl: optCommandImpl{working: optCtx.current}}
}

// NewCmd creates a wrapped Command object as OptCmd
func NewCmd() (opt OptCmd) {
	current := &Command{}
	return &cmdOpt{optCommandImpl: optCommandImpl{working: current}}
}

// NewSubCmd creates a wrapped Command object as OptCmd, and append it into the current working item.
func NewSubCmd() (opt OptCmd) {
	cmd := &Command{}
	optCtx.current.SubCommands = uniAddCmd(optCtx.current.SubCommands, cmd)
	return &subCmdOpt{optCommandImpl: optCommandImpl{working: cmd}}
}

func uniAddCmd(cmds []*Command, cmd *Command) []*Command {
	for _, f := range cmds {
		if f == cmd {
			return cmds
		}
	}
	return append(cmds, cmd)
}

func uniAddFlg(flags []*Flag, flg *Flag) []*Flag {
	for _, f := range flags {
		if f == flg {
			return flags
		}
	}
	return append(flags, flg)
}

func uniAddStr(a []string, s string) []string {
	for _, f := range a {
		if f == s {
			return a
		}
	}
	return append(a, s)
}

func uniAddStrs(a []string, ss ...string) []string {
	for _, s := range ss {
		found := false
		for _, f := range a {
			if f == s {
				found = true
				break
			}
		}
		if !found {
			a = append(a, s)
		}
	}
	return a
}

// NewBool creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewBool(defaultValue ...bool) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, workingFlag)
	opt = &boolOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv bool
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewString creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewString(defaultValue ...string) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &stringOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv string
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewStringSlice creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
//
// Sample:
//
//    cmdr.NewStringSlice("quick", "fox", "jumps").Titles("int-slice", "is").Group("").
//        AttachTo(parentCmdOpt)
func NewStringSlice(defaultValue ...string) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &stringSliceOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewIntSlice creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
//
// Sample:
//
//    cmdr.NewIntSlice(1, 2, 3).Titles("int-slice", "is").Group("").
//        AttachTo(parentCmdOpt)
func NewIntSlice(defaultValue ...int) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &intSliceOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewInt creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewInt(defaultValue ...int) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &intOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv int
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewUintSlice creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
//
// Sample:
//
//    cmdr.NewUintSlice(1, 2, 3).Titles("uint-slice", "us").Group("").
//        AttachTo(parentCmdOpt)
func NewUintSlice(defaultValue ...uint) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &uintSliceOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	opt.DefaultValue(defaultValue, "")
	return
}

// NewUint creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewUint(defaultValue ...uint) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &uintOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv uint
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewInt64 creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewInt64(defaultValue ...int64) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &int64Opt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv int64
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewUint64 creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewUint64(defaultValue ...uint64) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &uint64Opt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv uint64
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewFloat32 creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewFloat32(defaultValue ...float32) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &float32Opt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv float32
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewFloat64 creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewFloat64(defaultValue ...float64) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &float64Opt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv float64
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewComplex64 creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewComplex64(defaultValue ...complex64) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &complex64Opt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv complex64
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewComplex128 creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewComplex128(defaultValue ...complex128) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &complex128Opt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv complex128
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewDuration creates a wrapped OptFlag, you can connect it to a OptCmd via OptFlag.AttachXXX later.
func NewDuration(defaultValue ...time.Duration) (opt OptFlag) {
	workingFlag := &Flag{}
	// optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	opt = &durationOpt{optFlagImpl: optFlagImpl{working: workingFlag}}
	var dv time.Duration
	for _, v := range defaultValue {
		dv = v
	}
	opt.DefaultValue(dv, "")
	return
}

// NewDurationFrom creates a wrapped OptFlag, and append it into the current working item.
func NewDurationFrom(flg *Flag) (opt OptFlag) {
	optCtx.workingFlag = flg
	optCtx.current.Flags = uniAddFlg(optCtx.current.Flags, optCtx.workingFlag)
	return &durationOpt{optFlagImpl: optFlagImpl{working: optCtx.workingFlag}}
}
