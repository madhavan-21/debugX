package debugX

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

// TODO: add give the overall function execution flow tree outputs
// TODO: give least recently given function execution flow tree output
// Log levels
const (
	INFO         = "INFO"
	ERROR        = "ERROR"
	WARN         = "WARN"
	RES          = "RES"
	ALL          = "ALL"
	FLOW_CHECKER = "FLOW_CHECKER"
)

// Log level enable states
var (
	allDebugEnabled    = true
	infoEnabled        = false
	errorEnabled       = false
	warnEnabled        = false
	resEnabled         = false
	flowCheckerEnabled = false
)

// Color printers for different log levels
var (
	infoPrinter  = color.New(color.FgBlue).Add(color.Bold).SprintFunc()
	errorPrinter = color.New(color.FgRed).Add(color.Bold).SprintFunc()
	warnPrinter  = color.New(color.FgYellow).Add(color.Bold).SprintFunc()
	resPrinter   = color.New(color.FgGreen).Add(color.Bold).SprintFunc()
	flowPrinter  = color.New(color.FgGreen).Add(color.Bold).Add(color.BgHiWhite).SprintFunc()
	timePrinter  = color.New(color.FgMagenta).Add(color.Bold).Add(color.BgHiBlue).SprintFunc()
	linePrinter  = color.New(color.FgHiCyan).Add(color.Bold).Add(color.BgHiYellow).SprintFunc()
	filePrinter  = color.New(color.FgHiWhite).Add(color.BgHiRed).SprintFunc()
)

// Get the current timestamp
func getTimestamp() string {
	return time.Now().Format("15:04:05")
}

// Get the caller's file name and line number
func getCallerInfo() (string, int) {
	_, file, line, _ := runtime.Caller(2)
	return filepath.Base(file), line
}

// Control which log levels are enabled or disabled
func DebugControl(levels ...string) {
	if levels == nil {
		return
	}
	if len(levels) > 0 {
		allDebugEnabled = false
		infoEnabled, errorEnabled, warnEnabled, resEnabled, flowCheckerEnabled = false, false, false, false, false
		for _, l := range levels {
			l = strings.ToUpper(l)
			switch l {
			case INFO:
				infoEnabled = true
			case ERROR:
				errorEnabled = true
			case WARN:
				warnEnabled = true
			case RES:
				resEnabled = true
			case FLOW_CHECKER:
				flowCheckerEnabled = true
			default:
			case ALL:
				// If 'ALL' is explicitly passed, disable all debug levels
				allDebugEnabled = false
				infoEnabled, errorEnabled, warnEnabled, resEnabled, flowCheckerEnabled = false, false, false, false, false
			}
		}
	}
}

// Log an info message if enabled
func Info(format string, args ...interface{}) {
	if allDebugEnabled || infoEnabled {
		logMessage(INFO, infoPrinter, format, args...)
	}
}

// Log an error message if enabled
func Error(format string, args ...interface{}) {
	if allDebugEnabled || errorEnabled {
		logMessage(ERROR, errorPrinter, format, args...)
	}
}

// Log a warning message if enabled
func Warn(format string, args ...interface{}) {
	if allDebugEnabled || warnEnabled {
		logMessage(WARN, warnPrinter, format, args...)
	}
}

// Log a result message if enabled
func Res(format string, args ...interface{}) {
	if allDebugEnabled || resEnabled {
		logMessage(RES, resPrinter, format, args...)
	}
}

// Log a message with the appropriate level and color
func logMessage(level string, colorPrinter func(a ...interface{}) string, format string, args ...interface{}) {
	formattedMessage := fmt.Sprintf(format, args...)
	file, line := getCallerInfo()
	time := getTimestamp()
	fmt.Println(timePrinter(fmt.Sprintf("Time: %s", time)), filePrinter(fmt.Sprintf("File: %s", file)), linePrinter(fmt.Sprintf("Line:%d", line)), colorPrinter(fmt.Sprintf("%s: %s", level, formattedMessage)))
}

// FlowChecker logs entry and exit points of functions
func FlowChecker(params ...interface{}) func() {
	if allDebugEnabled || flowCheckerEnabled {
		pc, file, line, _ := runtime.Caller(1)
		function := filepath.Base(runtime.FuncForPC(pc).Name())
		shortFile := filepath.Base(file)
		time := getTimestamp()

		paramString := fmt.Sprintf("%v", params)
		flowFunction := fmt.Sprintf("Function %s entered with params: %s", function, paramString)
		fmt.Println(timePrinter(fmt.Sprintf("Time: %s", time)), filePrinter(fmt.Sprintf("File: %s", shortFile)), linePrinter(fmt.Sprintf("Line: %d", line)), flowPrinter(flowFunction))

		return func() {
			pc, file, line, _ := runtime.Caller(1)
			function := filepath.Base(runtime.FuncForPC(pc).Name())
			shortFile := filepath.Base(file)
			time := getTimestamp()
			flowFunction := fmt.Sprintf("Function %s exited", function)
			fmt.Println(timePrinter(fmt.Sprintf("Time: %s", time)), filePrinter(fmt.Sprintf("File: %s", shortFile)), linePrinter(fmt.Sprintf("Line: %d", line)), flowPrinter(flowFunction))
		}
	}

	return func() {}
}
