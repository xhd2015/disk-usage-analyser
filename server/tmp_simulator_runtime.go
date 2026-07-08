package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/xhd2015/xgo/support/cmd"
)

type SimulatorRuntimeCommandRunner func(name string, args ...string) ([]byte, error)

var simulatorRuntimeCommandRunner SimulatorRuntimeCommandRunner

func SetSimulatorRuntimeCommandRunner(runner SimulatorRuntimeCommandRunner) {
	simulatorRuntimeCommandRunner = runner
}

func defaultSimulatorRuntimeCommandRunner(name string, args ...string) ([]byte, error) {
	out, err := cmd.Debug().Output(name, args...)
	return []byte(out), err
}

func getSimulatorRuntimeCommandRunner() SimulatorRuntimeCommandRunner {
	if simulatorRuntimeCommandRunner != nil {
		return simulatorRuntimeCommandRunner
	}
	return defaultSimulatorRuntimeCommandRunner
}

type simulatorRuntimeRecord struct {
	Version    string `json:"version"`
	Identifier string `json:"identifier"`
	MountPath  string `json:"mountPath"`
	Deletable  bool   `json:"deletable"`
	State      string `json:"state"`
}

func ParseSimulatorRuntimeJSON(output string, sizeFn func(mountPath string) int64) ([]TmpRuntimeItem, error) {
	var records map[string]simulatorRuntimeRecord
	if err := json.Unmarshal([]byte(output), &records); err != nil {
		return nil, fmt.Errorf("parse simulator runtime json: %w", err)
	}

	var items []TmpRuntimeItem
	for _, rec := range records {
		if rec.MountPath == "" {
			continue
		}

		itemType := rec.Version
		if itemType == "" {
			itemType = rec.Identifier
		}

		size := sizeFn(rec.MountPath)

		var activeCount int64
		if rec.State == "Ready" {
			activeCount = 1
		}

		reclaimable := int64(0)
		if rec.Deletable {
			reclaimable = size
		}

		items = append(items, TmpRuntimeItem{
			Type:        itemType,
			TotalCount:  1,
			ActiveCount: activeCount,
			Size:        size,
			Reclaimable: reclaimable,
		})
	}

	return items, nil
}

func simulatorRuntimeMountSize(mountPath string) int64 {
	if mountPath == "" {
		return 0
	}
	out, err := cmd.Debug().Output("du", "-sb", mountPath)
	if err == nil {
		if size, parseErr := ParseDuSBOutput(string(out)); parseErr == nil {
			return size
		}
	}
	info, err := os.Stat(mountPath)
	if err != nil || !info.IsDir() {
		return 0
	}
	size, _, err := CalculateSize(os.DirFS(mountPath), ".")
	if err != nil {
		return 0
	}
	return size
}

func CollectSimulatorRuntimeStats() ([]TmpRuntimeItem, error) {
	if runtime.GOOS != "darwin" {
		return nil, nil
	}

	runner := getSimulatorRuntimeCommandRunner()
	output, err := runner("xcrun", "simctl", "runtime", "list", "-j")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, nil
	}

	items, err := ParseSimulatorRuntimeJSON(string(output), simulatorRuntimeMountSize)
	if err != nil {
		return nil, nil
	}
	return items, nil
}