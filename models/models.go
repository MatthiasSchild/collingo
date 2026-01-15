package models

type InfoModel struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type AccessModel struct {
	User string `json:"user"`
	Role string `json:"role"`
}

type ProjectModel struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Owner                string   `json:"owner"`
	BaseLanguage         string   `json:"baseLanguage"`
	TranslationLanguages []string `json:"translationLanguages"`
}

type GroupModel struct {
	ID            string       `json:"id"`
	Project       string       `json:"project"`
	Parent        *string      `json:"parent"`
	DisplayName   string       `json:"displayName"`
	TechnicalName string       `json:"technicalName"`
	Entries       []EntryModel `json:"entries"`
}

type GroupSummaryModel struct {
	ID              string   `json:"id"`
	DisplayName     string   `json:"displayName"`
	TechnicalName   string   `json:"technicalName"`
	Breadcrumbs     []string `json:"breadcrumbs"`
	BreadcrumbNames []string `json:"breadcrumbNames"`
}

type EntryModel struct {
	TechnicalName  string             `json:"technicalName"`
	BaseTerm       string             `json:"baseTerm"`
	ContextInfo    string             `json:"contextInfo"`
	TranslatedTerm []TranslationModel `json:"translatedTerm"`
}

type TranslationModel struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
	Type        string `json:"type"`
	Translator  string `json:"translator"`
}

type ResultModel[T any] struct {
	Success bool `json:"success"`
	Result  T    `json:"result"`
}

type ManyResultModel[T any] struct {
	Success bool  `json:"success"`
	Result  []T   `json:"result"`
	Total   int32 `json:"total"`
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
	HasMore bool  `json:"hasMore"`
}
