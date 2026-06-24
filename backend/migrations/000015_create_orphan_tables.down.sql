DROP TABLE IF EXISTS `bpbl`;
DROP TABLE IF EXISTS `gtentire`;
DROP TABLE IF EXISTS `rtl-tl1`;
DROP TABLE IF EXISTS `recipe_material`;

DELETE FROM resource_group_items WHERE resource_name IN ('bpbl', 'gtentire', 'rtl-tl1', 'recipe_material');
