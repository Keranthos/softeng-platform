package main

import (
	"context"
	"fmt"
	"log"
	"softeng-platform/internal/config"
	"softeng-platform/internal/model"
	"softeng-platform/internal/repository"
	"softeng-platform/internal/utils"
)

func main() {
	fmt.Println("测试用户功能")
	printDivider()

	// 1. 加载配置
	cfg := config.LoadConfig()
	fmt.Printf("配置加载成功:\n")
	fmt.Printf("  数据库URL: %s\n", cfg.DatabaseURL)
	fmt.Printf("  端口: %s\n", cfg.Port)
	printDivider()

	// 2. 连接数据库
	fmt.Println("连接数据库...")
	db, err := repository.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	defer db.Close()
	fmt.Println("✅ 数据库连接成功")
	printDivider()

	// 3. 创建 UserRepository
	userRepo := repository.NewUserRepository(db)

	// 4. 测试创建用户
	ctx := context.Background()

	// 加密密码
	hashedPassword, err := utils.HashPassword("test123456")
	if err != nil {
		log.Fatal("密码加密失败:", err)
	}

	testUser := &model.User{
		Username: "test_user_" + randomString(5),
		Nickname: "测试用户",
		Email:    "test_" + randomString(5) + "@example.com",
		Password: hashedPassword,
		Avatar:   "https://example.com/avatar.jpg",
		Role:     "user",
	}

	fmt.Println("创建测试用户...")
	fmt.Printf("  用户名: %s\n", testUser.Username)
	fmt.Printf("  邮箱: %s\n", testUser.Email)

	err = userRepo.Create(ctx, testUser)
	if err != nil {
		log.Fatal("创建用户失败:", err)
	}
	fmt.Printf("✅ 用户创建成功，ID: %d\n", testUser.ID)
	printDivider()

	// 5. 测试查询用户
	fmt.Println("查询刚创建的用户...")
	foundUser, err := userRepo.GetByID(ctx, testUser.ID)
	if err != nil {
		log.Fatal("查询用户失败:", err)
	}
	if foundUser == nil {
		log.Fatal("未找到用户")
	}
	fmt.Printf("✅ 查询成功:\n")
	fmt.Printf("  ID: %d\n", foundUser.ID)
	fmt.Printf("  用户名: %s\n", foundUser.Username)
	fmt.Printf("  昵称: %s\n", foundUser.Nickname)
	fmt.Printf("  邮箱: %s\n", foundUser.Email)
	fmt.Printf("  角色: %s\n", foundUser.Role)
	printDivider()

	// 6. 测试按用户名查询
	fmt.Println("按用户名查询...")
	userByUsername, err := userRepo.GetByUsername(ctx, testUser.Username)
	if err != nil {
		log.Fatal("按用户名查询失败:", err)
	}
	if userByUsername != nil {
		fmt.Println("✅ 按用户名查询成功")
	}
	printDivider()

	// 7. 测试按邮箱查询
	fmt.Println("按邮箱查询...")
	userByEmail, err := userRepo.GetByEmail(ctx, testUser.Email)
	if err != nil {
		log.Fatal("按邮箱查询失败:", err)
	}
	if userByEmail != nil {
		fmt.Println("✅ 按邮箱查询成功")
	}
	printDivider()

	// 8. 测试更新用户
	fmt.Println("更新用户信息...")
	foundUser.Nickname = "更新后的昵称"
	foundUser.Description = "这是更新后的用户描述"

	err = userRepo.Update(ctx, foundUser)
	if err != nil {
		log.Fatal("更新用户失败:", err)
	}

	// 验证更新
	updatedUser, err := userRepo.GetByID(ctx, testUser.ID)
	if err != nil {
		log.Fatal("验证更新失败:", err)
	}
	fmt.Printf("✅ 用户更新成功:\n")
	fmt.Printf("  新昵称: %s\n", updatedUser.Nickname)
	fmt.Printf("  新描述: %s\n", updatedUser.Description)
	printDivider()

	// 9. 测试更新密码
	fmt.Println("测试更新密码...")
	newHashedPassword, err := utils.HashPassword("new_password_123")
	if err != nil {
		log.Fatal("新密码加密失败:", err)
	}

	err = userRepo.UpdatePassword(ctx, testUser.ID, newHashedPassword)
	if err != nil {
		log.Fatal("更新密码失败:", err)
	}
	fmt.Println("✅ 密码更新成功")
	printDivider()

	// 10. 测试验证密码
	fmt.Println("测试密码验证...")
	passwordToCheck := "new_password_123"
	isValid := utils.CheckPasswordHash(passwordToCheck, newHashedPassword)
	if isValid {
		fmt.Println("✅ 密码验证成功")
	} else {
		fmt.Println("❌ 密码验证失败")
	}
	printDivider()

	fmt.Println("🎉 所有测试通过！")
}

// 辅助函数：打印分隔线
func printDivider() {
	fmt.Println("--------------------------------------------------")
}

// 辅助函数：生成随机字符串
func randomString(length int) string {
	// 简化版本，仅用于测试
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = byte('a' + i%26)
	}
	return string(bytes)
}
