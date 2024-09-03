package md2tgmd

import (
    "fmt"
    "regexp"
    "strings"
)

// Converter 结构体用于存储转换器的状态
type Converter struct {
    // 可能需要添加一些字段来存储状态
}

// NewConverter 创建一个新的Converter实例
func NewConverter() *Converter {
    return &Converter{}
}

// Convert 方法用于将Markdown转换为Telegram Markdown
func (c *Converter) Convert(markdown string) string {
    return c.post_process(c.pre_process(markdown))
}

func (c *Converter) pre_process(text string) string {
    // 移除文本中的零宽度空格
    text = strings.ReplaceAll(text, "\u200b", "")

    // 替换 HTML 实体
    text = strings.ReplaceAll(text, "<", "<")
    text = strings.ReplaceAll(text, ">", ">")
    text = strings.ReplaceAll(text, "&amp;", "&")

    // 处理代码块
    text = c.fenced_code_blocks(text)
    text = c.inline_code(text)

    // 处理其他元素
    text = c.heading(text)
    text = c.image(text)
    text = c.link(text)
    text = c.list(text)
    text = c.quote(text)
    text = c.emphasis(text)

    return text
}

func (c *Converter) post_process(text string) string {
    // 移除多余的换行符
    text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

    // 移除行末空格
    text = regexp.MustCompile(`[ \t]+\n`).ReplaceAllString(text, "\n")

    return strings.TrimSpace(text)
}

func (c *Converter) escape_chars(text string) string {
    escapeChars := []string{"_", "*", "`", "["}
    for _, char := range escapeChars {
        text = strings.ReplaceAll(text, char, "\\"+char)
    }
    return text
}

func (c *Converter) unescape_chars(text string) string {
    unescapeChars := []string{"_", "*", "`", "["}
    for _, char := range unescapeChars {
        text = strings.ReplaceAll(text, "\\"+char, char)
    }
    return text
}

func (c *Converter) fenced_code_blocks(text string) string {
    re := regexp.MustCompile("(?ms)```(?:([a-zA-Z0-9]+)\n)?(.+?)```")
    return re.ReplaceAllStringFunc(text, func(match string) string {
        parts := re.FindStringSubmatch(match)
        lang := parts[1]
        code := strings.TrimSpace(parts[2])
        if lang != "" {
            return fmt.Sprintf("```%s\n%s\n```", lang, code)
        }
        return fmt.Sprintf("```\n%s\n```", code)
    })
}

func (c *Converter) inline_code(text string) string {
    re := regexp.MustCompile("`([^`\n]+)`")
    return re.ReplaceAllString(text, "`$1`")
}

func (c *Converter) heading(text string) string {
    re := regexp.MustCompile(`(?m)^(#{1,3})\s*(.+)$`)
    return re.ReplaceAllStringFunc(text, func(match string) string {
        parts := re.FindStringSubmatch(match)
        level := len(parts[1])
        content := strings.TrimSpace(parts[2])
        switch level {
        case 1:
            return fmt.Sprintf("<b>%s</b>\n", content)
        case 2:
            return fmt.Sprintf("<u>%s</u>\n", content)
        case 3:
            return fmt.Sprintf("<i>%s</i>\n", content)
        default:
            return match
        }
    })
}

func (c *Converter) image(text string) string {
    re := regexp.MustCompile(`!\[([^\]]*)\]\(([^\s\)]+)(?:\s["']([^"']*)["'])?\)`)
    return re.ReplaceAllString(text, "")
}

func (c *Converter) link(text string) string {
    re := regexp.MustCompile(`\[([^\]]+)\]\(([^\s\)]+)(?:\s["']([^"']*)["'])?\)`)
    return re.ReplaceAllStringFunc(text, func(match string) string {
        parts := re.FindStringSubmatch(match)
        text := parts[1]
        url := parts[2]
        return fmt.Sprintf(`<a href="%s">%s</a>`, url, text)
    })
}

func (c *Converter) list(text string) string {
    lines := strings.Split(text, "\n")
    var result []string
    inList := false
    for _, line := range lines {
        if strings.HasPrefix(strings.TrimSpace(line), "- ") || strings.HasPrefix(strings.TrimSpace(line), "* ") {
            if !inList {
                inList = true
            }
            line = strings.TrimSpace(line)[2:]
            result = append(result, "• "+line)
        } else {
            if inList {
                inList = false
                result = append(result, "")
            }
            result = append(result, line)
        }
    }
    return strings.Join(result, "\n")
}

func (c *Converter) quote(text string) string {
    lines := strings.Split(text, "\n")
    var result []string
    for _, line := range lines {
        if strings.HasPrefix(strings.TrimSpace(line), "> ") {
            line = strings.TrimSpace(line)[2:]
            result = append(result, line)
        } else {
            result = append(result, line)
        }
    }
    return strings.Join(result, "\n")
}

func (c *Converter) emphasis(text string) string {
    // 粗体
    re := regexp.MustCompile(`\*\*(.+?)\*\*`)
    text = re.ReplaceAllString(text, "<b>$1</b>")

    // 斜体
    re = regexp.MustCompile(`\*(.+?)\*`)
    text = re.ReplaceAllString(text, "<i>$1</i>")

    return text
}
