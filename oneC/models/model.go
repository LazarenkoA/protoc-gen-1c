package models

type ServiceMethod struct {
	Name                      string
	HttpMethod                string
	Url                       string
	RequiredBodyParams        []string
	RequiredQueryParamsParams []string
	BodyParams                []string // параметры которые передаются через тело запроса в жсоне
	QueryParams               []string // параметры строки запроса /v1/customers?page=10&page_size=20
	PathParams                []string // параметры пути запроса /v1/customers/{id}
}

type Service struct {
	Name    string
	BaseUrl string
	Methods []*ServiceMethod
}
