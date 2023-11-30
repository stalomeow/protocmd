package main

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", filepath.Base(os.Args[0]), err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) > 1 {
		return fmt.Errorf("unknown argument %q (this program should be run by protoc, not directly)", os.Args[1])
	}

	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		return err
	}

	resp, err := response(req)
	if err != nil {
		return err
	}

	out, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	if _, err := os.Stdout.Write(out); err != nil {
		return err
	}
	return nil
}

func response(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	context := generateContext{}
	err := context.init(req)
	if err != nil {
		return nil, err
	}

	if lang, ok := context.popArg("lang"); ok {
		generator, err := getGeneratorByLang(lang)
		if err != nil {
			return nil, err
		}

		err = generator.generate(&context)
		if err != nil {
			return nil, err
		}
		return context.rsp, nil
	}
	return nil, errors.New("lang is not specified")
}
