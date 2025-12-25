CREATE DATABASE IF NOT EXISTS pet_service CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE pet_service;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    email VARCHAR(100) UNIQUE COMMENT '邮箱',
    phone VARCHAR(20) UNIQUE COMMENT '手机号',
    nickname VARCHAR(50) COMMENT '昵称',
    avatar VARCHAR(255) COMMENT '头像',
    status TINYINT DEFAULT 1 COMMENT '状态:0禁用,1正常',
    is_deleted TINYINT DEFAULT 0 COMMENT '是否删除:0否,1是',
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_phone (phone),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 插入测试数据
INSERT INTO users (username, password, email, phone, nickname, avatar, status) VALUES
('admin', '123456', 'admin@example.com', '13800138000', '管理员', 'https://example.com/avatar/admin.png', 1),
('testuser', '123456', 'test@example.com', '13800138001', '测试用户', 'https://example.com/avatar/test.png', 1);
