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

        MessageParser parser = Helpers.GetMessageParserByCmdId(TestReq.CmdId);
        ICmdMessage msg = parser.ParseFrom(bytes) as ICmdMessage;

        print(msg.CmdId);
        print(msg.CmdName);
        print((TestReq)msg);
    }
}
