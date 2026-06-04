package services

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"

	"paw-api/models"
	"paw-api/repositories"
)

type DocsService struct {
	collectionRepo *repositories.CollectionRepository
	requestRepo    *repositories.RequestRepository
}

func NewDocsService(
	collectionRepo *repositories.CollectionRepository,
	requestRepo *repositories.RequestRepository,
) *DocsService {
	return &DocsService{
		collectionRepo: collectionRepo,
		requestRepo:    requestRepo,
	}
}

func (s *DocsService) GenerateMarkdown(projectID int64) (string, error) {
	collections, err := s.collectionRepo.ListByProject(projectID)
	if err != nil {
		return "", err
	}

	collMap := make(map[int64]*models.Collection)
	childrenMap := make(map[int64][]*models.Collection)
	var roots []*models.Collection

	for i := range collections {
		c := &collections[i]
		collMap[c.ID] = c
		if c.ParentID == nil {
			roots = append(roots, c)
		} else {
			childrenMap[*c.ParentID] = append(childrenMap[*c.ParentID], c)
		}
	}

	sort.Slice(roots, func(i, j int) bool { return roots[i].SortOrder < roots[j].SortOrder })
	for _, children := range childrenMap {
		sort.Slice(children, func(i, j int) bool { return children[i].SortOrder < children[j].SortOrder })
	}

	_ = collMap

	var sb strings.Builder
	sb.WriteString("# API Documentation\n\n")

	for _, root := range roots {
		s.writeCollectionMD(&sb, root, childrenMap, 1)
	}

	return sb.String(), nil
}

func (s *DocsService) writeCollectionMD(sb *strings.Builder, coll *models.Collection, childrenMap map[int64][]*models.Collection, level int) {
	prefix := strings.Repeat("#", level+1)
	sb.WriteString(fmt.Sprintf("%s %s\n\n", prefix, coll.Name))

	requests, err := s.requestRepo.ListByCollection(coll.ID)
	if err == nil {
		sort.Slice(requests, func(i, j int) bool { return requests[i].SortOrder < requests[j].SortOrder })
		for _, req := range requests {
			s.writeRequestMD(sb, &req, level+2)
		}
	}

	children := childrenMap[coll.ID]
	sort.Slice(children, func(i, j int) bool { return children[i].SortOrder < children[j].SortOrder })
	for _, child := range children {
		s.writeCollectionMD(sb, child, childrenMap, level+1)
	}
}

func (s *DocsService) writeRequestMD(sb *strings.Builder, req *models.Request, level int) {
	prefix := strings.Repeat("#", level)
	sb.WriteString(fmt.Sprintf("%s `%s %s`\n\n", prefix, req.Method, req.URL))

	if req.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", req.Description))
	}

	var headers []struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.Unmarshal([]byte(req.Headers), &headers); err == nil && len(headers) > 0 {
		hasEnabled := false
		for _, h := range headers {
			if h.Enabled {
				hasEnabled = true
				break
			}
		}
		if hasEnabled {
			sb.WriteString("**Headers:**\n\n")
			sb.WriteString("| Key | Value |\n")
			sb.WriteString("| --- | --- |\n")
			for _, h := range headers {
				if h.Enabled {
					sb.WriteString(fmt.Sprintf("| %s | %s |\n", h.Key, h.Value))
				}
			}
			sb.WriteString("\n")
		}
	}

	var params []struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Enabled bool   `json:"enabled"`
	}
	if err := json.Unmarshal([]byte(req.Params), &params); err == nil && len(params) > 0 {
		hasEnabled := false
		for _, p := range params {
			if p.Enabled {
				hasEnabled = true
				break
			}
		}
		if hasEnabled {
			sb.WriteString("**Query Parameters:**\n\n")
			sb.WriteString("| Key | Value |\n")
			sb.WriteString("| --- | --- |\n")
			for _, p := range params {
				if p.Enabled {
					sb.WriteString(fmt.Sprintf("| %s | %s |\n", p.Key, p.Value))
				}
			}
			sb.WriteString("\n")
		}
	}

	if req.BodyType != "" && req.BodyType != "none" {
		sb.WriteString(fmt.Sprintf("**Body Type:** `%s`\n\n", req.BodyType))
		if req.Body != "" && req.Body != "{}" && req.Body != "[]" {
			var rawBody struct {
				Content string `json:"content"`
			}
			if json.Unmarshal([]byte(req.Body), &rawBody) == nil && rawBody.Content != "" {
				sb.WriteString("```\n")
				sb.WriteString(rawBody.Content)
				sb.WriteString("\n```\n\n")
			} else {
				sb.WriteString("```json\n")
				sb.WriteString(req.Body)
				sb.WriteString("\n```\n\n")
			}
		}
	}
}

func (s *DocsService) GenerateRequestMarkdown(requestID int64) (string, error) {
	req, err := s.requestRepo.GetByID(requestID)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	// Title with request name
	if req.Name != "" {
		sb.WriteString(fmt.Sprintf("# %s\n\n", req.Name))
	}

	s.writeRequestMD(&sb, req, 2)
	return sb.String(), nil
}

func (s *DocsService) GenerateRequestHTML(requestID int64) (string, error) {
	md, err := s.GenerateRequestMarkdown(requestID)
	if err != nil {
		return "", err
	}

	body := markdownToHTML(md)
	htmlOut := renderHTMLPage(body)
	return htmlOut, nil
}

func (s *DocsService) GenerateHTML(projectID int64) (string, error) {
	md, err := s.GenerateMarkdown(projectID)
	if err != nil {
		return "", err
	}

	body := markdownToHTML(md)
	htmlOut := renderHTMLPage(body)
	return htmlOut, nil
}

const htmlPageTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>API Documentation</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 900px; margin: 0 auto; padding: 20px; background: #fff; color: #333; line-height: 1.6; }
h1 { border-bottom: 2px solid #18a058; padding-bottom: 8px; }
h2 { border-bottom: 1px solid #ddd; padding-bottom: 6px; margin-top: 32px; }
h3 { margin-top: 24px; }
h4, h5, h6 { margin-top: 20px; }
code { background: #f5f5f5; padding: 2px 6px; border-radius: 3px; font-size: 90%%; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace; }
pre { background: #f5f5f5; padding: 16px; border-radius: 6px; overflow-x: auto; }
pre code { background: none; padding: 0; font-size: 85%%; }
table { border-collapse: collapse; width: 100%%; margin: 12px 0; }
th, td { border: 1px solid #ddd; padding: 8px 12px; text-align: left; }
th { background: #f2f2f2; font-weight: 600; }
p { margin: 8px 0; }
</style>
</head>
<body>
%s
</body>
</html>`

func renderHTMLPage(body string) string {
	return fmt.Sprintf(htmlPageTemplate, body)
}

func markdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var sb strings.Builder
	inCodeBlock := false
	var codeBlock strings.Builder
	codeBlockLang := ""
	inTable := false
	var tableRows []string

	flushTable := func() {
		if !inTable || len(tableRows) == 0 {
			return
		}
		sb.WriteString("<table>\n")
		for i, row := range tableRows {
			cells := strings.Split(strings.Trim(row, "|"), "|")
			if i == 0 {
				sb.WriteString("<thead><tr>")
				for _, cell := range cells {
					sb.WriteString("<th>" + inlineMarkdown(strings.TrimSpace(cell)) + "</th>")
				}
				sb.WriteString("</tr></thead>\n<tbody>\n")
			} else if i == 1 {
				// separator row, skip
				continue
			} else {
				sb.WriteString("<tr>")
				for _, cell := range cells {
					sb.WriteString("<td>" + inlineMarkdown(strings.TrimSpace(cell)) + "</td>")
				}
				sb.WriteString("</tr>\n")
			}
		}
		sb.WriteString("</tbody>\n</table>\n")
		tableRows = nil
		inTable = false
	}

	for _, line := range lines {
		// Code block toggle
		if strings.HasPrefix(line, "```") {
			if !inCodeBlock {
				flushTable()
				inCodeBlock = true
				codeBlockLang = strings.TrimPrefix(line, "```")
				codeBlock.Reset()
			} else {
				inCodeBlock = false
				langAttr := ""
				if codeBlockLang != "" {
					langAttr = fmt.Sprintf(` class="language-%s"`, html.EscapeString(codeBlockLang))
				}
				sb.WriteString(fmt.Sprintf("<pre><code%s>%s</code></pre>\n", langAttr, html.EscapeString(codeBlock.String())))
			}
			continue
		}
		if inCodeBlock {
			if codeBlock.Len() > 0 {
				codeBlock.WriteString("\n")
			}
			codeBlock.WriteString(line)
			continue
		}

		// Table row
		if strings.HasPrefix(strings.TrimSpace(line), "|") && strings.Contains(line, "|") {
			flushTablePending := false
			_ = flushTablePending
			if !inTable {
				inTable = true
			}
			tableRows = append(tableRows, line)
			continue
		}

		// Flush table if we're no longer in one
		flushTable()

		// Empty line
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Headings
		if matched, _ := regexp.MatchString(`^#{1,6}\s`, line); matched {
			level := 0
			for _, ch := range line {
				if ch == '#' {
					level++
				} else {
					break
				}
			}
			text := strings.TrimSpace(line[level:])
			sb.WriteString(fmt.Sprintf("<h%d>%s</h%d>\n", level, inlineMarkdown(text), level))
			continue
		}

		// Paragraph
		sb.WriteString("<p>" + inlineMarkdown(line) + "</p>\n")
	}

	// Flush any remaining
	if inCodeBlock {
		sb.WriteString(fmt.Sprintf("<pre><code>%s</code></pre>\n", html.EscapeString(codeBlock.String())))
	}
	flushTable()

	return sb.String()
}

var (
	reBold       = regexp.MustCompile(`\*\*(.+?)\*\*`)
	reInlineCode = regexp.MustCompile("`([^`]+)`")
)

func inlineMarkdown(s string) string {
	s = html.EscapeString(s)
	s = reInlineCode.ReplaceAllString(s, "<code>$1</code>")
	s = reBold.ReplaceAllString(s, "<strong>$1</strong>")
	return s
}
