// <auto-generated>
//     Generated by protoc-gen-cmd v1.0.0.  DO NOT EDIT!
// </auto-generated>

using pb = global::Google.Protobuf;

namespace Examples.CSharp
{
    public partial class MessageHelpers : pb::BaseMessageHelpers
    {
        public MessageHelpers()
        {
            this.Register(Examples.CSharp.Protos.TestReq.CmdId, Examples.CSharp.Protos.TestReq.CmdName, Examples.CSharp.Protos.TestReq.Parser, () => Examples.CSharp.Protos.TestReq.Descriptor);
            this.Register(Examples.CSharp.Protos.TestRsp.CmdId, Examples.CSharp.Protos.TestRsp.CmdName, Examples.CSharp.Protos.TestRsp.Parser, () => Examples.CSharp.Protos.TestRsp.Descriptor);
            this.Register(Examples.CSharp.Protos.TestRsp.Types.TransformInfo.CmdId, Examples.CSharp.Protos.TestRsp.Types.TransformInfo.CmdName, Examples.CSharp.Protos.TestRsp.Types.TransformInfo.Parser, () => Examples.CSharp.Protos.TestRsp.Types.TransformInfo.Descriptor);
        }
    }
}
