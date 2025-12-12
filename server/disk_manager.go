package server

import (
	"bytes"
	"disk-usage-analyser/server/disk"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/xgo/support/cmd"
)

func handleListDisks(w http.ResponseWriter, r *http.Request) {
	disks, err := disk.ListDisks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(disks)
}

type MountRequest struct {
	DeviceID string `json:"deviceID"`
	Password string `json:"password"`
}

func handleMountDisk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MountRequest
	// Try to decode body
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&req)
	}
	// Fallback to query param for deviceID
	if req.DeviceID == "" {
		req.DeviceID = r.URL.Query().Get("deviceID")
	}

	if req.DeviceID == "" {
		http.Error(w, "deviceID is required", http.StatusBadRequest)
		return
	}

	// Fetch disk info
	info, err := disk.GetDiskInfo(req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.MountPoint != "" {
		w.Write([]byte("ok"))
		return
	}

	// Check if it's ExFAT or Windows_NTFS (sometimes mislabeled)
	isExFAT := strings.EqualFold(info.FilesystemType, "exfat") ||
		strings.Contains(strings.ToLower(info.Content), "exfat") ||
		(strings.Contains(strings.ToLower(info.Content), "windows_ntfs") && strings.Contains(strings.ToLower(info.Content), "exfat")) ||
		strings.EqualFold(info.FilesystemType, "exfat")

	if isExFAT {
		volName := info.VolumeName
		if volName == "" {
			volName = req.DeviceID
		}

		// Use ~/Volumes/<Name> instead of /Volumes/<Name> to avoid permission issues
		homeDir, err := os.UserHomeDir()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get user home dir: %v", err), http.StatusInternalServerError)
			return
		}
		mountPoint := filepath.Join(homeDir, "Volumes", volName)

		// Create mount point as user
		if err := cmd.Debug().Run("mkdir", "-p", mountPoint); err != nil {
			http.Error(w, fmt.Sprintf("failed to create mount point: %v", err), http.StatusInternalServerError)
			return
		}

		// If password is not provided, try sudo -n first
		var outBuf bytes.Buffer
		var mountErr error

		if req.Password == "" {
			// Try non-interactive first
			mountErr = cmd.Debug().Stdout(&outBuf).Stderr(&outBuf).Run("sudo", "-n", "mount", "-t", "exfat", "/dev/"+req.DeviceID, mountPoint)
			if mountErr != nil {
				// Check if it failed due to missing password
				// sudo -n exits with 1 and usually prints something
				// But simpler is to assume if it fails we might need password
				// We can return 401 to prompt user
				// However, check output content to be sure?
				// "sudo: a password is required" is typical output
				outputStr := outBuf.String()
				if strings.Contains(outputStr, "password is required") || strings.Contains(outputStr, "sudo:") {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Sudo password required"))
					return
				}
				// Other error
				http.Error(w, fmt.Sprintf("failed to mount exfat disk (sudo -n): %v\nOutput: %s", mountErr, outputStr), http.StatusInternalServerError)
				return
			}
		} else {
			// Use sudo -S with password
			mountErr = cmd.Debug().Stdin(strings.NewReader(req.Password+"\n")).Stdout(&outBuf).Stderr(&outBuf).Run("sudo", "-S", "mount", "-t", "exfat", "/dev/"+req.DeviceID, mountPoint)
			if mountErr != nil {
				outputStr := outBuf.String()
				if strings.Contains(outputStr, "incorrect password") || strings.Contains(outputStr, "try again") {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Incorrect password"))
					return
				}
				http.Error(w, fmt.Sprintf("failed to mount exfat disk (sudo -S): %v\nOutput: %s", mountErr, outputStr), http.StatusInternalServerError)
				return
			}
		}

		w.Write([]byte("ok"))
		return
	}

	// Default behavior
	var outBuf bytes.Buffer
	err = cmd.Debug().Stdout(&outBuf).Stderr(&outBuf).Run("diskutil", "mount", req.DeviceID)
	if err != nil {
		outputStr := outBuf.String()
		if strings.Contains(outputStr, "SUIS premount dissented") {
			http.Error(w, "System is verifying the disk (fsck). Please wait until it finishes.", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("failed to mount disk: %v\nOutput: %s", err, outputStr), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}

func handleUnmountDisk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceID := r.URL.Query().Get("deviceID")
	if deviceID == "" {
		http.Error(w, "deviceID is required", http.StatusBadRequest)
		return
	}

	var outBuf bytes.Buffer
	err := cmd.Debug().Stdout(&outBuf).Stderr(&outBuf).Run("diskutil", "unmount", deviceID)
	if err != nil {
		outputStr := outBuf.String()
		http.Error(w, fmt.Sprintf("failed to unmount disk: %v\nOutput: %s", err, outputStr), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}

func handleOpenDisk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path is required", http.StatusBadRequest)
		return
	}

	var outBuf bytes.Buffer
	err := cmd.Debug().Stdout(&outBuf).Stderr(&outBuf).Run("open", path)
	if err != nil {
		outputStr := outBuf.String()
		http.Error(w, fmt.Sprintf("failed to open path: %v\nOutput: %s", err, outputStr), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("ok"))
}
