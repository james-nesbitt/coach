package log

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Configuration struct for a CliLog
type CliLogSettings struct {
	writer    io.Writer // a log writing target
	stack     []string  // patent name stack
	verbosity int       // current verbosity for this log object
	hush      bool      // if true, and verbosity is standard, then make the log quieter
}

// CliLog default logging handler
type CliLog struct {
	CliLogSettings
}

func (log *CliLog) Name() string {
	return log.stack[len(log.stack)]
}
func (log *CliLog) Verbosity() int {
	return log.verbosity
}
func (log *CliLog) SetVerbosity(verbosity int) {
	log.verbosity = verbosity
}
func (log *CliLog) MakeChild(target string) Log {
	return Log(&CliLog{
		CliLogSettings: CliLogSettings{
			writer:    log.writer,
			stack:     append(log.stack, target),
			verbosity: log.verbosity,
			hush:      log.hush,
		},
	})
}

func (log *CliLog) IsHushed() bool {
	return log.hush
}
func (log *CliLog) Hush() {
	log.hush = true
}
func (log *CliLog) UnHush() {
	log.hush = false
}

// Implement a Critical error
func (log *CliLog) Fatal(messages ...string) {
	log.writeLog(VERBOSITY_FATAL, messages...)
	panic("Execution halted on FATAL error")
}

// Implement an error
func (log *CliLog) Error(messages ...string) {
	log.writeLog(VERBOSITY_ERROR, messages...)
}

// Register a Warning error
func (log *CliLog) Warning(messages ...string) {
	log.writeLog(VERBOSITY_WARNING, messages...)
}

// Register a message
func (log *CliLog) Message(messages ...string) {
	log.writeLog(VERBOSITY_MESSAGE, messages...)
}

// Register an information verbose message
func (log *CliLog) Info(messages ...string) {
	log.writeLog(VERBOSITY_INFO, messages...)
}

// Debug message and data
func (log *CliLog) Debug(verbosity int, message string, objects ...interface{}) {
	log.writeLog(verbosity, message)
	if +verbosity <= log.verbosity && len(objects) > 0 && objects[0] != nil {
		fmt.Print("	")
		fmt.Fprintln(log, objects...)
	}
}

// internal logging writer
func (log *CliLog) writeLog(verbosity int, messages ...string) {

	elements := []string{} // a slice of string message elements

	if log.hush {
		switch verbosity {
		case VERBOSITY_WARNING:
			elements = append(elements, "[HUSHED WARNING]")
			verbosity = VERBOSITY_INFO
		case VERBOSITY_MESSAGE:
			elements = append(elements, "[HUSHED MESSAGE]")
			verbosity = VERBOSITY_INFO
		}
	}

	if verbosity > log.verbosity {
		return
	}

	switch verbosity {
	case VERBOSITY_FATAL:
		elements = append(elements, "[FATAL]", log.joinStack())
	case VERBOSITY_ERROR:
		elements = append(elements, "[ERROR]", log.joinStack())

	case VERBOSITY_WARNING:
		elements = append(elements, "[WARNING]")

	case VERBOSITY_MESSAGE:
		prefix := log.stack[len(log.stack)-1] + ": "
		if length := utf8.RuneCountInString(prefix); length < 15 {
			prefix += strings.Repeat("-", 15-length)
		}
		elements = append(elements, prefix)
	case VERBOSITY_INFO:
		elements = append(elements, "-->")

	default:
		elements = append(elements, "("+strconv.Itoa(verbosity)+")", log.joinStack())
	}

	elements = append(elements, messages...)

	output := strings.Join(elements, " ") + "\n"
	log.Write([]byte(output))

}

// joins the log targets into a printable string for message prefixing
func (log *CliLog) joinStack() string {
	output := ""
	if len(log.stack) > 0 {
		output += "[" + strings.Join(log.stack, "][") + "]"
	}
	length := utf8.RuneCountInString(output)
	if length < 25 {
		output += strings.Repeat("-", 25-length)
	}
	return output
}

// Implement io.writer
// Direct write a string of Bytes
func (log *CliLog) Write(message []byte) (int, error) {
	fmt.Print(string(message))
	return len(message), nil
}
