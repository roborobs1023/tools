package logs

import (
	"io"
	"log"
	"os"
)

type Log struct {
	Name        string
	File        *os.File
	MultiWriter bool
}

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

// InitLog creates log handlers for Log.
// To log to both Log and Console pass true to inialization function
func (l *Log) Init() error {

	var err error
	l.File, err = os.OpenFile(l.Name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	Trace = log.New(l.File,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(l.File,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(l.File,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(l.File,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	if l.MultiWriter {
		mw := io.MultiWriter(os.Stdout, l.File)
		Trace.SetOutput(mw)
		Info.SetOutput(mw)
		Warning.SetOutput(mw)
		Error.SetOutput(mw)
	}
	return nil
}
