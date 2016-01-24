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
 * CliLog is the default coach logger
 */

// Log factory method
func MakeCliLog(name string, writer io.Writer, verbosity int) Log {
	return Log(&CliLog{
		CliLogSettings: CliLogSettings{
			writer:    writer,
			stack:     []string{name},
			verbosity: verbosity,
			hush:      false,
		},
	})
}
