package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"disk-usage-analyser/tmpfiles"

	"github.com/xhd2015/dot-pkgs/go-pkgs/file/remotefs"
)

type DeleteBinariesRequest struct {
	Paths []string `json:"paths"`
}

type DeleteFailure struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

type DeleteBinariesResult struct {
	Deleted    []string        `json:"deleted"`
	Failed     []DeleteFailure `json:"failed"`
	FreedSize  int64           `json:"freedSize"`
	FreedHuman string          `json:"freedHuman"`
}

func HandleTmpBinariesDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, `{"error":"home directory unavailable"}`, http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read request body"}`, http.StatusBadRequest)
		return
	}

	var req DeleteBinariesRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	result := DeleteBinariesResult{
		Deleted: []string{},
		Failed:  []DeleteFailure{},
	}

	for _, displayPath := range req.Paths {
		absPath := displayPath
		if strings.HasPrefix(displayPath, "~/") || displayPath == "~" {
			absPath = resolveTildePath(displayPath, homeDir)
		}

		if info, statErr := os.Stat(absPath); statErr == nil && info.IsDir() {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: "not a regular file: directory",
			})
			continue
		}

		entry, ok := lookupBinarySession(homeDir, displayPath)
		if !ok {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: "path not in current scan results",
			})
			continue
		}

		absPath = entry.AbsPath
		if strings.HasPrefix(displayPath, "~/") {
			absPath = resolveTildePath(displayPath, homeDir)
		}

		if remote, err := remotefs.IsRemoteBackedPath(absPath); err != nil {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: err.Error(),
			})
			continue
		} else if remote {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: "remote-backed filesystem path",
			})
			continue
		}

		info, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				result.Failed = append(result.Failed, DeleteFailure{
					Path:  displayPath,
					Error: "file not found",
				})
				continue
			}
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: err.Error(),
			})
			continue
		}

		if info.IsDir() {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: "not a regular file: directory",
			})
			continue
		}
		if !info.Mode().IsRegular() {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: "not a regular file",
			})
			continue
		}

		if _, ok := tmpfiles.ClassifyFile(absPath, info.Size(), entry.Repo, homeDir); !ok {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: "file is not a binary",
			})
			continue
		}

		if err := os.Remove(absPath); err != nil {
			result.Failed = append(result.Failed, DeleteFailure{
				Path:  displayPath,
				Error: err.Error(),
			})
			continue
		}

		result.Deleted = append(result.Deleted, displayPath)
		result.FreedSize += entry.Size
	}

	result.FreedHuman = tmpfiles.FormatHumanSize(result.FreedSize)
	json.NewEncoder(w).Encode(result)
}