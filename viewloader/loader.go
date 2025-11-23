package viewloader

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
)

// LoadTemplates 根据配置加载所有模板文件并返回 multitemplate.Renderer
// 这是主要的入口函数,用于初始化 Gin 的 HTML 渲染器
func LoadTemplates(config *Config) (multitemplate.Renderer, error) {
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 创建渲染器
	render := multitemplate.NewRenderer()

	// 创建基础模板函数映射
	funcMap := template.FuncMap{
		// 内置函数: 判断菜单是否激活
		"isActive": func(current, target string) bool {
			return current == target
		},
		// 内置函数: 返回激活状态的 class 值
		"activeClass": func(current, target string) string {
			if current == target {
				return "active"
			}
			return ""
		},
		// 内置函数: 返回完整的 class 属性
		"activeAttr": func(current, target string) template.HTMLAttr {
			if current == target {
				return template.HTMLAttr(`class="active"`)
			}
			return ""
		},
	}

	// 合并用户自定义函数
	if config.Funcs != nil {
		for name, fn := range config.Funcs {
			funcMap[name] = fn
		}
	}

	// 加载模板组(需要继承布局的模板)
	for _, group := range config.TemplateGroups {
		if err := loadTemplateGroup(render, group, funcMap); err != nil {
			return nil, err
		}
	}

	// 加载独立模板(不需要继承布局的模板)
	if err := loadStandaloneTemplates(render, config.StandalonePaths, funcMap); err != nil {
		return nil, err
	}

	return render, nil
}

// loadTemplateGroup 加载一个模板组
func loadTemplateGroup(render multitemplate.Renderer, group TemplateGroup, funcMap template.FuncMap) error {
	// 检查布局模板是否存在
	if !fileExists(group.LayoutPath) {
		return &LoadError{
			Path:    group.LayoutPath,
			Message: fmt.Sprintf("模板组 [%s] 的布局模板文件不存在", group.Name),
			Err:     ErrLayoutNotFound,
		}
	}

	// 遍历视图目录
	for _, viewDir := range group.ViewDirs {
		// 检查路径是文件还是目录
		info, err := os.Stat(viewDir)
		if err != nil {
			return &LoadError{
				Path:    viewDir,
				Message: fmt.Sprintf("模板组 [%s] 的视图路径无法访问", group.Name),
				Err:     err,
			}
		}

		// 如果是目录,递归扫描
		if info.IsDir() {
			if err := loadViewDir(render, group.Name, group.LayoutPath, viewDir, funcMap); err != nil {
				return err
			}
		} else {
			// 如果是文件,直接加载
			if err := loadViewFile(render, group.LayoutPath, viewDir, funcMap); err != nil {
				return err
			}
		}
	}

	return nil
}

// loadViewDir 递归扫描目录并加载所有 HTML 视图文件
func loadViewDir(render multitemplate.Renderer, groupName, layoutPath, viewDir string, funcMap template.FuncMap) error {
	return filepath.WalkDir(viewDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return &LoadError{
				Path:    path,
				Message: fmt.Sprintf("模板组 [%s] 的视图目录遍历失败", groupName),
				Err:     err,
			}
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		// 只处理 HTML 文件
		if !strings.HasSuffix(strings.ToLower(path), ".html") {
			return nil
		}

		// 加载视图文件
		return loadViewFile(render, layoutPath, path, funcMap)
	})
}

// loadViewFile 加载单个视图文件(继承布局)
func loadViewFile(render multitemplate.Renderer, layoutPath, viewPath string, funcMap template.FuncMap) error {
	// 使用相对路径作为模板名称,并统一使用正斜杠
	templateName := normalizeTemplateName(viewPath)

	// 解析模板: 布局 + 视图
	tmpl, err := template.New(filepath.Base(layoutPath)).Funcs(funcMap).ParseFiles(layoutPath, viewPath)
	if err != nil {
		return &LoadError{
			Path:    viewPath,
			Message: "模板解析失败",
			Err:     err,
		}
	}

	// 注册到渲染器
	render.Add(templateName, tmpl)

	return nil
}

// loadStandaloneTemplates 加载独立模板(不继承布局)
func loadStandaloneTemplates(render multitemplate.Renderer, paths []string, funcMap template.FuncMap) error {
	for _, path := range paths {
		// 检查路径是否存在
		info, err := os.Stat(path)
		if err != nil {
			return &LoadError{
				Path:    path,
				Message: "独立视图路径无法访问",
				Err:     err,
			}
		}

		// 如果是目录,递归扫描
		if info.IsDir() {
			if err := loadStandaloneDir(render, path, funcMap); err != nil {
				return err
			}
		} else {
			// 如果是文件,直接加载
			if err := loadStandaloneFile(render, path, funcMap); err != nil {
				return err
			}
		}
	}

	return nil
}

// loadStandaloneDir 递归扫描目录并加载所有独立 HTML 文件
func loadStandaloneDir(render multitemplate.Renderer, dir string, funcMap template.FuncMap) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return &LoadError{
				Path:    path,
				Message: "独立视图目录遍历失败",
				Err:     err,
			}
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		// 只处理 HTML 文件
		if !strings.HasSuffix(strings.ToLower(path), ".html") {
			return nil
		}

		// 加载独立文件
		return loadStandaloneFile(render, path, funcMap)
	})
}

// loadStandaloneFile 加载单个独立视图文件(不继承布局)
func loadStandaloneFile(render multitemplate.Renderer, path string, funcMap template.FuncMap) error {
	// 使用相对路径作为模板名称,并统一使用正斜杠
	templateName := normalizeTemplateName(path)

	// 解析模板
	tmpl, err := template.New(filepath.Base(path)).Funcs(funcMap).ParseFiles(path)
	if err != nil {
		return &LoadError{
			Path:    path,
			Message: "独立模板解析失败",
			Err:     err,
		}
	}

	// 注册到渲染器
	render.Add(templateName, tmpl)

	return nil
}

// normalizeTemplateName 规范化模板名称
// 将反斜杠转换为正斜杠,确保跨平台兼容性
func normalizeTemplateName(path string) string {
	return filepath.ToSlash(path)
}

// fileExists 检查文件是否存在
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
