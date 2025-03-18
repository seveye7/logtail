package main

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hpcloud/tail"
)

type LogConfig struct {
	ConfigName string
	Config     *Config
	// Offset        int64
	// WatchFileName string
	Close chan struct{}
	wg    *sync.WaitGroup

	Offsets map[string]int64
	sync.Mutex
}

var (
	logFiles = map[string]*LogConfig{}
	lw       = &LogSwaper{}
)

func (l *LogConfig) Run(m map[string]int64, index int) {
	fileName := l.Config.Files[index].Name
	if strings.Contains(fileName, "{date}") {
		now := time.Now()
		fileName = strings.ReplaceAll(fileName, "{date}", now.Format("2006-01-02"))
	}
	l.Lock()
	offset := m[fileName]
	l.Offsets[fileName] = offset
	l.Unlock()

	writeTicker := time.NewTicker(time.Second * 3)
	now := time.Now()
	nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
	nextTimer := time.NewTimer(nextDay.Sub(now))
	defer func() {
		nextTimer.Stop()
		writeTicker.Stop()
		<-l.Close
		l.wg.Done()
	}()

	sinkWriter := NewSinkWriter(l.Config, index)
	parser := NewParser(l.Config, index)

	textLogger.Info("tail file", "file", fileName, "offset", offset)
	if fileName == "" {
		<-l.Close
		l.wg.Done()
		return
	}

	t, err := tail.TailFile(fileName, tail.Config{
		Follow:   true,
		Location: &tail.SeekInfo{Offset: offset},
		Logger:   lw,
		Poll:     true,
	})
	if err != nil {
		textLogger.Warn("tail file error: ", "err", err, "file", fileName)
		return
	}
	defer t.Stop()

	for {
		select {
		case line := <-t.Lines:
			parsed, err := parser.Parse(line.Text)
			if err != nil {
				textLogger.Warn("parse text error", "err", err)
				continue
			}

			n, total, err := sinkWriter.Write(parsed)
			if err != nil {
				textLogger.Warn("write to sink error, exit task", "err", err)
				return
			}

			l.Lock()
			// add \n length
			l.Offsets[fileName] += n + total
			l.Unlock()
		case <-writeTicker.C:
			n, total, err := sinkWriter.Write("")
			if err != nil {
				textLogger.Warn("time write to sink error, exit task", "err", err)
				return
			}
			if n == 0 {
				continue
			}
			l.Lock()
			// add \n length
			l.Offsets[fileName] += n + total
			l.Unlock()
		case <-l.Close:
			if l.wg != nil {
				l.wg.Done()
			}
			return
		case <-nextTimer.C:
			// 第二天
			now := time.Now()
			nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
			nextTimer.Reset(nextDay.Sub(now))

			fileName = l.Config.Files[index].Name
			if strings.Contains(fileName, "{date}") {
				fileName = strings.ReplaceAll(fileName, "{date}", now.Format("2006-01-02"))
				textLogger.Info("reset tail file", "file", fileName)

				l.Lock()
				l.Offsets[fileName] = 0
				l.Unlock()
				t.Stop()
				t, err = tail.TailFile(fileName, tail.Config{
					Follow:   true,
					Location: &tail.SeekInfo{Offset: 0},
					Logger:   lw,
					Poll:     true,
				})
				if err != nil {
					textLogger.Warn("tail file error: ", "err", err, "file", fileName)
					return
				}
			}
		}
	}
}

func rangeLogFiles(m map[string]int64) {
	for _, logFile := range logFiles {
		for i := range logFile.Config.Files {
			go logFile.Run(m, i)
		}
	}
}

func searchDir(dirName string, wg *sync.WaitGroup) error {
	// 遍历目录
	dir, err := os.ReadDir(dirName)
	if err != nil {
		return err
	}

	for _, file := range dir {
		fileName := file.Name()
		configName := dirName + "/" + fileName
		if file.IsDir() {
			searchDir(configName, wg)
		}
		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		// 读取文件
		config, err := LoadConfig(configName)
		log.Println(configName, config != nil, err)
		if err != nil {
			continue
		}

		f := &LogConfig{
			ConfigName: fileName,
			Config:     config,
			wg:         wg,
			Close:      make(chan struct{}, 1),
			// sinkWriter: NewSinkWriter(config),
			// parser:     NewParser(config),
			Offsets: map[string]int64{},
		}

		logFiles[fileName] = f
	}

	return nil
}
