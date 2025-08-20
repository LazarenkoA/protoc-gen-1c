package models

import "encoding/xml"

// MetaDataObject represents the root XML element
type MetaDataObject struct {
	XMLName        xml.Name        `xml:"MetaDataObject"`
	HTTPService    *HTTPService    `xml:"HTTPService,omitempty"`
	CommonTemplate *CommonTemplate `xml:"CommonTemplate,omitempty"`
	CommonModule   *CommonModule   `xml:"CommonModule,omitempty"`
	Version        string          `xml:"version,attr"`
	Xmlns          string          `xml:"xmlns,attr"`
	XmlnsApp       string          `xml:"xmlns:app,attr"`
	XmlnsCfg       string          `xml:"xmlns:cfg,attr"`
	XmlnsCmi       string          `xml:"xmlns:cmi,attr"`
	XmlnsEnt       string          `xml:"xmlns:ent,attr"`
	XmlnsLf        string          `xml:"xmlns:lf,attr"`
	XmlnsStyle     string          `xml:"xmlns:style,attr"`
	XmlnsSys       string          `xml:"xmlns:sys,attr"`
	XmlnsV8        string          `xml:"xmlns:v8,attr"`
	XmlnsV8ui      string          `xml:"xmlns:v8ui,attr"`
	XmlnsWeb       string          `xml:"xmlns:web,attr"`
	XmlnsWin       string          `xml:"xmlns:win,attr"`
	XmlnsXen       string          `xml:"xmlns:xen,attr"`
	XmlnsXpr       string          `xml:"xmlns:xpr,attr"`
	XmlnsXr        string          `xml:"xmlns:xr,attr"`
	XmlnsXs        string          `xml:"xmlns:xs,attr"`
	XmlnsXsi       string          `xml:"xmlns:xsi,attr"`
}

type CommonModule struct {
	UUID       string     `xml:"uuid,attr"`
	Properties Properties `xml:"Properties"`
}

type CommonTemplate struct {
	UUID       string     `xml:"uuid,attr"`
	Properties Properties `xml:"Properties"`
	Comment    string     `xml:"Comment,omitempty"`
}

type HTTPService struct {
	UUID         string       `xml:"uuid,attr"`
	Properties   Properties   `xml:"Properties"`
	ChildObjects ChildObjects `xml:"ChildObjects"`
}

type Properties struct {
	Name    string  `xml:"Name"`
	Synonym Synonym `xml:"Synonym"`
	Comment string  `xml:"Comment,omitempty"`

	// поля для HTTPService
	RootURL       string `xml:"RootURL,omitempty"`
	ReuseSessions string `xml:"ReuseSessions,omitempty"`
	SessionMaxAge int    `xml:"SessionMaxAge,omitempty"`

	// поля для CommonTemplate
	TemplateType string `xml:"TemplateType,omitempty"`

	// поля для CommonModule
	Global                    *bool  `xml:"Global,omitempty"`
	ClientManagedApplication  *bool  `xml:"ClientManagedApplication,omitempty"`
	Server                    *bool  `xml:"Server,omitempty"`
	ExternalConnection        *bool  `xml:"ExternalConnection,omitempty"`
	ClientOrdinaryApplication *bool  `xml:"ClientOrdinaryApplication,omitempty"`
	ServerCall                *bool  `xml:"ServerCall,omitempty"`
	Privileged                *bool  `xml:"Privileged,omitempty"`
	ReturnValuesReuse         string `xml:"ReturnValuesReuse,omitempty"`
}

type Synonym struct {
	Items []SynonymItem `xml:"v8:item"`
}

type SynonymItem struct {
	Lang    string `xml:"v8:lang,omitempty"`
	Content string `xml:"v8:content,omitempty"`
}

type ChildObjects struct {
	URLTemplates []URLTemplate `xml:"URLTemplate"`
}

type URLTemplate struct {
	UUID         string                `xml:"uuid,attr"`
	Properties   TemplateProperties    `xml:"Properties"`
	ChildObjects *TemplateChildObjects `xml:"ChildObjects"`
}

type TemplateProperties struct {
	Name     string `xml:"Name"`
	Synonym  string `xml:"Synonym,omitempty"`
	Comment  string `xml:"Comment,omitempty"`
	Template string `xml:"Template,omitempty"`
}

type TemplateChildObjects struct {
	Methods []Method `xml:"Method"`
}

type Method struct {
	UUID       string           `xml:"uuid,attr"`
	Properties MethodProperties `xml:"Properties"`
}

type MethodProperties struct {
	Name       string `xml:"Name"`
	Synonym    string `xml:"Synonym,omitempty"`
	Comment    string `xml:"Comment,omitempty"`
	HTTPMethod string `xml:"HTTPMethod,omitempty"`
	Handler    string `xml:"Handler"`
}
