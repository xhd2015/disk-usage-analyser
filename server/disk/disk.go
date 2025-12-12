package disk

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/xhd2015/xgo/support/cmd"
)

type Info struct {
	DeviceID   string `json:"deviceID"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Available  int64  `json:"available"`
	MountPoint string `json:"mountPoint"`
	Content    string `json:"content"`
	IsInternal bool   `json:"isInternal"`
	Status     string `json:"status"`
	Children   []Info `json:"children,omitempty"`
}

type ListOutput struct {
	AllDisksAndPartitions []Disk `json:"AllDisksAndPartitions"`
}

type Disk struct {
	DeviceIdentifier string      `json:"DeviceIdentifier"`
	Content          string      `json:"Content"`
	Size             int64       `json:"Size"`
	VolumeName       string      `json:"VolumeName"`
	MountPoint       string      `json:"MountPoint"`
	OSInternal       bool        `json:"OSInternal"`
	Partitions       []Partition `json:"Partitions"`
}

type Partition struct {
	DeviceIdentifier string `json:"DeviceIdentifier"`
	Content          string `json:"Content"`
	Size             int64  `json:"Size"`
	VolumeName       string `json:"VolumeName"`
	MountPoint       string `json:"MountPoint"`
}

type DetailInfo struct {
	FilesystemType            string `json:"FilesystemType"`
	FilesystemName            string `json:"FilesystemName"`
	VolumeName                string `json:"VolumeName"`
	MountPoint                string `json:"MountPoint"`
	Content                   string `json:"Content"`
	FilesystemUserVisibleName string `json:"FilesystemUserVisibleName"`
}

func GetDiskUsage() (map[string]int64, error) {
	output, err := cmd.Debug().Output("df", "-k")
	if err != nil {
		return nil, err
	}

	usage := make(map[string]int64)
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Filesystem") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}
		// Available is at index 3
		availStr := fields[3]
		availKB, err := strconv.ParseInt(availStr, 10, 64)
		if err != nil {
			continue
		}

		// Mount point starts at index 8. Join remaining fields.
		mountPoint := strings.Join(fields[8:], " ")
		usage[mountPoint] = availKB * 1024
	}
	return usage, nil
}

func GetDiskInfo(deviceID string) (*DetailInfo, error) {
	plistOutput, err := cmd.Debug().Output("diskutil", "info", "-plist", deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk info: %v", err)
	}
	jsonOutput, err := cmd.Debug().Stdin(strings.NewReader(plistOutput)).Output("plutil", "-convert", "json", "-r", "-o", "-", "--", "-")
	if err != nil {
		return nil, fmt.Errorf("failed to parse disk info (plutil): %v", err)
	}

	var info DetailInfo
	if err := json.Unmarshal([]byte(jsonOutput), &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal disk info: %v", err)
	}
	return &info, nil
}

func ListDisks() ([]Info, error) {
	// Execute diskutil list -plist
	plistOutput, err := cmd.Debug().Output("diskutil", "list", "-plist")
	if err != nil {
		return nil, fmt.Errorf("failed to run diskutil: %v", err)
	}

	// Execute plutil -convert json ... using plistOutput as stdin
	jsonOutput, err := cmd.Debug().Stdin(strings.NewReader(plistOutput)).Output("plutil", "-convert", "json", "-r", "-o", "-", "--", "-")
	if err != nil {
		return nil, fmt.Errorf("failed to run plutil: %v", err)
	}

	var data ListOutput
	if err := json.Unmarshal([]byte(jsonOutput), &data); err != nil {
		return nil, fmt.Errorf("failed to parse diskutil output: %v", err)
	}

	usage, _ := GetDiskUsage()

	// Check for running fsck processes
	psOut, _ := cmd.Debug().Output("ps", "aux")

	getStatus := func(deviceID string) string {
		if strings.Contains(psOut, "fsck") && strings.Contains(psOut, deviceID) {
			return "Checking"
		}
		return ""
	}

	var disks []Info
	for _, disk := range data.AllDisksAndPartitions {
		// Add the main disk
		content := disk.Content
		if content == "Windows_NTFS" {
			info, err := GetDiskInfo(disk.DeviceIdentifier)
			if err == nil && info.FilesystemUserVisibleName != "" {
				content = info.FilesystemUserVisibleName
			} else if err == nil && info.FilesystemType != "" {
				content = info.FilesystemType
			}
		}

		parent := Info{
			DeviceID:   disk.DeviceIdentifier,
			Name:       disk.VolumeName,
			Size:       disk.Size,
			Available:  usage[disk.MountPoint],
			MountPoint: disk.MountPoint,
			Content:    content,
			IsInternal: disk.OSInternal,
			Status:     getStatus(disk.DeviceIdentifier),
		}

		// Add partitions
		var children []Info
		for _, part := range disk.Partitions {
			partContent := part.Content
			if partContent == "Windows_NTFS" {
				info, err := GetDiskInfo(part.DeviceIdentifier)
				if err == nil && info.FilesystemUserVisibleName != "" {
					partContent = info.FilesystemUserVisibleName
				} else if err == nil && info.FilesystemType != "" {
					partContent = info.FilesystemType
				}
			}

			children = append(children, Info{
				DeviceID:   part.DeviceIdentifier,
				Name:       part.VolumeName,
				Size:       part.Size,
				Available:  usage[part.MountPoint],
				MountPoint: part.MountPoint,
				Content:    partContent,
				IsInternal: disk.OSInternal, // Inherit from parent
				Status:     getStatus(part.DeviceIdentifier),
			})
		}
		parent.Children = children
		disks = append(disks, parent)
	}
	return disks, nil
}
