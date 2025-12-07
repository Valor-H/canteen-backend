-- 餐厅管理系统数据库初始化脚本

-- 用户表
CREATE TABLE IF NOT EXISTS `sys_user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `dept_id` int(11) DEFAULT NULL,
  `nick_name` varchar(50) DEFAULT NULL,
  `count` int(11) DEFAULT '0',
  `card_no` varchar(50) DEFAULT NULL,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `idx_card_no` (`card_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 订餐记录表
CREATE TABLE IF NOT EXISTS `order_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `meal_id` int(11) DEFAULT NULL,
  `status` varchar(20) DEFAULT '已报餐',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `meal_type` varchar(20) DEFAULT NULL,
  `week_number` varchar(20) DEFAULT NULL,
  `order_date` date DEFAULT NULL,
  `weekday` varchar(10) DEFAULT NULL,
  `setmeal_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_week_number` (`week_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 周套餐表
CREATE TABLE IF NOT EXISTS `weekly_setmeal` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `week_number` varchar(20) NOT NULL,
  `weekday` varchar(10) NOT NULL,
  `meal_type` varchar(20) NOT NULL,
  `setmeal_id` int(11) DEFAULT NULL,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `create_user` int(11) DEFAULT NULL,
  `remark` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_week_number` (`week_number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 套餐表
CREATE TABLE IF NOT EXISTS `setmeal` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) DEFAULT NULL,
  `code` varchar(50) DEFAULT NULL,
  `description` text,
  `status` varchar(20) DEFAULT '启用',
  `sort` int(11) DEFAULT '0',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_user` int(11) DEFAULT NULL,
  `update_user` int(11) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 套餐菜品关联表
CREATE TABLE IF NOT EXISTS `setmeal_dish` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `setmeal_id` int(11) NOT NULL,
  `dish_id` int(11) NOT NULL,
  `sort` int(11) DEFAULT '0',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_setmeal_id` (`setmeal_id`),
  KEY `idx_dish_id` (`dish_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 菜品表
CREATE TABLE IF NOT EXISTS `dish` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `category_id` int(11) DEFAULT NULL,
  `code` varchar(50) DEFAULT NULL,
  `description` text,
  `status` varchar(20) DEFAULT '启用',
  `sort` int(11) DEFAULT '0',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_user` int(11) DEFAULT NULL,
  `update_user` int(11) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `image` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 菜品分类表
CREATE TABLE IF NOT EXISTS `dish_category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `code` varchar(50) DEFAULT NULL,
  `sort` int(11) DEFAULT '0',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 系统配置表
CREATE TABLE IF NOT EXISTS `canteen_config` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `config_key` varchar(50) NOT NULL,
  `config_value` varchar(200) DEFAULT NULL,
  `description` varchar(200) DEFAULT NULL,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_config_key` (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


----------------- TEST ---------------
-- -- 插入一些基础配置数据
-- INSERT INTO `canteen_config` (`config_key`, `config_value`, `description`) VALUES
-- ('flexible_dept_id', '1', '灵活部门ID'),
-- ('fixed_dept_id', '2', '固定部门ID'),
-- ('flexible_dinner_start_time', '17:00', '灵活晚餐开始时间'),
-- ('fixed_dinner_start_time', '18:00', '固定晚餐开始时间')
-- ON DUPLICATE KEY UPDATE config_value=VALUES(config_value);

-- -- 插入一些基础菜品分类
-- INSERT INTO `dish_category` (`name`, `code`, `sort`) VALUES
-- ('主食', 'STAPLE', 1),
-- ('荤菜', 'MEAT', 2),
-- ('素菜', 'VEGETABLE', 3),
-- ('汤类', 'SOUP', 4),
-- ('饮品', 'DRINK', 5)
-- ON DUPLICATE KEY UPDATE name=VALUES(name);

-- -- 插入一些基础菜品
-- INSERT INTO `dish` (`name`, `category_id`, `code`, `sort`, `create_user`) VALUES
-- ('米饭', 1, 'RICE', 1, 1),
-- ('馒头', 1, 'BUN', 2, 1),
-- ('红烧肉', 2, 'BRAISED_PORK', 1, 1),
-- ('糖醋排骨', 2, 'SWEET_SOUR_RIBS', 2, 1),
-- ('清炒时蔬', 3, 'STIR_FRIED_VEGETABLES', 1, 1),
-- ('麻婆豆腐', 3, 'MAP_TOFU', 2, 1),
-- ('番茄鸡蛋汤', 4, 'TOMATO_EGG_SOUP', 1, 1),
-- ('紫菜蛋花汤', 4, 'SEAWEED_SOUP', 2, 1),
-- ('豆浆', 5, 'SOY_MILK', 1, 1),
-- ('果汁', 5, 'JUICE', 2, 1)
-- ON DUPLICATE KEY UPDATE name=VALUES(name);

-- -- 插入一些基础套餐
-- INSERT INTO `setmeal` (`name`, `code`, `description`, `sort`, `create_user`) VALUES
-- ('经济套餐', 'ECONOMY', '经济实惠套餐', 1, 1),
-- ('标准套餐', 'STANDARD', '营养均衡标准套餐', 2, 1),
-- ('豪华套餐', 'DELUXE', '丰富搭配豪华套餐', 3, 1)
-- ON DUPLICATE KEY UPDATE name=VALUES(name);

-- -- 为套餐添加菜品
-- INSERT INTO `setmeal_dish` (`setmeal_id`, `dish_id`, `sort`) VALUES
-- -- 经济套餐
-- (1, 1, 1),  -- 米饭
-- (1, 3, 2),  -- 红烧肉
-- (1, 5, 3),  -- 清炒时蔬
-- (1, 7, 4),  -- 番茄鸡蛋汤
-- -- 标准套餐
-- (2, 1, 1),  -- 米饭
-- (2, 4, 2),  -- 糖醋排骨
-- (2, 6, 3),  -- 麻婆豆腐
-- (2, 8, 4),  -- 紫菜蛋花汤
-- (2, 9, 5),  -- 豆浆
-- -- 豪华套餐
-- (3, 1, 1),  -- 米饭
-- (3, 4, 2),  -- 糖醋排骨
-- (3, 3, 3),  -- 红烧肉
-- (3, 5, 4),  -- 清炒时蔬
-- (3, 7, 5),  -- 番茄鸡蛋汤
-- (3, 9, 6)   -- 豆浆
-- ON DUPLICATE KEY UPDATE sort=VALUES(sort);

-- -- 插入一个测试用户
-- INSERT INTO `sys_user` (`user_id`, `dept_id`, `nick_name`, `count`, `card_no`) VALUES
-- (1, 1, '测试用户', 10, 'TEST001')
-- ON DUPLICATE KEY UPDATE nick_name=VALUES(nick_name);