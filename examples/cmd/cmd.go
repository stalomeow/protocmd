package cmd

import (
	"google.golang.org/protobuf/proto"
)

type CmdMessage interface {
	proto.Message
	CmdName() string
	CmdId() uint16
}

type cmdInfo struct {
	name    string
	factory func() CmdMessage
}

var cmdMap = make(map[uint16]*cmdInfo)

func RegisterCmd(factory func() CmdMessage) {
	msg := factory()
	cmdId := msg.CmdId()

	cmdMap[cmdId] = &cmdInfo{
		name:    msg.CmdName(),
		factory: factory,
	}
}

func RegisteredCmdCount() int {
	return len(cmdMap)
}

func NewCmdById(cmdId uint16) CmdMessage {
	info, ok := cmdMap[cmdId]
	if !ok {
		return nil
	}
	return info.factory()
}

func GetCmdNameById(cmdId uint16) string {
	info, ok := cmdMap[cmdId]
	if !ok {
		return "<Unknown Cmd>"
	}
	return info.name
}
