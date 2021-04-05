// Copyright © 2020 Hedzr Yeh.

package cmdr

func defaultOnSwitchCharHit(parsed *Command, switchChar string, args []string) (err error) {
	if parsed != nil {
		printInDevMode("the last parsed command is %q - %q\n", parsed.GetTitleNames(), parsed.Description)
	}
	printInDevMode("SwitchChar FOUND: %v\nremains: %v\n\n", switchChar, args)
	return nil // ErrShouldBeStopException
}

func defaultOnPassThruCharHit(parsed *Command, switchChar string, args []string) (err error) {
	if parsed != nil {
		printInDevMode("the last parsed command is %q - %q\n", parsed.GetTitleNames(), parsed.Description)
	}
	printInDevMode("PassThrough flag FOUND: %v\nremains: %v\n\n", switchChar, args)
	return nil // ErrShouldBeStopException
}

func emptyUnknownOptionHandler(isFlag bool, title string, cmd *Command, args []string) (fallbackToDefaultDetector bool) {
	return false
}
