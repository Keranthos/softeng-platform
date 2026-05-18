# 工具使用快速指南

## 代码生成工具

### 快速开始

```bash
# 1. 进入项目根目录
cd softeng-platform/softeng-platform

# 2. 运行代码生成工具（示例：为 tools 表生成代码）
python tools/generate_repository.py tools

# 3. 查看生成的文件
# 文件位置: internal/repository/tools_repository_generated.go
```

### 生成的文件说明

生成的文件包含：
- Repository接口定义
- Repository基础实现（GetByID, Create, Update, Delete, GetList）

**注意**：生成的是模板代码，需要根据业务需求调整。

---

## API文档生成 (Swagger)

### 首次设置

```bash
# 1. 安装swag工具（如果还没安装）
go install github.com/swaggo/swag/cmd/swag@latest

# 2. 确保已安装依赖
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files

# 3. 生成文档
swag init -g cmd/server/main.go

# 4. 启动服务器
go run cmd/server/main.go

# 5. 访问文档
# 浏览器打开: http://localhost:8080/swagger/index.html
```

### 更新文档

每次修改Handler注释后，重新运行：
```bash
swag init -g cmd/server/main.go
```

### 添加新API注释示例

在Handler方法上添加注释：

```go
// GetTools 获取工具列表
// @Summary 获取工具列表
// @Description 根据分类、标签等条件获取工具列表，支持分页和排序
// @Tags tools
// @Accept json
// @Produce json
// @Param catagory query array false "工具分类"
// @Param tag query array false "工具标签"
// @Param sort query string false "排序方式"
// @Success 200 {object} map[string]interface{}
// @Router /tools [get]
func (h *ToolHandler) GetTools(c *gin.Context) {
    // ...
}
```

---

## 常见问题

### 代码生成工具报错？

1. 确保Python 3已安装
2. 确保schema.sql文件存在
3. 检查表名是否正确

### Swagger文档不显示？

1. 检查是否运行了 `swag init`
2. 检查main.go中是否添加了swagger路由
3. 检查注释格式是否正确

### 生成代码有语法错误？

这是正常的，生成的是模板代码，需要手动调整：
- 字段类型
- 业务逻辑
- 错误处理

---

## 下一步

- 查看 `tools/README.md` 获取详细文档
- 根据业务需求调整生成的代码
- 为所有API添加Swagger注释

