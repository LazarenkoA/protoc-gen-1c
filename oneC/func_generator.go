package oneC

import (
	"embed"
	"github.com/LazarenkoA/protoc-gen-1c/pkg/utils"
	"github.com/pkg/errors"
	"strings"
	"text/template"
)

type HandlerInfo struct {
	Method                    string
	HandlerName               string
	ServiceName               string
	RequiredBodyParams        []string
	RequiredQueryParamsParams []string
	BodyParams                []string // параметры которые передаются через тело запроса в жсоне
	QueryParams               []string // параметры строки запроса /v1/customers?page=10&page_size=20
	PathParams                []string // параметры пути запроса /v1/customers/{id}
	TmplName                  string
	RespCodes                 map[int32]string
	Funcs                     template.FuncMap
}

//go:embed templates/*
var templatesFS embed.FS

func getCommonModuleContent(info *HandlerInfo, log *utils.Logger) string {
	content, err := executeTemplate(info, "template_common_module")
	if err != nil {
		log.Error(err.Error(), "template_name", "template_common_module")
	}

	return content
}

func getHttpHandler(info *HandlerInfo, log *utils.Logger) string {
	content, err := executeTemplate(info, "template_method_handler")
	if err != nil {
		log.Error(err.Error(), "template_name", "template_method_handler")
	}

	return content
}

func swaggerHandler(info *HandlerInfo, log *utils.Logger) string {
	content, err := executeTemplate(info, "swagger_service_handler")
	if err != nil {
		log.Error(err.Error(), "template_name", "swagger_service_handler")
	}

	return content
}

func checkRequestFields(info *HandlerInfo, log *utils.Logger) string {
	content, err := executeTemplate(info, "template_request_check")
	if err != nil {
		log.Error(err.Error(), "template_name", "template_request_check")
	}

	return content
}

func checkResponseFields(info *HandlerInfo, log *utils.Logger) string {
	content, err := executeTemplate(info, "template_response_check")
	if err != nil {
		log.Error(err.Error(), "template_name", "template_response_check")
	}

	return content
}

func executeTemplate(info *HandlerInfo, name string) (string, error) {
	t, err := template.New("").Funcs(info.Funcs).ParseFS(templatesFS, "templates/"+name)
	if err != nil {
		return "", errors.Wrap(err, "parse template error")
	}

	var sb strings.Builder
	err = t.ExecuteTemplate(&sb, name, info)
	if err != nil {
		return "", errors.Wrap(err, "execute template error")
	}

	return sb.String(), nil
}

// boilerplate вспомогательные функции которые нужны каждому http сервису
func boilerplate() string {
	return `Функция JSON_ВЗначение(ТекстJSON, Ошибка = Ложь, ИменаСвойствСоЗначениямиДата = Неопределено) Экспорт
	Результат = Неопределено;
	
	ЧтениеJSON = Новый ЧтениеJSON;
	ЧтениеJSON.УстановитьСтроку(ТекстJSON);
	
	Результат = ПрочитатьJSON(ЧтениеJSON, Ложь, ИменаСвойствСоЗначениямиДата);	
	ЧтениеJSON.Закрыть();
	
	Возврат Результат;
КонецФункции           

Функция СформироватьОтветСервисаREST(КодСостояния = 200, СтрокаДляОтвета = "", Ошибка = Ложь, ДвоичныеДанные = Неопределено) Экспорт
	Заголовки = Новый Соответствие; 
	Заголовки.Вставить("Accept-Charset", "utf-8");
	
	Ответ = Новый HTTPСервисОтвет(КодСостояния, , Заголовки);
	
	Если Ошибка Тогда
		Заголовки.Вставить("Content-Type", "text/html; charset=utf-8");
		Ответ.УстановитьТелоИзСтроки(СтрокаДляОтвета, КодировкаТекста.UTF8, ИспользованиеByteOrderMark.НеИспользовать);
	 	Возврат Ответ;
	КонецЕсли;	
	
	Если ДвоичныеДанные = Неопределено Тогда
		Заголовки.Вставить("Content-Type", "application/json");
		Ответ.УстановитьТелоИзСтроки(СтрокаДляОтвета, КодировкаТекста.UTF8, ИспользованиеByteOrderMark.НеИспользовать);
	Иначе	
		Заголовки.Вставить("Content-Type", "application/octet-stream");
		Ответ.УстановитьТелоИзДвоичныхДанных(ДвоичныеДанные);
	КонецЕсли;
	
	Возврат Ответ;	

КонецФункции`
}

func disclaimer(comment string) string {
	return comment + " НЕ РЕДАКТИРОВАТЬ! Код сгенерированный программой protoc-gen-1c."
}
