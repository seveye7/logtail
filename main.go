package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"sync"
	"syscall"

	"github.com/judwhite/go-svc"
)

var (
	configDir  = flag.String("d", "/var/lib/logtail", "config dir")
	isConsole  = flag.Bool("c", false, "print log to console")
	textLogger *slog.Logger
)

func initLog(console bool) error {
	var textHandler *slog.TextHandler
	if console {
		textHandler = slog.NewTextHandler(os.Stdout, nil)
	} else {
		w, err := os.OpenFile("/var/log/logtail.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return err
		}
		textHandler = slog.NewTextHandler(w, nil)
	}

	textLogger = slog.New(textHandler)

	return nil
}

type LogTail struct {
	wg *sync.WaitGroup
}

func main() {
	logtail := &LogTail{
		wg: &sync.WaitGroup{},
	}
	if err := svc.Run(logtail, syscall.SIGINT, syscall.SIGTERM); err != nil {
		log.Fatalf("logtail exit, returned: %s", err)
	}
}

func (l *LogTail) Init(env svc.Environment) error {
	flag.Parse()

	if err := initLog(*isConsole); err != nil {
		return err
	}

	if _, err := os.Stat(logtailDir); os.IsNotExist(err) {
		err = os.MkdirAll(logtailDir, 0o755)
		if err != nil {
			textLogger.Warn("create logtail dir error", "err", err)
			return err
		}
	}

	textLogger.Info("logtail init")
	return nil
}

func (l *LogTail) Start() error {
	textLogger.Info("logtail start")

	m := loadOffset()
	//
	err := searchDir(*configDir, l.wg)
	if err != nil {
		return err
	}

	rangeLogFiles(m)
	return nil
}

func (l *LogTail) Stop() error {
	textLogger.Info("logtail close log files", "l", len(logFiles))
	for _, f := range logFiles {
		l.wg.Add(1)
		f.Close <- struct{}{}
	}
	l.wg.Wait()

	m := map[string]int64{}
	for _, f := range logFiles {
		// if f.Offset > 0 {
		// 	m[f.WatchFileName] = f.Offset
		// }
		for k, v := range f.Offsets {
			m[k] = v
		}
	}

	textLogger.Info("logtail save offset")
	saveOffset(m)
	textLogger.Info("logtail stop")

	return nil
}
