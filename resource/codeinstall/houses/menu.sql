-- 房源管理菜单/按钮权限初始化（business_auth_rule）
-- 说明：
-- 1) 本脚本可重复执行，不会产生重复菜单/按钮
-- 2) id 使用自增；pid 通过查询父菜单id写入

SET @uid := 1;
SET @now := NOW();

-- 1. 创建“房源管理”菜单（type=1）
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '房源管理', '', 'menu.houses', 6, 1, 0, 'icon-home', '/houses/houses', 'houses', '/houses/houses/index', '', '', NULL,
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 1 AND routepath = '/houses/houses' AND routename = 'houses'
);

SET @houses_menu_id := (
  SELECT id FROM business_auth_rule
  WHERE type = 1 AND routepath = '/houses/houses' AND routename = 'houses'
  ORDER BY id DESC LIMIT 1
);

-- 2. 按钮权限（type=2）
-- 查看列表
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '查看', '', '', 1, 2, @houses_menu_id, '', '', '', '', '', '/business/houses/houses/getList', 'view',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @houses_menu_id AND path = '/business/houses/houses/getList'
);

-- 添加/编辑
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '添加/编辑', '', '', 2, 2, @houses_menu_id, '', '', '', '', '', '/business/houses/houses/save', 'add',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @houses_menu_id AND path = '/business/houses/houses/save'
);

-- 删除
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '删除', '', '', 3, 2, @houses_menu_id, '', '', '', '', '', '/business/houses/houses/del', 'del',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @houses_menu_id AND path = '/business/houses/houses/del'
);

-- 状态/销售状态/排序（统一 upStatus）
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '状态', '', '', 4, 2, @houses_menu_id, '', '', '', '', '', '/business/houses/houses/upStatus', 'upStatus',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @houses_menu_id AND path = '/business/houses/houses/upStatus'
);

-- 详情
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '详情', '', '', 5, 2, @houses_menu_id, '', '', '', '', '', '/business/houses/houses/getContent', 'view',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @houses_menu_id AND path = '/business/houses/houses/getContent'
);

-- ============================
-- 3. 智能锁管理（TTLock）菜单/按钮
-- ============================

-- 3.1 创建“智能锁管理”子菜单（type=1）
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '智能锁管理', '', 'menu.ttlock', 7, 1, @houses_menu_id, 'icon-lock', '/houses/ttlock', 'ttlock', '/houses/ttlock/index', '', '', NULL,
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 1 AND routepath = '/houses/ttlock' AND routename = 'ttlock'
);

SET @ttlock_menu_id := (
  SELECT id FROM business_auth_rule
  WHERE type = 1 AND routepath = '/houses/ttlock' AND routename = 'ttlock'
  ORDER BY id DESC LIMIT 1
);

-- 3.2 按钮权限（type=2）
-- 状态/配置
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '状态', '', '', 1, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/status', 'view',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/status'
);

-- 同步锁列表
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '同步锁', '', '', 2, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/syncLocks', 'sync',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/syncLocks'
);

-- 本地锁列表
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '锁列表', '', '', 3, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/lockList', 'view',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/lockList'
);

-- 绑定房源
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '绑定', '', '', 4, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/bindProperty', 'bind',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/bindProperty'
);

-- 解绑房源
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '解绑', '', '', 5, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/unbindProperty', 'unbind',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/unbindProperty'
);

-- 查询房源绑定锁
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '绑定查询', '', '', 6, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/propertyLock', 'view',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/propertyLock'
);

-- 远程开锁
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '远程开锁', '', '', 7, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/remoteUnlock', 'unlock',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/remoteUnlock'
);

-- 网关下发密码
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '下发密码', '', '', 8, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/addKeyboardPwd', 'passcode',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/addKeyboardPwd'
);

-- 下发电子钥匙
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '下发钥匙', '', '', 9, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/sendEkey', 'ekey',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/sendEkey'
);

-- 云端锁详情
INSERT INTO business_auth_rule
  (uid, title, des, locale, weigh, type, pid, icon, routepath, routename, component, redirect, path, permission,
   status, isext, keepalive, requiresauth, hideinmenu, hidechildreninmenu, activemenu, noaffix, onlypage, createtime)
SELECT
  @uid, '锁详情', '', '', 10, 2, @ttlock_menu_id, '', '', '', '', '', '/business/houses/ttlock/lockDetail', 'view',
  0, 0, 0, 1, 0, 0, 0, 0, 0, @now
WHERE NOT EXISTS (
  SELECT 1 FROM business_auth_rule
  WHERE type = 2 AND pid = @ttlock_menu_id AND path = '/business/houses/ttlock/lockDetail'
);

