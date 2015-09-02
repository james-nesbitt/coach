package main

import (
	"io"
	"strings"
	"strconv"
	"unicode/utf8"

	"fmt"
)

const (
	LOG_SEVERITY_CRITICAL=1
	LOG_SEVERITY_ERROR=3
	LOG_SEVERITY_MESSAGE=5
	LOG_SEVERITY_WARNING=6
	LOG_SEVERITY_DEBUG=7
	LOG_SEVERITY_DEBUG_LOTS=8
	LOG_SEVERITY_DEBUG_WOAH=9
	LOG_SEVERITY_DEBUG_STAAAP=10
)

func GetLog(writer io.Writer, severity int) Log {
	return Log{
		parents: []string{},
		writer: writer,
	  severity: severity,
	}
}

type Log struct {
	parents []string
	writer io.Writer
	severity int
}
func (log *Log) ChildLog(topic string) Log {
	return Log{
		parents: append(log.parents, topic),
		writer: log.writer,
	  severity: log.severity,
	}
}

func (log *Log) Message(message string) {
	if (log.severity>=LOG_SEVERITY_MESSAGE) {
		log.write(
			LOG_SEVERITY_MESSAGE,
			message,
		)
	}
}
func (log *Log) Note(message string) {
	if (log.severity>=LOG_SEVERITY_MESSAGE) {
		log.write(
			LOG_SEVERITY_MESSAGE,
			message,
		)
	}
}

func (log *Log) Warning(message string) {
	if (log.severity>=LOG_SEVERITY_WARNING) {
		log.write(
			LOG_SEVERITY_WARNING,
			message,
		)
	}
}
func (log *Log) Error(message string) {
	if (log.severity>=LOG_SEVERITY_ERROR) {
		log.write(
			LOG_SEVERITY_ERROR,
			"ERROR",
			message,
		)
	}
}
func (log *Log) Fatal(message string) {
	if (log.severity>=LOG_SEVERITY_CRITICAL) {
		log.write(
			LOG_SEVERITY_CRITICAL,
			"CRITICAL",
			message,
		)
	}
}
func (log *Log) Debug(level int, message string) {
	if (log.severity>=level) {
		log.write(
			level,
			message,
		)
	}
}

func (log *Log) DebugObject(severity int, message string, elements ...interface{}) {
	if (log.severity>=severity) {
		fmt.Println("("+strconv.Itoa(severity)+")"+log.joinParents(), "{DEBUG}"+message, elements)
	}
}

func (log *Log) write(severity int, elements ...string) {
	elements = append( []string{log.joinParents()}, elements...)
	output := "("+strconv.Itoa(severity)+")"+strings.Join(elements, "\t")+ "\n"
	log.writer.Write( []byte(output) )
}
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

func (log Log) Write(message []byte) (int, error) {
	fmt.Print(string(message))
	return len(message), nil
}
