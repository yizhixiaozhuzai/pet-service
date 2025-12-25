package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// 数据库连接
	dsn := "root:123456@tcp(localhost:3306)/pet_service?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 创建生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:           "./biz/query",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery,
		FieldNullable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
	})

	g.UseDB(db)

	// 生成users表的模型
	g.ApplyBasic(
		g.GenerateModel("users"),
	)

	// 执行生成
	g.Execute()

	fmt.Println("GORM代码生成成功!")
}
