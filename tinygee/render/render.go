package render

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/xrjjing/Learn4Go/tinygee"
)

// TemplateRenderer 支持 FuncMap 与模板加载。
type TemplateRenderer struct {
	templates *template.Template
}

// New 创建模板渲染器。
func New(glob string, funcMap template.FuncMap) (*TemplateRenderer, error) {
	t := template.New(filepath.Base(glob)).Funcs(funcMap)
	parsed, err := t.ParseGlob(glob)
	if err != nil {
		return nil, err
	}
	return &TemplateRenderer{templates: parsed}, nil
}

// HTML 渲染模板。
func (tr *TemplateRenderer) HTML(c *tinygee.Context, code int, name string, data any) {
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	c.Status(code)
	_ = tr.templates.ExecuteTemplate(c.Writer, name, data)
}

// Static 返回一个处理静态文件的 HandlerFunc。
func Static(relative string, root http.FileSystem) tinygee.HandlerFunc {
	fileServer := http.StripPrefix(relative, http.FileServer(root))
	return func(c *tinygee.Context) {
		// 简单安全防护：仅允许 GET
		if c.Method != http.MethodGet {
			c.Status(http.StatusMethodNotAllowed)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
