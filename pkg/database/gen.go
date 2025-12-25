package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"pet-service/config"
)

// GenDB 生成数据库操作代码
func GenDB(cfg *config.Config) error {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./biz/query",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery,
		FieldNullable: true,
	})

	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN))
	if err != nil {
		return err
	}

	g.UseDB(db)

	// 自动生成所有表
	g.ApplyBasic(
	// g.GenerateModel("users"),
	// g.GenerateModel("products"),
	// 在这里添加需要生成的表
	)

	// 执行生成
	g.Execute()

	return nil
}

// Connect 连接数据库
func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	return db, nil
}
