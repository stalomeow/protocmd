// Code generated by protoc-gen-cmdid. DO NOT EDIT.
// versions:
// 	protoc-gen-cmdid v1.0.0
// 	protoc           v4.24.2
// source: test.proto

package protos

import (
	protocmd "github.com/stalomeow/protocmd"
)

const (
	SearchRequest_CmdId   uint16 = 1
	SearchRequest_CmdName string = "SearchRequest"
)

func (*SearchRequest) CmdId() uint16   { return SearchRequest_CmdId }
func (*SearchRequest) CmdName() string { return SearchRequest_CmdName }

const (
	SearchResponse_CmdId   uint16 = 2
	SearchResponse_CmdName string = "SearchResponse"
)

func (*SearchResponse) CmdId() uint16   { return SearchResponse_CmdId }
func (*SearchResponse) CmdName() string { return SearchResponse_CmdName }

const (
	SearchResponse_NestedData_CmdId   uint16 = 3
	SearchResponse_NestedData_CmdName string = "SearchResponse_NestedData"
)

func (*SearchResponse_NestedData) CmdId() uint16   { return SearchResponse_NestedData_CmdId }
func (*SearchResponse_NestedData) CmdName() string { return SearchResponse_NestedData_CmdName }

func init() {
	protocmd.Register(func() protocmd.CmdMessage { return new(SearchRequest) })
	protocmd.Register(func() protocmd.CmdMessage { return new(SearchResponse) })
	protocmd.Register(func() protocmd.CmdMessage { return new(SearchResponse_NestedData) })
}
