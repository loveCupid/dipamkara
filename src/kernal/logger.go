package kernal

import (
	"os"
	// "io"
	"context"
	"fmt"
	"log"
)

type log_item struct {
	file_name string
	circuler  int

	w *os.File
	l *log.Logger

	close chan int
	msg   chan string
}

type Logger struct {
	debug *log_item
	info  *log_item
	err   *log_item

	status   int    // 0 - running, 1 - stop
	env      string // online, test, pre
	buf_size int
}

const (
	prefix = "/tmp/"
)

func async(item *log_item) {
	for {
		select {
		case <-item.close:
			{
				item.w.Sync()
				item.w.Close()
				return
			}
		case m := <-item.msg:
			{
				// item.l.Println(m)
                fmt.Println(m)
			}
		}
	}
}

func NewLogger(name string, env string) *Logger {
	var err error
	lgr := new(Logger)

	lgr.env = env
	lgr.status = 0
	lgr.buf_size = 9999

    // Mkdir
    os.Mkdir(prefix + name + "/", 0777)

	{
		// Debug item
		lgr.debug = new(log_item)
		lgr.debug.file_name = prefix + name + "/" + name + ".debug.log"
		lgr.debug.circuler = 0
		lgr.debug.w, err = os.OpenFile(lgr.debug.file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		ErrorPanic(err)
		lgr.debug.l = log.New(lgr.debug.w, "", log.LstdFlags|log.Lshortfile)
		lgr.debug.msg = make(chan string, lgr.buf_size)
		lgr.debug.close = make(chan int)
		go async(lgr.debug)
	}
	{
		// info item
		lgr.info = new(log_item)
		lgr.info.file_name = prefix + name + "/" + name + ".info.log"
		lgr.info.circuler = 0
		lgr.info.w, err = os.OpenFile(lgr.info.file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		ErrorPanic(err)
		lgr.info.l = log.New(lgr.info.w, "", log.LstdFlags|log.Lshortfile)
		lgr.info.msg = make(chan string, lgr.buf_size)
		lgr.info.close = make(chan int)
		go async(lgr.info)
	}
	{
		// err item
		lgr.err = new(log_item)
		lgr.err.file_name = prefix + name + "/" + name + ".err.log"
		lgr.err.circuler = 0
		lgr.err.w, err = os.OpenFile(lgr.err.file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		ErrorPanic(err)
		lgr.err.l = log.New(lgr.err.w, "", log.LstdFlags|log.Lshortfile)
		lgr.err.msg = make(chan string, lgr.buf_size)
		lgr.err.close = make(chan int)
		go async(lgr.err)
	}

	return lgr
}

func ReleaseLogger(s *Server) {
	s.logger.debug.close <- 1
	s.logger.info.close <- 1
	s.logger.err.close <- 1
}

func (l *Logger) _debug(format string, v ...interface{}) {
	if l.status == 0 && l.env != "test" {
		l.debug.msg <- fmt.Sprintf(format, v...)
	}
}
func (l *Logger) _info(format string, v ...interface{}) {
	if l.status == 0 {
		l.info.msg <- fmt.Sprintf(format, v...)
	}
}
func (l *Logger) _err(format string, v ...interface{}) {
	if l.status == 0 {
		l.err.msg <- fmt.Sprintf(format, v...)
	}
}
func Debug(ctx context.Context, format string, v ...interface{}) {
	s := ctx.Value(Skey).(*Server)
	s.logger._debug(format, v...)
}
func Info(ctx context.Context, format string, v ...interface{}) {
	s := ctx.Value(Skey).(*Server)
	s.logger._info(format, v...)
}
func Error(ctx context.Context, format string, v ...interface{}) {
	s := ctx.Value(Skey).(*Server)
	s.logger._err(format, v...)
}
