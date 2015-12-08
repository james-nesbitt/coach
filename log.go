package main

import (
	"io"
	"strings"
	"strconv"
	"unicode/utf8"

	"fmt"
)

// log severity
const (
	LOG_SEVERITY_CRITICAL=1        // STOP and DIE
	LOG_SEVERITY_SEVERE=2          // Notify that a serious error has occured
	LOG_SEVERITY_ERROR=3           // Inform the user that an error has occured
	LOG_SEVERITY_WARNING=4         // Warn the user about something, usually a mild failure
	LOG_SEVERITY_MESSAGE=5         // Display a message to the user
	LOG_SEVERITY_INFO=6            // a superflous message
	LOG_SEVERITY_DEBUG=7           // Debug: start to display lots of debug details
	LOG_SEVERITY_DEBUG_LOTS=8      // Debug: display lots of debug stuff
	LOG_SEVERITY_DEBUG_WOAH=9      // Debug: this is getting pretty verbose
	LOG_SEVERITY_DEBUG_STAAAP=10   // Debug: ridiculous debugging verbosity
)

// Log factory
func GetLog(writer io.Writer, severity int) Log {
	return Log{
		parents: []string{},
		writer: writer,
		severity: severity,
		store: false,
	}
}

// a signle log entry
type LogEntry struct {
	severity int
	message []string
}

// a collection of log entries that knows how to write itself
type Log struct {
	parents []string
	writer io.Writer

	severity int

	record []LogEntry

	hush bool // if set to TRUE, then WARNINGS, MESSAGES and INFO are reduced in SEVERITY
	store bool // if set to TRUE, then this Log will collect but not write entries (use .flush())
}

// create a child log of the current log, and track the parents 
// as a string that can be used for better labelling
func (log *Log) ChildLog(topic string) Log {
	return Log{
		parents: append(log.parents, topic),
		writer: log.writer,
		severity: log.severity,
		store: log.store, // store entries if the parent is storing

		record: []LogEntry{},
	}
}

func (log *Log) Severity() int {
	return log.severity
}
func (log *Log) SetSverity(severity int) {
	log.severity = severity
}

// hush and unhush this log
func (log *Log) Hush() {
	log.hush=true
}
func (log *Log) UnHush() {
	log.hush=false
}

// enable and disable the log store
func (log *Log) Store(store bool) {
	log.store=store
}

// a standard single line message
func (log *Log) Message(message... string) {
	log.Record(
		LOG_SEVERITY_MESSAGE,
		message...
	)
}
// a standard single line message
func (log *Log) Info(message... string) {
	log.Record(
		LOG_SEVERITY_INFO,
		message...
	)
}
// A long message, which should not be formatted
func (log *Log) Note(message string) {
	log.Record(
		LOG_SEVERITY_MESSAGE,
		message,
	)
}
// warn the user that something non-serious has occured
func (log *Log) Warning(message... string) {
	log.Record(
		LOG_SEVERITY_WARNING,
		message...
	)
}
// warn the user that a significant error has occured
func (log *Log) Error(message... string) {
	message = append([]string{"ERROR"}, message...)
	log.Record(
		LOG_SEVERITY_ERROR,
		message...
	)
}
// a fatal error has occured, stop
func (log *Log) Fatal(message... string) {
	message = append([]string{"CRITICAL"}, message...)
	log.Record(
		LOG_SEVERITY_CRITICAL,
		message...
	)
}
// record a debugging message
func (log *Log) Debug(level int, message... string) {
	log.Record(
		level,
		message...
	)
}

// write or record a message, dependending on log.store
func (log *Log) Record(severity int, elements ...string) {
	newEntry := LogEntry{severity:severity, message:elements}
	if log.store {
		log.record = append(log.record, newEntry)
	} else {
		log.write([]LogEntry{ newEntry }, log.severity)
	}
}

// write all recorded messages, and then remove them
func (log *Log) Flush(severity int) {
	log.write( log.record, log.severity )
	log.record = []LogEntry{}
}

// write passed entries to the log writer, if they are severe enough
func (log *Log) write(record []LogEntry, severity int) {
	for _, entry := range record {
		entrySeverity := entry.severity

		elements := []string{}  // a slice of string message elements

		if log.hush {
			switch entrySeverity {
				case LOG_SEVERITY_WARNING:
					elements = append(elements, "[HUSHED WARNING]")
					entrySeverity = LOG_SEVERITY_INFO
				case LOG_SEVERITY_MESSAGE:
					elements = append(elements, "[HUSHED MESSAGE]")
					entrySeverity = LOG_SEVERITY_INFO
			}
		}

		if entrySeverity > severity {
			continue
		}

		switch entrySeverity {
			case LOG_SEVERITY_CRITICAL:
				elements = append(elements, "[CRITICAL]", log.joinParents())
			case LOG_SEVERITY_SEVERE:
				elements = append(elements, "[SEVERE]", log.joinParents())
			case LOG_SEVERITY_ERROR:
				elements = append(elements, "[ERROR]", log.joinParents())

			case LOG_SEVERITY_WARNING:
				elements = append(elements, "[WARNING]")

			case LOG_SEVERITY_MESSAGE:

			case LOG_SEVERITY_INFO:  		
				elements = append(elements, "-->")

			default:
				elements = append( elements, "("+strconv.Itoa(severity)+")", log.joinParents())
		}

		elements = append(elements, entry.message...)

		output := strings.Join(elements, " ")+ "\n"
		log.writer.Write( []byte(output) )

	}
}

// make a message prefix out of the log parents, to make it easier to track messages
func (log *Log) joinParents() string {
	output := ""
	if len(log.parents)>0 {
		output += "["+strings.Join(log.parents, "][")+"]"
	}
	length := utf8.RuneCountInString(output)
	if length<25 {
		output += strings.Repeat("-", 25-length)
	}
	return output
}

/**
 * Direct FMT methods
 *
 * Sometimes a write is not good enough
 */

// debug print an object
func (log *Log) DebugObject(severity int, message string, elements ...interface{}) {
	entrySeverity := severity

	if log.hush {
		switch entrySeverity {
			case LOG_SEVERITY_WARNING:
				elements = append(elements, "[HUSHED WARNING]")
				entrySeverity = LOG_SEVERITY_INFO
			case LOG_SEVERITY_MESSAGE:
				elements = append(elements, "[HUSHED MESSAGE]")
				entrySeverity = LOG_SEVERITY_INFO
		}
	}

	if entrySeverity > log.severity {
		return
	}

	fmt.Println("("+strconv.Itoa(severity)+") "+log.joinParents(), "{DEBUG}"+message+"\n\t", elements)
}

// Implemen io.writer
// Direct write a string of Bytes
func (log Log) Write(message []byte) (int, error) {
	fmt.Print(string(message))
	return len(message), nil
}
