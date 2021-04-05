/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"github.com/hedzr/cmdr/tool"
	"github.com/hedzr/log"
	"sync"
)

const (
	appNameDefault = "cmdr"

	// UnsortedGroup for commands and flags
	UnsortedGroup = "zzzg.unsorted"
	// AddonsGroup for commands and flags
	AddonsGroup = "zzzh.Addons"
	// ExtGroup for commands and flags
	ExtGroup = "zzzi.Extensions"
	// AliasesGroup for commands and flags
	AliasesGroup = "zzzj.Aliases"
	// SysMgmtGroup for commands and flags
	SysMgmtGroup = "zzz9.Misc"

	// DefaultEditor is 'vim'
	DefaultEditor = "vim"

	// ExternalToolEditor environment variable name, EDITOR is fit for most of shells.
	ExternalToolEditor = "EDITOR"

	// ExternalToolPasswordInput enables secure password input without echo.
	ExternalToolPasswordInput = "PASSWD"
)

type (
	// BaseOpt is base of `Command`, `Flag`
	BaseOpt struct {
		Name string `yaml:"name,omitempty" json:"name,omitempty"`
		// Short rune. short option/command name.
		// single char. example for flag: "a" -> "-a"
		Short string `yaml:"short-name,omitempty" json:"short-name,omitempty"`
		// Full full/long option/command name.
		// word string. example for flag: "addr" -> "--addr"
		Full string `yaml:"title,omitempty" json:"title,omitempty"`
		// Aliases are the more synonyms
		Aliases []string `yaml:"aliases,flow,omitempty" json:"aliases,flow,omitempty"`
		// Group group name
		Group string `yaml:"group,omitempty" json:"group,omitempty"`

		owner *Command
		// strHit keeps the matched title string from user input in command line
		strHit string

		Description     string `yaml:"desc,omitempty" json:"desc,omitempty"`
		LongDescription string `yaml:"long-desc,omitempty" json:"long-desc,omitempty"`
		Examples        string `yaml:"examples,omitempty" json:"examples,omitempty"`
		Hidden          bool   `yaml:"hidden,omitempty" json:"hidden,omitempty"`

		// Deprecated is a version string just like '0.5.9' or 'v0.5.9', that means this command/flag was/will be deprecated since `v0.5.9`.
		Deprecated string `yaml:"deprecated,omitempty" json:"deprecated,omitempty"`

		// Action is callback for the last recognized command/sub-command.
		// return: ErrShouldBeStopException will break the following flow and exit right now
		// cmd 是 flag 被识别时已经得到的子命令
		Action Handler `yaml:"-" json:"-"`
	}

	// Handler handles the event on a subcommand matched
	Handler func(cmd *Command, args []string) (err error)
	// Invoker is a Handler but without error returns
	Invoker func(cmd *Command, args []string)

	// Command holds the structure of commands and sub-commands
	Command struct {
		BaseOpt `yaml:",inline"`

		Flags []*Flag `yaml:"flags,omitempty" json:"flags,omitempty"`

		SubCommands []*Command `yaml:"subcmds,omitempty" json:"subcmds,omitempty"`

		// return: ErrShouldBeStopException will break the following flow and exit right now
		PreAction Handler `yaml:"-" json:"-"`
		// PostAction will be run after Action() invoked.
		PostAction Invoker `yaml:"-" json:"-"`
		// be shown at tail of command usages line. Such as for TailPlaceHolder="<host-fqdn> <ipv4/6>":
		// austr dns add <host-fqdn> <ipv4/6> [Options] [Parent/Global Options]
		TailPlaceHolder string `yaml:"tail-placeholder,omitempty" json:"tail-placeholder,omitempty"`
		// TailArgsText string
		// TailArgsDesc string

		root            *RootCommand
		allCmds         map[string]map[string]*Command // key1: Commnad.Group, key2: Command.Full
		allFlags        map[string]map[string]*Flag    // key1: Command.Flags[#].Group, key2: Command.Flags[#].Fullui
		plainCmds       map[string]*Command
		plainShortFlags map[string]*Flag
		plainLongFlags  map[string]*Flag
		headLikeFlag    *Flag

		presetCmdLines []string
		// Invoke is just for importing from a file.
		// invoke a command from the command tree in this app
		Invoke string `yaml:"invoke,omitempty" json:"invoke,omitempty"`
		// InvokeProc is just for importing from a file.
		// invoke the external commands (via: executable)
		InvokeProc string `yaml:"invoke-proc,omitempty" json:"invoke-proc,omitempty"`
		// InvokeShell is just for importing from a file.
		// invoke the external commands (via: shell)
		InvokeShell string `yaml:"invoke-sh,omitempty" json:"invoke-sh,omitempty"`
		// Shell is just for importing from a file.
		// invoke a command under this shell, or /usr/bin/env bash|zsh|...
		Shell string `yaml:"shell,omitempty" json:"shell,omitempty"`
	}

	// RootCommand holds some application information
	RootCommand struct {
		Command `yaml:",inline"`

		AppName    string `yaml:"appname,omitempty" json:"appname,omitempty"`
		Version    string `yaml:"version,omitempty" json:"version,omitempty"`
		VersionInt uint32 `yaml:"version-int,omitempty" json:"version-int,omitempty"`

		Copyright string `yaml:"copyright,omitempty" json:"copyright,omitempty"`
		Author    string `yaml:"author,omitempty" json:"author,omitempty"`
		Header    string `yaml:"header,omitempty" json:"header,omitempty"` // using `Header` for header and ignore built with `Copyright` and `Author`, and no usage lines too.

		PreActions  []Handler `yaml:"-" json:"-"`
		PostActions []Invoker `yaml:"-" json:"-"`

		ow   *bufio.Writer
		oerr *bufio.Writer
	}

	// Flag means a flag, a option, or a opt.
	Flag struct {
		BaseOpt `yaml:",inline"`

		// ToggleGroup for Toggle Group
		ToggleGroup string `yaml:"toggle-group,omitempty" json:"toggle-group,omitempty"`
		// DefaultValuePlaceholder for flag
		DefaultValuePlaceholder string `yaml:"default-placeholder,omitempty" json:"default-placeholder,omitempty"`
		// DefaultValue default value for flag
		DefaultValue interface{} `yaml:"default,omitempty" json:"default,omitempty"`
		// DefaultValueType is a string to indicate the data-type of DefaultValue.
		// It's only available in loading flag definition from a config-file (yaml/json/...).
		// Never used in writing your codes manually.
		DefaultValueType string `yaml:"type,omitempty" json:"type,omitempty"`
		// ValidArgs for enum flag
		ValidArgs []string `yaml:"valid-args,omitempty" json:"valid-args,omitempty"`
		// Required to-do
		Required bool `yaml:"required,omitempty" json:"required,omitempty"`

		// ExternalTool to get the value text by invoking external tool.
		// It's an environment variable name, such as: "EDITOR" (or cmdr.ExternalToolEditor)
		ExternalTool string `yaml:"external-tool,omitempty" json:"external-tool,omitempty"`

		// EnvVars give a list to bind to environment variables manually
		// it'll take effects since v1.6.9
		EnvVars []string `yaml:"envvars,flow,omitempty" json:"envvars,flow,omitempty"`

		// HeadLike enables a free-hand option like `head -3`.
		//
		// When a free-hand option presents, it'll be treated as a named option with an integer value.
		//
		// For example, option/flag = `{{Full:"line",Short:"l"},HeadLike:true}`, the command line:
		// `app -3`
		// is equivalent to `app -l 3`, and so on.
		//
		// HeadLike assumed an named option with an integer value, that means, Min and Max can be applied on it too.
		// NOTE: Only one head-like option can be defined in a command/sub-command chain.
		HeadLike bool `yaml:"head-like,omitempty" json:"head-like,omitempty"`

		// Min minimal value of a range.
		Min int64 `yaml:"min,omitempty" json:"min,omitempty"`
		// Max maximal value of a range.
		Max int64 `yaml:"max,omitempty" json:"max,omitempty"`

		onSet func(keyPath string, value interface{})

		// times how many times this flag was triggered.
		// To access it with `Flag.GetTriggeredTimes()`.
		times int

		// PostAction treat this flag as a command!
		// PostAction Handler

		// by default, a flag is always `optional`.
	}

	// Options is a holder of all options
	Options struct {
		entries   map[string]interface{}
		hierarchy map[string]interface{}
		rw        *sync.RWMutex

		usedConfigFile      string
		usedConfigSubDir    string
		usedAlterConfigFile string
		configFiles         []string
		filesWatching       []string
		batchMerging        bool

		onConfigReloadedFunctions map[ConfigReloaded]bool
		rwlCfgReload              *sync.RWMutex
		rwCB                      sync.RWMutex
		onMergingSet              OnOptionSetCB
		onSet                     OnOptionSetCB
	}

	// OptOne struct {
	// 	Children map[string]*OptOne `yaml:"c,omitempty"`
	// 	Value    interface{}        `yaml:"v,omitempty"`
	// }

	// ConfigReloaded for config reloaded
	ConfigReloaded interface {
		OnConfigReloaded()
	}

	// OnOptionSetCB is a callback function while an option is being set (or merged)
	OnOptionSetCB func(keyPath string, value, oldVal interface{})
	// OnSwitchCharHitCB is a callback function ...
	OnSwitchCharHitCB func(parsed *Command, switchChar string, args []string) (err error)
	// OnPassThruCharHitCB is a callback function ...
	OnPassThruCharHitCB func(parsed *Command, switchChar string, args []string) (err error)

	// HookFunc the hook function prototype for SetBeforeXrefBuilding and SetAfterXrefBuilt
	HookFunc func(root *RootCommand, args []string)

	// HookOptsFunc the hook function prototype
	HookOptsFunc func(root *RootCommand, opts *Options)

	// HookHelpScreenFunc the hook function prototype
	HookHelpScreenFunc func(w *ExecWorker, p Painter, cmd *Command, justFlags bool)
)

var (
	//
	// doNotLoadingConfigFiles = false

	// // rootCommand the root of all commands
	// rootCommand *RootCommand
	// // rootOptions *Opt
	// rxxtOptions = newOptions()

	// usedConfigFile
	// usedConfigFile            string
	// usedConfigSubDir          string
	// configFiles               []string
	// onConfigReloadedFunctions map[ConfigReloaded]bool
	//
	// predefinedLocations = []string{
	// 	"./ci/etc/%s/%s.yml",
	// 	"/etc/%s/%s.yml",
	// 	"/usr/local/etc/%s/%s.yml",
	// 	os.Getenv("HOME") + "/.%s/%s.yml",
	// }

	//
	// defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
	// defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)

	//
	// currentHelpPainter Painter

	// CurrentDescColor the print color for description line
	CurrentDescColor = FgDarkGray
	// CurrentDefaultValueColor the print color for default value line
	CurrentDefaultValueColor = FgDarkGray
	// CurrentGroupTitleColor the print color for titles
	CurrentGroupTitleColor = DarkColor

	// globalShowVersion   func()
	// globalShowBuildInfo func()

	// beforeXrefBuilding []HookFunc
	// afterXrefBuilt     []HookFunc

	// getEditor sets callback to get editor program
	// getEditor func() (string, error)

	defaultStringMetric = tool.JaroWinklerDistance(tool.JWWithThreshold(similarThreshold))
)

const similarThreshold = 0.6666666666666666

// GetStrictMode enables error when opt value missed. such as:
// xxx a b --prefix''   => error: prefix opt has no value specified.
// xxx a b --prefix'/'  => ok.
//
// ENV: use `CMDR_APP_STRICT_MODE=true` to enable strict-mode.
// NOTE: `CMDR_APP_` prefix could be set by user (via: `EnvPrefix` && `RxxtPrefix`).
//
// the flag value of `--strict-mode`.
func GetStrictMode() bool {
	return GetBoolR("strict-mode")
}

// GetTraceMode returns the flag value of `--trace`/`-tr`
//
// NOTE
//     log.GetTraceMode()/SetTraceMode() have higher universality
//
// the flag value of `--trace` or `-tr` is always stored
// in cmdr Option Store, so you can retrieved it by
// GetBoolR("trace") and set it by Set("trace", true).
// You could also set it with SetTraceMode(b bool).
//
// The `--trace` is not enabled in default, so you have to
// add it manually:
//
//     import "github.com/hedzr/cmdr-addons/pkg/plugins/trace"
//     cmdr.Exec(buildRootCmd(),
//         trace.WithTraceEnable(true),
//     )
func GetTraceMode() bool {
	return GetBoolR("trace") || log.GetTraceMode()
}

// SetTraceMode setup the tracing mode status in Option Store
func SetTraceMode(b bool) {
	Set("trace", b)
}

// GetDebugMode returns the flag value of `--debug`/`-D`
//
// NOTE
//     log.GetDebugMode()/SetDebugMode() have higher universality
//
// the flag value of `--debug` or `-D` is always stored
// in cmdr Option Store, so you can retrieved it by
// GetBoolR("debug") and set it by Set("debug", true).
// You could also set it with SetDebugMode(b bool).
func GetDebugMode() bool {
	return GetBoolR("debug") || log.GetDebugMode()
}

// SetDebugMode setup the debug mode status in Option Store
func SetDebugMode(b bool) {
	Set("debug", b)
}

// NewLoggerConfig returns a default LoggerConfig
func NewLoggerConfig() *log.LoggerConfig {
	lc := NewLoggerConfigWith(false, "sugar", "error")
	return lc
}

// NewLoggerConfigWith returns a default LoggerConfig
func NewLoggerConfigWith(enabled bool, backend, level string) *log.LoggerConfig {
	log.SetTraceMode(GetTraceMode())
	log.SetDebugMode(GetDebugMode())
	lc := log.NewLoggerConfigWith(enabled, backend, level)
	_ = GetSectionFrom("logger", &lc)
	return lc
}

// GetVerboseMode returns the flag value of `--verbose`/`-v`
func GetVerboseMode() bool {
	return GetBoolR("verbose")
}

// GetQuietMode returns the flag value of `--quiet`/`-q`
func GetQuietMode() bool {
	return GetBoolR("quiet")
}

// GetNoColorMode return the flag value of `--no-color`
func GetNoColorMode() bool {
	return GetBoolR("no-color")
}

// func init() {
// 	// onConfigReloadedFunctions = make(map[ConfigReloaded]bool)
// 	// SetCurrentHelpPainter(new(helpPainter))
// }
