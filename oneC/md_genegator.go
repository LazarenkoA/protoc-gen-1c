package oneC

import (
	"1c-grpc-gateway/oneC/models"
	"1c-grpc-gateway/pkg/utils"
	"encoding/xml"
	"fmt"
	"github.com/beevik/etree"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"golang.org/x/text/encoding/charmap"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CreateSwaggerHttpService(root string, openApiData []byte, log *utils.Logger) error {
	log.Info("create swagger http service")

	const (
		methodName = "MethodGet"
		tmplName   = "swagger"
	)

	service := &models.Service{
		Name:    "Swagger",
		BaseUrl: "swagger",
		Methods: []*models.ServiceMethod{
			{
				Name:       methodName,
				HttpMethod: http.MethodGet,
				Url:        "/doc",
			},
		},
	}

	httpServicesDir, err := createHttpService(root, service, log)
	if err != nil {
		return err
	}

	// создаем обработчик
	f, _ := os.Create(filepath.Join(httpServicesDir, "Ext", "Module.bsl"))
	defer f.Close()

	builder := strings.Builder{}
	builder.WriteString(disclaimer("//"))
	builder.WriteString("\n\n")

	content := swaggerHandler(&HandlerInfo{
		HandlerName: methodName,
		ServiceName: service.Name,
		TmplName:    tmplName,
	}, log)

	builder.WriteString(content)
	builder.WriteString("\n\n")

	_, _ = f.WriteString(builder.String())

	// создаем общий макет
	return errors.Wrap(CommonTemplate(root, tmplName, openApiData, log), "create common template error")
}

func CommonModules(root string, name string, service *models.Service, log *utils.Logger) error {
	log.Info("create common modules")

	confFile, dir := prepareDirectoryStruct(filepath.Join(root, "CommonModules"), name, false)
	bslPath := filepath.Join(dir, "Ext", "Module.bsl")
	if _, err := os.Stat(bslPath); err == nil {
		return nil // что б не перетирать
	}

	m, err := newCommonModule(root, name)
	if err != nil {
		return err
	}

	log.Debug("store metaDataObject", "confFile", confFile)
	if err := storeMetaDataObject(m, confFile); err != nil {
		return errors.Wrap(err, "storeMetaDataObject")
	}

	f, _ := os.Create(bslPath)

	builder := strings.Builder{}

	// создаем обработчик для каждого метода
	for _, method := range service.Methods {
		content := getCommonModuleContent(&HandlerInfo{
			HandlerName: method.Name,
			ServiceName: service.Name,
			QueryParams: method.QueryParams,
			PathParams:  method.PathParams,
			BodyParams:  method.BodyParams,
			RespCodes:   method.RespCodes,
		}, log)

		builder.WriteString(content)
		builder.WriteString("\n\n")
	}

	_, _ = f.WriteString(builder.String())

	return f.Close()
}

func CommonTemplate(root, name string, data []byte, log *utils.Logger) error {
	log.Info("create common template")

	confFile, dir := prepareDirectoryStruct(filepath.Join(root, "CommonTemplates"), name, true)

	tmpl, err := newCommonTemplate(root, name)
	if err != nil {
		return err
	}

	log.Debug("store metaDataObject", "confFile", confFile)
	if err := storeMetaDataObject(tmpl, confFile); err != nil {
		return errors.Wrap(err, "storeMetaDataObject")
	}

	f, _ := os.Create(filepath.Join(dir, "Ext", "Template.txt"))
	builder := strings.Builder{}
	builder.WriteString(disclaimer("#"))
	builder.WriteString("\n\n")
	builder.WriteString(string(data))

	writer := charmap.Windows1251.NewEncoder().Writer(f)
	writer.Write([]byte(builder.String()))

	return f.Close()
}

func CreateHttpService(root string, service *models.Service, log *utils.Logger) error {
	log.Info("create new http service", "name", service.Name)

	httpServicesDir, err := createHttpService(root, service, log)
	if err != nil {
		return err
	}

	f, _ := os.Create(filepath.Join(httpServicesDir, "Ext", "Module.bsl"))
	builder := strings.Builder{}
	builder.WriteString(disclaimer("//"))
	builder.WriteString("\n\n")

	// создаем обработчик для каждого метода
	for _, method := range service.Methods {
		content := getHttpHandler(&HandlerInfo{
			HandlerName: method.Name,
			ServiceName: service.Name,
		}, log)

		builder.WriteString(content)
		builder.WriteString("\n\n")

		content = checkRequestFields(&HandlerInfo{
			HandlerName:               method.Name,
			ServiceName:               service.Name,
			RequiredQueryParamsParams: method.RequiredQueryParamsParams,
			RequiredBodyParams:        method.RequiredBodyParams,
			BodyParams:                method.BodyParams,
			QueryParams:               method.QueryParams,
			PathParams:                method.PathParams,
			Funcs:                     map[string]any{"join": strings.Join},
		}, log)
		builder.WriteString(content)
		builder.WriteString("\n")

		content = checkResponseFields(&HandlerInfo{
			HandlerName: method.Name,
			ServiceName: service.Name,
			RespCodes:   method.RespCodes,
		}, log)
		builder.WriteString(content)
		builder.WriteString("\n")
	}

	builder.WriteString(boilerplate())
	builder.WriteString("\n\n")

	_, _ = f.WriteString(builder.String())
	_ = f.Close()

	return CommonModules(root, fmt.Sprintf("%sПереопределяемый", service.Name), service, log)
}

func prepareDirectoryStruct(root, name string, rewrite bool) (string, string) {
	newDir := filepath.Join(root, name)
	confFile := filepath.Join(root, fmt.Sprintf("%s.xml", name))

	if rewrite {
		_ = os.Remove(confFile)
		_ = os.RemoveAll(newDir)
	}

	_ = os.MkdirAll(filepath.Join(newDir, "Ext"), os.ModeDir)

	return confFile, newDir
}

// createHttpService создает сервис без обработчиков
func createHttpService(root string, service *models.Service, log *utils.Logger) (string, error) {
	srv, err := newHttpService(root, service)
	if err != nil {
		return "", err
	}

	confFile, httpServicesDir := prepareDirectoryStruct(filepath.Join(root, "HTTPServices"), service.Name, true)

	log.Debug("store metaDataObject", "confFile", confFile)
	if err := storeMetaDataObject(srv, confFile); err != nil {
		return "", errors.Wrap(err, "storeMetaDataObject")
	}

	return httpServicesDir, nil
}

func storeMetaDataObject(md *models.MetaDataObject, path string) error {
	data, _ := xml.MarshalIndent(md, "", "    ")
	f, _ := os.Create(path)
	_, _ = f.Write(data)

	return f.Close()
}

func newHttpService(root string, service *models.Service) (*models.MetaDataObject, error) {
	if err := registerMetaObject(root, "HTTPService", service.Name); err != nil {
		return nil, err
	}

	srv := &models.MetaDataObject{
		Version:    "2.16",
		Xmlns:      "http://v8.1c.ru/8.3/MDClasses",
		XmlnsApp:   "http://v8.1c.ru/8.2/managed-application/core",
		XmlnsCfg:   "http://v8.1c.ru/8.1/data/enterprise/current-config",
		XmlnsCmi:   "http://v8.1c.ru/8.2/managed-application/cmi",
		XmlnsEnt:   "http://v8.1c.ru/8.1/data/enterprise",
		XmlnsLf:    "http://v8.1c.ru/8.2/managed-application/logform",
		XmlnsStyle: "http://v8.1c.ru/8.1/data/ui/style",
		XmlnsSys:   "http://v8.1c.ru/8.1/data/ui/fonts/system",
		XmlnsV8:    "http://v8.1c.ru/8.1/data/core",
		XmlnsV8ui:  "http://v8.1c.ru/8.1/data/ui",
		XmlnsWeb:   "http://v8.1c.ru/8.1/data/ui/colors/web",
		XmlnsWin:   "http://v8.1c.ru/8.1/data/ui/colors/windows",
		XmlnsXen:   "http://v8.1c.ru/8.3/xcf/enums",
		XmlnsXpr:   "http://v8.1c.ru/8.3/xcf/predef",
		XmlnsXr:    "http://v8.1c.ru/8.3/xcf/readable",
		XmlnsXs:    "http://www.w3.org/2001/XMLSchema",
		XmlnsXsi:   "http://www.w3.org/2001/XMLSchema-instance",

		HTTPService: &models.HTTPService{
			UUID: uuid.NewString(),
			Properties: models.Properties{
				Name: service.Name,
				Synonym: models.Synonym{
					Items: []models.SynonymItem{
						{
							Lang:    "ru",
							Content: service.Name,
						},
					},
				},
				Comment:       "",
				RootURL:       service.BaseUrl,
				ReuseSessions: "AutoUse",
				SessionMaxAge: 20,
			},
			ChildObjects: models.ChildObjects{},
		},
	}

	appendMethod(srv.HTTPService, service)
	return srv, nil
}

func newCommonTemplate(root string, name string) (*models.MetaDataObject, error) {
	if err := registerMetaObject(root, "CommonTemplate", name); err != nil {
		return nil, err
	}

	tmpl := &models.MetaDataObject{
		Version:    "2.16",
		Xmlns:      "http://v8.1c.ru/8.3/MDClasses",
		XmlnsApp:   "http://v8.1c.ru/8.2/managed-application/core",
		XmlnsCfg:   "http://v8.1c.ru/8.1/data/enterprise/current-config",
		XmlnsCmi:   "http://v8.1c.ru/8.2/managed-application/cmi",
		XmlnsEnt:   "http://v8.1c.ru/8.1/data/enterprise",
		XmlnsLf:    "http://v8.1c.ru/8.2/managed-application/logform",
		XmlnsStyle: "http://v8.1c.ru/8.1/data/ui/style",
		XmlnsSys:   "http://v8.1c.ru/8.1/data/ui/fonts/system",
		XmlnsV8:    "http://v8.1c.ru/8.1/data/core",
		XmlnsV8ui:  "http://v8.1c.ru/8.1/data/ui",
		XmlnsWeb:   "http://v8.1c.ru/8.1/data/ui/colors/web",
		XmlnsWin:   "http://v8.1c.ru/8.1/data/ui/colors/windows",
		XmlnsXen:   "http://v8.1c.ru/8.3/xcf/enums",
		XmlnsXpr:   "http://v8.1c.ru/8.3/xcf/predef",
		XmlnsXr:    "http://v8.1c.ru/8.3/xcf/readable",
		XmlnsXs:    "http://www.w3.org/2001/XMLSchema",
		XmlnsXsi:   "http://www.w3.org/2001/XMLSchema-instance",

		CommonTemplate: &models.CommonTemplate{
			UUID: uuid.NewString(),
			Properties: models.Properties{
				Name: name,
				Synonym: models.Synonym{
					Items: []models.SynonymItem{
						{
							Lang:    "ru",
							Content: name,
						},
					},
				},
				Comment:      "",
				TemplateType: "TextDocument", // для наших нужд всегда такой
			},
		},
	}

	return tmpl, nil
}

func newCommonModule(root string, name string) (*models.MetaDataObject, error) {
	if err := registerMetaObject(root, "CommonModule", name); err != nil {
		return nil, err
	}

	m := &models.MetaDataObject{
		Version:    "2.16",
		Xmlns:      "http://v8.1c.ru/8.3/MDClasses",
		XmlnsApp:   "http://v8.1c.ru/8.2/managed-application/core",
		XmlnsCfg:   "http://v8.1c.ru/8.1/data/enterprise/current-config",
		XmlnsCmi:   "http://v8.1c.ru/8.2/managed-application/cmi",
		XmlnsEnt:   "http://v8.1c.ru/8.1/data/enterprise",
		XmlnsLf:    "http://v8.1c.ru/8.2/managed-application/logform",
		XmlnsStyle: "http://v8.1c.ru/8.1/data/ui/style",
		XmlnsSys:   "http://v8.1c.ru/8.1/data/ui/fonts/system",
		XmlnsV8:    "http://v8.1c.ru/8.1/data/core",
		XmlnsV8ui:  "http://v8.1c.ru/8.1/data/ui",
		XmlnsWeb:   "http://v8.1c.ru/8.1/data/ui/colors/web",
		XmlnsWin:   "http://v8.1c.ru/8.1/data/ui/colors/windows",
		XmlnsXen:   "http://v8.1c.ru/8.3/xcf/enums",
		XmlnsXpr:   "http://v8.1c.ru/8.3/xcf/predef",
		XmlnsXr:    "http://v8.1c.ru/8.3/xcf/readable",
		XmlnsXs:    "http://www.w3.org/2001/XMLSchema",
		XmlnsXsi:   "http://www.w3.org/2001/XMLSchema-instance",

		CommonModule: &models.CommonModule{
			UUID: uuid.NewString(),
			Properties: models.Properties{
				Name:              name,
				Server:            utils.Ptr(true),
				ReturnValuesReuse: "DontUse", // для наших нужд всегда такой
			},
		},
	}

	return m, nil
}

// registerMetaObject фиксация объекта в Configuration.xml
func registerMetaObject(root, path, name string) error {
	confPath := filepath.Join(root, "Configuration.xml")
	f, err := os.Open(confPath)
	if err != nil {
		return err
	}

	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(f); err != nil {
		return errors.Wrap(err, "xml read")
	}
	_ = f.Close()

	childObjects := doc.FindElement("//ChildObjects")
	if childObjects == nil {
		return errors.New("parse xml. ChildObjects not found")
	}

	httpService := childObjects.FindElements(path)
	_, exist := lo.Find(httpService, func(item *etree.Element) bool {
		return item.Text() == name
	})
	if exist {
		return nil
	}

	newNode := childObjects.CreateElement(path)
	newNode.SetText(name)

	doc.Indent(3)
	newXml, _ := doc.WriteToString()
	os.Remove(f.Name())

	f, _ = os.Create(confPath)
	f.WriteString(newXml)
	return f.Close()
}

func appendMethod(httpService *models.HTTPService, service *models.Service) {
	// группируем по url
	methodsByURL := lo.GroupBy(service.Methods, func(item *models.ServiceMethod) string {
		return item.Url
	})

	for url, methods := range methodsByURL {
		httpService.ChildObjects.URLTemplates = append(httpService.ChildObjects.URLTemplates, models.URLTemplate{
			UUID: uuid.NewString(),
			Properties: models.TemplateProperties{
				Name:     fmt.Sprintf("Шаблон%s-%d", httpService.Properties.Name, len(httpService.ChildObjects.URLTemplates)+1),
				Template: url,
			},
			ChildObjects: &models.TemplateChildObjects{Methods: make([]models.Method, 0, len(methods))},
		})

		last := httpService.ChildObjects.URLTemplates[len(httpService.ChildObjects.URLTemplates)-1]
		for _, method := range methods {
			last.ChildObjects.Methods = append(last.ChildObjects.Methods, models.Method{
				UUID: uuid.NewString(),
				Properties: models.MethodProperties{
					Name:       fmt.Sprintf("%s_%s", method.HttpMethod, method.Name),
					HTTPMethod: method.HttpMethod,
					Handler:    fmt.Sprintf("Обработчик_%s%s", service.Name, method.Name),
				},
			})
		}
	}
}
