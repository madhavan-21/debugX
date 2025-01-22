package debugX

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// TODO: add give the overall function execution flow tree outputs
// TODO: give least recently given function execution flow tree output

// Log levels
const (
	info         = "INFO"
	err          = "ERROR"
	warn         = "WARN"
	res          = "RES"
	all          = "ALL"
	flow_checker = "FLOW_CHECKER"
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

var (
	callTree   = make(map[string][]string)
	callTreeMu sync.Mutex
	rootFunc   string
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
			case info:
				infoEnabled = true
			case err:
				errorEnabled = true
			case warn:
				warnEnabled = true
			case res:
				resEnabled = true
			case flow_checker:
				flowCheckerEnabled = true
			case all:
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
		formattedMessage := fmt.Sprintf(format, args...)
		file, line := getCallerInfo()
		time := getTimestamp()
		fmt.Println(timePrinter(fmt.Sprintf("Time: %s", time)), filePrinter(fmt.Sprintf("File: %s", file)), linePrinter(fmt.Sprintf("Line:%d", line)), infoPrinter(fmt.Sprintf("%s: %s", "INFO", formattedMessage)))

	}
}

// Log an error message if enabled
func Error(format string, args ...interface{}) {
	if allDebugEnabled || errorEnabled {
		formattedMessage := fmt.Sprintf(format, args...)
		file, line := getCallerInfo()
		time := getTimestamp()
		fmt.Println(timePrinter(fmt.Sprintf("Time: %s", time)), filePrinter(fmt.Sprintf("File: %s", file)), linePrinter(fmt.Sprintf("Line:%d", line)), errorPrinter(fmt.Sprintf("%s: %s", "ERROR", formattedMessage)))
	}
}

// Log a warning message if enabled
func Warn(format string, args ...interface{}) {
	if allDebugEnabled || warnEnabled {
		formattedMessage := fmt.Sprintf(format, args...)
		file, line := getCallerInfo()
		time := getTimestamp()
		fmt.Println(timePrinter(fmt.Sprintf("Time: %s", time)), filePrinter(fmt.Sprintf("File: %s", file)), linePrinter(fmt.Sprintf("Line:%d", line)), warnPrinter(fmt.Sprintf("%s: %s", "INFO", formattedMessage)))
	}
}

// Log a result message if enabled
func Res(format string, args ...interface{}) {
	if allDebugEnabled || resEnabled {
		formattedMessage := fmt.Sprintf(format, args...)
		file, line := getCallerInfo()
		time := getTimestamp()
		fmt.Println(timePrinter(fmt.Sprintf("Time: %s", time)), filePrinter(fmt.Sprintf("File: %s", file)), linePrinter(fmt.Sprintf("Line:%d", line)), resPrinter(fmt.Sprintf("%s: %s", "INFO", formattedMessage)))
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
		flowFunction := fmt.Sprintf("Function %s entered with parammeter: %s", function, paramString)
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

// Function to capture the current function's name
func getCurrentFunctionName(skip int) string {
	pc, _, _, _ := runtime.Caller(skip)
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// Function to log function calls dynamically
func logFunctionCall(parent, child string) {
	callTreeMu.Lock()
	defer callTreeMu.Unlock()
	callTree[parent] = append(callTree[parent], child)
}

// Dynamic wrapper for root function
func InvokeAndTrack(f func()) string {
	// Set the root function
	rootFunc = getCurrentFunctionName(2)

	// Wrap and track function execution
	callTree = make(map[string][]string)
	wrappedFunc := func() {
		trackFunctionCalls(func() {
			defer func() { recover() }() // Recover from errors
			f()
		})
	}
	wrappedFunc()

	// Generate and return the flow tree
	return generateFlowTree(rootFunc)
}

// Function to track calls dynamically using runtime stack
func trackFunctionCalls(f func()) {
	pc, _, _, _ := runtime.Caller(1) // Get caller info
	_ = runtime.FuncForPC(pc).Name()
	callStack := []string{}

	defer func() {
		for i := 1; ; i++ {
			pc, _, _, ok := runtime.Caller(i)
			if !ok {
				break
			}
			fn := runtime.FuncForPC(pc).Name()
			callStack = append(callStack, fn)
		}
		for i := len(callStack) - 1; i > 0; i-- {
			logFunctionCall(callStack[i], callStack[i-1])
		}
	}()

	f() // Execute the root function
}

// Function to generate a flow tree as a string
func generateFlowTree(root string) string {
	callTreeMu.Lock()
	defer callTreeMu.Unlock()

	visited := make(map[string]bool)
	var result string
	var buildTree func(node string, depth int)
	buildTree = func(node string, depth int) {
		if visited[node] {
			return
		}
		visited[node] = true
		indentation := ""
		for i := 0; i < depth; i++ {
			indentation += "  "
		}
		result += fmt.Sprintf("%s- %s\n", indentation, node)
		for _, child := range callTree[node] {
			buildTree(child, depth+1)
		}
	}
	buildTree(root, 0)
	return result
}
