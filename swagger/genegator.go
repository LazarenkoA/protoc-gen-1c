package swagger

import (
	"1c-grpc-gateway/pkg/utils"
	custompb "1c-grpc-gateway/proto/gen"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/gnostic/cmd/protoc-gen-openapi/generator"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"gopkg.in/yaml.v3"
	"maps"
	"slices"
)

func Generate(plugin *protogen.Plugin, log *utils.Logger) []byte {
	conf := generator.Configuration{
		Version:         utils.Ptr("0.0.1"),
		Title:           utils.Ptr(""),
		Description:     utils.Ptr(""),
		Naming:          utils.Ptr("proto"),
		FQSchemaNaming:  utils.Ptr(false),
		EnumType:        utils.Ptr("integer"),
		CircularDepth:   utils.Ptr(2),
		DefaultResponse: utils.Ptr(false),
		OutputMode:      utils.Ptr("merged"),
	}

	log.Info("generate openAPI")

	outputFile := plugin.NewGeneratedFile("", "")
	generator.NewOpenAPIv3Generator(plugin, conf, plugin.Files).Run(outputFile)
	data, _ := outputFile.Content()

	return setRequired(plugin, data, log)
}

// setRequired установка признака "обязательный"
func setRequired(plugin *protogen.Plugin, data []byte, log *utils.Logger) []byte {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		log.Error(errors.Wrap(err, "openapi3 load data error").Error())
		return []byte{}
	}

	required := getRequiredFields(plugin)

	updateParam := func(op *openapi3.Operation) {
		if op == nil {
			return
		}

		for _, param := range op.Parameters {
			if _, ok := required[param.Value.Name]; ok {
				param.Value.Required = true
			}
		}
	}

	for _, path := range doc.Paths.Map() {
		updateParam(path.Get)
		updateParam(path.Post)
		updateParam(path.Put)
		updateParam(path.Patch)
		updateParam(path.Options)
		updateParam(path.Delete)
	}

	for _, schema := range doc.Components.Schemas {
		schema.Value.Required = slices.Collect(maps.Keys(required))
	}

	newData, _ := doc.MarshalYAML()
	outData, _ := yaml.Marshal(&newData)
	return outData
}

func getRequiredFields(plugin *protogen.Plugin) map[string]struct{} {
	result := map[string]struct{}{}

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		for _, m := range file.Messages {
			for _, field := range m.Fields {
				opts := field.Desc.Options().(*descriptorpb.FieldOptions)
				if proto.HasExtension(opts, custompb.E_Required) && proto.GetExtension(opts, custompb.E_Required).(bool) {
					result[field.Desc.TextName()] = struct{}{}
				}
			}
		}
	}

	return result
}
