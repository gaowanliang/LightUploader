/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

import (
	"bufio"
	"github.com/hedzr/cmdr/conf"
	"gopkg.in/hedzr/errors.v2"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
)

// WithXrefBuildingHooks sets the hook before and after building xref indices.
// It's replacers for AddOnBeforeXrefBuilding, and AddOnAfterXrefBuilt.
//
// By using beforeXrefBuildingX, you could append, modify, or remove the
// builtin commands and options.
func WithXrefBuildingHooks(beforeXrefBuildingX, afterXrefBuiltX HookFunc) ExecOption {
	return func(w *ExecWorker) {
		if beforeXrefBuildingX != nil {
			w.beforeXrefBuilding = append(w.beforeXrefBuilding, beforeXrefBuildingX)
		}
		if afterXrefBuiltX != nil {
			w.afterXrefBuilt = append(w.afterXrefBuilt, afterXrefBuiltX)
		}
	}
}

// WithAutomaticEnvHooks sets the hook after building automatic environment.
//
// At this point, you could override the option default values.
func WithAutomaticEnvHooks(hook HookOptsFunc) ExecOption {
	return func(w *ExecWorker) {
		if hook != nil {
			w.afterAutomaticEnv = append(w.afterAutomaticEnv, hook)
		}
	}
}

// WithEnvVarMap adds a (envvar-name, value) map, which will be applied
// to string option value, string-slice option values, ....
// For example, you could define a key-value entry in your `<app>.yml` file:
//    app:
//      test-value: "$THIS/$APPNAME.yml"
//      home-dir: "$HOME"
// it will be expanded by mapping to OS environment and this map (WithEnvVarMap).
// That is, $THIS will be expanded to the directory path of this
// executable, $APPNAME to the app name.
// And of course, $HOME will be mapped to os home directory path.
func WithEnvVarMap(varToValue map[string]func() string) ExecOption {
	return func(w *ExecWorker) {
		if varToValue == nil {
			varToValue = make(map[string]func() string)
		}
		w.envVarToValueMap = varToValue
		testAndSetMap(w.envVarToValueMap, "THIS", func() string { return GetExecutableDir() })
		testAndSetMap(w.envVarToValueMap, "APPNAME", func() string { return conf.AppName })
		testAndSetMap(w.envVarToValueMap, "CFG_DIR", func() string { return path.Dir(GetUsedConfigFile()) })
	}
}

func testAndSetMap(m map[string]func() string, key string, value func() string) {
	if _, ok := m[key]; !ok {
		m[key] = value
	}
}

// WithEnvPrefix sets the environment variable text prefixes.
// cmdr will lookup envvars for a key.
// Default env-prefix is array ["CMDR"], ie 'CMDR_'
func WithEnvPrefix(prefix ...string) ExecOption {
	return func(w *ExecWorker) {
		w.envPrefixes = prefix
	}
}

// WithOptionsPrefix create a top-level namespace, which contains all normalized `Flag`s.
// =WithRxxtPrefix
// Default Options Prefix is array ["app"], ie 'app.xxx'
func WithOptionsPrefix(prefix ...string) ExecOption {
	return func(w *ExecWorker) {
		w.rxxtPrefixes = prefix
	}
}

// WithRxxtPrefix create a top-level namespace, which contains all normalized `Flag`s.
// cmdr will lookup envvars for a key.
// Default Options Prefix is array ["app"], ie 'app.xxx'
func WithRxxtPrefix(prefix ...string) ExecOption {
	return func(w *ExecWorker) {
		w.rxxtPrefixes = prefix
	}
}

// WithPluginLocations sets the addon locations.
//
// An addon is a golang plugin for cmdr.
//
// Default locations are:
//
//     []string{
//		   "./ci/local/share/$APPNAME/addons",
//		   "$HOME/.local/share/$APPNAME/addons",
//		   "$HOME/.$APPNAME/addons",
//		   "/usr/local/share/$APPNAME/addons",
//		   "/usr/share/$APPNAME/addons",
//     },
//
// See also internalResetWorkerNoLock()
func WithPluginLocations(locations ...string) ExecOption {
	return func(w *ExecWorker) {
		w.pluginsLocations = locations
	}
}

// WithExtensionsLocations sets the extension locations.
//
// A extension is an executable (shell scripts, binary executable, ,,,).
//
// Default locations are:
//
//     []string{
//		   "./ci/local/share/$APPNAME/ext",
//		   "$HOME/.local/share/$APPNAME/ext",
//		   "$HOME/.$APPNAME/ext",
//		   "/usr/local/share/$APPNAME/ext",
//		   "/usr/share/$APPNAME/ext",
//     },
//
// See also internalResetWorkerNoLock()
func WithExtensionsLocations(locations ...string) ExecOption {
	return func(w *ExecWorker) {
		w.extensionsLocations = locations
	}
}

// WithPredefinedLocations sets the main config file locations.
//
// Default is:
//
//     []string{
//			"./ci/etc/$APPNAME/$APPNAME.yml",       // for developer
//			"/etc/$APPNAME/$APPNAME.yml",           // regular location
//			"/usr/local/etc/$APPNAME/$APPNAME.yml", // regular macOS HomeBrew location
//			"/opt/etc/$APPNAME/$APPNAME.yml",       // regular location
//			"$HOME/.config/$APPNAME/$APPNAME.yml",  // per user
//			"$HOME/.$APPNAME/$APPNAME.yml",         // ext location per user
//			"$THIS/$APPNAME.yml", // executable's directory
//			"$APPNAME.yml",       // current directory
//     },
//
// See also internalResetWorkerNoLock()
func WithPredefinedLocations(locations ...string) ExecOption {
	return func(w *ExecWorker) {
		w.predefinedLocations = locations
	}
}

// WithAlterLocations sets the alter config file locations.
//
// Default is:
//
//     alterLocations: []string{
//         "./bin/$APPNAME.yml", // for developer, current bin directory
//         "/var/lib/$APPNAME",
//         "$THIS/$APPNAME.yml", // executable's directory
//     },
//
// NOTE that just one config file will be loaded, the child `conf.d` folder not supports.
//
// cmdr will SAVE the changes in this alter config file automatically once loaded.
//
func WithAlterLocations(locations ...string) ExecOption {
	return func(w *ExecWorker) {
		w.alterLocations = locations
	}
}

// WithWatchMainConfigFileToo enables the watcher on main config file.
//
// By default, cmdr watches all files in the sub-directory `conf.d` of
// which folder contains the main config file. But as a feature,
// cmdr ignore the changes in main config file.
//
// WithWatchMainConfigFileToo can switch this feature.
//
// envvars:
//
//    CFG_DIR: will be set to the folder contains the main config file
//    no-watch-conf-dir: if set true, the watcher will be disabled.
//
func WithWatchMainConfigFileToo(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.watchMainConfigFileToo = b
	}
}

// WithNoLoadConfigFiles true means no loading config files
func WithNoLoadConfigFiles(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.doNotLoadingConfigFiles = b
	}
}

// WithNoWatchConfigFiles true means no watching the config files
func WithNoWatchConfigFiles(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.doNotWatchingConfigFiles = b
	}
}

// WithConfigLoadedListener adds a functor on config loaded and merged
func WithConfigLoadedListener(c ConfigReloaded) ExecOption {
	return func(w *ExecWorker) {
		AddOnConfigLoadedListener(c)
	}
}

// WithConfigSubDirAutoName specify an alternate folder name instead `conf.d`.
//
// By default, cmdr watches all files in the sub-directory `conf.d` of
// which folder contains the main config file.
//
func WithConfigSubDirAutoName(folderName string) ExecOption {
	return func(w *ExecWorker) {
		w.confDFolderName = folderName
	}
}

// WithSearchAlterConfigFiles adds CURRDIR/.<appname>.yml and CURRDIR/.<appname>/*.yml
// to the assumed config files and folders
func WithSearchAlterConfigFiles(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.watchChildConfigFiles = b
	}
}

// WithOptionMergeModifying adds a callback which invoked on new
// configurations was merging into, typically on loading the
// modified external config file(s).
// NOTE that MergeWith(...) can make it triggered too.
// A onMergingSet callback will be enabled after first loaded
// any external configuration files and environment variables
// merged.
func WithOptionMergeModifying(onMergingSet func(keyPath string, value, oldVal interface{})) ExecOption {
	return func(w *ExecWorker) {
		w.onOptionMergingSet = onMergingSet
	}
}

// WithOptionModifying adds a callback which invoked at each time
// any option was modified.
// It will be triggered after any external config files first loaded
// and the env variables had been merged.
func WithOptionModifying(onOptionSet func(keyPath string, value, oldVal interface{})) ExecOption {
	return func(w *ExecWorker) {
		w.onOptionSet = onOptionSet
	}
}

// WithIgnoreWrongEnumValue will be put into `cmdrError.Ignorable`
// while wrong enumerable value found in parsing command-line
// options.
// The default is true.
//
// Main program might decide whether it's a warning or error.
//
// See also
//
// [Flag.ValidArgs]
func WithIgnoreWrongEnumValue(ignored bool) ExecOption {
	return func(w *ExecWorker) {
		w.shouldIgnoreWrongEnumValue = ignored
		// ShouldIgnoreWrongEnumValue = ignored
	}
}

// WithWarnForUnknownCommand warns the end user if unknown command found.
//
// By default, cmdr ignore the first unknown command and treat them as
// remained arguments.
//
func WithWarnForUnknownCommand(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.treatUnknownCommandAsArgs = !b
	}
}

// WithBuiltinCommands enables/disables those builtin predefined commands. Such as:
//
//  - versionsCmds / EnableVersionCommands supports injecting the default `--version` flags and commands
//  - helpCmds / EnableHelpCommands supports injecting the default `--help` flags and commands
//  - verboseCmds / EnableVerboseCommands supports injecting the default `--verbose` flags and commands
//  - generalCmdrCmds / EnableCmdrCommands support these flags: `--strict-mode`, `--no-env-overrides`, and `--no-color`
//  - generateCmds / EnableGenerateCommands supports injecting the default `generate` commands and sub-commands
//
func WithBuiltinCommands(versionsCmds, helpCmds, verboseCmds, generateCmds, generalCmdrCmds bool) ExecOption {
	return func(w *ExecWorker) {
		w.enableVersionCommands = versionsCmds
		w.enableHelpCommands = helpCmds
		w.enableVerboseCommands = verboseCmds
		w.enableCmdrCommands = generalCmdrCmds
		w.enableGenerateCommands = generateCmds
	}
}

// WithInternalOutputStreams sets the internal output streams for debugging
func WithInternalOutputStreams(out, err *bufio.Writer) ExecOption {
	return func(w *ExecWorker) {
		w.defaultStdout = out
		w.defaultStderr = err

		if w.defaultStdout == nil {
			w.defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
		}
		if w.defaultStderr == nil {
			w.defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)
		}
	}
}

// WithConfigFileLoadingHooks adds the hook function to the front and back of trying to load config files.
//
// These two hooks always are triggered whatever WithNoLoadConfigFiles is enabled or not.
//
func WithConfigFileLoadingHooks(before, after HookFunc) ExecOption {
	return func(w *ExecWorker) {
		if before != nil {
			w.beforeConfigFileLoading = append(w.beforeConfigFileLoading, before)
		}
		if after != nil {
			w.afterConfigFileLoading = append(w.afterConfigFileLoading, after)
		}
	}
}

// WithHelpScreenHooks adds the hook function to the front and back of printing help screen
func WithHelpScreenHooks(before, after HookHelpScreenFunc) ExecOption {
	return func(w *ExecWorker) {
		if before != nil {
			w.beforeHelpScreen = append(w.beforeHelpScreen, before)
		}
		if after != nil {
			w.afterHelpScreen = append(w.afterHelpScreen, after)
		}
	}
}

// WithPagerEnabled transfer cmdr stdout to OS pager.
// The environment variable `PAGER` will be checkout, otherwise `less`.
//
// NOTICE ONLY the outputs of cmdr (such as help screen) will be paged.
func WithPagerEnabled(enabled ...bool) ExecOption {
	return func(w *ExecWorker) {
		for _, b := range enabled {
			enableShellPager(w, b)
			return
		}
		enableShellPager(w, true) // if params `enabled` are missed
	}
}

// EnableShellPager transfer cmdr stdout to OS pager.
// The environment variable `PAGER` will be checkout, otherwise `less`.
//
func EnableShellPager(enabled bool) {
	w := internalGetWorker()
	enableShellPager(w, enabled)
}

func enableShellPager(w *ExecWorker, enabled bool) {
	if enabled {
		w.closers = append(w.closers, openPager(w))
		return
	}

	for _, c := range w.closers {
		if c != nil {
			c()
		}
	}
	w.closers = nil
}

func closePager(w *ExecWorker, cmd *exec.Cmd, pager io.WriteCloser) func() {
	return func() {
		var err = errors.NewContainer("closePager errors")
		err.Attach(pager.Close())
		err.Attach(cmd.Wait())
		if !err.IsEmpty() {
			Logger.Errorf("Close pager errors: %v", err.Error())
		}
		// log.Println("pager closed.")

		w.defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
		w.defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)
		if w.bufferedStdio {
			w.rootCommand.ow = w.defaultStdout
			w.rootCommand.oerr = w.defaultStderr
		} else {
			w.rootCommand.ow = nil
			w.rootCommand.oerr = nil
		}
	}
}

func openPager(w *ExecWorker) (closer func()) {
	var err error
	var cmd *exec.Cmd
	var pager io.WriteCloser
	pagerApp := os.Getenv("PAGER")
	if pagerApp == "" {
		pagerApp = "less"

		// NOTE: here is another pager with column mode supports:
		// https://github.com/noborus/ov
	}
	cmd = exec.Command(pagerApp)
	pager, err = cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if pagerApp == "less" {
		cmd.Args = []string{pagerApp, "-SEX"}
	} else if runtime.GOOS == "darwin" {
		cmd.Args = []string{pagerApp, "-SEX", "-R-"}
	}
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// log.Printf("run %q %v....", pagerApp, cmd.Args)
	w.defaultStdout = bufio.NewWriterSize(pager, 32768)
	return closePager(w, cmd, pager)
}

// WithCustomShowVersion supports your `ShowVersion()` instead of internal `showVersion()`
func WithCustomShowVersion(fn func()) ExecOption {
	return func(w *ExecWorker) {
		w.globalShowVersion = fn
	}
}

// WithCustomShowBuildInfo supports your `ShowBuildInfo()` instead of internal `showBuildInfo()`
func WithCustomShowBuildInfo(fn func()) ExecOption {
	return func(w *ExecWorker) {
		w.globalShowBuildInfo = fn
	}
}

// WithNoDefaultHelpScreen true to disable printing help screen if without any arguments
func WithNoDefaultHelpScreen(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.noDefaultHelpScreen = b
	}
}

// WithHelpPainter allows to change the behavior and facade of help screen.
func WithHelpPainter(painter Painter) ExecOption {
	return func(w *ExecWorker) {
		w.currentHelpPainter = painter
	}
}

// WithHelpTabStop sets the tab-stop position in the help screen
// Default tabstop is 48
//
// Deprecated since v1.7.8 because the tab-stop position will be autosize from then on.
func WithHelpTabStop(tabStop int) ExecOption {
	return func(w *ExecWorker) {
		initTabStop(tabStop)
	}
}

// WithUnknownOptionHandler enables your customized wrong command/flag processor.
// internal processor supports smart suggestions for those wrong commands and flags.
func WithUnknownOptionHandler(handler UnknownOptionHandler) ExecOption {
	return func(w *ExecWorker) {
		unknownOptionHandler = handler
	}
}

// WithSimilarThreshold defines a threshold for command/option similar detector.
// Default threshold is 0.6666666666666666.
// See also JaroWinklerDistance
func WithSimilarThreshold(similarThreshold float64) ExecOption {
	return func(w *ExecWorker) {
		w.similarThreshold = similarThreshold
	}
}

// WithNoColor make console outputs plain and without ANSI escape colors
//
// Since v1.6.2+
func WithNoColor(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.noColor = b
	}
}

// WithNoEnvOverrides enables the internal no-env-overrides mode
//
// Since v1.6.2+
//
// In this mode, cmdr do NOT find and transfer equivalent envvar
// value into cmdr options store.
func WithNoEnvOverrides(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.noEnvOverrides = b
	}
}

// WithStrictMode enables the internal strict mode
//
// Since v1.6.2+
//
// In this mode, any warnings will be treat as an error and cause app
// fatal exit.
//
// In normal mode, these cases are assumed as warnings:
// - flag name not found
// - command or sub-command name not found
// - value extracting failed
// - ...
func WithStrictMode(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.strictMode = b
	}
}

// WithAfterArgsParsed sets a callback point after command-line args parsed by cmdr internal exec().
//
// Your callback func will be invoked before invoking the matched command `cmd`.
// At this time, all command-line args parsed and a command found.
//
// If program was launched with empty or wrong arguments, your callback func won't be triggered.
//
// When empty argument or `--help` found, cmdr will display help screen. To customize it
// see also cmdr.WithCustomShowVersion and cmdr.WithCustomShowBuildInfo.
//
// When any wrong/warn arguments found, cmdr will display some tip messages. To customize it
// see also cmdr.WithUnknownOptionHandler.
//
func WithAfterArgsParsed(hookFunc Handler) ExecOption {
	return func(w *ExecWorker) {
		w.afterArgsParsed = hookFunc
	}
}

// WithHelpTailLine setup the tail line in help screen
//
// Default line is:
//   "\nType '-h' or '--help' to get command help screen."
func WithHelpTailLine(line string) ExecOption {
	return func(w *ExecWorker) {
		w.helpTailLine = line
	}
}

// WithUnhandledErrorHandler handle the panics or exceptions generally
func WithUnhandledErrorHandler(handler UnhandledErrorHandler) ExecOption {
	return func(w *ExecWorker) {
		unhandledErrorHandler = handler
	}
}

type (
	// UnhandledErrorHandler for WithUnhandledErrorHandler
	UnhandledErrorHandler func(err interface{})
)

var (
	unhandledErrorHandler UnhandledErrorHandler
)

// WithNoCommandAction do NOT run the action of the matched command.
func WithNoCommandAction(b bool) ExecOption {
	return func(w *ExecWorker) {
		w.noCommandAction = b
	}
}

// WithOnSwitchCharHit handle the exact single switch-char (such as '-', '/', '~') matched.
// For example, type `bin/fluent mx -d - --help` will trigger this callback at the 2nd flag '-'.
func WithOnSwitchCharHit(fn func(parsed *Command, switchChar string, args []string) (err error)) ExecOption {
	return func(w *ExecWorker) {
		w.onSwitchCharHit = fn
	}
}

// WithOnPassThruCharHit handle the passthrough char(s) (i.e. '--') matched
// For example, type `bin/fluent mx -d -- --help` will trigger this callback at the 2nd flag '--'.
func WithOnPassThruCharHit(fn func(parsed *Command, switchChar string, args []string) (err error)) ExecOption {
	return func(w *ExecWorker) {
		w.onPassThruCharHit = fn
	}
}
