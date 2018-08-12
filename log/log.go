//-----------------------------------------------------------------------------
/*

Wrapper on standard logging.

*/
//-----------------------------------------------------------------------------

package log

import (
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

//-----------------------------------------------------------------------------

type LogWriter struct{}

var (
	Info  = log.New(LogWriter{}, "INFO ", 0)
	Debug = log.New(LogWriter{}, "DEBUG ", 0)
	Error = log.New(LogWriter{}, "ERROR ", 0)
)

func (f LogWriter) Write(p []byte) (n int, err error) {
	pc, file, line, ok := runtime.Caller(4)
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	log.Printf("%s:%d %s: %s", filepath.Base(file), line, fnName, p)
	return len(p), nil
}

//-----------------------------------------------------------------------------
