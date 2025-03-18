package models

type CodeData struct {
	Parts      []Part
	Sections   []Section
	Chapters   []Chapter
	Paragraphs []Paragraph
	Articles   []Article
	Clauses    []Clause
	SubClauses []SubClause
}

type Part struct {
	ID     int    `json:"id"`
	NameRu string `json:"nameRu"`
	NameKz string `json:"nameKz"`
}

type Section struct {
	ID           int    `json:"id"`
	ParentPartID int    `json:"parentPartId"`
	NameRu       string `json:"nameRu"`
	NameKz       string `json:"nameKz"`
}

type Chapter struct {
	ID              int    `json:"id"`
	ParentSectionID int    `json:"parentSectionId"`
	ParentPartID    int    `json:"parentPartId"`
	NameRu          string `json:"nameRu"`
	NameKz          string `json:"nameKz"`
}

type Paragraph struct {
	ID              int    `json:"id"`
	ParentChapterID int    `json:"parentChapterId"`
	ParentSectionID int    `json:"parentSectionId"`
	ParentPartID    int    `json:"parentPartId"`
	NameRu          string `json:"nameRu"`
	NameKz          string `json:"nameKz"`
}

type Article struct {
	ID                int    `json:"id"`
	ParentParagraphID int    `json:"parentParagraphId"`
	ParentChapterID   int    `json:"parentChapterId"`
	ParentSectionID   int    `json:"parentSectionId"`
	ParentPartID      int    `json:"parentPartId"`
	NameRu            string `json:"nameRu"`
	NameKz            string `json:"nameKz"`
}

type Clause struct {
	ID                int    `json:"id"`
	ParentArticleID   int    `json:"parentArticleId"`
	ParentParagraphID int    `json:"parentParagraphId"`
	ParentChapterID   int    `json:"parentChapterId"`
	ParentSectionID   int    `json:"parentSectionId"`
	ParentPartID      int    `json:"parentPartId"`
	NameRu            string `json:"nameRu"`
	NameKz            string `json:"nameKz"`
}

type SubClause struct {
	ID                int    `json:"id"`
	ParentClauseID    int    `json:"parentClauseId"`
	ParentArticleID   int    `json:"parentArticleId"`
	ParentParagraphID int    `json:"parentParagraphId"`
	ParentChapterID   int    `json:"parentChapterId"`
	ParentSectionID   int    `json:"parentSectionId"`
	ParentPartID      int    `json:"parentPartId"`
	NameRu            string `json:"nameRu"`
	NameKz            string `json:"nameKz"`
}

type ParsedData struct {
	Parts      []Part      `json:"parts"`
	Sections   []Section   `json:"sections"`
	Chapters   []Chapter   `json:"chapters"`
	Paragraphs []Paragraph `json:"paragraphs"`
	Articles   []Article   `json:"articles"`
	Clauses    []Clause    `json:"clauses"`
	SubClauses []SubClause `json:"subClauses"`
}

type DocumentResult struct {
	ParsedData ParsedData        `json:"parsedData"`
	SQLQueries []string          `json:"sqlQueries"`
	CSVFiles   map[string]string `json:"csvFiles"`
}
