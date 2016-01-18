package log

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Log verbosity enum
const (
	VERBOSITY_FATAL   = iota // STOP and DIE
	VERBOSITY_ERROR          // Inform the user that an error has occured
	VERBOSITY_WARNING        // Warn the user about something, usually a mild failure
	VERBOSITY_MESSAGE        // Display a message to the user
	VERBOSITY_INFO           // a superflous message

	VERBOSITY_DEBUG        // Debug: start to display lots of debug details
	VERBOSITY_DEBUG_LOTS   // Debug: display lots of debug stuff
	VERBOSITY_DEBUG_WOAH   // Debug: this is getting pretty verbose
	VERBOSITY_DEBUG_STAAAP // Debug: ridiculous debugging verbosity
)

// Log output handler for coach
type Log interface {
	Name() string // get the name of the log

	Verbosity() int             // get the verbosity of the log
	SetVerbosity(verbosity int) // set a new verbosity for the log

	MakeChild(target string) Log

	Hush()          // Hush a log to make warnings, and messages less verbose
	UnHush()        // Un hush the log
	IsHushed() bool // is the log hushed

	Fatal(messages ...string)                                   // Fatal error has occured. Execution should be stopped
	Error(messages ...string)                                   // A signigicant error has occured, and the user should be warned
	Warning(messages ...string)                                 // An error has occured, but execution can continue without worry
	Message(messages ...string)                                 // A user notification
	Info(messages ...string)                                    // Verbose gratuitous information messages for a user
	Debug(verbosity int, message string, object ...interface{}) // Debugging output that should only be shown if requested

	Write(message []byte) (int, error)	
}

/**
 * CoachLog is the default coach logger
 */

// Log factory method
func MakeCoachLog(name string, writer io.Writer, verbosity int) Log {
	return Log(&CoachLog{
		CoachLogSettings: CoachLogSettings{
			writer:    writer,
			stack:     []string{name},
			verbosity: verbosity,
			hush:      false,
		},
	})
}

// Configuration struct for a CoachLog
type CoachLogSettings struct {
	writer    io.Writer // a log writing target
	stack     []string  // patent name stack
	verbosity int       // current verbosity for this log object
	hush      bool      // if true, and verbosity is standard, then make the log quieter
}

// CoachLog default logging handler
type CoachLog struct {
	CoachLogSettings
}

func (log *CoachLog) Name() string {
	return log.stack[len(log.stack)]
}
func (log *CoachLog) Verbosity() int {
	return log.verbosity
}
func (log *CoachLog) SetVerbosity(verbosity int) {
	log.verbosity = verbosity
}
func (log *CoachLog) MakeChild(target string) Log {
	return Log(&CoachLog{
		CoachLogSettings: CoachLogSettings{
			writer:    log.writer,
			stack:     append(log.stack, target),
			verbosity: log.verbosity,
			hush:      log.hush,
		},
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
func (log *CoachLog) Fatal(messages ...string) {
	log.writeLog(VERBOSITY_FATAL, messages...)
	panic("Execution halted on FATAL error")
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
	if +verbosity <= log.verbosity && len(objects)>0 && objects[0]!=nil {
		fmt.Print("	")
		fmt.Fprintln(log, objects...)
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
	case VERBOSITY_FATAL:
		elements = append(elements, "[FATAL]", log.joinStack())
	case VERBOSITY_ERROR:
		elements = append(elements, "[ERROR]", log.joinStack())

	case VERBOSITY_WARNING:
		elements = append(elements, "[WARNING]")

	case VERBOSITY_MESSAGE:
		prefix := log.stack[len(log.stack)-1]+": "
		if length := utf8.RuneCountInString(prefix); length < 25 {
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
func (log *CoachLog) joinStack() string {
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
func (log *CoachLog) Write(message []byte) (int, error) {
	fmt.Print(string(message))
	return len(message), nil
}
