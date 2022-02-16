package qga

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/libvirt/libvirt-go"
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

	data, err := execStatus(dom, retExecObj.Return.Pid)
	if err != nil {
		return nil, err
	}
	if data != nil {
		return data, err
	}
	return nil, errors.New("failed to get exec return data")
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
	_ = json.Unmarshal([]byte(retExecStatus), &retExecStatusObj)

	if !retExecStatusObj.Return.Exited {
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
