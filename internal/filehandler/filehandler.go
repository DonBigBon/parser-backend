package filehandler

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/DonBigBon/parser-backend/internal/models"
	"github.com/xuri/excelize/v2"
)

var supportedFormats = map[string]bool{
	".docx": true,
	".doc":  true,
	".txt":  true,
	".rtf":  true,
}

func SaveUploadedFile(file io.Reader, filename string) (string, error) {
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			return "", err
		}
	}

	ext := filepath.Ext(filename)
	if !supportedFormats[strings.ToLower(ext)] {
		return "", errors.New("unsupported file format")
	}

	filePath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func GenerateCSV(codeData *models.CodeData) (map[string]string, error) {
	csvFiles := make(map[string]string)

	csvDir := "./csv_output"
	if _, err := os.Stat(csvDir); os.IsNotExist(err) {
		err = os.MkdirAll(csvDir, 0755)
		if err != nil {
			return nil, err
		}
	}

	f := excelize.NewFile()

	sheetName := "Parts"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheetName, "A1", "PartNumber")
	f.SetCellValue(sheetName, "B1", "NameRu")
	f.SetCellValue(sheetName, "C1", "NameKz")

	for i, part := range codeData.Parts {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), part.Number)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), part.NameRu)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), part.NameKz)
	}

	sheetName = "Sections"
	index, err = f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheetName, "A1", "PartNumber")
	f.SetCellValue(sheetName, "B1", "SectionNumber")
	f.SetCellValue(sheetName, "C1", "NameRu")
	f.SetCellValue(sheetName, "D1", "NameKz")

	for i, section := range codeData.Sections {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), section.PartNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), section.Number)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), section.NameRu)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), section.NameKz)
	}

	sheetName = "Chapters"
	index, err = f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheetName, "A1", "PartNumber")
	f.SetCellValue(sheetName, "B1", "SectionNumber")
	f.SetCellValue(sheetName, "C1", "ChapterNumber")
	f.SetCellValue(sheetName, "D1", "NameRu")
	f.SetCellValue(sheetName, "E1", "NameKz")

	for i, chapter := range codeData.Chapters {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), chapter.PartNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), chapter.SectionNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), chapter.Number)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), chapter.NameRu)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), chapter.NameKz)
	}

	sheetName = "Paragraphs"
	index, err = f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheetName, "A1", "PartNumber")
	f.SetCellValue(sheetName, "B1", "SectionNumber")
	f.SetCellValue(sheetName, "C1", "ChapterNumber")
	f.SetCellValue(sheetName, "D1", "ParagraphNumber")
	f.SetCellValue(sheetName, "E1", "NameRu")
	f.SetCellValue(sheetName, "F1", "NameKz")

	for i, paragraph := range codeData.Paragraphs {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), paragraph.PartNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), paragraph.SectionNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), paragraph.ChapterNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), paragraph.Number)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), paragraph.NameRu)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), paragraph.NameKz)
	}

	sheetName = "Articles"
	index, err = f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheetName, "A1", "PartNumber")
	f.SetCellValue(sheetName, "B1", "SectionNumber")
	f.SetCellValue(sheetName, "C1", "ChapterNumber")
	f.SetCellValue(sheetName, "D1", "ParagraphNumber")
	f.SetCellValue(sheetName, "E1", "ArticleNumber")
	f.SetCellValue(sheetName, "F1", "NameRu")
	f.SetCellValue(sheetName, "G1", "NameKz")

	for i, article := range codeData.Articles {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), article.PartNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), article.SectionNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), article.ChapterNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), article.ParagraphNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), article.Number)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), article.NameRu)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), article.NameKz)
	}

	sheetName = "Clauses"
	index, err = f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheetName, "A1", "PartNumber")
	f.SetCellValue(sheetName, "B1", "SectionNumber")
	f.SetCellValue(sheetName, "C1", "ChapterNumber")
	f.SetCellValue(sheetName, "D1", "ParagraphNumber")
	f.SetCellValue(sheetName, "E1", "ArticleNumber")
	f.SetCellValue(sheetName, "F1", "ClauseNumber")
	f.SetCellValue(sheetName, "G1", "NameRu")
	f.SetCellValue(sheetName, "H1", "NameKz")

	for i, clause := range codeData.Clauses {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), clause.PartNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), clause.SectionNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), clause.ChapterNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), clause.ParagraphNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), clause.ArticleNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), clause.Number)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), clause.NameRu)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), clause.NameKz)
	}

	sheetName = "SubClauses"
	index, err = f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}
	f.SetCellValue(sheetName, "A1", "PartNumber")
	f.SetCellValue(sheetName, "B1", "SectionNumber")
	f.SetCellValue(sheetName, "C1", "ChapterNumber")
	f.SetCellValue(sheetName, "D1", "ParagraphNumber")
	f.SetCellValue(sheetName, "E1", "ArticleNumber")
	f.SetCellValue(sheetName, "F1", "ClauseNumber")
	f.SetCellValue(sheetName, "G1", "SubClauseNumber")
	f.SetCellValue(sheetName, "H1", "NameRu")
	f.SetCellValue(sheetName, "I1", "NameKz")

	for i, subClause := range codeData.SubClauses {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), subClause.PartNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), subClause.SectionNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), subClause.ChapterNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), subClause.ParagraphNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), subClause.ArticleNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), subClause.ClauseNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), subClause.Number)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), subClause.NameRu)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), subClause.NameKz)
	}

	f.SetActiveSheet(index)

	csvPath := filepath.Join(csvDir, "code_data.xlsx")
	if err := f.SaveAs(csvPath); err != nil {
		return nil, err
	}

	csvFiles["excel"] = csvPath

	return csvFiles, nil
}

func GenerateSQLDump(codeData *models.CodeData) (string, error) {
	sqlDir := "./sql_output"
	if _, err := os.Stat(sqlDir); os.IsNotExist(err) {
		err = os.MkdirAll(sqlDir, 0755)
		if err != nil {
			return "", err
		}
	}

	sqlPath := filepath.Join(sqlDir, "code_data.sql")
	file, err := os.Create(sqlPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	sql := `
-- Очистка таблиц
DELETE FROM SubClauses;
DELETE FROM Clauses;
DELETE FROM Articles;
DELETE FROM Paragraphs;
DELETE FROM Chapters;
DELETE FROM Sections;
DELETE FROM Parts;
DELETE FROM Codes;

-- Создаем новый кодекс
INSERT INTO Codes (Name) VALUES ('Новый кодекс');
DECLARE @CodeID INT = SCOPE_IDENTITY();

-- Вставка частей
`

	partIDMap := make(map[string]string)
	for i, part := range codeData.Parts {
		partID := fmt.Sprintf("@PartID_%s", part.Number)
		partIDMap[part.Number] = partID

		sql += fmt.Sprintf("INSERT INTO Parts (CodeID, Number, NameRu, NameKz) VALUES (@CodeID, '%s', '%s', '%s');\n",
			part.Number, escapeSQLString(part.NameRu), escapeSQLString(part.NameKz))
		sql += fmt.Sprintf("DECLARE %s INT = SCOPE_IDENTITY();\n\n", partID)
	}

	sectionIDMap := make(map[string]string)
	for i, section := range codeData.Sections {
		sectionKey := fmt.Sprintf("%s_%s", section.PartNumber, section.Number)
		sectionID := fmt.Sprintf("@SectionID_%s", sectionKey)
		sectionIDMap[sectionKey] = sectionID

		partID := partIDMap[section.PartNumber]
		if partID == "" {
			partID = "NULL"
		}

		sql += fmt.Sprintf("INSERT INTO Sections (PartID, Number, NameRu, NameKz) VALUES (%s, '%s', '%s', '%s');\n",
			partID, section.Number, escapeSQLString(section.NameRu), escapeSQLString(section.NameKz))
		sql += fmt.Sprintf("DECLARE %s INT = SCOPE_IDENTITY();\n\n", sectionID)
	}

	chapterIDMap := make(map[string]string)
	for i, chapter := range codeData.Chapters {
		chapterKey := fmt.Sprintf("%s_%s_%s", chapter.PartNumber, chapter.SectionNumber, chapter.Number)
		chapterID := fmt.Sprintf("@ChapterID_%s", chapterKey)
		chapterIDMap[chapterKey] = chapterID

		sectionKey := fmt.Sprintf("%s_%s", chapter.PartNumber, chapter.SectionNumber)
		sectionID := sectionIDMap[sectionKey]
		if sectionID == "" {
			sectionID = "NULL"
		}

		sql += fmt.Sprintf("INSERT INTO Chapters (SectionID, Number, NameRu, NameKz) VALUES (%s, '%s', '%s', '%s');\n",
			sectionID, chapter.Number, escapeSQLString(chapter.NameRu), escapeSQLString(chapter.NameKz))
		sql += fmt.Sprintf("DECLARE %s INT = SCOPE_IDENTITY();\n\n", chapterID)
	}

	paragraphIDMap := make(map[string]string)
	for i, paragraph := range codeData.Paragraphs {
		paragraphKey := fmt.Sprintf("%s_%s_%s_%s", paragraph.PartNumber, paragraph.SectionNumber, paragraph.ChapterNumber, paragraph.Number)
		paragraphID := fmt.Sprintf("@ParagraphID_%s", paragraphKey)
		paragraphIDMap[paragraphKey] = paragraphID

		chapterKey := fmt.Sprintf("%s_%s_%s", paragraph.PartNumber, paragraph.SectionNumber, paragraph.ChapterNumber)
		chapterID := chapterIDMap[chapterKey]
		if chapterID == "" {
			chapterID = "NULL"
		}

		sql += fmt.Sprintf("INSERT INTO Paragraphs (ChapterID, Number, NameRu, NameKz) VALUES (%s, '%s', '%s', '%s');\n",
			chapterID, paragraph.Number, escapeSQLString(paragraph.NameRu), escapeSQLString(paragraph.NameKz))
		sql += fmt.Sprintf("DECLARE %s INT = SCOPE_IDENTITY();\n\n", paragraphID)
	}

	articleIDMap := make(map[string]string)
	for i, article := range codeData.Articles {
		articleKey := fmt.Sprintf("%s_%s_%s_%s_%s", article.PartNumber, article.SectionNumber, article.ChapterNumber, article.ParagraphNumber, article.Number)
		articleID := fmt.Sprintf("@ArticleID_%s", articleKey)
		articleIDMap[articleKey] = articleID

		paragraphKey := fmt.Sprintf("%s_%s_%s_%s", article.PartNumber, article.SectionNumber, article.ChapterNumber, article.ParagraphNumber)
		paragraphID := paragraphIDMap[paragraphKey]
		if paragraphID == "" {
			paragraphID = "NULL"
		}

		sql += fmt.Sprintf("INSERT INTO Articles (ParagraphID, Number, NameRu, NameKz) VALUES (%s, '%s', '%s', '%s');\n",
			paragraphID, article.Number, escapeSQLString(article.NameRu), escapeSQLString(article.NameKz))
		sql += fmt.Sprintf("DECLARE %s INT = SCOPE_IDENTITY();\n\n", articleID)
	}

	clauseIDMap := make(map[string]string)
	for i, clause := range codeData.Clauses {
		clauseKey := fmt.Sprintf("%s_%s_%s_%s_%s_%s", clause.PartNumber, clause.SectionNumber, clause.ChapterNumber, clause.ParagraphNumber, clause.ArticleNumber, clause.Number)
		clauseID := fmt.Sprintf("@ClauseID_%s", clauseKey)
		clauseIDMap[clauseKey] = clauseID

		articleKey := fmt.Sprintf("%s_%s_%s_%s_%s", clause.PartNumber, clause.SectionNumber, clause.ChapterNumber, clause.ParagraphNumber, clause.ArticleNumber)
		articleID := articleIDMap[articleKey]
		if articleID == "" {
			articleID = "NULL"
		}

		sql += fmt.Sprintf("INSERT INTO Clauses (ArticleID, Number, NameRu, NameKz) VALUES (%s, '%s', '%s', '%s');\n",
			articleID, clause.Number, escapeSQLString(clause.NameRu), escapeSQLString(clause.NameKz))
		sql += fmt.Sprintf("DECLARE %s INT = SCOPE_IDENTITY();\n\n", clauseID)
	}

	for i, subClause := range codeData.SubClauses {
		clauseKey := fmt.Sprintf("%s_%s_%s_%s_%s_%s", subClause.PartNumber, subClause.SectionNumber, subClause.ChapterNumber, subClause.ParagraphNumber, subClause.ArticleNumber, subClause.ClauseNumber)
		clauseID := clauseIDMap[clauseKey]
		if clauseID == "" {
			clauseID = "NULL"
		}

		sql += fmt.Sprintf("INSERT INTO SubClauses (ClauseID, Number, NameRu, NameKz) VALUES (%s, '%s', '%s', '%s');\n",
			clauseID, subClause.Number, escapeSQLString(subClause.NameRu), escapeSQLString(subClause.NameKz))
	}

	_, err = file.WriteString(sql)
	if err != nil {
		return "", err
	}

	return sqlPath, nil
}

func escapeSQLString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
