package main

import (
	"fmt"
	"github.com/stalomeow/protocmd"
	"github.com/stalomeow/protocmd/examples/go/protos" // load all generated protos
	"google.golang.org/protobuf/proto"
)

func main() {
	// Get loaded cmd count
	fmt.Printf("Load %v cmd\n", protocmd.CmdCount())
	fmt.Println(protocmd.AllCmdIds())

	// Get cmdName by id
	if cmdName, ok := protocmd.CmdName(protos.TestRsp_TransformInfo_CmdId); ok {
		fmt.Println(cmdName)
	}

	buf, err := proto.Marshal(&protos.TestReq{Uid: "123321"})
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}

	// Create message by cmdId
	msg, err := protocmd.NewMessageByCmdId(protos.TestReq_CmdId)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	err = proto.Unmarshal(buf, msg)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}

	fmt.Printf("Cmd (id: %v, name: %s)\n", msg.CmdId(), msg.CmdName())
	fmt.Println(msg.(*protos.TestReq))
}
