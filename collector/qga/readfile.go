package qga

import (
	"encoding/base64"
	"encoding/json"
	"github.com/libvirt/libvirt-go"
)

func ReadFile(dom *libvirt.Domain, path string) ([]byte, error) {
	cpuStatAll := []byte{}

	// open file, return file handle
	fileOpenObj := fileOpen{"guest-file-open",
		struct {
			Path string `json:"path"`
			Mode string `json:"mode"`
		}{Path: path, Mode: "r"}}

	retOpen, err := qemuAgentCommand(dom, fileOpenObj)
	if err != nil {
		return []byte{}, nil
	}

	retOpenObj := fileOpenRet{}
	err = json.Unmarshal([]byte(retOpen), &retOpenObj)
	if err != nil {
		return []byte{}, nil
	}

	// close file
	defer func(handle int) {
		fileCloseObj := fileClose{"guest-file-close",
			struct {
				Handle int `json:"handle"`
			}{Handle: handle},
		}

		_, err := qemuAgentCommand(dom, fileCloseObj)
		if err != nil {
			return
		}
	}(retOpenObj.Return)

	// read file, return content and eof
	fileReadObj := fileRead{"guest-file-read",
		struct {
			Handle int `json:"handle"`
			Count  int `json:"count"`
		}{Handle: retOpenObj.Return, Count: 2048},
	}
	for i := 0; i < 10; i++ {
		retRead, err := qemuAgentCommand(dom, fileReadObj)
		if err != nil {
			return []byte{}, nil
		}
		retReadObj := fileReadRet{}
		err = json.Unmarshal([]byte(retRead), &retReadObj)
		if err != nil {
			return []byte{}, nil
		}
		cpuStat, err := base64.StdEncoding.DecodeString(retReadObj.Return.Bufb64)
		cpuStatAll = append(cpuStatAll, cpuStat...)
		if retReadObj.Return.Eof {
			break
		}
	}

	return cpuStatAll, nil
}
