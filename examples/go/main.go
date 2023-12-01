package main

import (
	"fmt"
	"github.com/stalomeow/protocmd"
	"github.com/stalomeow/protocmd/examples/go/protos" // load all protos
)

func main() {
	// Get loaded cmd count
	fmt.Printf("Load %v cmd\n", protocmd.CmdCount())
	fmt.Println(protocmd.AllCmdIds())

	// Get cmd name by id
	fmt.Println(protocmd.CmdName(protos.TestRsp_TransformInfo_CmdId))

	// Create cmd by id
	msg, err := protocmd.NewMessageByCmdId(protos.TestReq_CmdId)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}

	// Get cmd id and name
	fmt.Printf("Cmd (id: %v, name: %s)\n", msg.CmdId(), msg.CmdName())

	req := msg.(*protos.TestReq)
	req.Uid = "123321"
	fmt.Println(req)
}
