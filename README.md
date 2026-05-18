# 详细代码分析

## 1. 程序入口 (main.go)
**功能：** 应用程序的启动入口，负责初始化所有组件和设置路由

### 关键流程：
- 加载配置（端口、数据库连接、JWT密钥）
- 连接数据库
- 初始化各层组件（Repository → Service → Handler）
- 配置 Gin 路由和中间件
- 启动 HTTP 服务器

### 路由分组：
- **/auth**：用户认证相关（注册、登录、忘记密码）
- **/users**：用户个人中心（需要认证）
- **/tools**：工具资源管理
- **/courses**：课程资源管理
- **/projects**：项目资源管理
- **/admin**：管理员功能（需要管理员权限）

---

## 2. 配置管理 (config.go)
使用环境变量或默认值配置应用参数

### 关键配置项：
- `PORT`：服务器端口（默认 8080）
- `DATABASE_URL`：MySQL 连接字符串
- `JWT_SECRET`：JWT 签名密钥

---

## 3. 数据模型层 (model/ 目录)
定义了应用的核心数据结构：

### user.go：用户模型及相关请求结构体
- `User`：用户基本信息
- `RegisterRequest`：注册请求
- `LoginRequest`：登录请求
- `UpdateProfileRequest`：更新资料请求

### common.go：通用资源模型
- `ResourceWeb`：网页资源
- `ResourceUpload`：上传资源
- `ResourceReview`：审核资源
- `Comment`：评论结构

### course.go：课程相关模型
- `Course`：课程基本信息
- `CourseDetail`：课程详情
- `CourseComment`：课程评论

### project.go：项目相关模型
- `Project`：项目基本信息
- `ProjectDetail`：项目详情
- `ProjectComment`：项目评论

### tool.go：工具相关模型
- `Tool`：工具详细信息
- `ToolReview`：工具审核信息

---

## 4. 处理器层 (handler/ 目录)
处理 HTTP 请求，接收参数并调用 Service 层：

### auth.go：处理认证相关请求
- `Register`：用户注册
- `Login`：用户登录
- `ForgotPassword`：重置密码

### user.go：处理用户个人中心请求
- `GetProfile`：获取用户资料
- `UpdateProfile`：更新资料
- `GetCollection`：获取收藏
- `DeleteCollection`：删除收藏
- `GetStatus`：获取审核状态
- `GetSummit`：获取个人提交

### tool.go：处理工具相关请求
- `GetTools`：获取工具列表
- `SearchTools`：搜索工具
- `SubmitTool`：提交新工具
- 点赞、收藏、评论等互动功能

### course.go：处理课程相关请求
- `GetCourses`：获取课程列表
- `UploadResource`：上传课程资源
- `DownloadTextbook`：下载教材
- 评论、收藏、点赞等功能

### project.go：处理项目相关请求
- `GetProjects`：获取项目列表
- `UploadProject`：上传项目
- `UpdateProject`：更新项目
- 互动功能

### admin.go：管理员功能
- `GetPending`：获取待审核内容
- `ReviewItem`：审核项目

---

## 5. 业务逻辑层 (service/ 目录)
实现核心业务逻辑：

### auth.go：认证业务逻辑
- 用户注册、登录、密码重置
- 邮箱验证码、邀请码验证

### user.go：用户业务逻辑
- 个人资料管理
- 收藏管理
- 审核状态查询

### tool.go：工具业务逻辑
- 工具提交、搜索、展示
- 互动功能处理

### course.go：课程业务逻辑
- 课程资源管理
- 教材下载
- 互动功能

### project.go：项目业务逻辑
- 项目管理
- 互动功能

### admin.go：管理业务逻辑
- 审核内容管理
- 审批操作

---

## 6. 数据访问层 (repository/ 目录)
数据库操作抽象层：

### database.go：数据库连接管理

### user.go：用户数据操作（实际数据库操作）
- CRUD 操作
- 支持 PostgreSQL（从代码中的 $1 参数占位符判断）

### tool.go、course.go、project.go：
- 目前返回模拟数据（硬编码）
- 定义了完整的接口，但未实现真实数据库操作

---

## 7. 中间件 (middleware/ 目录)

### auth.go：
- `AuthMiddleware`：JWT 认证中间件
- `AdminMiddleware`：管理员权限检查

### cors.go：跨域资源共享配置

---

## 8. 工具函数 (utils/ 目录)

### jwt.go：JWT 令牌生成和验证

### password.go：密码加密和验证（使用 bcrypt）

### validator.go：数据验证工具

### response.go：统一响应格式

---

## 9. 依赖管理 (go.sum)
项目依赖的主要包：
- Web 框架：gin-gonic/gin
- 数据库驱动：database/sql + 疑似 PostgreSQL 驱动
- JWT 处理：golang-jwt/jwt/v4
- 配置管理：joho/godotenv
- 密码加密：golang.org/x/crypto/bcrypt
- 数据验证：go-playground/validator/v10

---

## 核心功能模块

### 1. 用户系统
- 注册、登录、密码重置
- 个人资料管理
- 收藏管理
- 提交记录跟踪

### 2. 资源管理三大模块
#### 工具 (Tools)
- 软件工具分享
- 分类标签管理
- 用户评价系统

#### 课程 (Courses)
- 课程资源共享
- 教材下载
- 教学资料上传

#### 项目 (Projects)
- 项目作品展示
- GitHub 链接管理
- 技术栈标签

### 3. 互动功能
- 点赞/取消点赞
- 收藏/取消收藏
- 评论/回复
- 浏览量统计

### 4. 审核系统
- 内容提交后进入待审核状态
- 管理员审批（通过/拒绝）
- 拒绝时提供理由

### 5. 权限控制
- 普通用户：浏览、提交、互动
- 认证用户：管理个人内容
- 管理员：内容审核、系统管理

---

## 技术特点

### ✅ 优点
- **架构清晰：** 严格的分层设计，职责分离明确
- **接口定义完整：** 各层之间有清晰的接口约定
- **错误处理统一：** 使用统一的响应格式
- **安全性考虑：** JWT 认证、密码加密、CORS 配置
- **可扩展性：** 模块化设计，易于添加新功能

---
