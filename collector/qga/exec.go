package qga

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/libvirt/libvirt-go"
	"time"
)

func Exec(dom *libvirt.Domain, cmd GuestExecArg) ([]byte, error) {
	// exec
	execObj := guestExec{
		"guest-exec",
		cmd,
	}
	retExec, err := qemuAgentCommand(dom, execObj)
	if err != nil {
		return nil, err
	}
	retExecObj := retGuestExec{}
	err = json.Unmarshal([]byte(retExec), &retExecObj)
	if err != nil {
		return []byte{}, err
	}

	// get exec return data, default wait 2s
	for i := 0; i < 10; i++ {
		data, err := execStatus(dom, retExecObj.Return.Pid)
		if err != nil {
			return nil, err
		}
		if data != nil {
			return data, err
		}
		time.Sleep(200)
	}
	return nil, errors.New("failed to get exec return data.")
}

func execStatus(dom *libvirt.Domain, pid int) ([]byte, error) {
	// get exec status
	execStatusObj := guestExecStatus{
		"guest-exec-status",
		struct {
			Pid int `json:"pid"`
		}{Pid: pid},
	}
	retExecStatus, err := qemuAgentCommand(dom, execStatusObj)
	if err != nil {
		return []byte{}, err
	}
	retExecStatusObj := retGuestExecStatus{}
	err = json.Unmarshal([]byte(retExecStatus), &retExecStatusObj)

	if retExecStatusObj.Return.Exited == false {
		if retExecStatusObj.Return.Exitcode != 0 {
			errData, _ := base64.StdEncoding.DecodeString(retExecStatusObj.Return.ErrData)
			return nil, errors.New(string(errData))
		}
		return nil, nil
	}
	ret, err := base64.StdEncoding.DecodeString(retExecStatusObj.Return.OutData)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
