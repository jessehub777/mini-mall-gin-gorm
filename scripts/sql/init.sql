-- ============================================
-- mini-mall-gin-gorm 初始化脚本（MySQL 8.0）
-- 说明：用于快速创建学习用数据库、表结构和演示数据
-- ============================================

-- 1) 创建数据库
CREATE DATABASE IF NOT EXISTS mini_mall
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_general_ci;

USE mini_mall;

-- 2) 用户表：登录 + 个人信息
CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  username VARCHAR(50) NOT NULL COMMENT '登录用户名，唯一',
  password VARCHAR(255) NOT NULL COMMENT '密码哈希',
  nickname VARCHAR(100) NOT NULL DEFAULT '' COMMENT '昵称',
  email VARCHAR(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  phone VARCHAR(30) NOT NULL DEFAULT '' COMMENT '手机号',
  address VARCHAR(255) NOT NULL DEFAULT '' COMMENT '地址',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (id),
  UNIQUE KEY uk_users_username (username)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 3) 商品表：商品信息 + 库存
CREATE TABLE IF NOT EXISTS products (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  name VARCHAR(120) NOT NULL COMMENT '商品名称',
  description VARCHAR(500) NOT NULL DEFAULT '' COMMENT '商品描述',
  price DECIMAL(10,2) NOT NULL COMMENT '商品单价',
  stock INT NOT NULL DEFAULT 0 COMMENT '库存',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- 4) 购买记录表：用户购买商品明细（用户-商品关系）
CREATE TABLE IF NOT EXISTS purchases (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
  quantity INT NOT NULL COMMENT '购买数量',
  unit_price DECIMAL(10,2) NOT NULL COMMENT '下单时商品单价快照',
  total_amount DECIMAL(10,2) NOT NULL COMMENT '总金额（unit_price * quantity）',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (id),
  KEY idx_purchases_user_id (user_id),
  KEY idx_purchases_product_id (product_id),
  CONSTRAINT fk_purchases_user_id FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_purchases_product_id FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购买记录表';

-- 5) 演示商品数据
INSERT INTO products (name, description, price, stock)
VALUES
('机械键盘', '87键热插拔机械键盘，支持RGB灯效', 299.00, 50),
('无线鼠标', '轻量化无线鼠标，办公与游戏通用', 129.00, 80),
('27寸显示器', '2K分辨率，支持75Hz刷新率', 1099.00, 20)
ON DUPLICATE KEY UPDATE
  description = VALUES(description),
  price = VALUES(price),
  stock = VALUES(stock);
