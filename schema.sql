-- HouseManager Database Schema (GoFly规范版本)
-- 根据 GoFly 框架数据表规范优化，合并现有表结构

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 1. 用户表 (基于现有business_user优化)
-- ----------------------------
DROP TABLE IF EXISTS `business_user`;
CREATE TABLE `business_user` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  `username` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '用户名',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '姓名',
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '昵称',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '备注',
  `password` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '密码',
  `salt` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '密码盐',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '电子邮箱',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '手机号',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '头像',
  `sex` tinyint(1) UNSIGNED NOT NULL DEFAULT 0 COMMENT '性别:1男,2女,0未知',
  `successions` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '连续登录天数',
  `maxsuccessions` int(10) UNSIGNED NOT NULL DEFAULT 1 COMMENT '最大连续登录天数',
  `prevtime` bigint(20) NULL DEFAULT NULL COMMENT '上次登录时间',
  `logintime` bigint(20) NULL DEFAULT NULL COMMENT '登录时间',
  `loginip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '登录IP',
  `loginfailure` tinyint(3) UNSIGNED NOT NULL DEFAULT 0 COMMENT '失败次数',
  `role` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'user' COMMENT '经纪人',
  `store_id` int(11) NOT NULL DEFAULT 0 COMMENT '所属门店ID',
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '头衔',
  `introduction` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '自我介绍',
  `openid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '微信OpenID',
  `unionid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '微信UnionID',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
  `createtime` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `updatetime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  `deletetime` datetime(0) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_username` (`username`) USING BTREE,
  INDEX `idx_email` (`email`) USING BTREE,
  INDEX `idx_mobile` (`mobile`) USING BTREE,
  INDEX `idx_openid` (`openid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 110 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '经纪人表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 2. 门店表 (含多图)
-- ----------------------------
DROP TABLE IF EXISTS `business_stores`;
CREATE TABLE `business_stores` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '门店名称',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '门店地址',
  `contact_phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '联系电话',
  `manager_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '店长姓名',
  `images` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '门店图片(多图,逗号分隔)',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
  `weigh` int(11) NOT NULL DEFAULT 0 COMMENT '排序权重',
  `createtime` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `updatetime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  `deletetime` datetime(0) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '门店表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 3. 房源表 (含多图)
-- ----------------------------
DROP TABLE IF EXISTS `business_properties`;
CREATE TABLE `business_properties` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '房源标题',
  `price` decimal(12, 2) NOT NULL DEFAULT 0.00 COMMENT '价格',
  `price_unit` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '万' COMMENT '价格单位',
  `area` decimal(10, 2) NOT NULL DEFAULT 0.00 COMMENT '面积(平米)',
  `rooms` tinyint(2) NOT NULL DEFAULT 0 COMMENT '室',
  `halls` tinyint(2) NOT NULL DEFAULT 0 COMMENT '厅',
  `bathrooms` tinyint(2) NOT NULL DEFAULT 0 COMMENT '卫',
  `floor_level` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '楼层',
  `total_floors` int(11) NOT NULL DEFAULT 0 COMMENT '总楼层',
  `orientation` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '朝向',
  `build_year` int(11) NOT NULL DEFAULT 0 COMMENT '建成年份',
  `property_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '物业类型',
  `decoration_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '装修类型',
  `community_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '小区名称',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '详细地址',
  `latitude` decimal(10, 6) NULL DEFAULT NULL COMMENT '纬度',
  `longitude` decimal(10, 6) NULL DEFAULT NULL COMMENT '经度',
  `tags` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '标签(逗号分隔)',
  `images` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '房源图片(多图,逗号分隔)',
  `cover_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '封面图',
  `has_smart_lock` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否有智能锁:0无,1有',
  `commission_rate` decimal(5, 2) NULL DEFAULT 0.00 COMMENT '佣金比例',
  `commission_reward` decimal(12, 2) NULL DEFAULT 0.00 COMMENT '成交奖励金',
  `owner_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '业主姓名',
  `owner_phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '业主电话',
  `receiver_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '收房人姓名',
  `receiver_phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '收房人电话',
  `receiver_price` decimal(12, 2) NULL DEFAULT NULL COMMENT '收房价格(支付给业主)',
  `agent_id` int(11) NOT NULL DEFAULT 0 COMMENT '维护经纪人ID',
  `sale_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'on_sale' COMMENT '销售状态:on_sale在售,sold已售,off_market下架',
  `hot_status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '推荐状态:0不推荐,1推荐',
  `view_count` int(11) NOT NULL DEFAULT 0 COMMENT '浏览次数',
  `follow_count` int(11) NOT NULL DEFAULT 0 COMMENT '关注人数',
  `showing_count` int(11) NOT NULL DEFAULT 0 COMMENT '带看次数',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
  `weigh` int(11) NOT NULL DEFAULT 0 COMMENT '排序权重',
  `createtime` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `updatetime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  `deletetime` datetime(0) NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_agent_id` (`agent_id`) USING BTREE,
  INDEX `idx_sale_status` (`sale_status`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '房源表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 4. 装修进度表 (含多图)
-- ----------------------------
DROP TABLE IF EXISTS `business_renovations`;
CREATE TABLE `business_renovations` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `property_id` int(11) NOT NULL DEFAULT 0 COMMENT '房源ID',
  `renovation_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'none' COMMENT '装修状态:none未装修,in_progress装修中,done已完成',
  `progress_percentage` tinyint(3) NOT NULL DEFAULT 0 COMMENT '进度百分比',
  `current_stage` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '当前阶段',
  `start_date` date NULL DEFAULT NULL COMMENT '开始日期',
  `estimated_finish_date` date NULL DEFAULT NULL COMMENT '预计完工日期',
  `actual_finish_date` date NULL DEFAULT NULL COMMENT '实际完工日期',
  `materials` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '材料列表(逗号分隔)',
  `images` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '装修图片(多图,逗号分隔)',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '施工说明',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
  `createtime` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `updatetime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_property_id` (`property_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '装修进度表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 5. 用户收藏表
-- ----------------------------
DROP TABLE IF EXISTS `business_favorites`;
CREATE TABLE `business_favorites` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '用户ID',
  `property_id` int(11) NOT NULL DEFAULT 0 COMMENT '房源ID',
  `createtime` datetime(0) NULL DEFAULT NULL COMMENT '关注时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_user_property` (`user_id`, `property_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户收藏表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 5.1 通通锁（TTLock）授权表
-- ----------------------------
DROP TABLE IF EXISTS `business_ttlock_accounts`;
CREATE TABLE `business_ttlock_accounts` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  `access_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'TTLock access_token',
  `refresh_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'TTLock refresh_token',
  `expire_at` bigint(20) NOT NULL DEFAULT 0 COMMENT 'access_token 过期时间(秒时间戳)',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
  `createtime` bigint(20) NOT NULL DEFAULT 0 COMMENT '创建时间(秒时间戳)',
  `updatetime` bigint(20) NOT NULL DEFAULT 0 COMMENT '更新时间(秒时间戳)',
  `deletetime` bigint(20) NOT NULL DEFAULT 0 COMMENT '删除时间(秒时间戳)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_business_id` (`business_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'TTLock 授权表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 5.2 智能锁表（TTLock 锁主数据）
-- ----------------------------
DROP TABLE IF EXISTS `business_smart_locks`;
CREATE TABLE `business_smart_locks` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  `ttlock_lock_id` bigint(20) NOT NULL DEFAULT 0 COMMENT 'TTLock LockId',
  `lock_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '锁名称',
  `lock_mac` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '锁MAC',
  `battery` int(11) NOT NULL DEFAULT 0 COMMENT '电量',
  `model_num` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '型号',
  `raw_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '原始数据(JSON)',
  `last_sync_at` bigint(20) NOT NULL DEFAULT 0 COMMENT '最近同步时间(秒时间戳)',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
  `createtime` bigint(20) NOT NULL DEFAULT 0 COMMENT '创建时间(秒时间戳)',
  `updatetime` bigint(20) NOT NULL DEFAULT 0 COMMENT '更新时间(秒时间戳)',
  `deletetime` bigint(20) NOT NULL DEFAULT 0 COMMENT '删除时间(秒时间戳)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_business_lock` (`business_id`, `ttlock_lock_id`) USING BTREE,
  INDEX `idx_business_id` (`business_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '智能锁表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 5.3 房源-锁绑定表（默认一房一锁）
-- ----------------------------
DROP TABLE IF EXISTS `business_property_locks`;
CREATE TABLE `business_property_locks` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  `property_id` int(11) NOT NULL DEFAULT 0 COMMENT '房源ID',
  `ttlock_lock_id` bigint(20) NOT NULL DEFAULT 0 COMMENT 'TTLock LockId',
  `bind_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '绑定状态:1已绑,0解绑',
  `bind_by_user_id` int(11) NOT NULL DEFAULT 0 COMMENT '绑定人(business_user.id)',
  `bind_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '绑定时间(秒时间戳)',
  `unbind_time` bigint(20) NULL DEFAULT NULL COMMENT '解绑时间(秒时间戳)',
  `createtime` bigint(20) NOT NULL DEFAULT 0 COMMENT '创建时间(秒时间戳)',
  `updatetime` bigint(20) NOT NULL DEFAULT 0 COMMENT '更新时间(秒时间戳)',
  `deletetime` bigint(20) NOT NULL DEFAULT 0 COMMENT '删除时间(秒时间戳)',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_business_property` (`business_id`, `property_id`) USING BTREE,
  INDEX `idx_lock_id` (`ttlock_lock_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '房源智能锁绑定表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 5.4 智能锁事件日志表
-- ----------------------------
DROP TABLE IF EXISTS `business_lock_events`;
CREATE TABLE `business_lock_events` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '操作人(business_user.id)',
  `property_id` int(11) NOT NULL DEFAULT 0 COMMENT '房源ID',
  `ttlock_lock_id` bigint(20) NOT NULL DEFAULT 0 COMMENT 'TTLock LockId',
  `event_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '事件类型',
  `result` tinyint(1) NOT NULL DEFAULT 0 COMMENT '结果:1成功,0失败',
  `err_code` int(11) NOT NULL DEFAULT 0 COMMENT '错误码',
  `err_msg` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '错误信息',
  `meta_json` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '扩展数据(JSON)',
  `createtime` bigint(20) NOT NULL DEFAULT 0 COMMENT '创建时间(秒时间戳)',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_business_id` (`business_id`) USING BTREE,
  INDEX `idx_property_id` (`property_id`) USING BTREE,
  INDEX `idx_lock_id` (`ttlock_lock_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '智能锁事件日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 6. 海报模板表
-- ----------------------------
DROP TABLE IF EXISTS `business_poster_templates`;
CREATE TABLE `business_poster_templates` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '模板名称',
  `preview_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '预览图URL',
  `template_config` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '配置信息(JSON)',
  `template_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'property' COMMENT '模板类型:property房源,agent经纪人,festive节日',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
  `weigh` int(11) NOT NULL DEFAULT 0 COMMENT '排序权重',
  `createtime` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `updatetime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '海报模板表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 7. 用户活动日志表
-- ----------------------------
DROP TABLE IF EXISTS `business_user_activity_logs`;
CREATE TABLE `business_user_activity_logs` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '用户ID',
  `property_id` int(11) NOT NULL DEFAULT 0 COMMENT '房源ID',
  `activity_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '活动类型:view浏览,share分享,call拨打,showing带看,unlock开锁',
  `meta_data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '元数据(JSON)',
  `createtime` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_type` (`user_id`, `activity_type`) USING BTREE,
  INDEX `idx_property` (`property_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户活动日志表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 8. 房源状态变更记录表
-- ----------------------------
DROP TABLE IF EXISTS `business_property_status_logs`;
CREATE TABLE `business_property_status_logs` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `business_id` int(11) NOT NULL DEFAULT 0 COMMENT '业务主账号id',
  `property_id` int(11) NOT NULL DEFAULT 0 COMMENT '房源ID',
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '操作人(business_user.id)',
  `field` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '变更字段:status/hot_status/sale_status/weigh等',
  `before_value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '变更前',
  `after_value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '变更后',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '备注/原因',
  `createtime` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_biz_property` (`business_id`, `property_id`) USING BTREE,
  INDEX `idx_property_time` (`property_id`, `createtime`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '房源状态变更记录表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- 9. 开锁记录表
-- ----------------------------
DROP TABLE IF EXISTS `business_unlock_requests`;
CREATE TABLE `business_unlock_requests` (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` int(11) NOT NULL DEFAULT 0 COMMENT '申请人ID',
  `property_id` int(11) NOT NULL DEFAULT 0 COMMENT '房源ID',
  `request_status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'pending' COMMENT '状态:pending待审核,approved已通过,rejected已拒绝,completed已完成,cancelled已取消',
  `request_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'password' COMMENT '开锁方式:password密码,bluetooth蓝牙',
  `request_time` datetime(0) NULL DEFAULT NULL COMMENT '申请时间',
  `expires_at` datetime(0) NULL DEFAULT NULL COMMENT '过期时间',
  `approval_remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '审批备注',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '状态:0正常,1禁用',
  `createtime` datetime(0) NULL DEFAULT NULL COMMENT '创建时间',
  `updatetime` datetime(0) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id` (`user_id`) USING BTREE,
  INDEX `idx_property_id` (`property_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '开锁记录表' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
