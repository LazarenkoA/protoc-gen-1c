package swagger

import (
	"1c-grpc-gateway/pkg/utils"
	custompb "1c-grpc-gateway/proto/gen"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/gnostic/cmd/protoc-gen-openapi/generator"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/compiler/protogen"
	"gopkg.in/yaml.v3"
	"maps"
	"net/http"
	"slices"
	"strconv"
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

	return AppendInfo(plugin, data, log)
}

func AppendInfo(plugin *protogen.Plugin, data []byte, log *utils.Logger) []byte {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		log.Error(errors.Wrap(err, "openapi3 load data error").Error())
		return []byte{}
	}

	setRequired(plugin, doc)
	setStatusCodes(plugin, doc)

	newData, _ := doc.MarshalYAML()
	outData, _ := yaml.Marshal(&newData)
	return outData
}

// setRequired установка признака "обязательный"
func setRequired(plugin *protogen.Plugin, doc *openapi3.T) {
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
}

// setStatusCodes установка кодов ответа
func setStatusCodes(plugin *protogen.Plugin, doc *openapi3.T) {
	respCodes := getResponseCodes(plugin)
	update := func(url, method string, op *openapi3.Operation) {
		if op == nil {
			return
		}

		codes, ok := respCodes[method+url]
		if !ok {
			return
		}

		for _, code := range codes {
			op.Responses.Set(strconv.Itoa(int(code.Code)), &openapi3.ResponseRef{Value: &openapi3.Response{
				Description: utils.Ptr(code.Comment),
			}})
		}
	}

	for url, path := range doc.Paths.Map() {
		update(url, http.MethodGet, path.Get)
		update(url, http.MethodPost, path.Post)
		update(url, http.MethodPut, path.Put)
		update(url, http.MethodPatch, path.Patch)
		update(url, http.MethodOptions, path.Options)
		update(url, http.MethodDelete, path.Delete)
	}
}

func getResponseCodes(plugin *protogen.Plugin) map[string][]*custompb.StatusCode {
	result := map[string][]*custompb.StatusCode{}

	for _, f := range plugin.Files {
		for _, svc := range f.Services {
			for _, m := range svc.Methods {
				method, url, _ := utils.GetMethodInfo(m)
				result[method+url] = utils.GetMethodOptions[[]*custompb.StatusCode](m, custompb.E_Codes)
			}
		}
	}

	return result
}

func getRequiredFields(plugin *protogen.Plugin) map[string]struct{} {
	result := map[string]struct{}{}

	for _, file := range plugin.Files {
		if !file.Generate {
			continue
		}

		for _, m := range file.Messages {
			for _, field := range m.Fields {
				if utils.GetFieldOptions[bool](field, custompb.E_Required) {
					result[field.Desc.TextName()] = struct{}{}
				}
			}
		}
	}

	return result
}
