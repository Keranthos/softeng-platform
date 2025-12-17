package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"softeng-platform/internal/config"
	"softeng-platform/internal/handler"
	"softeng-platform/internal/middleware"
	"softeng-platform/internal/repository"
	"softeng-platform/internal/service"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

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
	userService := service.NewUserService(userRepo)
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

	// 设置路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())

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
		tools.GET("/profile", toolHandler.GetTools)
		tools.GET("/search", toolHandler.SearchTools)
		tools.GET("/:resourceId", toolHandler.GetTool)
		tools.POST("/submit", middleware.AuthMiddleware(), toolHandler.SubmitTool)
		tools.POST("/:resourceId/views", toolHandler.AddView)
		tools.POST("/:resourceId/collections", middleware.AuthMiddleware(), toolHandler.CollectTool)
		tools.DELETE("/:resourceId/collections", middleware.AuthMiddleware(), toolHandler.UncollectTool)
		tools.POST("/:resourceId/comments", middleware.AuthMiddleware(), toolHandler.AddComment)
		tools.DELETE("/:resourceId/comments", middleware.AuthMiddleware(), toolHandler.DeleteComment)
		tools.POST("/:resourceId/comments/:commentId/reply", middleware.AuthMiddleware(), toolHandler.ReplyComment)
		tools.DELETE("/:resourceId/comments/:commentId/reply", middleware.AuthMiddleware(), toolHandler.DeleteReply)
		tools.POST("/:resourceId/like", middleware.AuthMiddleware(), toolHandler.LikeTool)
		tools.DELETE("/:resourceId/like", middleware.AuthMiddleware(), toolHandler.UnlikeTool)
	}

	// 课程路由
	courses := r.Group("/courses")
	{
		courses.GET("/profile", courseHandler.GetCourses)
		courses.GET("/search", courseHandler.SearchCourses)
		courses.GET("/:courseId", courseHandler.GetCourse)
		courses.POST("/:courseId/upload", middleware.AuthMiddleware(), courseHandler.UploadResource)
		courses.GET("/:courseId/textbooks/:textbookId/download", middleware.AuthMiddleware(), courseHandler.DownloadTextbook)
		courses.POST("/:courseId/comments", middleware.AuthMiddleware(), courseHandler.AddComment)
		courses.DELETE("/:courseId/comments", middleware.AuthMiddleware(), courseHandler.DeleteComment)
		courses.POST("/:courseId/comments/:commentId/reply", middleware.AuthMiddleware(), courseHandler.ReplyComment)
		courses.DELETE("/:courseId/comments/:commentId/reply", middleware.AuthMiddleware(), courseHandler.DeleteReply)
		courses.POST("/:courseId/view", courseHandler.AddView)
		courses.POST("/:courseId/collected", middleware.AuthMiddleware(), courseHandler.CollectCourse)
		courses.DELETE("/:courseId/collected", middleware.AuthMiddleware(), courseHandler.UncollectCourse)
		courses.POST("/:courseId/like", middleware.AuthMiddleware(), courseHandler.LikeCourse)
		courses.DELETE("/:courseId/like", middleware.AuthMiddleware(), courseHandler.UnlikeCourse)
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
		projects.POST("/:projectId/comments", middleware.AuthMiddleware(), projectHandler.AddComment)
		projects.DELETE("/:projectId/comments", middleware.AuthMiddleware(), projectHandler.DeleteComment)
		projects.POST("/:projectId/comments/:commentId/reply", middleware.AuthMiddleware(), projectHandler.ReplyComment)
		projects.DELETE("/:projectId/comments/:commentId/reply", middleware.AuthMiddleware(), projectHandler.DeleteReply)
		projects.POST("/:projectId/view", projectHandler.AddView)
		projects.POST("/:projectId/collected", middleware.AuthMiddleware(), projectHandler.CollectProject)
		projects.DELETE("/:projectId/collected", middleware.AuthMiddleware(), projectHandler.UncollectProject)
	}

	// 管理员路由
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware()) // 先验证身份
	admin.Use(middleware.AdminMiddleware()) // 再验证管理员权限
	{
		admin.GET("/pending", adminHandler.GetPending)
		admin.POST("/review/:itemId", adminHandler.ReviewItem) // 改为POST方法以支持requestBody
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
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
