# protocmd

This project provides an easy way to manage a mapping between uint16 ids (a.k.a. CmdId/MsgId) and protobuf Messages. A typical use is encoding and decoding binary network packet.

Currently, only support Go and Unity (C#).

This project is comprised of two components:

- Runtime library for Go: This provides the ability to create a protobuf Message using a predefined uint16 id.

- Code generator: The `protoc-gen-cmd` tool is a compiler plugin to protoc, the protocol buffer compiler. It generates code that maps ids to messages and vice versa.

The runtime library for Unity is [github.com/stalomeow/Protobuf-Unity](https://github.com/stalomeow/Protobuf-Unity).

## Installation

```
go get github.com/stalomeow/protocmd
```

## Code Generation

### Configuration

Create a yaml file to configure the mapping between messages and ids.

``` yaml
# format:
# message-full-name: uint16-cmd-id

# First Group
protocmd.examples.TestReq: 1010
protocmd.examples.TestRsp: 1011

# Second Group
protocmd.examples.TestRsp.TransformInfo: 2010
# protocmd.examples.Vector3: 2055
```

### Generate Go code

Options:

- `lang`: Output language. Must be `go` here.
- `config`: The configuration file name. The default value is `cmd.yaml`.
- Options used by `protoc-gen-go` are also supported.

Example:

```
protoc --cmd_out=go --cmd_opt=lang=go,module=github.com/stalomeow/protocmd/examples/go protos/*.proto -I=protos
```

### Generate C# code

Options:

- `lang`: Output language. Must be `csharp` here.
- `config`: The configuration file name. The default value is `cmd.yaml`.
- [`base_namespace`](https://protobuf.dev/reference/csharp/csharp-generated/#compiler_options): When this option is specified, the generator creates a directory hierarchy for generated source code corresponding to the namespaces of the generated classes, using the value of the option to indicate which part of the namespace should be considered as the "base" for the output directory.
- `msg_helpers_name`: The name of a generated class (`msg_helpers`) holding all messages that have cmdIds. The default value is `MessageHelpers`.
- `msg_helpers_ns`: The namespace of `msg_helpers`. If `base_namespace` is specified, the default value is `base_namespace` otherwise empty.

Example:

```
protoc --cmd_out=unity --cmd_opt=lang=csharp,base_namespace=Examples.CSharp protos/*.proto -I=protos
```

## Runtime Usage

[Examples](/examples/)

### Go

``` go
package main

import (
	"fmt"
	"github.com/stalomeow/protocmd"
	"github.com/stalomeow/protocmd/examples/go/protos" // load all generated protos
	"google.golang.org/protobuf/proto"
)

func main() {
	cmdId := protos.TestReq_CmdId

	buf, err := proto.Marshal(&protos.TestReq{Uid: "123321"})
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}

	// Create message by cmdId
	msg, err := protocmd.NewMessageByCmdId(cmdId)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	err = proto.Unmarshal(buf, msg)
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}

	// Get cmdId and name
	fmt.Printf("Cmd (id: %v, name: %s)\n", msg.CmdId(), msg.CmdName())
	fmt.Println(msg.(*protos.TestReq))
}
```

### Unity

``` c#
using Examples.CSharp.Protos;
using Examples.CSharp;
using Google.Protobuf;
using UnityEngine;

public class Example : MonoBehaviour
{
    public static MessageHelpers Helpers = new();

    private void Start()
    {
        byte[] bytes = new TestReq() { Uid = "123321" }.ToByteArray();

		// Create message by cmdId
        MessageParser parser = Helpers.GetMessageParserByCmdId(TestReq.CmdId);
        ICmdMessage msg = parser.ParseFrom(bytes) as ICmdMessage;

		// Get cmdId and name
        print(msg.CmdId);
        print(msg.CmdName);
        print((TestReq)msg);
    }
}
```
