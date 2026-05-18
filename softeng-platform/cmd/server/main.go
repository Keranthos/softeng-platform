package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"softeng-platform/internal/config"
	"softeng-platform/internal/handler"
	"softeng-platform/internal/middleware"
	"softeng-platform/internal/repository"
	"softeng-platform/internal/service"
	"softeng-platform/pkg/response"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	
	// Swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title 软件工程平台 API
// @version 1.0
// @description 软件工程平台后端API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 使用 "Bearer {token}" 格式，token通过登录接口获取

func main() {
	// 初始化配置
	cfg := config.LoadConfig()

	// 初始化数据库
	db, err := repository.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}()

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	toolRepo := repository.NewToolRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo, toolRepo, projectRepo)
	toolService := service.NewToolService(toolRepo)
	courseService := service.NewCourseService(courseRepo)
	projectService := service.NewProjectService(projectRepo)
	adminService := service.NewAdminService(toolRepo, courseRepo, projectRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	toolHandler := handler.NewToolHandler(toolService)
	courseHandler := handler.NewCourseHandler(courseService)
	projectHandler := handler.NewProjectHandler(projectService)
	adminHandler := handler.NewAdminHandler(adminService)
	uploadHandler := handler.NewUploadHandler()

	// 设置路由
	r := gin.Default()
	
	// 强制输出调试信息，确认代码已重新编译
	log.Println("=== Server starting with updated routes ===")

	// 中间件
	r.Use(middleware.CORS())
	
	// 添加请求日志中间件（用于调试）
	r.Use(func(c *gin.Context) {
		log.Printf("[DEBUG] Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	// 静态文件服务（图片上传目录）
	r.Static("/uploads", "./uploads")
	
	// Swagger API文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 认证路由
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
	}

	// 用户路由
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.POST("/logout", userHandler.Logout)
		users.GET("/profile", userHandler.GetProfile)
		users.GET("/status", userHandler.GetStatus)
		users.GET("/collection", userHandler.GetCollection)
		users.POST("/update", userHandler.UpdateProfile)
		users.DELETE("/collection/:resourceType/:resourceId/", userHandler.DeleteCollection)
		users.GET("/summit", userHandler.GetSummit)
		users.PUT("/status/:resourceType/:resourceId/statu", userHandler.UpdateResourceStatus)
		users.POST("/profile/new_email", userHandler.UpdateEmail)
		users.POST("/profile/new_passward", userHandler.UpdatePassword) // 保持与API文档一致（即使拼写错误）
	}

	// 工具路由
	tools := r.Group("/tools")
	{
		// 具体路由放在前面
		tools.GET("/profile", toolHandler.GetTools)                                                    // 获取工具列表
		tools.GET("/search", toolHandler.SearchTools)                                                  // 搜索工具
		tools.POST("/submit", middleware.AuthMiddleware(), toolHandler.SubmitTool)                    // 提交工具
		
		// 更具体的参数路由放在前面
		tools.GET("/:resourceId/comments", toolHandler.GetComments)                                   // 获取工具评论（新增）
		tools.POST("/:resourceId/comments", middleware.AuthMiddleware(), toolHandler.AddComment)      // 发表评论
		tools.DELETE("/:resourceId/comments/:commentId", middleware.AuthMiddleware(), toolHandler.DeleteComment) // 删除评论（修正路径）
		tools.POST("/:resourceId/comments/:commentId/like", middleware.AuthMiddleware(), toolHandler.LikeComment) // 点赞评论（新增）
		tools.POST("/:resourceId/comments/:commentId/reply", middleware.AuthMiddleware(), toolHandler.ReplyComment) // 回复评论
		tools.DELETE("/:resourceId/comments/:commentId/reply", middleware.AuthMiddleware(), toolHandler.DeleteReply) // 删除回复
		tools.POST("/:resourceId/views", toolHandler.AddView)                                         // 增加浏览量
		tools.POST("/:resourceId/collections", middleware.AuthMiddleware(), toolHandler.CollectTool)  // 收藏工具
		tools.DELETE("/:resourceId/collections", middleware.AuthMiddleware(), toolHandler.UncollectTool) // 取消收藏
		tools.POST("/:resourceId/like", middleware.AuthMiddleware(), toolHandler.LikeTool)            // 点赞工具
		tools.DELETE("/:resourceId/like", middleware.AuthMiddleware(), toolHandler.UnlikeTool)        // 取消点赞
		
		// 最通用的参数路由放在最后
		tools.GET("/:resourceId", toolHandler.GetTool)                                                 // 获取工具详情
		tools.PUT("/:resourceId", middleware.AuthMiddleware(), toolHandler.UpdateTool)                // 更新工具
	}
	
	// 添加 NoRoute handler 用于调试
	r.NoRoute(func(c *gin.Context) {
		log.Printf("[DEBUG] NoRoute: Method=%s, Path=%s", c.Request.Method, c.Request.URL.Path)
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Route not found: %s %s", c.Request.Method, c.Request.URL.Path))
	})

	// 课程路由（使用单数course以匹配前端）
	course := r.Group("/course")
	{
		course.GET("", courseHandler.GetCourses)                         // 获取课程列表
		course.GET("/search", courseHandler.SearchCourses)               // 搜索课程
		course.GET("/:courseId", courseHandler.GetCourse)                // 获取课程详情
		course.POST("/submit", middleware.AuthMiddleware(), courseHandler.SubmitCourse) // 提交课程（新增）
		course.POST("/:courseId/view", courseHandler.AddView)            // 增加浏览量
		course.POST("/:courseId/collections", middleware.AuthMiddleware(), courseHandler.CollectCourse) // 收藏课程
		course.DELETE("/:courseId/collections", middleware.AuthMiddleware(), courseHandler.UncollectCourse) // 取消收藏
		course.POST("/:courseId/like", middleware.AuthMiddleware(), courseHandler.LikeCourse) // 点赞课程
		course.DELETE("/:courseId/like", middleware.AuthMiddleware(), courseHandler.UnlikeCourse) // 取消点赞
		course.GET("/:courseId/comments", courseHandler.GetComments)     // 获取课程评论（新增）
		course.POST("/:courseId/comments", middleware.AuthMiddleware(), courseHandler.AddComment) // 发表评论
		course.DELETE("/:courseId/comments/:commentId", middleware.AuthMiddleware(), courseHandler.DeleteComment) // 删除评论
		course.POST("/:courseId/comments/:commentId/like", middleware.AuthMiddleware(), courseHandler.LikeComment) // 点赞评论（新增）
		course.POST("/:courseId/comments/:commentId/reply", middleware.AuthMiddleware(), courseHandler.ReplyComment) // 回复评论
		course.DELETE("/:courseId/comments/:commentId/reply", middleware.AuthMiddleware(), courseHandler.DeleteReply) // 删除回复
		course.GET("/:courseId/resources", courseHandler.GetResources)   // 获取课程资源（新增）
		course.POST("/:courseId/resources", middleware.AuthMiddleware(), courseHandler.UploadResource) // 上传资源（改为resources）
		course.GET("/:courseId/textbooks/:textbookId/download", middleware.AuthMiddleware(), courseHandler.DownloadTextbook) // 下载课本
		// 以下为前端定义但可能暂时不实现的功能
		// course.POST("/:courseId/learning-plan", ...) // 加入学习计划
		// course.DELETE("/:courseId/learning-plan", ...) // 从学习计划移除
		// course.GET("/:courseId/progress", ...) // 获取学习进度
		// course.PUT("/:courseId/progress/:chapterId", ...) // 更新学习进度
		// course.POST("/:courseId/rating", ...) // 添加课程评分
		// course.POST("/analyze", ...) // 分析课程链接
	}

	// 项目路由
	projects := r.Group("/projects")
	{
		projects.GET("/profile", projectHandler.GetProjects)
		projects.GET("/search", projectHandler.SearchProjects)
		projects.GET("/:projectId", projectHandler.GetProject)
		projects.PUT("/:projectId", middleware.AuthMiddleware(), projectHandler.UpdateProject)
		projects.POST("/upload", middleware.AuthMiddleware(), projectHandler.UploadProject)
		projects.POST("/:projectId/like", middleware.AuthMiddleware(), projectHandler.LikeProject)
		projects.DELETE("/:projectId/like", middleware.AuthMiddleware(), projectHandler.UnlikeProject)
		projects.GET("/:projectId/comments", projectHandler.GetComments)                                    // 获取项目评论列表
		projects.POST("/:projectId/comments", middleware.AuthMiddleware(), projectHandler.AddComment)      // 发表评论
		projects.DELETE("/:projectId/comments/:commentId", middleware.AuthMiddleware(), projectHandler.DeleteComment) // 删除评论
		projects.POST("/:projectId/comments/:commentId/like", middleware.AuthMiddleware(), projectHandler.LikeComment) // 点赞评论
		projects.POST("/:projectId/comments/:commentId/reply", middleware.AuthMiddleware(), projectHandler.ReplyComment) // 回复评论
		projects.DELETE("/:projectId/comments/:commentId/reply", middleware.AuthMiddleware(), projectHandler.DeleteReply) // 删除回复
		projects.POST("/:projectId/view", projectHandler.AddView)
		projects.POST("/:projectId/collected", middleware.AuthMiddleware(), projectHandler.CollectProject)
		projects.DELETE("/:projectId/collected", middleware.AuthMiddleware(), projectHandler.UncollectProject)
	}

	// 管理员路由
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware()) // 先验证身份
	admin.Use(middleware.AdminMiddleware()) // 再验证管理员权限
	{
		admin.GET("/pending", adminHandler.GetPending)              // 获取待审核内容
		admin.POST("/review/:itemId", adminHandler.ReviewItem)      // 审核项目（支持POST和GET，前端使用GET但需要requestBody，所以用POST）
		admin.GET("/review/:itemId", adminHandler.ReviewItem)       // 也支持GET方法（前端调用的是GET）
	}

	// 上传路由
	upload := r.Group("/api/upload")
	upload.Use(middleware.AuthMiddleware()) // 需要登录才能上传
	{
		upload.POST("/image", uploadHandler.UploadImage)           // 上传图片文件
		upload.POST("/process", uploadHandler.ProcessImageURL)    // 处理图片URL（自动本地化）
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// 在goroutine中启动服务器
	serverError := make(chan error, 1)
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
			serverError <- err
		}
	}()

	// 等待一小段时间确保服务器启动成功
	time.Sleep(500 * time.Millisecond)
	
	// 检查服务器是否成功启动
	select {
	case err := <-serverError:
		log.Fatalf("Failed to start server: %v", err)
	default:
		log.Println("Server is running on port", cfg.Port)
		log.Println("Press Ctrl+C to stop the server")
	}

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// 监听中断信号（在Windows下使用os.Interrupt，Unix下使用syscall.SIGTERM）
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	
	// 等待信号或服务器错误
	select {
	case sig := <-quit:
		log.Printf("Received signal: %v", sig)
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	}
	
	log.Println("Shutting down server...")

	// 设置5秒的超时时间用于优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
