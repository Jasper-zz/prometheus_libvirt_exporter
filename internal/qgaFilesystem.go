package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type filesystemLabels struct {
	MountPoint, Device, FsType string
}

type filesystemStats struct {
	Labels filesystemLabels
	Inodes float64
	IAvail float64
	Size   float64
	Avail  float64
}

func GetFilesystem(data []byte) ([]filesystemStats, error) {
	var fsStats []filesystemStats
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		partsNum := []int64{}
		parts := strings.Fields(scanner.Text())
		if !strings.HasPrefix(parts[0], "/dev/") {
			continue
		}
		if len(parts) < 7 {
			return nil, fmt.Errorf("malformed mount point information: %q", scanner.Text())
		}
		for _, p := range parts[3:] {
			ret, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse mount point information failed: %q", scanner.Text())
			}
			partsNum = append(partsNum, ret*1024)
		}
		fsStats = append(fsStats, filesystemStats{
			filesystemLabels{
				Device: parts[0],
				FsType:     parts[1],
				MountPoint:     parts[2],
			},
			float64(partsNum[0]),
			float64(partsNum[1]),
			float64(partsNum[2]),
			float64(partsNum[3]),
		})
	}

	return fsStats, scanner.Err()
}
