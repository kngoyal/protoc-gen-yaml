package components

import (
	"sync"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

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
	sync.WaitGroup `yaml:"wg,omitempty"`
	sync.Mutex     `yaml:"mutex,omitempty"`
	Services       []*Service
}

func (srvcs *Services) append(srvc *Service) {
	srvcs.Lock()
	defer srvcs.Unlock()
	// log.Println(srvc)
	srvcs.Services = append(srvcs.Services, srvc)
}

func (services *Services) generate(file *protogen.File) {
	for _, service := range file.Services {
		services.Add(1)
		go func(service *protogen.Service) {
			defer services.Done()
			services.buildFromDesc(service.Desc)
		}(service)
	}
	services.Wait()
}

func (services *Services) buildFromDesc(desc protoreflect.ServiceDescriptor) {
	var methods []*Method
	n := desc.Methods().Len()

	for i := 0; i < n; i++ {
		method := desc.Methods().Get(i)
		methods = append(methods, &Method{
			Name:   string(method.Name()),
			Input:  string(method.Input().FullName()),
			Output: string(method.Output().FullName()),
		})
	}
	services.append(&Service{
		Name:    string(desc.FullName()),
		Methods: methods,
	})
}
