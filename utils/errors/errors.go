package errors

import (
	builtin "errors"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
)

const (
	calldepth          int    = 2
	funtionCallerDepth int    = 3
	defaultTitle       string = "There is something wrong, please try again"
)

var mutex *sync.Mutex
var goPath string
var callerParamRegexp *regexp.Regexp

func init() {
	if mutex == nil {
		mutex = &sync.Mutex{}
	}

	if goPath == "" {
		goPath = strings.ReplaceAll(os.Getenv("GOPATH"), "\\", "/")
	}

	if callerParamRegexp == nil {
		callerParamRegexp = regexp.MustCompile(`(0x(([a-f]|[0-9])+)(,?)(\s?))+`)
	}
}

type Error struct {
	file     string
	title    string
	errorObj error
}

func (e *Error) Error() string {
	return e.errorObj.Error()
}

func (e *Error) GetSource() string {
	return e.file
}

func (e *Error) GetTitle() string {
	return e.title
}

func getCallerFunctionName() string {
	mutex.Lock() // need to lock it, it's expensive
	defer mutex.Unlock()

	callstack := debug.Stack()
	lines := strings.Split(string(callstack), "\n")

	return string(callerParamRegexp.ReplaceAll([]byte(lines[funtionCallerDepth+4]), []byte("")))
}

func New(errorMsg string, title ...string) *Error {

	newTitle := defaultTitle
	if len(title) > 0 {
		newTitle = title[0]
	}
	return &Error{
		file:     getCallerFunctionName(),
		title:    newTitle,
		errorObj: builtin.New(errorMsg),
	}
}
func NewFromError(err error, title ...string) *Error {

	newTitle := defaultTitle
	if len(title) > 0 {
		newTitle = title[0]
	}
	return &Error{
		file:     getCallerFunctionName(),
		title:    newTitle,
		errorObj: err,
	}
}
