package server

import (
	"encoding/json"
	"errors"
	"runtime"
	"strconv"
	"strings"

	"github.com/xhd2015/xgo/support/cmd"
)

type TmpVmStorageItem struct {
	Label string `json:"label"`
	Path  string `json:"path"`
	Size  int64  `json:"size"`
}

type TmpVmInternal struct {
	Items          []TmpVmStorageItem `json:"items"`
	TotalSize      int64              `json:"totalSize"`
	MachineRunning bool               `json:"machineRunning"`
}

const (
	podmanVmStoragePath = "/home/core/.local/share/containers/storage"
	podmanVmOverlayPath = "/home/core/.local/share/containers/storage/overlay"
)

type PodmanMachineRunner func(args ...string) ([]byte, error)

var (
	podmanMachineRunner   PodmanMachineRunner
	podmanVmGOOSOverride  string
)

func SetPodmanMachineRunner(runner PodmanMachineRunner) {
	podmanMachineRunner = runner
}

func SetPodmanVmGOOSOverride(goos string) {
	podmanVmGOOSOverride = goos
}

func getPodmanVmGOOS() string {
	if podmanVmGOOSOverride != "" {
		return podmanVmGOOSOverride
	}
	return runtime.GOOS
}

func getPodmanMachineRunner() PodmanMachineRunner {
	if podmanMachineRunner != nil {
		return podmanMachineRunner
	}
	return defaultPodmanMachineRunner
}

func defaultPodmanMachineRunner(args ...string) ([]byte, error) {
	out, err := cmd.Debug().Output("podman", args...)
	return []byte(out), err
}

func ParseDuSBOutput(output string) (int64, error) {
	var lastValid int64
	var found bool
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		bytes, ok := parseDuSBLine(line)
		if ok {
			lastValid = bytes
			found = true
		}
	}
	if !found {
		return 0, errors.New("no valid du -sb output")
	}
	return lastValid, nil
}

func parseDuSBLine(line string) (int64, bool) {
	end := 0
	for end < len(line) && line[end] >= '0' && line[end] <= '9' {
		end++
	}
	if end == 0 {
		return 0, false
	}
	bytes, err := strconv.ParseInt(line[:end], 10, 64)
	if err != nil {
		return 0, false
	}
	return bytes, true
}

type podmanMachineInfo struct {
	Running bool `json:"Running"`
}

func isPodmanMachineRunning() (bool, error) {
	runner := getPodmanMachineRunner()
	output, err := runner("machine", "list", "--format", "json")
	if err != nil {
		return false, err
	}
	var machines []podmanMachineInfo
	if err := json.Unmarshal(output, &machines); err != nil {
		return false, err
	}
	for _, m := range machines {
		if m.Running {
			return true, nil
		}
	}
	return false, nil
}

func podmanMachineSSH(args ...string) ([]byte, error) {
	sshArgs := append([]string{"machine", "ssh", "--"}, args...)
	return getPodmanMachineRunner()(sshArgs...)
}

// podmanMachineDu runs du -sb inside the Podman VM. du often exits non-zero when
// some overlay dirs are unreadable, but still prints a usable total on stdout.
func podmanMachineDu(path string) (int64, error) {
	out, err := podmanMachineSSH("du", "-sb", path)
	size, parseErr := ParseDuSBOutput(string(out))
	if parseErr == nil {
		return size, nil
	}
	if err != nil {
		return 0, err
	}
	return 0, parseErr
}

func CollectPodmanVmInternal() (*TmpVmInternal, error) {
	if getPodmanVmGOOS() != "darwin" {
		return nil, nil
	}

	running, err := isPodmanMachineRunning()
	if err != nil || !running {
		return nil, nil
	}

	storageSize, err := podmanMachineDu(podmanVmStoragePath)
	if err != nil {
		return nil, nil
	}

	overlaySize, err := podmanMachineDu(podmanVmOverlayPath)
	if err != nil {
		return nil, nil
	}

	return &TmpVmInternal{
		Items: []TmpVmStorageItem{
			{Label: "Container storage", Path: podmanVmStoragePath, Size: storageSize},
			{Label: "Overlay layers", Path: podmanVmOverlayPath, Size: overlaySize},
		},
		TotalSize:      storageSize,
		MachineRunning: true,
	}, nil
}

func CollectPodmanRuntimeViaSSH() ([]TmpRuntimeItem, error) {
	if getPodmanVmGOOS() != "darwin" {
		return nil, nil
	}

	output, err := podmanMachineSSH("podman", "system", "df", "--format", "json")
	if err != nil {
		return nil, nil
	}
	items, err := ParseSystemDFJSON(string(output))
	if err != nil {
		return nil, nil
	}
	filtered := FilterRuntimeItems(items, "Images", "Build Cache")
	if len(filtered) == 0 {
		return nil, nil
	}
	return filtered, nil
}