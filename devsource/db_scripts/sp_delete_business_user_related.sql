/*
  GoFly Admin - 删除 business_user 用户关联数据（安全版：不删除业务资产）

  适用场景：
  - 只输入 business_user.id，即可清理该用户的“个人关联数据”
  - 不会删除门店/房源/智能锁等业务资产数据
  - 对业务资产表中引用该用户的字段，会做“解除引用”（置 0）

  使用方式：
  1）在目标数据库执行本脚本（创建存储过程）
  2）执行：CALL sp_delete_business_user_related(用户ID);

  说明：
  - 为了兼容不同版本表结构，过程内部会先判断表/字段是否存在，再执行对应 SQL
  - business_user 采用“脱敏 + 软删除（deletetime）”方式，避免其它未知表仍引用该 ID 时产生问题
*/

DROP PROCEDURE IF EXISTS `sp_delete_business_user_related`;
DELIMITER $$
CREATE PROCEDURE `sp_delete_business_user_related`(IN p_user_id INT)
BEGIN
  DECLARE v_db VARCHAR(64);
  DECLARE v_user_exists INT DEFAULT 0;

  -- 影响行数统计（便于执行后核对）
  DECLARE v_upd_properties_agent_id INT DEFAULT 0;
  DECLARE v_upd_property_locks_bind_by_user_id INT DEFAULT 0;
  DECLARE v_del_favorites INT DEFAULT 0;
  DECLARE v_del_user_activity_logs INT DEFAULT 0;
  DECLARE v_del_unlock_requests INT DEFAULT 0;
  DECLARE v_del_lock_events INT DEFAULT 0;
  DECLARE v_del_property_status_logs INT DEFAULT 0;
  DECLARE v_upd_business_user INT DEFAULT 0;
  DECLARE v_upd_business_user_openid INT DEFAULT 0;
  DECLARE v_upd_business_user_unionid INT DEFAULT 0;

  -- 出错自动回滚，并把原始错误抛出
  DECLARE EXIT HANDLER FOR SQLEXCEPTION
  BEGIN
    ROLLBACK;
    RESIGNAL;
  END;

  SET v_db = DATABASE();

  IF p_user_id IS NULL OR p_user_id <= 0 THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '参数错误：p_user_id 必须为正整数';
  END IF;

  -- 必须存在 business_user 表
  IF NOT EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = v_db AND table_name = 'business_user'
  ) THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '当前数据库不存在 business_user 表，无法执行';
  END IF;

  -- 用户必须存在（避免误操作）
  SELECT COUNT(*) INTO v_user_exists
  FROM `business_user`
  WHERE `id` = p_user_id;

  IF v_user_exists = 0 THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'business_user 中不存在该用户ID，已终止执行';
  END IF;

  START TRANSACTION;

  -- 1) 解除业务资产引用（不删除资产）
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = v_db AND table_name = 'business_properties' AND column_name = 'agent_id'
  ) THEN
    UPDATE `business_properties`
    SET `agent_id` = 0
    WHERE `agent_id` = p_user_id;
    SET v_upd_properties_agent_id = ROW_COUNT();
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = v_db AND table_name = 'business_property_locks' AND column_name = 'bind_by_user_id'
  ) THEN
    UPDATE `business_property_locks`
    SET `bind_by_user_id` = 0
    WHERE `bind_by_user_id` = p_user_id;
    SET v_upd_property_locks_bind_by_user_id = ROW_COUNT();
  END IF;

  -- 2) 删除个人关联数据
  IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = v_db AND table_name = 'business_favorites'
  ) THEN
    DELETE FROM `business_favorites` WHERE `user_id` = p_user_id;
    SET v_del_favorites = ROW_COUNT();
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = v_db AND table_name = 'business_user_activity_logs'
  ) THEN
    DELETE FROM `business_user_activity_logs` WHERE `user_id` = p_user_id;
    SET v_del_user_activity_logs = ROW_COUNT();
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = v_db AND table_name = 'business_unlock_requests'
  ) THEN
    DELETE FROM `business_unlock_requests` WHERE `user_id` = p_user_id;
    SET v_del_unlock_requests = ROW_COUNT();
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = v_db AND table_name = 'business_lock_events'
  ) THEN
    DELETE FROM `business_lock_events` WHERE `user_id` = p_user_id;
    SET v_del_lock_events = ROW_COUNT();
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.tables
    WHERE table_schema = v_db AND table_name = 'business_property_status_logs'
  ) THEN
    DELETE FROM `business_property_status_logs` WHERE `user_id` = p_user_id;
    SET v_del_property_status_logs = ROW_COUNT();
  END IF;

  -- 3) business_user：脱敏 + 软删除（与 gform 的 deletetime 机制保持一致）
  UPDATE `business_user`
  SET
    `status` = 1,
    `deletetime` = NOW(),
    `updatetime` = NOW(),
    `username` = CONCAT('deleted_', `id`),
    `name` = '',
    `nickname` = '',
    `remark` = CONCAT(IFNULL(`remark`, ''), '[已删除]'),
    `password` = '',
    `salt` = '',
    `email` = '',
    `mobile` = '',
    `avatar` = ''
  WHERE `id` = p_user_id;
  SET v_upd_business_user = ROW_COUNT();

  -- 可选字段：如果存在则一并清空
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = v_db AND table_name = 'business_user' AND column_name = 'openid'
  ) THEN
    UPDATE `business_user` SET `openid` = '' WHERE `id` = p_user_id;
    SET v_upd_business_user_openid = ROW_COUNT();
  END IF;

  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = v_db AND table_name = 'business_user' AND column_name = 'unionid'
  ) THEN
    UPDATE `business_user` SET `unionid` = '' WHERE `id` = p_user_id;
    SET v_upd_business_user_unionid = ROW_COUNT();
  END IF;

  COMMIT;

  -- 执行结果汇总
  SELECT
    p_user_id AS `user_id`,
    v_upd_properties_agent_id AS `updated_business_properties_agent_id`,
    v_upd_property_locks_bind_by_user_id AS `updated_business_property_locks_bind_by_user_id`,
    v_del_favorites AS `deleted_business_favorites`,
    v_del_user_activity_logs AS `deleted_business_user_activity_logs`,
    v_del_unlock_requests AS `deleted_business_unlock_requests`,
    v_del_lock_events AS `deleted_business_lock_events`,
    v_del_property_status_logs AS `deleted_business_property_status_logs`,
    v_upd_business_user AS `updated_business_user_masked_and_soft_deleted`,
    v_upd_business_user_openid AS `updated_business_user_openid`,
    v_upd_business_user_unionid AS `updated_business_user_unionid`;
END$$
DELIMITER ;
