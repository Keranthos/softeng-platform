package safety

import (
	"strings"
)

// 与前端 contentSafety.js 保持一致的演示词库（课程报告：双端校验）
var defaultBlockedWords = []string{
	"法轮功", "赌博", "色情", "代开发票", "枪支", "炸药",
}

var blockedURLPrefixes = []string{
	"javascript:",
	"data:text/html",
	"vbscript:",
}

func textHits(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	ls := strings.ToLower(s)
	var out []string
	for _, w := range defaultBlockedWords {
		if w == "" {
			continue
		}
		if strings.Contains(s, w) || strings.Contains(ls, strings.ToLower(w)) {
			out = append(out, "包含敏感词演示命中：「"+w+"」")
		}
	}
	return out
}

func urlIssues(u string) []string {
	if strings.TrimSpace(u) == "" {
		return nil
	}
	low := strings.ToLower(strings.TrimSpace(u))
	var out []string
	for _, p := range blockedURLPrefixes {
		if strings.HasPrefix(low, p) {
			out = append(out, "链接协议不允许："+p)
		}
	}
	return out
}

func joinReasons(parts [][]string) string {
	var flat []string
	for _, p := range parts {
		flat = append(flat, p...)
	}
	return strings.Join(flat, "；")
}

// ValidateToolSubmit 工具提交 / 更新正文
func ValidateToolSubmit(name, link, description, detail string) string {
	return joinReasons([][]string{
		textHits(name),
		textHits(description),
		textHits(detail),
		urlIssues(link),
	})
}

// ValidateProjectBody 项目上传 / 更新
func ValidateProjectBody(name, description, detail, github string) string {
	return joinReasons([][]string{
		textHits(name),
		textHits(description),
		textHits(detail),
		urlIssues(github),
	})
}

// ValidateCourseResource 课程资料上传
func ValidateCourseResource(resourceURL, description string) string {
	return joinReasons([][]string{
		textHits(description),
		urlIssues(resourceURL),
	})
}

// ValidateComment 评论正文
func ValidateComment(content string) string {
	return joinReasons([][]string{textHits(content)})
}
