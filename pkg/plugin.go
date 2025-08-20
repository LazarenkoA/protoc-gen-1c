package pkg

import (
	"1c-grpc-gateway/oneC"
	"1c-grpc-gateway/pkg/utils"
	custompb "1c-grpc-gateway/proto/gen"
	"1c-grpc-gateway/swagger"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
	"log/slog"
	"os"
	"strings"
)

const (
	tmpFile = "csdcvfdv"
)

type Plugin struct {
	logger *utils.Logger
}

func NewPlugin() *Plugin {
	return &Plugin{
		logger: utils.NewLogger(slog.LevelDebug),
	}
}

func (p *Plugin) Process(plugin *protogen.Plugin) (err error) {
	defer func() {
		if err != nil {
			p.logger.Error(err.Error())
		}
	}()

	args := parseParams(plugin)
	if !needLogger(args) {
		p.logger.Disable()
	}

	defer func() {
		os.Remove(tmpFile)
	}()

	p.logger.Info("start process")

	plugin.NewGeneratedFile(tmpFile, "") // нужно потому что просит protoc, без этого он начнет выдавать "First file chunk returned by plugin did not specify a file name"

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		root := ""
		opts := file.Proto.GetOptions()
		if proto.HasExtension(opts, custompb.E_ConfigurationCatalog) {
			root = proto.GetExtension(opts, custompb.E_ConfigurationCatalog).(string)
		}

		if needSwagger(args) {
			openApiData := swagger.Generate(plugin, p.logger)
			if err := oneC.CreateSwaggerHttpService(root, openApiData, p.logger); err != nil {
				return errors.Wrap(err, "create swagger http service error")
			}

			//todo удалить
			//f, _ := os.Create("swagger.yaml")
			//f.Write(openApiData)
			//f.Close()
		}

		// Для всех сервисов
		for _, service := range file.Services {
			if err := oneC.CreateHttpService(root, protogenServiceToOneC(service), p.logger); err != nil {
				return errors.Wrap(err, "create http service error")
			}
		}
	}
	return nil
}

func getServiceOptions[T any](service *protogen.Service, extInfo *protoimpl.ExtensionInfo) T {
	opts := service.Desc.Options().(*descriptorpb.ServiceOptions)

	var value T
	if proto.HasExtension(opts, extInfo) {
		value, _ = proto.GetExtension(opts, extInfo).(T)
	}

	return value
}

func getMethodOptions[T any](method *protogen.Method, extInfo *protoimpl.ExtensionInfo) T {
	opts := method.Desc.Options().(*descriptorpb.MethodOptions)

	var value T
	if proto.HasExtension(opts, extInfo) {
		value, _ = proto.GetExtension(opts, extInfo).(T)
	}

	return value
}

func getFieldOptions[T any](field *protogen.Field, extInfo *protoimpl.ExtensionInfo) T {
	opts := field.Desc.Options().(*descriptorpb.FieldOptions)

	var value T
	if proto.HasExtension(opts, extInfo) {
		value, _ = proto.GetExtension(opts, extInfo).(T)
	}

	return value
}

func getRequiredFields(method *protogen.Method) []string {
	var result []string
	for _, fld := range method.Input.Fields {
		if getFieldOptions[bool](fld, custompb.E_Required) {
			result = append(result, fld.Desc.TextName())
		}
	}

	return result
}

func parseParams(plugin *protogen.Plugin) map[string]string {
	params := plugin.Request.GetParameter()
	result := map[string]string{}
	for _, param := range strings.Split(params, ",") {
		if kv := strings.Split(param, "="); len(kv) == 2 {
			result[kv[0]] = kv[1]
		}
	}
	return result
}

func needSwagger(args map[string]string) bool {
	v, _ := args["swagger"]
	return v == "1"
}

func needLogger(args map[string]string) bool {
	v, _ := args["logger"]
	return v == "1"
}
