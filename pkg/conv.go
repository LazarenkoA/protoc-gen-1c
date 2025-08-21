package pkg

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/compiler/protogen"
	onecmodels "protoc-gen-1c/oneC/models"
	"protoc-gen-1c/pkg/utils"
	custompb "protoc-gen-1c/proto/gen"
)

func protogenServiceToOneC(service *protogen.Service) *onecmodels.Service {
	oneCService := &onecmodels.Service{
		Name:    service.GoName,
		BaseUrl: utils.GetServiceOptions[string](service, custompb.E_BaseUrl),
		Methods: make([]*onecmodels.ServiceMethod, 0, len(service.Methods)),
	}

	for _, method := range service.Methods {
		httpMethod, url, body := utils.GetMethodInfo(method)
		requiredFields := utils.GetRequiredFields(method)

		oneCService.Methods = append(oneCService.Methods, &onecmodels.ServiceMethod{
			Name:        method.GoName,
			HttpMethod:  httpMethod,
			Url:         url,
			BodyParams:  utils.GetBodyParams(method, body),
			QueryParams: utils.GetQueryParams(method, body),
			PathParams:  utils.GetPathParams(url),
			RespCodes:   utils.GetRespCodes(method),
		})

		last := oneCService.Methods[len(oneCService.Methods)-1]
		last.QueryParams, _ = lo.Difference(last.QueryParams, last.PathParams) // PathParams побеждают
		last.RequiredQueryParamsParams = lo.Intersect(requiredFields, last.QueryParams)
		last.RequiredBodyParams = lo.Intersect(requiredFields, last.BodyParams)
	}

	return oneCService
}
