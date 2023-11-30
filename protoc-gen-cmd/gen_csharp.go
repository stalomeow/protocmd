package main

import (
	"flag"
	"google.golang.org/protobuf/reflect/protoreflect"
	"path"
	"strings"
)

type csharpGenerator struct {
	allTypeFullNames []string
	initClassName    string
	initClassNs      string
}

func init() {
	registerLangGenerator(&csharpGenerator{})
}

func (*csharpGenerator) langName() string {
	return "csharp"
}

func (gen *csharpGenerator) initGenerator(context *generateContext) error {
	gen.allTypeFullNames = make([]string, 0)

	flags := flag.FlagSet{}
	flags.StringVar(&gen.initClassName, "init_class_name", "CmdMessageLoader", "")
	flags.StringVar(&gen.initClassNs, "init_class_ns", "", "")
	return context.writeArgsToFlagSet(&flags)
}

func (gen *csharpGenerator) generate(context *generateContext) error {
	err := gen.initGenerator(context)
	if err != nil {
		return err
	}

	files, err := context.filterFilesToGenerate()
	if err != nil {
		return err
	}

	for _, f := range files {
		gf := genFile{}
		// TODO: Add Header
		gf.println("using Google.Protobuf;")
		gf.println()

		ns := getCSharpFileNamespace(f)
		if ns != "" {
			gf.println("namespace ", ns, " {")
			gf.indent(1)
		}

		prevTypeCount := len(gen.allTypeFullNames)
		gen.writeMsg(context, f.Messages(), &gf, ns)
		if len(gen.allTypeFullNames) <= prevTypeCount {
			continue
		}

		if ns != "" {
			gf.indent(-1)
			gf.println("}")
		}

		filename := getCSharpFileNameWithoutExt(f)
		context.addGenFile(filename+".cmd.cs", &gf)
	}

	gen.writeInitClass(context)
	return nil
}

func (gen *csharpGenerator) writeMsg(context *generateContext, messages protoreflect.MessageDescriptors, gf *genFile, ns string) {
	for i := 0; i < messages.Len(); i++ {
		msg := messages.Get(i)

		cmdId, ok := context.config.CmdIdMap[string(msg.FullName())]
		if !ok {
			continue
		}

		cmdName := string(msg.Name())
		typeFullName := cmdName
		if ns != "" {
			typeFullName = ns + "." + typeFullName
		}
		gen.allTypeFullNames = append(gen.allTypeFullNames, typeFullName)

		gf.println("partial class ", cmdName, " : ICmdMessage")
		gf.println("{")
		gf.indent(1)
		gf.println("public static ushort CmdId => ", cmdId, ";")
		gf.println("ushort ICmdMessage.CmdId => ", cmdId, ";")
		gf.println()
		gf.println("public static string CmdName => \"", cmdName, "\";")
		gf.println("string ICmdMessage.CmdName => \"", cmdName, "\";")

		if nestedMessages := msg.Messages(); nestedMessages.Len() > 0 {
			gf.println()
			gf.println("partial class Types")
			gf.println("{")
			gf.indent(1)

			gen.writeMsg(context, nestedMessages, gf, typeFullName+".Types")

			gf.indent(-1)
			gf.println("}")
		}

		gf.indent(-1)
		gf.println("}")
	}
}

func (gen *csharpGenerator) writeInitClass(context *generateContext) {
	if len(gen.allTypeFullNames) <= 0 {
		return
	}

	gf := genFile{}

	if gen.initClassNs != "" {
		gf.println("namespace ", gen.initClassNs)
		gf.println("{")
		gf.indent(1)
	}

	gf.println("internal static partial class ", gen.initClassName)
	gf.println("{")
	gf.indent(1)
	gf.println("[UnityEngine.RuntimeInitializeOnLoadMethod]")
	gf.println("private static void InitCmdMessages()")
	gf.println("{")
	gf.indent(1)

	for _, t := range gen.allTypeFullNames {
		gf.println("Google.Protobuf.CmdMessageUtility.RegisterCmd(", t, ".CmdId, ", t, ".CmdName, ", t, ".Parser);")
	}

	gf.indent(-1)
	gf.println("}")
	gf.indent(-1)
	gf.println("}")

	if gen.initClassNs != "" {
		gf.indent(-1)
		gf.println("}")
	}

	context.addGenFile(gen.initClassName+".cs", &gf)
}

func getCSharpFileNameWithoutExt(desc protoreflect.FileDescriptor) string {
	filename := string(desc.Name())
	filename = strings.TrimSuffix(filename, path.Ext(filename))
	return underscoresToCamelCase(filename, true, true)
}

func getCSharpFileNamespace(desc protoreflect.FileDescriptor) string {
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