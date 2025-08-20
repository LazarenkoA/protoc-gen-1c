package pkg

import (
	onecmodels "1c-grpc-gateway/oneC/models"
	custompb "1c-grpc-gateway/proto/gen"
	annotationspb "google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"net/http"
)

func protogenServiceToOneC(service *protogen.Service) *onecmodels.Service {
	oneCService := &onecmodels.Service{
		Name:    service.GoName,
		BaseUrl: getServiceOptions[string](service, custompb.E_BaseUrl),
		Methods: make([]onecmodels.ServiceMethod, 0, len(service.Methods)),
	}

	for _, method := range service.Methods {
		httpMethod, url := "", ""
		methodOpts := method.Desc.Options().(*descriptorpb.MethodOptions)
		if proto.HasExtension(methodOpts, annotationspb.E_Http) {
			httpRule := proto.GetExtension(methodOpts, annotationspb.E_Http).(*annotationspb.HttpRule)
			switch pattern := httpRule.Pattern.(type) {
			case *annotationspb.HttpRule_Get:
				httpMethod = http.MethodGet
				url = pattern.Get
			case *annotationspb.HttpRule_Post:
				httpMethod = http.MethodPost
				url = pattern.Post
			case *annotationspb.HttpRule_Delete:
				httpMethod = http.MethodDelete
				url = pattern.Delete
			case *annotationspb.HttpRule_Patch:
				httpMethod = http.MethodPatch
				url = pattern.Patch
			case *annotationspb.HttpRule_Put:
				httpMethod = http.MethodPut
				url = pattern.Put
			}
		}

		oneCService.Methods = append(oneCService.Methods, onecmodels.ServiceMethod{
			Name:           method.GoName,
			HttpMethod:     httpMethod,
			Url:            url,
			RequiredFields: getRequiredFields(method),
		})
	}

	return oneCService
}
