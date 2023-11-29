package main

import (
	"bytes"
	"flag"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/pluginpb"
	"path"
	"strings"
)

type csGenContext struct {
	initClassName    string
	initClassNs      string
	config           *cmdIdConfig
	genTypeFullNames []string
}

type genFile struct {
	buf         bytes.Buffer
	indentCount int
}

func genCSharp(req *pluginpb.CodeGeneratorRequest, config *cmdIdConfig) (*pluginpb.CodeGeneratorResponse, error) {
	context := &csGenContext{
		config:           config,
		genTypeFullNames: make([]string, 0),
	}
	flags := flag.FlagSet{}
	flags.StringVar(&context.initClassName, "init_class_name", "CmdMessageInitializer", "")
	flags.StringVar(&context.initClassNs, "init_class_ns", "", "")

	err := context.config.writeFlags(&flags)
	if err != nil {
		return nil, err
	}

	rsp := &pluginpb.CodeGeneratorResponse{}
	fileReg := new(protoregistry.Files)
	genFileMap := make(map[string]interface{})
	for _, fileName := range req.FileToGenerate {
		genFileMap[fileName] = nil
	}

	for _, f := range req.ProtoFile {
		if _, ok := genFileMap[f.GetName()]; !ok {
			continue
		}

		filename := strings.TrimSuffix(f.GetName(), path.Ext(f.GetName()))
		filename = underscoresToCamelCase(filename, true, true)
		filename += ".cmd.cs"

		desc, err := protodesc.NewFile(f, fileReg)
		if err != nil {
			return nil, fmt.Errorf("invalid FileDescriptorProto %q: %v", f.GetName(), err)
		}
		if err := fileReg.RegisterFile(desc); err != nil {
			return nil, fmt.Errorf("cannot register descriptor %q: %v", f.GetName(), err)
		}

		g := genFile{}
		// TODO: Add Header
		g.P("using Google.Protobuf;")
		g.P()

		ns := getFileNamespace(desc)
		if ns != "" {
			g.P("namespace ", ns, " {")
			g.Indent(1)
		}

		if writeCount := writeCSMessages(context, desc.Messages(), &g, ns); writeCount <= 0 {
			continue
		}

		if ns != "" {
			g.Indent(-1)
			g.P("}")
		}

		rsp.File = append(rsp.File, &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(filename),
			Content: proto.String(g.buf.String()),
		})
	}

	writeInitClass(context, rsp)
	return rsp, nil
}

func writeInitClass(context *csGenContext, rsp *pluginpb.CodeGeneratorResponse) {
	if len(context.genTypeFullNames) <= 0 {
		return
	}

	g := genFile{}

	if context.initClassNs != "" {
		g.P("namespace ", context.initClassNs)
		g.P("{")
		g.Indent(1)
	}

	g.P("internal static partial class ", context.initClassName)
	g.P("{")
	g.Indent(1)
	g.P("[UnityEngine.RuntimeInitializeOnLoadMethod]")
	g.P("private static void Initialize()")
	g.P("{")
	g.Indent(1)

	for _, t := range context.genTypeFullNames {
		g.P("Google.Protobuf.CmdMessageUtility.RegisterCmd(", t, ".CmdId, ", t, ".CmdName, ", t, ".Parser);")
	}

	g.Indent(-1)
	g.P("}")
	g.Indent(-1)
	g.P("}")

	if context.initClassNs != "" {
		g.Indent(-1)
		g.P("}")
	}

	rsp.File = append(rsp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(context.initClassName + ".cs"),
		Content: proto.String(g.buf.String()),
	})
}

func writeCSMessages(context *csGenContext, messages protoreflect.MessageDescriptors, g *genFile, ns string) int {
	var writeCount = 0

	for i := 0; i < messages.Len(); i++ {
		msg := messages.Get(i)
		cmdId, ok := context.config.cmdIdMap[string(msg.FullName())]
		if !ok {
			continue
		}
		writeCount++

		cmdName := string(msg.Name())
		typeFullName := cmdName
		if ns != "" {
			typeFullName = ns + "." + typeFullName
		}
		context.genTypeFullNames = append(context.genTypeFullNames, typeFullName)

		g.P("partial class ", cmdName, " : ICmdMessage")
		g.P("{")
		g.Indent(1)
		g.P("public static ushort CmdId => ", cmdId, ";")
		g.P("ushort ICmdMessage.CmdId => ", cmdId, ";")
		g.P()
		g.P("public static string CmdName => \"", cmdName, "\";")
		g.P("string ICmdMessage.CmdName => \"", cmdName, "\";")

		if nestedMessages := msg.Messages(); nestedMessages.Len() > 0 {
			g.P()
			g.P("partial class Types")
			g.P("{")
			g.Indent(1)

			writeCount += writeCSMessages(context, nestedMessages, g, typeFullName+".Types")

			g.Indent(-1)
			g.P("}")
		}

		g.Indent(-1)
		g.P("}")
	}

	return writeCount
}

func getFileNamespace(desc protoreflect.FileDescriptor) string {
	options := desc.Options().ProtoReflect()
	fields := options.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		if f.Name() == "csharp_namespace" {
			if ns := options.Get(f).String(); ns != "" {
				return ns
			}
		}
	}

	return underscoresToCamelCase(string(desc.Package()), true, true)
}

func underscoresToCamelCase(input string, capNextLetter, preservePeriod bool) string {
	var result strings.Builder

	for i := 0; i < len(input); i++ {
		if 'a' <= input[i] && input[i] <= 'z' {
			if capNextLetter {
				result.WriteByte(input[i] + 'A' - 'a')
			} else {
				result.WriteByte(input[i])
			}
			capNextLetter = false
		} else if 'A' <= input[i] && input[i] <= 'Z' {
			if i == 0 && !capNextLetter {
				result.WriteByte(input[i] + ('a' - 'A'))
			} else {
				result.WriteByte(input[i])
			}
			capNextLetter = false
		} else if '0' <= input[i] && input[i] <= '9' {
			result.WriteByte(input[i])
			capNextLetter = true
		} else {
			capNextLetter = true
			if input[i] == '.' && preservePeriod {
				result.WriteByte('.')
			}
		}
	}

	if len(input) > 0 && input[len(input)-1] == '#' {
		result.WriteByte('_')
	}

	if result.Len() > 0 && ('0' <= result.String()[0] && result.String()[0] <= '9') &&
		len(input) > 0 && input[0] == '_' {
		result.WriteString("_")
	}

	return result.String()
}

func (g *genFile) P(v ...interface{}) {
	if len(v) > 0 {
		for i := 0; i < g.indentCount; i++ {
			fmt.Fprint(&g.buf, "    ")
		}
		for _, x := range v {
			fmt.Fprint(&g.buf, x)
		}
	}
	fmt.Fprintln(&g.buf)
}

func (g *genFile) Indent(count int) {
	g.indentCount += count

	if g.indentCount < 0 {
		g.indentCount = 0
	}
}
