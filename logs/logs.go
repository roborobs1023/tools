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

	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

// InitLog creates log handlers for Log.
// To log to both Log and Console pass true to inialization function
func (l *Log) InitLog() error {
	var err error
	l.File, err = os.OpenFile(l.Name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	l.Trace = log.New(l.File,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.Info = log.New(l.File,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.Warning = log.New(l.File,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	l.Error = log.New(l.File,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	if l.MultiWriter {
		mw := io.MultiWriter(os.Stdout, l.File)
		l.Trace.SetOutput(mw)
		l.Info.SetOutput(mw)
		l.Warning.SetOutput(mw)
		l.Error.SetOutput(mw)
	}
	return nil
}
