DELETE FROM business_properties WHERE title LIKE 'TEST-PROP-20260214-IMPORTSMALL-%';
INSERT INTO business_properties (business_id,title,sale_status,agent_id,status,createtime,updatetime) VALUES (1,'TEST-PROP-20260214-IMPORTSMALL-001','on_sale',112,0,'2026-02-14 12:10:00','2026-02-14 12:10:00');
INSERT INTO business_properties (business_id,title,sale_status,agent_id,status,createtime,updatetime) VALUES (1,'TEST-PROP-20260214-IMPORTSMALL-002','in_sale',112,0,'2026-02-14 12:11:00','2026-02-14 12:11:00');
