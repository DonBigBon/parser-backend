package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/DonBigBon/parser-backend/internal/models"
)

type DocumentNode struct {
	Type      string
	ID        int
	NameRu    string
	NameKz    string
	ParentIDs map[string]int
	Children  []*DocumentNode
}

type Parser struct {
	rootNode *DocumentNode
	patterns map[string]*regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		rootNode: &DocumentNode{
			Type:     "ROOT",
			Children: make([]*DocumentNode, 0),
		},
		patterns: map[string]*regexp.Regexp{
			"PART":      regexp.MustCompile(`(?m)^(?:\s*)ЧАСТЬ\s+(\d+)[\.\s]+(.+?)(?:\n|$)`),
			"SECTION":   regexp.MustCompile(`(?m)^(?:\s*)РАЗДЕЛ\s+(\d+)[\.\s]+(.+?)(?:\n|$)`),
			"CHAPTER":   regexp.MustCompile(`(?m)^(?:\s*)Глава\s+(\d+)[\.\s]+(.+?)(?:\n|$)`),
			"PARAGRAPH": regexp.MustCompile(`(?m)^(?:\s*)Параграф\s+(\d+)[\.\s]+(.+?)(?:\n|$)`),
			"ARTICLE":   regexp.MustCompile(`(?m)^(?:\s*)Статья\s+(\d+)[\.\s]+(.+?)(?:\n|$)`),
			"CLAUSE":    regexp.MustCompile(`(?m)^(?:\s*)(\d+)\)\s+(.+?)(?:\n|$)`),
			"SUBCLAUSE": regexp.MustCompile(`(?m)^(?:\s*)([a-zа-яA-ZА-Я])\)\s+(.+?)(?:\n|$)`),
		},
	}
}

func (p *Parser) ParseDocument(content string) *DocumentNode {
	lines := strings.Split(content, "\n")

	context := make(map[string]*DocumentNode)
	context["ROOT"] = p.rootNode

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		p.processLineForType(line, "PART", context)
		p.processLineForType(line, "SECTION", context)
		p.processLineForType(line, "CHAPTER", context)
		p.processLineForType(line, "PARAGRAPH", context)
		p.processLineForType(line, "ARTICLE", context)
		p.processLineForType(line, "CLAUSE", context)
		p.processLineForType(line, "SUBCLAUSE", context)
	}

	return p.rootNode
}

func (p *Parser) processLineForType(line, nodeType string, context map[string]*DocumentNode) bool {
	if match := p.patterns[nodeType].FindStringSubmatch(line); match != nil {
		nodeID := parseIntID(match[1])
		nodeName := match[2]
		nameRu, nameKz := splitNames(nodeName)

		newNode := &DocumentNode{
			Type:      nodeType,
			ID:        nodeID,
			NameRu:    nameRu,
			NameKz:    nameKz,
			ParentIDs: make(map[string]int),
			Children:  make([]*DocumentNode, 0),
		}

		parentTypes := getParentTypes(nodeType)

		for _, parentType := range parentTypes {
			if parent, ok := context[parentType]; ok {
				newNode.ParentIDs[parentType] = parent.ID
				parent.Children = append(parent.Children, newNode)
			} else {
				newNode.ParentIDs[parentType] = 0
			}
		}

		if len(parentTypes) == 0 {
			context["ROOT"].Children = append(context["ROOT"].Children, newNode)
		}

		context[nodeType] = newNode

		for _, childType := range getChildTypes(nodeType) {
			delete(context, childType)
		}

		return true
	}

	return false
}

func (p *Parser) ConvertToFlatData() models.ParsedData {
	var data models.ParsedData

	p.traverseTree(p.rootNode, &data)

	return data
}

func (p *Parser) traverseTree(node *DocumentNode, data *models.ParsedData) {
	switch node.Type {
	case "PART":
		data.Parts = append(data.Parts, models.Part{
			ID:     node.ID,
			NameRu: node.NameRu,
			NameKz: node.NameKz,
		})
	case "SECTION":
		data.Sections = append(data.Sections, models.Section{
			ID:           node.ID,
			ParentPartID: node.ParentIDs["PART"],
			NameRu:       node.NameRu,
			NameKz:       node.NameKz,
		})
	case "CHAPTER":
		data.Chapters = append(data.Chapters, models.Chapter{
			ID:              node.ID,
			ParentSectionID: node.ParentIDs["SECTION"],
			ParentPartID:    node.ParentIDs["PART"],
			NameRu:          node.NameRu,
			NameKz:          node.NameKz,
		})
	case "PARAGRAPH":
		data.Paragraphs = append(data.Paragraphs, models.Paragraph{
			ID:              node.ID,
			ParentChapterID: node.ParentIDs["CHAPTER"],
			ParentSectionID: node.ParentIDs["SECTION"],
			ParentPartID:    node.ParentIDs["PART"],
			NameRu:          node.NameRu,
			NameKz:          node.NameKz,
		})
	case "ARTICLE":
		data.Articles = append(data.Articles, models.Article{
			ID:                node.ID,
			ParentParagraphID: node.ParentIDs["PARAGRAPH"],
			ParentChapterID:   node.ParentIDs["CHAPTER"],
			ParentSectionID:   node.ParentIDs["SECTION"],
			ParentPartID:      node.ParentIDs["PART"],
			NameRu:            node.NameRu,
			NameKz:            node.NameKz,
		})
	case "CLAUSE":
		data.Clauses = append(data.Clauses, models.Clause{
			ID:                node.ID,
			ParentArticleID:   node.ParentIDs["ARTICLE"],
			ParentParagraphID: node.ParentIDs["PARAGRAPH"],
			ParentChapterID:   node.ParentIDs["CHAPTER"],
			ParentSectionID:   node.ParentIDs["SECTION"],
			ParentPartID:      node.ParentIDs["PART"],
			NameRu:            node.NameRu,
			NameKz:            node.NameKz,
		})
	case "SUBCLAUSE":
		data.SubClauses = append(data.SubClauses, models.SubClause{
			ID:                node.ID,
			ParentClauseID:    node.ParentIDs["CLAUSE"],
			ParentArticleID:   node.ParentIDs["ARTICLE"],
			ParentParagraphID: node.ParentIDs["PARAGRAPH"],
			ParentChapterID:   node.ParentIDs["CHAPTER"],
			ParentSectionID:   node.ParentIDs["SECTION"],
			ParentPartID:      node.ParentIDs["PART"],
			NameRu:            node.NameRu,
			NameKz:            node.NameKz,
		})
	}

	for _, child := range node.Children {
		p.traverseTree(child, data)
	}
}

func parseIntID(idStr string) int {
	var id int
	fmt.Sscanf(idStr, "%d", &id)
	return id
}

func splitNames(fullName string) (string, string) {
	parts := strings.Split(fullName, "/")

	nameRu := strings.TrimSpace(parts[0])
	nameKz := ""

	if len(parts) > 1 {
		nameKz = strings.TrimSpace(parts[1])
	}

	return nameRu, nameKz
}

func getParentTypes(nodeType string) []string {
	switch nodeType {
	case "PART":
		return []string{}
	case "SECTION":
		return []string{"PART"}
	case "CHAPTER":
		return []string{"SECTION", "PART"}
	case "PARAGRAPH":
		return []string{"CHAPTER", "SECTION", "PART"}
	case "ARTICLE":
		return []string{"PARAGRAPH", "CHAPTER", "SECTION", "PART"}
	case "CLAUSE":
		return []string{"ARTICLE", "PARAGRAPH", "CHAPTER", "SECTION", "PART"}
	case "SUBCLAUSE":
		return []string{"CLAUSE", "ARTICLE", "PARAGRAPH", "CHAPTER", "SECTION", "PART"}
	default:
		return []string{}
	}
}

func getChildTypes(nodeType string) []string {
	switch nodeType {
	case "PART":
		return []string{"SECTION", "CHAPTER", "PARAGRAPH", "ARTICLE", "CLAUSE", "SUBCLAUSE"}
	case "SECTION":
		return []string{"CHAPTER", "PARAGRAPH", "ARTICLE", "CLAUSE", "SUBCLAUSE"}
	case "CHAPTER":
		return []string{"PARAGRAPH", "ARTICLE", "CLAUSE", "SUBCLAUSE"}
	case "PARAGRAPH":
		return []string{"ARTICLE", "CLAUSE", "SUBCLAUSE"}
	case "ARTICLE":
		return []string{"CLAUSE", "SUBCLAUSE"}
	case "CLAUSE":
		return []string{"SUBCLAUSE"}
	default:
		return []string{}
	}
}
