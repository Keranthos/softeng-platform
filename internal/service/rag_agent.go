package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"softeng-platform/internal/config"
	"softeng-platform/internal/repository"
	"strings"
	"time"
)

// RagSource 检索片段（报告：RAG 引用溯源）
type RagSource struct {
	Kind    string `json:"kind"`
	ID      string `json:"id"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

// RagChatResult Agent 输出
type RagChatResult struct {
	Answer      string      `json:"answer"`
	Sources     []RagSource `json:"sources"`
	Mode        string      `json:"mode"` // gemini | llm | demo
	Model       string      `json:"model,omitempty"`
	UserMessage string      `json:"userMessage"`
}

type RagAgentService interface {
	Chat(ctx context.Context, userMessage string) (*RagChatResult, error)
}

type ragAgentService struct {
	db  *repository.Database
	cfg *config.Config
	hc  *http.Client
}

func NewRagAgentService(db *repository.Database, cfg *config.Config) RagAgentService {
	return &ragAgentService{
		db:  db,
		cfg: cfg,
		hc:  &http.Client{Timeout: 90 * time.Second},
	}
}

// redactClientErr 避免把密钥、代理口令写进返回给前端的文案（例如 http 错误里会带完整 URL 含 ?key=）
func (s *ragAgentService) redactClientErr(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	if k := strings.TrimSpace(s.cfg.GoogleAPIKey); k != "" {
		msg = strings.ReplaceAll(msg, k, "(已脱敏)")
	}
	if k := strings.TrimSpace(s.cfg.OpenAIAPIKey); k != "" {
		msg = strings.ReplaceAll(msg, k, "(已脱敏)")
	}
	if k := strings.TrimSpace(s.cfg.GeminiProxyToken); k != "" {
		msg = strings.ReplaceAll(msg, k, "(已脱敏)")
	}
	return msg
}

func (s *ragAgentService) Chat(ctx context.Context, userMessage string) (*RagChatResult, error) {
	q := strings.TrimSpace(userMessage)
	if q == "" {
		return nil, fmt.Errorf("empty message")
	}
	sources, err := s.retrieve(ctx, q)
	if err != nil {
		sources = nil
	}
	ctxBlock := s.buildContextBlock(sources)

	// 1) 与 blog/agent 一致：优先 Gemini（GOOGLE_API_KEY / GEMINI_API_KEY）
	if strings.TrimSpace(s.cfg.GoogleAPIKey) != "" {
		ans, model, err := s.callGemini(ctx, q, ctxBlock)
		if err == nil {
			return &RagChatResult{
				Answer:      ans,
				Sources:     sources,
				Mode:        "gemini",
				Model:       model,
				UserMessage: q,
			}, nil
		}
		if strings.TrimSpace(s.cfg.OpenAIAPIKey) != "" {
			if ans2, model2, err2 := s.callOpenAI(ctx, q, ctxBlock); err2 == nil {
				return &RagChatResult{
					Answer:      ans2,
					Sources:     sources,
					Mode:        "llm",
					Model:       model2,
					UserMessage: q,
				}, nil
			}
		}
		return &RagChatResult{
			Answer:      s.answerDemo(q, sources, ctxBlock) + "\n\n（Gemini 调用失败，已回退：" + s.redactClientErr(err) + "）",
			Sources:     sources,
			Mode:        "demo",
			UserMessage: q,
		}, nil
	}

	// 2) OpenAI 兼容网关
	if strings.TrimSpace(s.cfg.OpenAIAPIKey) != "" {
		ans, model, err := s.callOpenAI(ctx, q, ctxBlock)
		if err != nil {
			return &RagChatResult{
				Answer:      s.answerDemo(q, sources, ctxBlock) + "\n\n（LLM 调用失败，已回退为离线综合：" + s.redactClientErr(err) + "）",
				Sources:     sources,
				Mode:        "demo",
				UserMessage: q,
			}, nil
		}
		return &RagChatResult{
			Answer:      ans,
			Sources:     sources,
			Mode:        "llm",
			Model:       model,
			UserMessage: q,
		}, nil
	}

	// 3) 离线演示
	return &RagChatResult{
		Answer:      s.answerDemo(q, sources, ctxBlock),
		Sources:     sources,
		Mode:        "demo",
		UserMessage: q,
	}, nil
}

func (s *ragAgentService) retrieve(ctx context.Context, q string) ([]RagSource, error) {
	pat := "%" + q + "%"
	var out []RagSource

	rowsT, err := s.db.QueryContext(ctx, `
		SELECT resource_id, resource_name, COALESCE(description,''), COALESCE(description_detail,'')
		FROM tools WHERE status = 'approved'
		  AND (resource_name LIKE ? OR description LIKE ? OR description_detail LIKE ?)
		ORDER BY views DESC LIMIT 6`, pat, pat, pat)
	if err == nil {
		defer rowsT.Close()
		for rowsT.Next() {
			var id int
			var name, d1, d2 string
			if rowsT.Scan(&id, &name, &d1, &d2) != nil {
				continue
			}
			sn := strings.TrimSpace(d1)
			if len(sn) > 220 {
				sn = sn[:220] + "…"
			}
			out = append(out, RagSource{Kind: "tool", ID: fmt.Sprintf("%d", id), Title: name, Snippet: sn})
		}
	}

	rowsC, err := s.db.QueryContext(ctx, `
		SELECT course_id, name, semester FROM courses
		WHERE name LIKE ? OR semester LIKE ?
		ORDER BY views DESC LIMIT 4`, pat, pat)
	if err == nil {
		defer rowsC.Close()
		for rowsC.Next() {
			var id int
			var name, sem string
			if rowsC.Scan(&id, &name, &sem) != nil {
				continue
			}
			out = append(out, RagSource{Kind: "course", ID: fmt.Sprintf("%d", id), Title: name, Snippet: "学期：" + sem})
		}
	}

	rowsP, err := s.db.QueryContext(ctx, `
		SELECT project_id, name, COALESCE(description,''), COALESCE(LEFT(detail, 400),'')
		FROM projects
		WHERE name LIKE ? OR description LIKE ? OR detail LIKE ?
		ORDER BY views DESC LIMIT 6`, pat, pat, pat)
	if err == nil {
		defer rowsP.Close()
		for rowsP.Next() {
			var id int
			var name, d1, d2 string
			if rowsP.Scan(&id, &name, &d1, &d2) != nil {
				continue
			}
			sn := strings.TrimSpace(d1)
			if sn == "" {
				sn = strings.TrimSpace(d2)
			}
			if len(sn) > 220 {
				sn = sn[:220] + "…"
			}
			out = append(out, RagSource{Kind: "project", ID: fmt.Sprintf("%d", id), Title: name, Snippet: sn})
		}
	}

	if len(out) > 14 {
		out = out[:14]
	}
	return out, nil
}

func (s *ragAgentService) buildContextBlock(src []RagSource) string {
	if len(src) == 0 {
		return "（知识库中未检索到与问题强相关的条目，请基于软件工程通识回答。）"
	}
	var b strings.Builder
	for i, x := range src {
		b.WriteString(fmt.Sprintf("[%d] (%s #%s) %s\n摘要：%s\n\n", i+1, x.Kind, x.ID, x.Title, x.Snippet))
	}
	return b.String()
}

func (s *ragAgentService) answerDemo(q string, src []RagSource, ctxBlock string) string {
	if len(src) == 0 {
		return fmt.Sprintf("【离线 RAG 演示】未在平台库中命中「%s」相关条目。建议在工具/课程/项目模块补充关键词更匹配的简介后再试；或在 .env 中配置 GOOGLE_API_KEY（Gemini）或 OPENAI_API_KEY。", q)
	}
	var titles []string
	for _, x := range src {
		if len(titles) < 5 {
			titles = append(titles, fmt.Sprintf("%s《%s》", x.Kind, x.Title))
		}
	}
	return fmt.Sprintf("【离线 RAG 演示】根据检索到的 %d 条知识片段，与您的问题「%s」相关的资源主要包括：%s。\n\n"+
		"以下为检索上下文摘要（供报告说明「引用溯源」）：\n%s\n"+
		"（配置 GOOGLE_API_KEY + GEMINI_CHAT_MODEL 将走与 blog/agent 相同的 Gemini generateContent；或配置 OPENAI_API_KEY 走 OpenAI 兼容接口。）",
		len(src), q, strings.Join(titles, "、"), ctxBlock)
}

type oaMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type oaReq struct {
	Model    string  `json:"model"`
	Messages []oaMsg `json:"messages"`
}

type oaResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func isGoogleGenerativeLanguageHost(base string) bool {
	b := strings.TrimSuffix(strings.TrimSpace(strings.ToLower(base)), "/")
	return strings.HasSuffix(b, "generativelanguage.googleapis.com")
}

type geminiTextPart struct {
	Text string `json:"text"`
}

type geminiUserContent struct {
	Role  string           `json:"role"`
	Parts []geminiTextPart `json:"parts"`
}

type geminiGenerateBody struct {
	SystemInstruction struct {
		Parts []geminiTextPart `json:"parts"`
	} `json:"systemInstruction"`
	Contents         []geminiUserContent `json:"contents"`
	GenerationConfig struct {
		Temperature float64 `json:"temperature"`
	} `json:"generationConfig"`
}

type geminiPartOut struct {
	Text string `json:"text"`
}

type geminiContentOut struct {
	Parts []geminiPartOut `json:"parts"`
}

type geminiCandidate struct {
	Content geminiContentOut `json:"content"`
}

type geminiGenerateResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
	Error        *struct {
		Message string `json:"message"`
	} `json:"error"`
}

// callGemini 调用与 blog/agent 相同的 Generative Language API（v1beta generateContent）
func (s *ragAgentService) callGemini(ctx context.Context, userQ, contextBlock string) (string, string, error) {
	base := strings.TrimRight(strings.TrimSpace(s.cfg.GeminiAPIBase), "/")
	model := strings.TrimSpace(s.cfg.GeminiChatModel)
	if model == "" {
		model = "gemini-2.5-flash"
	}
	apiURL := base + "/v1beta/models/" + model + ":generateContent"
	u, err := url.Parse(apiURL)
	if err != nil {
		return "", "", fmt.Errorf("gemini url: %w", err)
	}

	useGoogleKey := isGoogleGenerativeLanguageHost(base)
	if useGoogleKey {
		q := u.Query()
		q.Set("key", strings.TrimSpace(s.cfg.GoogleAPIKey))
		u.RawQuery = q.Encode()
	}

	system := "你是软件工程学习资源平台的助教助手。请严格根据「知识库检索片段」回答问题；若片段不足以回答，请明确说明并给出学习建议。回答尽量使用中文，条理清晰。"
	userText := "用户问题：\n" + userQ + "\n\n知识库检索片段：\n" + contextBlock

	var body geminiGenerateBody
	body.SystemInstruction.Parts = []geminiTextPart{{Text: system}}
	body.Contents = []geminiUserContent{{Role: "user", Parts: []geminiTextPart{{Text: userText}}}}
	body.GenerationConfig.Temperature = s.cfg.GeminiTemperature

	rawBody, err := json.Marshal(body)
	if err != nil {
		return "", "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(rawBody))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if !useGoogleKey {
		tok := strings.TrimSpace(s.cfg.GeminiProxyToken)
		if tok == "" {
			return "", "", fmt.Errorf("GEMINI_API_BASE 非 Google 官方时需在 .env 设置 GEMINI_PROXY_TOKEN（与 blog/agent 一致）")
		}
		req.Header.Set("X-Proxy-Token", tok)
	}

	res, err := s.hc.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", "", fmt.Errorf("gemini http %d: %s", res.StatusCode, string(raw))
	}
	var parsed geminiGenerateResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", "", err
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return "", "", fmt.Errorf("gemini api: %s", parsed.Error.Message)
	}
	if len(parsed.Candidates) == 0 {
		return "", "", fmt.Errorf("gemini empty candidates: %s", string(raw))
	}
	var texts []string
	for _, p := range parsed.Candidates[0].Content.Parts {
		texts = append(texts, p.Text)
	}
	out := strings.TrimSpace(strings.Join(texts, ""))
	if out == "" {
		return "", "", fmt.Errorf("gemini empty text body")
	}
	return out, model, nil
}

func (s *ragAgentService) callOpenAI(ctx context.Context, userQ, contextBlock string) (string, string, error) {
	base := strings.TrimRight(strings.TrimSpace(s.cfg.OpenAIBaseURL), "/")
	url := base + "/chat/completions"
	model := strings.TrimSpace(s.cfg.OpenAIModel)
	if model == "" {
		model = "gpt-4o-mini"
	}
	system := "你是软件工程学习资源平台的助教助手。请严格根据「知识库检索片段」回答问题；若片段不足以回答，请明确说明并给出学习建议。回答尽量使用中文，条理清晰。"
	user := "用户问题：\n" + userQ + "\n\n知识库检索片段：\n" + contextBlock

	body, _ := json.Marshal(oaReq{
		Model: model,
		Messages: []oaMsg{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.cfg.OpenAIAPIKey)

	res, err := s.hc.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", "", fmt.Errorf("openai http %d: %s", res.StatusCode, string(raw))
	}
	var parsed oaResp
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", "", err
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return "", "", fmt.Errorf("openai error: %s", parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		return "", "", fmt.Errorf("empty choices")
	}
	return strings.TrimSpace(parsed.Choices[0].Message.Content), model, nil
}
