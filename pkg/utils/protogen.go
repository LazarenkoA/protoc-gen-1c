package utils

import (
	custompb "github.com/LazarenkoA/protoc-gen-1c/proto/gen"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/descriptorpb"
	"net/http"
	"regexp"
)

func GetServiceOptions[T any](service *protogen.Service, extInfo *protoimpl.ExtensionInfo) T {
	opts := service.Desc.Options().(*descriptorpb.ServiceOptions)

	var value T
	if proto.HasExtension(opts, extInfo) {
		value, _ = proto.GetExtension(opts, extInfo).(T)
	}

	return value
}

func GetMethodOptions[T any](method *protogen.Method, extInfo *protoimpl.ExtensionInfo) T {
	opts := method.Desc.Options().(*descriptorpb.MethodOptions)

	var value T
	if proto.HasExtension(opts, extInfo) {
		value, _ = proto.GetExtension(opts, extInfo).(T)
	}

	return value
}

func GetMethodInfo(method *protogen.Method) (httpMethod string, url string, body string) {
	httpRule := GetMethodOptions[*annotations.HttpRule](method, annotations.E_Http)
	if httpRule != nil {
		switch pattern := httpRule.Pattern.(type) {
		case *annotations.HttpRule_Get:
			httpMethod = http.MethodGet
			url = pattern.Get
		case *annotations.HttpRule_Post:
			httpMethod = http.MethodPost
			url = pattern.Post
		case *annotations.HttpRule_Delete:
			httpMethod = http.MethodDelete
			url = pattern.Delete
		case *annotations.HttpRule_Patch:
			httpMethod = http.MethodPatch
			url = pattern.Patch
		case *annotations.HttpRule_Put:
			httpMethod = http.MethodPut
			url = pattern.Put
		}
	}

	return httpMethod, url, httpRule.GetBody()
}

func GetFieldOptions[T any](field *protogen.Field, extInfo *protoimpl.ExtensionInfo) T {
	opts := field.Desc.Options().(*descriptorpb.FieldOptions)

	var value T
	if proto.HasExtension(opts, extInfo) {
		value, _ = proto.GetExtension(opts, extInfo).(T)
	}

	return value
}

func GetRequiredFields(method *protogen.Method) []string {
	var result []string
	for _, fld := range method.Input.Fields {
		if GetFieldOptions[bool](fld, custompb.E_Required) {
			result = append(result, fld.Desc.TextName())
		}
	}

	return result
}

func GetRespCodes(method *protogen.Method) map[int32]string {
	result := map[int32]string{}

	codes := GetMethodOptions[[]*custompb.StatusCode](method, custompb.E_Codes)
	for _, cod := range codes {
		result[cod.Code] = cod.Comment
	}

	return result
}

func GetBodyParams(method *protogen.Method, body string) []string {
	var result []string

	if body == "" {
		return result
	}

	for _, fld := range method.Input.Fields {
		result = append(result, fld.Desc.TextName())
	}

	return result
}

func GetQueryParams(method *protogen.Method, body string) []string {
	var result []string

	if body == "*" {
		return result
	}

	for _, fld := range method.Input.Fields {
		result = append(result, fld.Desc.TextName())
	}

	return result
}

func GetPathParams(url string) []string {
	var result []string

	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(url, -1)
	for _, m := range matches {
		result = append(result, m[1])
	}

	return result
}
