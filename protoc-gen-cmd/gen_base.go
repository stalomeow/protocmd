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
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strings"
)

const genVersion = "1.0.0"

type cmdYamlConfig struct {
	CmdIdMap map[string]uint16
}

type generateContext struct {
	req     *pluginpb.CodeGeneratorRequest
	rsp     *pluginpb.CodeGeneratorResponse
	rawArgs map[string]string
	config  *cmdYamlConfig
}

func (context *generateContext) init(req *pluginpb.CodeGeneratorRequest) error {
	context.req = req
	context.rsp = new(pluginpb.CodeGeneratorResponse)
	context.rawArgs = parseArgs(req.GetParameter())
	context.config = new(cmdYamlConfig)
	return loadYamlConfig(context)
}

func parseArgs(rawParameter string) map[string]string {
	args := make(map[string]string)
	for _, param := range strings.Split(rawParameter, ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}

		if param != "" {
			args[param] = value
		}
	}
	return args
}

func loadYamlConfig(context *generateContext) error {
	yamlName := "cmd.yaml"
	if v, ok := context.popArg("config"); ok {
		yamlName = v
	}

	buf, err := os.ReadFile(yamlName)
	if err != nil {
		return err
	}
	context.config.CmdIdMap = make(map[string]uint16)
	return yaml.Unmarshal(buf, &context.config.CmdIdMap)
}

func (context *generateContext) popArg(key string) (string, bool) {
	value, ok := context.rawArgs[key]
	delete(context.rawArgs, key)
	return value, ok
}

func (context *generateContext) writeArgsToFlagSet(flags *flag.FlagSet) error {
	for k, v := range context.rawArgs {
		err := flags.Set(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (context *generateContext) filterFilesToGenerate() ([]protoreflect.FileDescriptor, error) {
	genFileMap := make(map[string]interface{})
	for _, fileName := range context.req.FileToGenerate {
		genFileMap[fileName] = nil
	}

	fileReg := new(protoregistry.Files)
	results := make([]protoreflect.FileDescriptor, 0)
	for _, f := range context.req.ProtoFile {
		if _, ok := genFileMap[f.GetName()]; !ok {
			continue
		}

		desc, err := protodesc.NewFile(f, fileReg)
		if err != nil {
			return nil, fmt.Errorf("invalid FileDescriptorProto %q: %v", f.GetName(), err)
		}
		if err := fileReg.RegisterFile(desc); err != nil {
			return nil, fmt.Errorf("cannot register descriptor %q: %v", f.GetName(), err)
		}

		results = append(results, desc)
	}
	return results, nil
}

func (context *generateContext) addGenFile(filename string, g *genFile) {
	context.rsp.File = append(context.rsp.File, &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(filename),
		Content: proto.String(g.buf.String()),
	})
}

type langGenerator interface {
	langName() string
	generate(context *generateContext) error
}

var langMap = make(map[string]reflect.Type)

func registerLangGenerator(gen langGenerator) {
	langMap[gen.langName()] = reflect.TypeOf(gen).Elem()
}

func getGeneratorByLang(langName string) (langGenerator, error) {
	if t, ok := langMap[langName]; ok {
		return reflect.New(t).Interface().(langGenerator), nil
	}
	return nil, fmt.Errorf("lang '%s' was not registered", langName)
}

type genFile struct {
	buf         bytes.Buffer
	indentCount int
}

func (g *genFile) println(v ...interface{}) {
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

func (g *genFile) indent(count int) {
	g.indentCount += count

	if g.indentCount < 0 {
		g.indentCount = 0
	}
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
