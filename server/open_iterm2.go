package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
)

type OpenIterm2Request struct {
	Path string `json:"path"`
}

func HandleOpenIterm2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to read request body"}`, http.StatusBadRequest)
		return
	}

	var req OpenIterm2Request
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Path) == "" {
		http.Error(w, `{"error":"path required"}`, http.StatusBadRequest)
		return
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		http.Error(w, `{"error":"home directory unavailable"}`, http.StatusInternalServerError)
		return
	}

	absPath := filepath.Clean(resolveTildePath(req.Path, homeDir))
	target := absPath
	if filepath.Base(target) == "node_modules" {
		target = filepath.Dir(target)
	}

	if err := iterm2.OpenConfig(target, &iterm2.Config{Mode: iterm2.ModeForceNew}); err != nil {
		msg, code := openIterm2ErrorStatus(err)
		http.Error(w, `{"error":"`+jsonEscape(msg)+`"}`, code)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "opened": target})
}

func openIterm2ErrorStatus(err error) (string, int) {
	switch {
	case err == iterm2.ErrUnsupportedPlatform:
		return err.Error(), http.StatusNotImplemented
	case err == iterm2.ErrNotInstalled:
		return err.Error(), http.StatusServiceUnavailable
	default:
		return err.Error(), http.StatusInternalServerError
	}
}

func jsonEscape(s string) string {
	b, _ := json.Marshal(s)
	if len(b) >= 2 {
		return string(b[1 : len(b)-1])
	}
	return s
}