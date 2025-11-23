package viewloader

// TemplateGroup 表示一组需要继承相同布局模板的视图配置
type TemplateGroup struct {
	// Name 模板组名称,用于标识和日志输出
	Name string

	// LayoutPath 布局模板文件路径
	// 示例: "views/layouts/admin.html"
	LayoutPath string

	// ViewDirs 需要继承该布局的视图文件目录列表
	// 支持多级目录,将递归扫描目录下所有 .html 文件
	// 示例: []string{"views/admin", "views/admin/users"}
	ViewDirs []string
}

// Config 视图加载器配置
type Config struct {
	// TemplateGroups 模板组列表
	// 每个模板组包含一个布局模板和多个需要继承该布局的视图文件
	TemplateGroups []TemplateGroup

	// StandalonePaths 独立视图路径列表
	// 这些视图不需要继承任何布局模板,可以是文件路径或目录路径
	// - 文件路径示例: "views/auth/login.html"
	// - 目录路径示例: "views/public" (会递归扫描该目录下所有 HTML 文件)
	StandalonePaths []string

	// Funcs 自定义模板函数映射
	// 可选配置,用于在模板中使用自定义函数
	Funcs map[string]interface{}
}

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	if len(c.TemplateGroups) == 0 && len(c.StandalonePaths) == 0 {
		return ErrEmptyConfig
	}

	// 验证每个模板组
	for i, group := range c.TemplateGroups {
		if group.Name == "" {
			return &ConfigError{
				Field:   "TemplateGroups[" + string(rune(i)) + "].Name",
				Message: "模板组名称不能为空",
			}
		}
		if group.LayoutPath == "" {
			return &ConfigError{
				Field:   "TemplateGroups[" + string(rune(i)) + "].LayoutPath",
				Message: "布局模板路径不能为空",
			}
		}
		if len(group.ViewDirs) == 0 {
			return &ConfigError{
				Field:   "TemplateGroups[" + string(rune(i)) + "].ViewDirs",
				Message: "视图目录列表不能为空",
			}
		}
	}

	return nil
}
