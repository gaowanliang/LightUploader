package cmdrplugin

// PluginEntry entry of an addon (golang plugin)
type PluginEntry interface {
	PluginCmd
	AddonTitle() string
	AddonDescription() string
	AddonCopyright() string
	AddonVersion() string
}

// PluginCompBase component for cmd and flag of an addon
type PluginCompBase interface {
	Name() string
	ShortName() string
	Aliases() []string
	Description() string
}

// PluginCmd a command of an addon
type PluginCmd interface {
	PluginCompBase
	SubCommands() []PluginCmd
	Flags() []PluginFlag
	Action(args []string) (err error)
}

// PluginFlag a flag of a command of an addon
type PluginFlag interface {
	PluginCompBase
	DefaultValue() interface{}
	PlaceHolder() string
	Action() (err error) // onSet
}
