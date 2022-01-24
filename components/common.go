package components

import (
	"google.golang.org/protobuf/compiler/protogen"
	"gopkg.in/yaml.v2"
)

type Component interface {
	generate(file *protogen.File)
}

func Get(file *protogen.File, componentType string) []byte {
	var component Component

	switch componentType {
	case "message":
		component = &Messages{}
	case "service":
		component = &Services{}
	}

	component.generate(file)
	out, _ := yaml.Marshal(&component)
	return out
}
