package main

import (
	"google.golang.org/protobuf/compiler/protogen"

	"github.com/protoc-gen-yaml/components"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for path, f := range gen.FilesByPath {
			if !f.Generate {
				continue
			}
			generateFile(gen, f, path)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File, path string) {
	filename := path + ".yaml"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	out_messages := components.Get(file, "message")
	out_services := components.Get(file, "service")

	out := append(out_messages, out_services...)
	// log.Println(string(out))
	g.P(string(out))
}
