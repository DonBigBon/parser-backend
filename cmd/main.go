package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

// Добавьте другие структуры (Paragraph, Article, Clause, SubClause)

type ParsedData struct {
	Parts    []Part    `json:"parts"`
	Sections []Section `json:"sections"`
	Chapters []Chapter `json:"chapters"`
	// Добавьте другие слайсы
}

// Конфигурация базы данных
var (
	dbUser     = "sa"
	dbPassword = "YourStrongPassword"
	dbHost     = "localhost"
	dbPort     = 1433
	dbName     = "ParserDB"
)

func main() {
	// Настройка подключения к базе данных
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		dbHost, dbUser, dbPassword, dbPort, dbName)

	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Error connecting to database:", err.Error())
	}
	defer db.Close()

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err.Error())
	}

	// Настройка Gin роутера
	router := gin.Default()

	// Настройка CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Маршруты API
	router.POST("/api/parse", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("document")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
			return
		}
		defer file.Close()

		// Сохранение загруженного файла
		tempFilePath := filepath.Join("uploads", header.Filename)
		os.MkdirAll("uploads", os.ModePerm)

		out, err := os.Create(tempFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file content"})
			return
		}

		// Чтение содержимого файла
		content, err := ioutil.ReadFile(tempFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
			return
		}

		// Парсинг документа
		parsedData := parseDocument(string(content))

		// Создание SQL запросов
		sqlQueries := generateSQLQueries(parsedData)

		// Создание CSV файлов
		csvFiles := generateCSVFiles(parsedData)

		// Ответ клиенту
		c.JSON(http.StatusOK, gin.H{
			"sqlQueries": sqlQueries,
			"csvFiles":   csvFiles,
		})
	})

	// Запуск сервера
	log.Println("Starting server on port 8080...")
	router.Run(":8080")
}

// Функция для парсинга документа
func parseDocument(content string) ParsedData {
	var parsedData ParsedData

	// Регулярные выражения для поиска частей, разделов и т.д.
	partRegex := regexp.MustCompile(`ЧАСТЬ\s+(\d+)\.\s+(.*?)(?:\n|$)`)
	sectionRegex := regexp.MustCompile(`РАЗДЕЛ\s+(\d+)\.\s+(.*?)(?:\n|$)`)
	chapterRegex := regexp.MustCompile(`Глава\s+(\d+)\.\s+(.*?)(?:\n|$)`)
	// Добавьте регулярные выражения для других элементов

	// Поиск частей
	partMatches := partRegex.FindAllStringSubmatch(content, -1)
	for _, match := range partMatches {
		if len(match) >= 3 {
			partID := parseID(match[1])
			partName := strings.TrimSpace(match[2])
			// Здесь необходимо разделить на русское и казахское название
			// Для простоты примера, предположим, что название на казахском указано в скобках
			nameRu, nameKz := splitNames(partName)

			parsedData.Parts = append(parsedData.Parts, Part{
				ID:     partID,
				NameRu: nameRu,
				NameKz: nameKz,
			})
		}
	}

	// Поиск разделов и других элементов...
	// Аналогично с частями, но с учетом родительских элементов

	return parsedData
}

// Вспомогательная функция для разделения названий на русском и казахском
func splitNames(fullName string) (string, string) {
	// Предполагаем, что название на казахском указано в скобках
	parts := strings.Split(fullName, "(")
	nameRu := strings.TrimSpace(parts[0])
	nameKz := ""

	if len(parts) > 1 {
		nameKz = strings.TrimSpace(strings.TrimSuffix(parts[1], ")"))
	}

	return nameRu, nameKz
}

// Функция для преобразования строки в число
func parseID(idStr string) int {
	var id int
	fmt.Sscanf(idStr, "%d", &id)
	return id
}

// Функция для генерации SQL запросов
func generateSQLQueries(data ParsedData) []string {
	var queries []string

	// Запросы для вставки частей
	for _, part := range data.Parts {
		query := fmt.Sprintf("INSERT INTO Parts (PartId, NameRu, NameKz) VALUES (%d, N'%s', N'%s');",
			part.ID, escapeSQLString(part.NameRu), escapeSQLString(part.NameKz))
		queries = append(queries, query)
	}

	// Запросы для вставки разделов
	for _, section := range data.Sections {
		query := fmt.Sprintf("INSERT INTO Sections (SectionId, ParentPartId, NameRu, NameKz) VALUES (%d, %d, N'%s', N'%s');",
			section.ID, section.ParentPartID, escapeSQLString(section.NameRu), escapeSQLString(section.NameKz))
		queries = append(queries, query)
	}

	// Аналогично для других элементов...

	return queries
}

// Функция для экранирования строк в SQL запросах
func escapeSQLString(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// Функция для генерации CSV файлов
func generateCSVFiles(data ParsedData) map[string]string {
	csvFiles := make(map[string]string)

	// Создание CSV для частей
	partsCSV := createCSV([]string{"PartId", "NameRu", "NameKz"})
	for _, part := range data.Parts {
		partsCSV = append(partsCSV, []string{
			fmt.Sprintf("%d", part.ID),
			part.NameRu,
			part.NameKz,
		})
	}

	csvFiles["parts"] = writeCSVToString(partsCSV)

	// Аналогично для других элементов...

	return csvFiles
}

// Функция для создания заголовка CSV файла
func createCSV(headers []string) [][]string {
	csv := make([][]string, 0)
	csv = append(csv, headers)
	return csv
}

// Функция для записи CSV в строку
func writeCSVToString(records [][]string) string {
	var csvString strings.Builder
	writer := csv.NewWriter(&csvString)

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			log.Println("Error writing CSV record:", err)
			return ""
		}
	}

	writer.Flush()
	return csvString.String()
}
