package log

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode/utf8"
	"strings"
)

// Log verbosity enum
const (
	VERBOSITY_CRITICAL = iota // STOP and DIE
	VERBOSITY_SEVERE          // Notify that a serious error has occured
	VERBOSITY_ERROR           // Inform the user that an error has occured
	VERBOSITY_WARNING         // Warn the user about something, usually a mild failure
	VERBOSITY_MESSAGE         // Display a message to the user
	VERBOSITY_INFO            // a superflous message

	VERBOSITY_DEBUG        // Debug: start to display lots of debug details
	VERBOSITY_DEBUG_LOTS   // Debug: display lots of debug stuff
	VERBOSITY_DEBUG_WOAH   // Debug: this is getting pretty verbose
	VERBOSITY_DEBUG_STAAAP // Debug: ridiculous debugging verbosity
)

// Log factory method
func GetLog(name string, verbosity int) Log {
	return Log(&CoachLog{
		writer:    os.Stdout,
		targets:   []string{name},
		verbosity: verbosity,
	})
}

// Log output handler for coach
type Log interface {
	Target() string // get the name of the log

	Verbosity() int             // get the verbosity of the log
	SetVerbosity(verbosity int) // set a new verbosity for the log

	MakeChild(target string) Log

	Hush()
	UnHush()
	IsHushed() bool

	Critical(messages ...string)
	Error(messages ...string)
	Warning(messages ...string)
	Message(messages ...string)
	Info(messages ...string)
	Debug(verbosity int, message string, object ...interface{})
}

// CoachLog default logging handler
type CoachLog struct {
	writer    io.Writer // a log writing target
	targets   []string  // Target name stack
	verbosity int       // current verbosity for this log object
	hush      bool      // if true, and verbosity is standard, then make the log quieter
}

func (log *CoachLog) Target() string {
	return log.targets[len(log.targets)]
}
func (log *CoachLog) Verbosity() int {
	return log.verbosity
}
func (log *CoachLog) SetVerbosity(verbosity int) {
	log.verbosity = verbosity
}
func (log *CoachLog) MakeChild(target string) Log {
	return Log(&CoachLog{
		targets:   append(log.targets, target),
		verbosity: log.verbosity,
	})
}

func (log *CoachLog) IsHushed() bool {
	return log.hush
}
func (log *CoachLog) Hush() {
	log.hush = true
}
func (log *CoachLog) UnHush() {
	log.hush = false
}

// Implement a Critical error
func (log *CoachLog) Critical(messages ...string) {
	log.writeLog(VERBOSITY_CRITICAL, messages...)
}

// Implement an error
func (log *CoachLog) Error(messages ...string) {
	log.writeLog(VERBOSITY_ERROR, messages...)
}

// Register a Warning error
func (log *CoachLog) Warning(messages ...string) {
	log.writeLog(VERBOSITY_WARNING, messages...)
}

// Register a message
func (log *CoachLog) Message(messages ...string) {
	log.writeLog(VERBOSITY_MESSAGE, messages...)
}

// Register an information verbose message
func (log *CoachLog) Info(messages ...string) {
	log.writeLog(VERBOSITY_INFO, messages...)
}

// Debug message and data
func (log *CoachLog) Debug(verbosity int, message string, objects ...interface{}) {
	log.writeLog(verbosity, message)
	if objects!=nil {
		fmt.Fprintln(log, objects)
	}
}

// internal logging writer
func (log *CoachLog) writeLog(verbosity int, messages ...string) {

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
	case VERBOSITY_CRITICAL:
		elements = append(elements, "[CRITICAL]", log.joinTargets())
	case VERBOSITY_SEVERE:
		elements = append(elements, "[SEVERE]", log.joinTargets())
	case VERBOSITY_ERROR:
		elements = append(elements, "[ERROR]", log.joinTargets())

	case VERBOSITY_WARNING:
		elements = append(elements, "[WARNING]")

	case VERBOSITY_MESSAGE:

	case VERBOSITY_INFO:
		elements = append(elements, "-->")

	default:
		elements = append(elements, "("+strconv.Itoa(verbosity)+")", log.joinTargets())
	}

	elements = append(elements, messages...)

	output := strings.Join(elements, " ") + "\n"
	log.Write([]byte(output))

}

// joins the log targets into a printable string for message prefixing
func (log *CoachLog) joinTargets() string {
	output := ""
	if len(log.targets) > 0 {
		output += "[" + strings.Join(log.targets, "][") + "]"
	}
	length := utf8.RuneCountInString(output)
	if length < 25 {
		output += strings.Repeat("-", 25-length)
	}
	return output
}

// Implement io.writer
// Direct write a string of Bytes
func (log *CoachLog) Write(message []byte) (int, error) {
	fmt.Print(string(message))
	return len(message), nil
}
