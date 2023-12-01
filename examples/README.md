# Examples

Directories and files:

- `go`: Go example.
- `protos`: Proto source files.
- `unity`: Unity example.
- `cmd.yaml`: The configuration file.

## Command generating Go code

```
protoc --go_out=go --go_opt=module=github.com/stalomeow/protocmd/examples/go --cmd_out=go --cmd_opt=lang=go,module=github.com/stalomeow/protocmd/examples/go protos/*.proto -I=protos
```

## Command generating C# code

```
protoc --csharp_out=unity --csharp_opt=base_namespace=Examples.CSharp --cmd_out=unity --cmd_opt=lang=csharp,base_namespace=Examples.CSharp protos/*.proto -I=protos
```
