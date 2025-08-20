package models

type ServiceMethod struct {
	Name           string
	HttpMethod     string
	Url            string
	RequiredFields []string
}

type Service struct {
	Name    string
	BaseUrl string
	Methods []ServiceMethod
}
