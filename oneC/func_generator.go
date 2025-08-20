package oneC

import (
	"1c-grpc-gateway/pkg/utils"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
	"text/template"
)

type HandlerInfo struct {
	Method         string
	HandlerName    string
	ServiceName    string
	RequiredFields []string
	TmplName       string
	Funcs          template.FuncMap
}

func getCommonModuleContent(info *HandlerInfo, log *utils.Logger) string {
	content, err := executeTemplate(info, "oneC/templates/template_common_module")
	if err != nil {
		log.Error(err.Error(), "template_name", "template_common_module")
	}

	return content
}

func getHttpHandler(info *HandlerInfo, log *utils.Logger) string {
	content, err := executeTemplate(info, "oneC/templates/template_method_handler")
	if err != nil {
		log.Error(err.Error(), "template_name", "template_method_handler")
	}

	return content
}

func swaggerHandler(info *HandlerInfo, log *utils.Logger) string {
	content, err := executeTemplate(info, "oneC/templates/swagger_service_handler")
	if err != nil {
		log.Error(err.Error(), "template_name", "swagger_service_handler")
	}

	return content
}

func executeTemplate(info *HandlerInfo, tmplPath string) (string, error) {
	tmplFile, err := os.Open(tmplPath)
	if err != nil {
		return "", err
	}

	tmpl, _ := io.ReadAll(tmplFile)
	t, err := template.New("message").Funcs(info.Funcs).Parse(string(tmpl))
	if err != nil {
		return "", errors.Wrap(err, "parse template error")
	}

	var sb strings.Builder
	err = t.Execute(&sb, info)
	if err != nil {
		return "", errors.Wrap(err, "execute template error")
	}

	return sb.String(), nil
}

func checkFields(info *HandlerInfo, log *utils.Logger) string {
	content, err := executeTemplate(info, "oneC/templates/template_required_check")
	if err != nil {
		log.Error(err.Error(), "template_name", "template_required_check")
	}
	return content
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
