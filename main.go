package main

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"gopkg.in/yaml.v2"
)

type Field struct {
	Name   string `yaml:"name"`
	Number int    `yaml:"number"`
}

type Message struct {
	Name   string   `yaml:"name"`
	Fields []*Field `yaml:"fields"`
}

type Messages struct {
	Messages []*Message
}

type Method struct {
	Name   string `yaml:"name"`
	Input  string `yaml:"input_type"`
	Output string `yaml:"output_type"`
}

type Service struct {
	Name    string    `yaml:"name"`
	Methods []*Method `yaml:"methods"`
}

type Services struct {
	Services []*Service
}

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for p, f := range gen.FilesByPath {
			if !f.Generate {
				continue
			}
			generateFile(gen, f, p)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File, path string) {
	filename := path + ".yaml"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	messages := generateMessages(file)
	services := generateServices(file)

	out_messages, _ := yaml.Marshal(&messages)
	out_services, _ := yaml.Marshal(&services)

	out := append(out_messages, out_services...)
	g.P(string(out))
}

func generateMessages(file *protogen.File) *Messages {
	messages := Messages{}
	for _, message := range file.Messages {

		desc := message.Desc
		messages.Messages = append(messages.Messages, generateMessage(desc))

		n := desc.Messages().Len()
		for i := 0; i < n; i++ {
			desc = desc.Messages().Get(i)
			messages.Messages = append(messages.Messages, generateMessage(desc))
		}
	}
	return &messages
}

func generateMessage(desc protoreflect.MessageDescriptor) *Message {
	var fields []*Field
	n := desc.Fields().Len()

	for i := 0; i < n; i++ {
		field := desc.Fields().ByNumber(protoreflect.FieldNumber(i + 1))
		fields = append(fields, &Field{
			Name:   string(field.Name()),
			Number: int(field.Number()),
		})
	}
	return &Message{
		Name:   string(desc.FullName()),
		Fields: fields,
	}
}

func generateServices(file *protogen.File) *Services {
	services := Services{}
	for _, service := range file.Services {
		var methods []*Method
		desc := service.Desc
		n := desc.Methods().Len()

		for i := 0; i < n; i++ {
			method := desc.Methods().Get(i)
			methods = append(methods, &Method{
				Name:   string(method.Name()),
				Input:  string(method.Input().FullName()),
				Output: string(method.Output().FullName()),
			})
		}
		services.Services = append(services.Services, &Service{
			Name:    string(desc.FullName()),
			Methods: methods,
		})
	}
	return &services
}
