package protocmd

import (
	"fmt"
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
var cmdMapDup = make(map[uint16]interface{})

func Register(factory func() CmdMessage) {
	msg := factory()
	cmdId := msg.CmdId()

	if _, ok := cmdMap[cmdId]; ok {
		cmdMapDup[cmdId] = nil
	}

	cmdMap[cmdId] = &cmdInfo{
		name:    msg.CmdName(),
		factory: factory,
	}
}

func Unregister(cmdId uint16) {
	delete(cmdMap, cmdId)
	delete(cmdMapDup, cmdId)
}

func New(cmdId uint16) (CmdMessage, error) {
	info, ok := cmdMap[cmdId]
	if !ok {
		return nil, fmt.Errorf("failed to new cmd with id '%v' which was not registered", cmdId)
	}
	return info.factory(), nil
}

func Count() int {
	return len(cmdMap)
}

func HasDuplicated() bool {
	return len(cmdMapDup) > 0
}

func DuplicatedCmdIds() []uint16 {
	keys := make([]uint16, len(cmdMapDup))
	i := 0
	for k := range cmdMapDup {
		keys[i] = k
		i++
	}
	return keys
}

func Name(cmdId uint16) string {
	return NameOrDefault(cmdId, "<UnknownCmd>")
}

func NameOrDefault(cmdId uint16, defaultName string) string {
	info, ok := cmdMap[cmdId]
	if !ok {
		return defaultName
	}
	return info.name
}
