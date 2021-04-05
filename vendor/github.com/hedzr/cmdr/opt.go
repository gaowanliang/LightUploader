/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

type (
	// // Opt never used?
	// Opt interface {
	// 	Titles(short, long string, aliases ...string) (opt Opt)
	// 	Short(short string) (opt Opt)
	// 	Long(long string) (opt Opt)
	// 	Aliases(ss ...string) (opt Opt)
	// 	Description(oneLine, long string) (opt Opt)
	// 	Examples(examples string) (opt Opt)
	// 	Group(group string) (opt Opt)
	// 	Hidden(hidden bool) (opt Opt)
	// 	Deprecated(deprecation string) (opt Opt)
	// 	Action(action Handler) (opt Opt)
	// }

	// OptFlag to support fluent api of cmdr.
	// see also: cmdr.Root().NewSubCommand()/.NewFlag()
	//
	// For an option, its default value must be declared with exact type as is
	OptFlag interface {
		// Titles: broken API since v1.6.39.
		//
		// If necessary, an order prefix can be attached to the long title.
		// The title with prefix will be set to Name field and striped to Long field.
		//
		// An order prefix is a dotted string with multiple alphabet and digit. Such as:
		// "zzzz.", "0001.", "700.", "A1." ...
		Titles(long, short string, aliases ...string) (opt OptFlag)
		Short(short string) (opt OptFlag)
		Long(long string) (opt OptFlag)
		// Name is an internal identity, and an order prefix is optional
		// An order prefix is a dotted string with multiple alphabet and digit. Such as:
		// "zzzz.", "0001.", "700.", "A1." ...
		Name(name string) (opt OptFlag)
		Aliases(ss ...string) (opt OptFlag)
		Description(oneLineDesc string, longDesc ...string) (opt OptFlag)
		Examples(examples string) (opt OptFlag)
		Group(group string) (opt OptFlag)
		Hidden(hidden bool) (opt OptFlag)
		Deprecated(deprecation string) (opt OptFlag)
		// Action will be triggered once being parsed ok
		Action(action Handler) (opt OptFlag)

		ToggleGroup(group string) (opt OptFlag)
		// DefaultValue needs an exact typed 'val'.
		// IMPORTANT: cmdr interprets value type of an option based on the underlying default value set.
		DefaultValue(val interface{}, placeholder string) (opt OptFlag)
		Placeholder(placeholder string) (opt OptFlag)
		ExternalTool(envKeyName string) (opt OptFlag)
		ValidArgs(list ...string) (opt OptFlag)
		// HeadLike enables `head -n` mode.
		// 'min', 'max' will be ignored at this version, its might be impl in the future.
		// There's only one head-like flag in one command and its parent and children commands.
		HeadLike(enable bool, min, max int64) (opt OptFlag)

		// EnvKeys is a list of env-var names of binding on this flag
		EnvKeys(keys ...string) (opt OptFlag)
		// Required flag.
		//
		// NOTE
		//
		//   Required() set the required flag to true while it's invoked with empty params.
		Required(required ...bool) (opt OptFlag)

		OwnerCommand() (opt OptCmd)
		SetOwner(opt OptCmd)

		RootCommand() *RootCommand

		ToFlag() *Flag

		// AttachTo attach as a flag of `opt` OptCmd object
		AttachTo(opt OptCmd)
		// AttachToCommand attach as a flag of *Command object
		AttachToCommand(cmd *Command)
		// AttachToRoot attach as a flag of *RootCommand object
		AttachToRoot(root *RootCommand)

		OnSet
	}

	// OptCmd to support fluent api of cmdr.
	// see also: cmdr.Root().NewSubCommand()/.NewFlag()
	OptCmd interface {
		// Titles: broken API since v1.6.39
		//
		// If necessary, an order prefix can be attached to the long title.
		// The title with prefix will be set to Name field and striped to Long field.
		//
		// An order prefix is a dotted string with multiple alphabet and digit. Such as:
		// "zzzz.", "0001.", "700.", "A1." ...
		Titles(long, short string, aliases ...string) (opt OptCmd)
		Short(short string) (opt OptCmd)
		Long(long string) (opt OptCmd)
		// Name is an internal identity, and an order prefix is optional
		// An order prefix is a dotted string with multiple alphabet and digit. Such as:
		// "zzzz.", "0001.", "700.", "A1." ...
		Name(name string) (opt OptCmd)
		Aliases(ss ...string) (opt OptCmd)
		Description(oneLine string, long ...string) (opt OptCmd)
		Examples(examples string) (opt OptCmd)
		Group(group string) (opt OptCmd)
		Hidden(hidden bool) (opt OptCmd)
		Deprecated(deprecation string) (opt OptCmd)
		// Action will be triggered after all command-line arguments parsed
		Action(action Handler) (opt OptCmd)

		// FlagAdd(flg *Flag) (opt OptCmd)
		// SubCommand(cmd *Command) (opt OptCmd)

		// PreAction will be invoked before running Action
		// NOTE that RootCommand.PreAction will be invoked too.
		PreAction(pre Handler) (opt OptCmd)
		// PostAction will be invoked after run Action
		// NOTE that RootCommand.PostAction will be invoked too.
		PostAction(post Invoker) (opt OptCmd)

		TailPlaceholder(placeholder string) (opt OptCmd)

		// NewFlag create a new flag object and return it for further operations.
		// Deprecated since v1.6.9, replace it with FlagV(defaultValue)
		//
		// Deprecated since v1.6.50, we recommend the new form:
		//    cmdr.NewBool(false).Titles(...)...AttachTo(ownerCmd)
		NewFlag(typ OptFlagType) (opt OptFlag)
		// NewFlagV create a new flag object and return it for further operations.
		// the titles in arguments MUST be: longTitle, [shortTitle, [aliasTitles...]]
		//
		// Deprecated since v1.6.50, we recommend the new form:
		//    cmdr.NewBool(false).Titles(...)...AttachTo(ownerCmd)
		NewFlagV(defaultValue interface{}, titles ...string) (opt OptFlag)
		// NewSubCommand make a new sub-command optcmd object with optional titles.
		// the titles in arguments MUST be: longTitle, [shortTitle, [aliasTitles...]]
		NewSubCommand(titles ...string) (opt OptCmd)

		OwnerCommand() (opt OptCmd)
		SetOwner(opt OptCmd)

		RootCommand() *RootCommand

		ToCommand() *Command

		AddOptFlag(flag OptFlag)
		AddFlag(flag *Flag)
		// AddOptCmd adds 'opt' OptCmd as a sub-command
		AddOptCmd(opt OptCmd)
		// AddCommand adds a *Command as a sub-command
		AddCommand(cmd *Command)
		// AttachTo attaches itself as a sub-command of 'opt' OptCmd object
		AttachTo(opt OptCmd)
		// AttachTo attaches itself as a sub-command of *Command object
		AttachToCommand(cmd *Command)
		// AttachTo attaches itself as a sub-command of *RootCommand object
		AttachToRoot(root *RootCommand)
	}

	// OnSet interface
	OnSet interface {
		// OnSet will be callback'd after this flag parsed
		OnSet(f func(keyPath string, value interface{})) (opt OptFlag)
	}

	// OptFlagType to support fluent api of cmdr.
	// see also: OptCmd.NewFlag(OptFlagType)
	//
	// Usage
	//
	//   root := cmdr.Root()
	//   co := root.NewSubCommand()
	//   co.NewFlag(cmdr.OptFlagTypeUint)
	//
	// See also those short-hand constructors: Bool(), Int(), ....
	OptFlagType int
)

const (
	// OptFlagTypeBool to create a new bool flag
	OptFlagTypeBool OptFlagType = iota
	// OptFlagTypeInt to create a new int flag
	OptFlagTypeInt OptFlagType = iota + 1
	// OptFlagTypeUint to create a new uint flag
	OptFlagTypeUint OptFlagType = iota + 2
	// OptFlagTypeInt64 to create a new int64 flag
	OptFlagTypeInt64 OptFlagType = iota + 3
	// OptFlagTypeUint64 to create a new uint64 flag
	OptFlagTypeUint64 OptFlagType = iota + 4
	// OptFlagTypeFloat32 to create a new int float32 flag
	OptFlagTypeFloat32 OptFlagType = iota + 8
	// OptFlagTypeFloat64 to create a new int float64 flag
	OptFlagTypeFloat64 OptFlagType = iota + 9
	// OptFlagTypeComplex64 to create a new int complex64 flag
	OptFlagTypeComplex64 OptFlagType = iota + 10
	// OptFlagTypeComplex128 to create a new int complex128 flag
	OptFlagTypeComplex128 OptFlagType = iota + 11
	// OptFlagTypeString to create a new string flag
	OptFlagTypeString OptFlagType = iota + 12
	// OptFlagTypeStringSlice to create a new string slice flag
	OptFlagTypeStringSlice OptFlagType = iota + 13
	// OptFlagTypeIntSlice to create a new int slice flag
	OptFlagTypeIntSlice OptFlagType = iota + 14
	// OptFlagTypeInt64Slice to create a new int slice flag
	OptFlagTypeInt64Slice OptFlagType = iota + 15
	// OptFlagTypeUint64Slice to create a new int slice flag
	OptFlagTypeUint64Slice OptFlagType = iota + 16
	// OptFlagTypeDuration to create a new duration flag
	OptFlagTypeDuration OptFlagType = iota + 17
	// OptFlagTypeHumanReadableSize to create a new human readable size flag
	OptFlagTypeHumanReadableSize OptFlagType = iota + 18
)

type optContext struct {
	current     *Command
	root        *RootCommand
	workingFlag *Flag
}

var optCtx *optContext

// Root for fluent api, to create a new [*RootCmdOpt] object.
func Root(appName, version string) (opt *RootCmdOpt) {
	root := &RootCommand{AppName: appName, Version: version, Command: Command{BaseOpt: BaseOpt{Name: appName}}}
	// rootCommand = root
	opt = RootFrom(root)
	return
}

// RootFrom for fluent api, to create the new [*RootCmdOpt] object from an existed [RootCommand]
func RootFrom(root *RootCommand) (opt *RootCmdOpt) {
	optCtx = &optContext{current: &root.Command, root: root, workingFlag: nil}

	opt = &RootCmdOpt{optCommandImpl: optCommandImpl{working: optCtx.current}}
	opt.parent = opt
	return
}
