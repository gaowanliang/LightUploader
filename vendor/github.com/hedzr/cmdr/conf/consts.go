/*
 * Copyright © 2019 Hedzr Yeh.
 */

// Package conf are used to store the app-level constants (app name/vaersion) for cmdr and your app
package conf

var (
	// CfgFile never used
	CfgFile string
	// AppName app name
	AppName string

	// these 3 variables will be rewrote when app had been building by ci-tool

	// Version app version
	Version = "0.2.1"
	// Buildstamp app built stamp
	Buildstamp = ""
	// Githash app git hash
	Githash = ""
	// GitShortVersion from `git describe --long` [NEVER USED]
	GitShortVersion = ""
	// GoVersion `go version` string
	GoVersion = ""

	// ServerTag app server tag names
	ServerTag = ""
	// ServerID app server id
	ServerID = ""
)
