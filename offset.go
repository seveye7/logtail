package main

import (
	"encoding/json"
	"os"
)

var (
	logtailDir = "/var/lib/logtail"
	offsetFile = logtailDir + "/offset.json"
)

func loadOffset() map[string]int64 {
	m := map[string]int64{}
	buff, err := os.ReadFile(offsetFile)
	if err != nil {
		textLogger.Warn("load offset error", "err", err)
		return m
	}

	err = json.Unmarshal(buff, &m)
	if err != nil {
		textLogger.Warn("load offset error", "err", err)
	}
	return m
}

func saveOffset(m map[string]int64) {
	buff, err := json.Marshal(m)
	if err != nil {
		textLogger.Warn("save offset error", "err", err)
		return
	}

	err = os.WriteFile(offsetFile, buff, 0o644)
	if err != nil {
		textLogger.Warn("save offset error", "err", err)
	}
}
