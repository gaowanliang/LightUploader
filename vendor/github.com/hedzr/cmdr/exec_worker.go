package cmdr

import (
	"bufio"
	cmdrBase "github.com/hedzr/cmdr-base"
	"os"
	"runtime"
	"sync"
)

// ExecWorker is a core logic worker and holder
type ExecWorker struct {
	switchCharset string

	// beforeXrefBuildingX, afterXrefBuiltX HookFunc
	beforeXrefBuilding      []HookFunc
	beforeConfigFileLoading []HookFunc
	afterConfigFileLoading  []HookFunc
	afterXrefBuilt          []HookFunc
	afterAutomaticEnv       []HookOptsFunc
	beforeHelpScreen        []HookHelpScreenFunc
	afterHelpScreen         []HookHelpScreenFunc

	envPrefixes         []string
	rxxtPrefixes        []string
	predefinedLocations []string // predefined config file locations
	alterLocations      []string // alter config file locations, so we can write back the changes
	pluginsLocations    []string
	extensionsLocations []string

	shouldIgnoreWrongEnumValue bool

	enableVersionCommands     bool
	enableHelpCommands        bool
	enableVerboseCommands     bool
	enableCmdrCommands        bool
	enableGenerateCommands    bool
	treatUnknownCommandAsArgs bool

	watchMainConfigFileToo   bool
	doNotLoadingConfigFiles  bool
	doNotWatchingConfigFiles bool
	confDFolderName          string
	watchChildConfigFiles    bool

	globalShowVersion   func()
	globalShowBuildInfo func()

	currentHelpPainter Painter

	bufferedStdio bool
	defaultStdout *bufio.Writer
	defaultStderr *bufio.Writer
	closers       []func()

	// rootCommand the root of all commands
	rootCommand *RootCommand
	// rootOptions *Opt
	rxxtOptions        *Options
	onOptionMergingSet OnOptionSetCB
	onOptionSet        OnOptionSetCB

	similarThreshold      float64
	noDefaultHelpScreen   bool
	noColor               bool
	noEnvOverrides        bool
	strictMode            bool
	noUnknownCmdTip       bool
	noCommandAction       bool
	noPluggableAddons     bool
	noPluggableExtensions bool

	logexInitialFunctor Handler
	logexPrefix         string
	logexSkipFrames     int

	afterArgsParsed Handler

	envVarToValueMap map[string]func() string

	helpTailLine string

	onSwitchCharHit   OnSwitchCharHitCB
	onPassThruCharHit OnPassThruCharHitCB

	addons []cmdrBase.PluginEntry

	lastPkg *ptpkg
}

// ExecOption is the functional option for Exec()
type ExecOption func(w *ExecWorker)

func internalGetWorker() (w *ExecWorker) {
	uniqueWorkerLock.RLock()
	w = uniqueWorker
	uniqueWorkerLock.RUnlock()
	return
}

func internalResetWorkerNoLock() (w *ExecWorker) {
	w = &ExecWorker{
		switchCharset: "-~",

		envPrefixes:  []string{"CMDR"},
		rxxtPrefixes: []string{"app"},

		predefinedLocations: []string{
			"./ci/etc/$APPNAME/$APPNAME.yml",       // for developer
			"/etc/$APPNAME/$APPNAME.yml",           // regular location
			"/usr/local/etc/$APPNAME/$APPNAME.yml", // regular macOS HomeBrew location
			"/opt/etc/$APPNAME/$APPNAME.yml",       // regular location
			"/var/lib/etc/$APPNAME/$APPNAME.yml",   // regular location
			"$HOME/.config/$APPNAME/$APPNAME.yml",  // per user
			"$HOME/.$APPNAME/$APPNAME.yml",         // ext location per user
			// "$XDG_CONFIG_HOME/$APPNAME/$APPNAME.yml", // ?? seldom defined | generally it's $HOME/.config
			"$THIS/$APPNAME.yml", // executable's directory
			"$APPNAME.yml",       // current directory
			// "./ci/etc/%s/%s.yml",
			// "/etc/%s/%s.yml",
			// "/usr/local/etc/%s/%s.yml",
			// "$HOME/.%s/%s.yml",
			// "$HOME/.config/%s/%s.yml",
		},

		alterLocations: []string{
			"./bin/$APPNAME.yml", // for developer, current bin directory
			"/var/lib/$APPNAME",  //
			"$THIS/$APPNAME.yml", // executable's directory
		},

		pluginsLocations: []string{
			"./ci/local/share/$APPNAME/addons",
			"$HOME/.local/share/$APPNAME/addons",
			"$HOME/.$APPNAME/addons",
			"/usr/local/share/$APPNAME/addons",
			"/usr/share/$APPNAME/addons",
		},
		extensionsLocations: []string{
			"./ci/local/share/$APPNAME/ext",
			"$HOME/.local/share/$APPNAME/ext",
			"$HOME/.$APPNAME/ext",
			"/usr/local/share/$APPNAME/ext",
			"/usr/share/$APPNAME/ext",
		},

		shouldIgnoreWrongEnumValue: true,

		enableVersionCommands:     true,
		enableHelpCommands:        true,
		enableVerboseCommands:     true,
		enableCmdrCommands:        true,
		enableGenerateCommands:    true,
		treatUnknownCommandAsArgs: true,

		doNotLoadingConfigFiles: false,

		currentHelpPainter: new(helpPainter),

		defaultStdout: bufio.NewWriterSize(os.Stdout, 16384),
		defaultStderr: bufio.NewWriterSize(os.Stderr, 16384),

		rxxtOptions: newOptions(),

		similarThreshold:    similarThreshold,
		noDefaultHelpScreen: false,

		helpTailLine: defaultTailLine,

		confDFolderName: confDFolderNameConst,
	}

	WithEnvVarMap(nil)(w)

	if runtime.GOOS == "windows" {
		w.switchCharset = "-/~"
	}

	uniqueWorker = w
	return
}

func init() {
	_ = internalResetWorkerNoLock()
}

var uniqueWorkerLock sync.RWMutex
var uniqueWorker *ExecWorker
var noResetWorker = true

const confDFolderNameConst = "conf.d"
