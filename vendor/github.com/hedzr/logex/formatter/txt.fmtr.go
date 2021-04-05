/*
 * Copyright © 2019 Hedzr Yeh.
 */

package formatter

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	defaultTimestampFormat = time.RFC3339

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97

	bgNormal       = 0
	bgBoldOrBright = 1
	bgDim          = 2
	bgUnderline    = 4
	bgUlink        = 5
	bgHidden       = 8

	darkColor = lightGray

	SKIP = "SKIP"
)

var baseTimestamp time.Time

func init() {
	baseTimestamp = time.Now()
}

// TextFormatter formats logs into text
type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &TextFormatter{
	//     FieldMap: FieldMap{
	//         FieldKeyTime:  "@timestamp",
	//         FieldKeyLevel: "@level",
	//         FieldKeyMsg:   "@message"}}
	FieldMap logrus.FieldMap

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	Skip       int
	EnableSkip bool

	terminalInitOnce sync.Once
}

func (f *TextFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)

		if f.isTerminal {
			initTerminal(entry.Logger.Out)
		}
	}
}

func (f *TextFormatter) isColored() bool {
	isColored := f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows"))

	if f.EnvironmentOverrideColors {
		if force, ok := os.LookupEnv("CLICOLOR_FORCE"); ok && force != "0" {
			isColored = true
		} else if ok && force == "0" {
			isColored = false
		} else if os.Getenv("CLICOLOR") == "0" {
			isColored = false
		}
	}

	return isColored && !f.DisableColors
}

var (
	callerInitOnce     sync.Once
	logrusPackage      string
	logexPackage       string
	logPackage         string
	minimumCallerDepth int
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

// getPackageName reduces a fully qualified function name to the package name
// There really ought to be to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller(skipFrames int) *runtime.Frame {

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, 2)
		_ = runtime.Callers(0, pcs)
		logPackage = "github.com/hedzr/log"
		logexPackage = getPackageName(runtime.FuncForPC(pcs[1]).Name())
		logrusPackage = "github.com/sirupsen/logrus"

		// now that we have the cache, we can skip a minimum count of known-logrus functions
		// XXX this is dubious, the number of frames may vary
		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth /*+skipFrames*/, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	skipped := 0
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logrusPackage && pkg != logexPackage && pkg != logPackage {
			if skipped < skipFrames {
				skipped++
				continue
			}
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	var funcVal, fileVal string

	fixedKeys := make([]string, 0, 4+len(data))
	if !f.DisableTimestamp {
		fixedKeys = append(fixedKeys, resolve(f.FieldMap, logrus.FieldKeyTime))
	}
	fixedKeys = append(fixedKeys, resolve(f.FieldMap, logrus.FieldKeyLevel))
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, resolve(f.FieldMap, logrus.FieldKeyMsg))
	}
	// NOTICE if entry.err != "" {
	// 	fixedKeys = append(fixedKeys, resolve(f.FieldMap, logrus.FieldKeyLogrusError))
	// }
	if entry.HasCaller() {
		fixedKeys = append(fixedKeys,
			resolve(f.FieldMap, logrus.FieldKeyFunc), resolve(f.FieldMap, logrus.FieldKeyFile))
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		} else {
			funcVal = entry.Caller.Function
			if ff, ok := data[resolve(f.FieldMap, logrus.FieldKeyFile)]; ok && strings.Contains(entry.Caller.File, "logrus.go") {
				fileVal = ff.(string)
			} else {
				var sf interface{}
				ok, fb := false, false
				if f.EnableSkip && f.Skip > 0 {
					sf, ok = f.Skip, true
					// } else if !f.EnableSkip {
					//	 sf, ok = 1, true
				}
				if v, yes := data[SKIP]; yes {
					sf, ok = sf.(int)+v.(int), yes
				}
				if ok {
					if skipFrames, ok := sf.(int); ok && skipFrames > 0 {
						// println("skipFrames: %v", skipFrames)
						fb = true
						delete(data, SKIP)
						for i, k := range keys {
							if k == SKIP {
								keys[i] = "via"
							}
						}
						// data["via"] = fmt.Sprintf("%s:%d (%v)", entry.Caller.File, entry.Caller.Line, funcVal)
						data["via"] = fmt.Sprintf("%v", funcVal)

						entryCaller := getCaller(skipFrames)
						fileVal = fmt.Sprintf("%s:%d", entryCaller.File, entryCaller.Line)
						// funcVal = entryCaller.Function
						entry.Caller = entryCaller
					}
				}
				// fmt.Printf("runtime-version: %v, skipFrames: %v. ok: %v, EnableSkip: %v, %v. fb: %v. logexPackage: %v.\n", runtime.Version(), sf, ok, f.EnableSkip, f.Skip, fb, logexPackage)
				if !fb {
					fileVal = fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
				}
			}
		}
	}

	if !f.DisableSorting {
		if f.SortingFunc == nil {
			sort.Strings(keys)
			fixedKeys = append(fixedKeys, keys...)
		} else {
			if !f.isColored() {
				fixedKeys = append(fixedKeys, keys...)
				f.SortingFunc(fixedKeys)
			} else {
				f.SortingFunc(keys)
			}
		}
	} else {
		fixedKeys = append(fixedKeys, keys...)
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.terminalInitOnce.Do(func() { f.init(entry) })

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}
	if f.isColored() {
		f.printColored(b, entry, keys, data, timestampFormat)
	} else {

		for _, key := range fixedKeys {
			var value interface{}
			switch {
			case key == resolve(f.FieldMap, logrus.FieldKeyTime):
				value = entry.Time.Format(timestampFormat)
			case key == resolve(f.FieldMap, logrus.FieldKeyLevel):
				value = entry.Level.String()
			case key == resolve(f.FieldMap, logrus.FieldKeyMsg):
				value = entry.Message
			// NOTICE case key == resolve(f.FieldMap, logrus.FieldKeyLogrusError):
			// 	value = entry.err
			case key == resolve(f.FieldMap, logrus.FieldKeyFunc) && entry.HasCaller():
				value = funcVal
			case key == resolve(f.FieldMap, logrus.FieldKeyFile) && entry.HasCaller():
				value = fileVal
			default:
				value = data[key]
			}
			f.appendKeyValue(b, key, value)
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, keys []string, data logrus.Fields, timestampFormat string) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = darkGray // lightGray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = cyan
	}

	levelText := strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation {
		levelText = levelText[0:4]
	}

	// Remove a single newline if it already exists in the message to keep
	// the behavior of logrus text_formatter the same as the stdlib log package
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	caller := ""
	skipFile := false

	if entry.HasCaller() {
		funcVal := fmt.Sprintf("\u001B[%dm%s\u001B[0m()", darkGray, entry.Caller.Function)
		var fileVal string // fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if f, ok := data[resolve(f.FieldMap, logrus.FieldKeyFile)]; ok && strings.Contains(entry.Caller.File, "logrus.go") {
			fileVal = f.(string)
			skipFile = true
		} else {
			fileVal = fmt.Sprintf("\u001B[%dm%s:%d\u001B[0m", lightBlue, entry.Caller.File, entry.Caller.Line)
		}

		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
		caller = fileVal + " " + funcVal
	}

	if f.DisableTimestamp {
		_, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m%-48s %s ", levelColor, levelText, entry.Message, caller)
	} else if !f.FullTimestamp {
		// echo -e "Normal \e[2mDim"
		_, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m\x1b[2m\x1b[%dm[%04d]\x1b[0m%-48s \x1b[2m\x1b[%dm%s\x1b[0m ",
			levelColor, levelText, darkColor, int(entry.Time.Sub(baseTimestamp)/time.Second), entry.Message, darkColor, caller)
	} else {
		_, _ = fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s]%-48s %s ", levelColor, levelText, entry.Time.Format(timestampFormat), entry.Message, caller)
	}
	for _, k := range keys {
		if skipFile {
			if strings.HasSuffix(k, resolve(f.FieldMap, logrus.FieldKeyFile)) {
				continue
			}
		}

		v := data[k]
		_, _ = fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
