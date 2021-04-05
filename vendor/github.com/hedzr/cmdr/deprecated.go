/*
 * Copyright © 2019 Hedzr Yeh.
 */

package cmdr

var (
// // EnableVersionCommands supports injecting the default `--version` flags and commands
// //
// // Deprecated: from v1.5.0
// EnableVersionCommands = true
// // EnableHelpCommands supports injecting the default `--help` flags and commands
// //
// // Deprecated: from v1.5.0
// EnableHelpCommands = true
// // EnableVerboseCommands supports injecting the default `--verbose` flags and commands
// //
// // Deprecated: from v1.5.0
// EnableVerboseCommands = true
// // EnableCmdrCommands support these flags: `--strict-mode`, `--no-env-overrides`
// //
// // Deprecated: from v1.5.0
// EnableCmdrCommands = true
// // EnableGenerateCommands supports injecting the default `generate` commands and subcommands
// //
// // Deprecated: from v1.5.0
// EnableGenerateCommands = true
//
// // EnvPrefix attaches a prefix to key to retrieve the option value.
// //
// // Deprecated: from v1.5.0
// EnvPrefix = []string{"CMDR"}
//
// // RxxtPrefix create a top-level namespace, which contains all normalized `Flag`s.
// //
// // Deprecated: from v1.5.0
// RxxtPrefix = []string{"app"}
//
// // ShouldIgnoreWrongEnumValue will be put into `cmdrError.Ignorable` while wrong enumerable value found in parsing command-line options.
// // main program might decide whether it's a warning or error.
// // see also: [Flag.ValidArgs]
// //
// // Deprecated: from v1.5.0
// ShouldIgnoreWrongEnumValue = false
)

// // AddOnBeforeXrefBuilding add hook func
// //
// // Deprecated: from v1.5.0
// func AddOnBeforeXrefBuilding(cb HookFunc) {
// 	uniqueWorker.AddOnBeforeXrefBuilding(cb)
// }
//
// // AddOnAfterXrefBuilt add hook func
// //
// // Deprecated: from v1.5.0
// func AddOnAfterXrefBuilt(cb HookFunc) {
// 	uniqueWorker.AddOnAfterXrefBuilt(cb)
// }
//
// // SetInternalOutputStreams sets the internal output streams for debugging
// //
// // Deprecated: from v1.5.0
// func SetInternalOutputStreams(out, err *bufio.Writer) {
// 	uniqueWorker.defaultStdout = out
// 	uniqueWorker.defaultStderr = err
//
// 	if uniqueWorker.defaultStdout == nil {
// 		uniqueWorker.defaultStdout = bufio.NewWriterSize(os.Stdout, 16384)
// 	}
// 	if uniqueWorker.defaultStderr == nil {
// 		uniqueWorker.defaultStderr = bufio.NewWriterSize(os.Stderr, 16384)
// 	}
// }
//
// // SetCustomShowVersion supports your `ShowVersion()` instead of internal `showVersion()`
// //
// // Deprecated: from v1.5.0
// //
// func SetCustomShowVersion(fn func()) {
// 	uniqueWorker.globalShowVersion = fn
// }
//
// // SetCustomShowBuildInfo supports your `ShowBuildInfo()` instead of internal `showBuildInfo()`
// //
// // Deprecated: from v1.5.0
// func SetCustomShowBuildInfo(fn func()) {
// 	uniqueWorker.globalShowBuildInfo = fn
// }
//
// // PrintBuildInfo print building information
// //
// // Deprecated: from v1.5.0
// func PrintBuildInfo() {
// 	uniqueWorker.showBuildInfo()
// }
//
// // SetNoLoadConfigFiles true means no loading config files
// //
// // Deprecated: from v1.5.0
// func SetNoLoadConfigFiles(b bool) {
// 	uniqueWorker.doNotLoadingConfigFiles = b
// }
//
// // SetCurrentHelpPainter allows to change the behavior and facade of help screen.
// //
// // Deprecated: from v1.5.0
// func SetCurrentHelpPainter(painter Painter) {
// 	uniqueWorker.currentHelpPainter = painter
// }
//
// // SetHelpTabStop sets the tab stop for help screen output.
// //
// // Deprecated: from v1.5.0, replaced with WithHelpTabStop(tabStop).
// func SetHelpTabStop(tabStop int) {
// 	initTabStop(tabStop)
// }
//
// // ExecWith is main entry of `cmdr`.
// //
// // Deprecated: from v1.5.0
// func ExecWith(rootCmdForTesting *RootCommand, beforeXrefBuildingX, afterXrefBuiltX HookFunc) (err error) {
// 	w := uniqueWorker
//
// 	if beforeXrefBuildingX != nil {
// 		w.beforeXrefBuilding = append(w.beforeXrefBuilding, beforeXrefBuildingX)
// 	}
// 	if afterXrefBuiltX != nil {
// 		w.afterXrefBuilt = append(w.afterXrefBuilt, afterXrefBuiltX)
// 	}
//
// 	err = w.InternalExecFor(rootCmdForTesting, os.Args)
// 	return
// }
