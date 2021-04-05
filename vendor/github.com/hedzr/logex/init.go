/*
 * Copyright © 2019 Hedzr Yeh.
 */

package logex

import (
	"github.com/hedzr/log"
	"github.com/hedzr/logex/formatter"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

func GetLevel() log.Level {
	return log.GetLevel()
}

func Enable() {
	logrus.SetFormatter(&formatter.TextFormatter{ForceColors: true})
	logrus.SetReportCaller(true)
	// logrus.AddHook(logex.hook.DefaultContextHook)
}

func EnableWith(lvl log.Level, opts ...Option) {
	if lvl == log.OffLevel {
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.SetOutput(ioutil.Discard)
	} else {
		logrus.SetLevel(logrus.Level(lvl))
		logrus.SetOutput(os.Stdout)
	}
	log.SetLevel(lvl)
	logrus.SetFormatter(&formatter.TextFormatter{ForceColors: true})
	logrus.SetReportCaller(true)
	// logrus.AddHook(logex.hook.DefaultContextHook)
	for _, opt := range opts {
		opt()
	}
}

func SetupLoggingFormat(format string, logexSkipFrames int) {
	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:  "2006-01-02 15:04:05.000",
			DisableTimestamp: false,
			PrettyPrint:      false,
		})
	default:
		e := false
		if logexSkipFrames > 0 {
			e = true
		}
		logrus.SetFormatter(&formatter.TextFormatter{
			ForceColors:               true,
			DisableColors:             false,
			FullTimestamp:             true,
			TimestampFormat:           "2006-01-02 15:04:05.000",
			Skip:                      logexSkipFrames,
			EnableSkip:                e,
			EnvironmentOverrideColors: true,
		})
	}
	if GetLevel() == log.OffLevel {
		logrus.SetLevel(logrus.ErrorLevel)
		logrus.SetOutput(ioutil.Discard)
	}
}

type Option func()

const SKIP = formatter.SKIP

// var level log.Level

// func Enable() {
//
// 	var foreground = vxconf.GetBoolR("app.foreground")
// 	var file = daemon.DefaultLogFile()
// 	var lvl = vxconf.GetStringR("app.logger.level")
//
// 	var target = vxconf.GetStringR("app.logger.target")
// 	var format = vxconf.GetStringR("app.logger.format")
// 	if len(target) == 0 {
// 		target = "default"
// 	}
// 	if len(format) == 0 {
// 		format = "text"
// 	}
// 	if target == "journal" {
// 		format = "text"
// 	}
// 	switch format {
// 	case "json":
// 		logrus.SetFormatter(&logrus.JSONFormatter{})
// 	default:
// 		logrus.SetFormatter(&fmtr.TextFormatter{ForceColors: true})
// 	}
// 	// Log as JSON instead of the default ASCII formatter.
//
// 	can_use_log_file, journal_mode := ij(target, foreground)
// 	l, _ := logrus.ParseLevel(lvl)
// 	if cli_common.Debug && l < logrus.DebugLevel {
// 		l = logrus.DebugLevel
// 	}
// 	logrus.SetLevel(l)
// 	logrus.Debugf("Using logger: format=%v, lvl=%v/%v, target=%v, journal.mode=%v, using.logrus.file=%v", format, lvl, l, target, journal_mode, can_use_log_file)
//
// 	if foreground == false && can_use_log_file {
// 		if len(file) == 0 {
// 			file = os.DevNull // "/dev/null"
// 		}
//
// 		logFile, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0440)
// 		if err != nil {
// 			logFile, err = os.Create(file)
// 			if err != nil {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			} else {
// 				fmt.Printf("[FG] Using new log file: %s\n\n", file)
// 			}
// 		} else {
// 			fmt.Printf("[FG] Using exists log file: %s\n\n", file)
// 		}
//
// 		// logrus.Infof("Using log file: %s", file)
// 		// fmt.Printf("Using log file: %s\n\n", file)
//
// 		// Output to stdout instead of the default stderr
// 		// Can be any io.Writer, see below for File example
// 		// logrus.SetOutput(os.Stdout)
// 		logrus.SetOutput(logFile)
// 	} else {
// 		logrus.SetReportCaller(true)
// 		// logrus.AddHook(DefaultContextHook)
// 	}
//
// 	// if hook, err := logrus_syslogrus.NewSyslogHook("udp", "localhost:514", syslogrus.LOG_INFO, ""); err == nil { logrus.Hooks.Add(hook) }
//
// 	// var lvl = vxconf.GetStringR("app.logger.level")
// 	// var file = daemon.DefaultLogFile()
// 	// var err error
// 	//
// 	// logFile, err := os.OpenFile(file, os.O_WRONLY,0400)
// 	// if err != nil{
// 	// 	fmt.Println(err)
// 	// }
// 	//
// 	// lf, err := logging.LogLevel(lvl)
// 	// backend1 := logging.NewLogBackend(logFile, "", 0)
// 	// backend1Leveled := logging.AddModuleLevel(backend1)
// 	// backend1Leveled.SetLevel(lf, "")
// 	//
// 	// if foreground {
// 	// 	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
// 	// 	backend2Formatter := logging.NewBackendFormatter(backend2, format)
// 	// 	logging.SetBackend(backend1Leveled, backend2Formatter)
// 	// }
//
// 	// logrus.Debugf("debug %s", Password("secret"))
// 	// logrus.Info("info")
// 	// logrus.Notice("notice")
// 	// logrus.Warning("warning")
// 	// logrus.Error("suweia.com")
// 	// logrus.Critical("太严重了")
// 	// os.Exit(0)
//
// 	// etcd.Warn("sssss")
// 	// etcd.Warn("sssss")
// 	// etcd.Warn("sssss")
// 	// etcd.Warn("sssss")
// }
