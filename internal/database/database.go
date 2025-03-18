package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/DonBigBon/parser-backend/internal/models"
	_ "github.com/denisenkom/go-mssqldb"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type DBHandler struct {
	db *sql.DB
}

func NewDBHandler(config DBConfig) (*DBHandler, error) {
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		config.Host, config.User, config.Password, config.Port, config.DBName)

	db, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return &DBHandler{db: db}, nil
}

func (h *DBHandler) Close() error {
	return h.db.Close()
}

func (h *DBHandler) GenerateSQLQueries(data models.ParsedData) []string {
	var queries []string

	for _, part := range data.Parts {
		query := fmt.Sprintf("INSERT INTO Parts (PartId, NameRu, NameKz) VALUES (%d, N'%s', N'%s');",
			part.ID, escapeSQLString(part.NameRu), escapeSQLString(part.NameKz))
		queries = append(queries, query)
	}

	for _, section := range data.Sections {
		query := fmt.Sprintf("INSERT INTO Sections (SectionId, ParentPartId, NameRu, NameKz) VALUES (%d, %d, N'%s', N'%s');",
			section.ID, section.ParentPartID, escapeSQLString(section.NameRu), escapeSQLString(section.NameKz))
		queries = append(queries, query)
	}

	for _, chapter := range data.Chapters {
		query := fmt.Sprintf("INSERT INTO Chapters (ChapterId, ParentSectionId, ParentPartId, NameRu, NameKz) VALUES (%d, %d, %d, N'%s', N'%s');",
			chapter.ID, chapter.ParentSectionID, chapter.ParentPartID, escapeSQLString(chapter.NameRu), escapeSQLString(chapter.NameKz))
		queries = append(queries, query)
	}

	for _, paragraph := range data.Paragraphs {
		query := fmt.Sprintf("INSERT INTO Paragraphs (ParagraphId, ParentChapterId, ParentSectionId, ParentPartId, NameRu, NameKz) VALUES (%d, %d, %d, %d, N'%s', N'%s');",
			paragraph.ID, paragraph.ParentChapterID, paragraph.ParentSectionID, paragraph.ParentPartID,
			escapeSQLString(paragraph.NameRu), escapeSQLString(paragraph.NameKz))
		queries = append(queries, query)
	}

	for _, article := range data.Articles {
		query := fmt.Sprintf("INSERT INTO Articles (ArticleId, ParentParagraphId, ParentChapterId, ParentSectionId, ParentPartId, NameRu, NameKz) VALUES (%d, %d, %d, %d, %d, N'%s', N'%s');",
			article.ID, article.ParentParagraphID, article.ParentChapterID, article.ParentSectionID, article.ParentPartID,
			escapeSQLString(article.NameRu), escapeSQLString(article.NameKz))
		queries = append(queries, query)
	}

	for _, clause := range data.Clauses {
		query := fmt.Sprintf("INSERT INTO Clauses (ClauseId, ParentArticleId, ParentParagraphId, ParentChapterId, ParentSectionId, ParentPartId, NameRu, NameKz) VALUES (%d, %d, %d, %d, %d, %d, N'%s', N'%s');",
			clause.ID, clause.ParentArticleID, clause.ParentParagraphID, clause.ParentChapterID, clause.ParentSectionID, clause.ParentPartID,
			escapeSQLString(clause.NameRu), escapeSQLString(clause.NameKz))
		queries = append(queries, query)
	}

	for _, subClause := range data.SubClauses {
		query := fmt.Sprintf("INSERT INTO SubClauses (SubClauseId, ParentClauseId, ParentArticleId, ParentParagraphId, ParentChapterId, ParentSectionId, ParentPartId, NameRu, NameKz) VALUES (%d, %d, %d, %d, %d, %d, %d, N'%s', N'%s');",
			subClause.ID, subClause.ParentClauseID, subClause.ParentArticleID, subClause.ParentParagraphID, subClause.ParentChapterID, subClause.ParentSectionID, subClause.ParentPartID,
			escapeSQLString(subClause.NameRu), escapeSQLString(subClause.NameKz))
		queries = append(queries, query)
	}

	return queries
}

func (h *DBHandler) ExecuteQueries(queries []string) error {
	for _, query := range queries {
		_, err := h.db.Exec(query)
		if err != nil {
			return fmt.Errorf("error executing query: %v", err)
		}
	}
	return nil
}

func escapeSQLString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
