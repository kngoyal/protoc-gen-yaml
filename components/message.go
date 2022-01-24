package components

import (
	"sync"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
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
	sync.WaitGroup `yaml:"wg,omitempty"`
	sync.Mutex     `yaml:"mutex,omitempty"`
	Messages       []*Message
}

func (msgs *Messages) append(msg *Message) {
	msgs.Lock()
	defer msgs.Unlock()
	// log.Println(msg)
	msgs.Messages = append(msgs.Messages, msg)
}

func (messages *Messages) generate(file *protogen.File) {
	for _, message := range file.Messages {
		messages.Add(1)
		go func(desc protoreflect.MessageDescriptor) {
			defer messages.Done()
			messages.getFromDesc(desc)
		}(message.Desc)
	}
	messages.Wait()
}

func (messages *Messages) getFromDesc(desc protoreflect.MessageDescriptor) {
	messages.Add(1)
	go func(desc protoreflect.MessageDescriptor) {
		defer messages.Done()
		messages.buildFromDesc(desc)
	}(desc)
	messages.Add(1)
	go func(desc protoreflect.MessageDescriptor) {
		defer messages.Done()
		messages.getFromNestedDesc(desc)
	}(desc)
}

func (messages *Messages) getFromNestedDesc(desc protoreflect.MessageDescriptor) {
	n := desc.Messages().Len()
	for i := 0; i < n; i++ {
		messages.Add(1)
		go func(desc protoreflect.MessageDescriptor, i int) {
			defer messages.Done()
			messages.getFromDesc(desc.Messages().Get(i))
		}(desc, i)
	}
}

func (messages *Messages) buildFromDesc(desc protoreflect.MessageDescriptor) {
	var fields []*Field
	n := desc.Fields().Len()

	for i := 0; i < n; i++ {
		field := desc.Fields().ByNumber(protoreflect.FieldNumber(i + 1))
		fields = append(fields, &Field{
			Name:   string(field.Name()),
			Number: int(field.Number()),
		})
	}
	messages.append(&Message{
		Name:   string(desc.FullName()),
		Fields: fields,
	})
}
