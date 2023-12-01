package protocmd

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type CmdMessage interface {
	proto.Message
	CmdName() string
	CmdId() uint16
}

type cmdInfo struct {
	name       string
	factory    func() CmdMessage
	descriptor protoreflect.MessageDescriptor
}

var cmdInfoMap = make(map[uint16]*cmdInfo)
var cmdIdMap = make(map[string]uint16)

func Register(factory func() CmdMessage) {
	msg := factory()
	cmdId := msg.CmdId()

	cmdInfoMap[cmdId] = &cmdInfo{
		name:       msg.CmdName(),
		factory:    factory,
		descriptor: msg.ProtoReflect().Descriptor(),
	}
}

func NewMessageByCmdId(cmdId uint16) (CmdMessage, error) {
	info, ok := cmdInfoMap[cmdId]
	if !ok {
		return nil, fmt.Errorf("failed to new Message with cmdId '%v' which was not registered", cmdId)
	}
	return info.factory(), nil
}

func MessageDescriptorByCmdId(cmdId uint16) (protoreflect.MessageDescriptor, error) {
	info, ok := cmdInfoMap[cmdId]
	if !ok {
		return nil, fmt.Errorf("failed to get MessageDescriptor with cmdId '%v' which was not registered", cmdId)
	}
	return info.descriptor, nil
}

func CmdCount() int {
	return len(cmdInfoMap)
}

func AllCmdIds() []uint16 {
	ids := make([]uint16, len(cmdInfoMap))
	i := 0
	for k := range cmdInfoMap {
		ids[i] = k
		i++
	}
	return ids
}

func CmdName(cmdId uint16) (string, bool) {
	info, ok := cmdInfoMap[cmdId]
	if !ok {
		return "", false
	}
	return info.name, true
}

func CmdId(cmdName string) (uint16, bool) {
	cmdId, ok := cmdIdMap[cmdName]
	return cmdId, ok
}
