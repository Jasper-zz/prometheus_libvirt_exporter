package qga

import (
	"encoding/json"
	"github.com/libvirt/libvirt-go"
)

// guest-read-file
type fileOpen struct {
	Execute   string `json:"execute"`
	Arguments struct {
		Path string `json:"path"`
		Mode string `json:"mode"`
	} `json:"arguments"`
}

type fileOpenRet struct {
	Return int `json:"return"`
}

type fileRead struct {
	Execute   string `json:"execute"`
	Arguments struct {
		Handle int `json:"handle"`
		Count  int `json:"count"`
	} `json:"arguments"`
}

type fileReadRet struct {
	Return struct {
		Count  int    `json:"count"`
		Bufb64 string `json:"buf-b64"`
		Eof    bool   `json:"eof"`
	}
}

type fileClose struct {
	Execute   string `json:"execute"`
	Arguments struct {
		Handle int `json:"handle"`
	} `json:"arguments"`
}

type fileCloseRet struct {
	Return struct{} `json:"return"`
}

// guest-exec
type GuestExecArg struct {
	Path          string   `json:"path"`
	Arg           []string `json:"arg"`
	CaptureOutput bool     `json:"capture-output"`
}

type guestExec struct {
	Execute   string       `json:"execute"`
	Arguments GuestExecArg `json:"arguments"`
}

type guestExecStatus struct {
	Execute   string `json:"execute"`
	Arguments struct {
		Pid int `json:"pid"`
	} `json:"arguments"`
}

type retGuestExec struct {
	Return struct {
		Pid int `json:"pid"`
	} `json:"return"`
}

type retGuestExecStatus struct {
	Return struct {
		Exitcode int    `json:"exitcode"`
		OutData  string `json:"out-data"`
		ErrData  string `json:"err-data"`
		Exited   bool   `json:"exited"`
	} `json:"return"`
}

// qemu-agent-command
func qemuAgentCommand(dom *libvirt.Domain, cmdObj interface{}) (string, error) {
	cmd, err := json.Marshal(cmdObj)
	if err != nil {
		return "", err
	}
	cmdRet, err := dom.QemuAgentCommand(string(cmd), 3, 0)
	if err != nil {
		return "", err
	}
	return cmdRet, nil
}
