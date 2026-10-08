-- SQLite fixed-ID equipment photo import. Copy media/equipment into dist/h5/media/equipment first.
-- Sources and checksums: equipment_sources.json. Images retain their original rights.
BEGIN IMMEDIATE;
CREATE TEMP TABLE fithub_equipment_photo_guard (valid INTEGER NOT NULL CHECK(valid = 1));
-- ID 1: 哑铃 ; https://shop.lifefitness.com/products/hammer-strength-4-sided-rubber-dumbbells
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=1 AND NOT ((lower(trim(name)) = lower('dumbbell') OR zh_name = '哑铃') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 1,'dumbbell','哑铃','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=1);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/01-equipment.jpg"]')) WHERE id=1;
-- ID 2: 杠铃 ; https://bellsofsteel.us/products/olympic-weightlifting-barbell
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=2 AND NOT ((lower(trim(name)) = lower('barbell') OR zh_name = '杠铃') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 2,'barbell','杠铃','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=2);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/02-equipment.jpg"]')) WHERE id=2;
-- ID 3: 跑步机 ; https://shop.lifefitness.com/products/atmos-treadmill-home-edition
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=3 AND NOT ((lower(trim(name)) = lower('treadmill') OR zh_name = '跑步机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 3,'treadmill','跑步机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=3);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/03-equipment.png"]')) WHERE id=3;
-- ID 4: 椭圆机 ; https://shop.lifefitness.com/products/atmos-elliptical-home-edition
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=4 AND NOT ((lower(trim(name)) = lower('elliptical') OR zh_name = '椭圆机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 4,'elliptical','椭圆机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=4);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/04-equipment.png"]')) WHERE id=4;
-- ID 5: 仰卧板 ; https://shop.lifefitness.com/products/life-fitness-abdominal-bench
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=5 AND NOT ((lower(trim(name)) = lower('ab_machine') OR zh_name = '仰卧板') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 5,'ab_machine','仰卧板','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=5);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/05-equipment.jpg"]')) WHERE id=5;
-- ID 6: 坐姿推胸器 ; https://shop.lifefitness.com/products/insignia-series-chest-press
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=6 AND NOT ((lower(trim(name)) = lower('chest_press_machine') OR zh_name = '坐姿推胸器') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 6,'chest_press_machine','坐姿推胸器','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=6);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/06-equipment.jpg"]')) WHERE id=6;
-- ID 7: 蝴蝶机 ; https://shop.lifefitness.com/products/insignia-series-pectoral-fly-rear-deltoid
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=7 AND NOT ((lower(trim(name)) = lower('butterfly_machine') OR zh_name = '蝴蝶机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 7,'butterfly_machine','蝴蝶机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=7);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/07-equipment.jpg"]')) WHERE id=7;
-- ID 8: 坐姿划船器 ; https://shop.lifefitness.com/products/insignia-series-row
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=8 AND NOT ((lower(trim(name)) = lower('row_machine') OR zh_name = '坐姿划船器') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 8,'row_machine','坐姿划船器','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=8);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/08-equipment.jpg"]')) WHERE id=8;
-- ID 9: 腿部推蹬机 ; https://shop.lifefitness.com/products/insignia-series-leg-press
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=9 AND NOT ((lower(trim(name)) = lower('leg_extension_machine') OR zh_name = '腿部推蹬机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 9,'leg_extension_machine','腿部推蹬机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=9);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/09-equipment.jpg"]')) WHERE id=9;
-- ID 10: 腿部伸展机 ; https://shop.lifefitness.com/products/insignia-series-leg-extension
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=10 AND NOT ((lower(trim(name)) = lower('leg_curl_machine') OR zh_name = '腿部伸展机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 10,'leg_curl_machine','腿部伸展机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=10);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/10-equipment.jpg"]')) WHERE id=10;
-- ID 11: 腿部弯举机 ; https://shop.lifefitness.com/products/insignia-series-seated-leg-curl
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=11 AND NOT ((lower(trim(name)) = lower('leg_curl_machine') OR zh_name = '腿部弯举机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 11,'leg_curl_machine','腿部弯举机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=11);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/11-equipment.jpg"]')) WHERE id=11;
-- ID 12: 史密斯机 ; https://shop.lifefitness.com/products/life-fitness-smith-machine
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=12 AND NOT ((lower(trim(name)) = lower('smith_machine') OR zh_name = '史密斯机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 12,'smith_machine','史密斯机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=12);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/12-equipment.jpg"]')) WHERE id=12;
-- ID 13: 龙门架 ; https://shop.lifefitness.com/products/life-fitness-dual-adjustable-pulley
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=13 AND NOT ((lower(trim(name)) = lower('cable_machine') OR zh_name = '龙门架') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 13,'cable_machine','龙门架','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=13);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/13-equipment.jpg"]')) WHERE id=13;
-- ID 14: 引体向上器 ; https://shop.lifefitness.com/products/hammer-strength-select-assist-dip-chin
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=14 AND NOT ((lower(trim(name)) = lower('pull_up_machine') OR zh_name = '引体向上器') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 14,'pull_up_machine','引体向上器','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=14);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/14-equipment.jpg"]')) WHERE id=14;
-- ID 15: 罗马椅 ; https://shop.lifefitness.com/products/signature-series-back-extension
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=15 AND NOT ((lower(trim(name)) = lower('roman_chair') OR zh_name = '罗马椅') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 15,'roman_chair','罗马椅','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=15);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/15-equipment.jpg"]')) WHERE id=15;
-- ID 16: 壶铃 ; https://shop.lifefitness.com/products/life-fitness-kettlebells
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=16 AND NOT ((lower(trim(name)) = lower('kettlebell') OR zh_name = '壶铃') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 16,'kettlebell','壶铃','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=16);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/16-equipment.jpg"]')) WHERE id=16;
-- ID 17: 弹力带 ; https://repfitness.com/products/resistance-bands-2-0
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=17 AND NOT ((lower(trim(name)) = lower('resistance_band') OR zh_name = '弹力带') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 17,'resistance_band','弹力带','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=17);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/17-equipment.png"]')) WHERE id=17;
-- ID 18: 泡沫轴 ; https://repfitness.com/products/high-density-foam-roller
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=18 AND NOT ((lower(trim(name)) = lower('foam_roller') OR zh_name = '泡沫轴') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 18,'foam_roller','泡沫轴','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=18);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/18-equipment.jpg"]')) WHERE id=18;
-- ID 19: 瑜伽球 ; https://repfitness.com/products/physio-ball
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=19 AND NOT ((lower(trim(name)) = lower('yoga_ball') OR zh_name = '瑜伽球') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 19,'yoga_ball','瑜伽球','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=19);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/19-equipment.jpg"]')) WHERE id=19;
-- ID 20: 健身战绳 ; https://bellsofsteel.us/products/battle-rope
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=20 AND NOT ((lower(trim(name)) = lower('battle_rope') OR zh_name = '健身战绳') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 20,'battle_rope','健身战绳','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=20);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/20-equipment.jpg"]')) WHERE id=20;
-- ID 21: TRX悬挂训练带 ; https://repfitness.com/products/trx-suspension-trainers
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=21 AND NOT ((lower(trim(name)) = lower('trx_suspension_trainer') OR zh_name = 'TRX悬挂训练带') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 21,'trx_suspension_trainer','TRX悬挂训练带','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=21);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/21-equipment.jpg"]')) WHERE id=21;
-- ID 22: 瑜伽垫 ; https://repfitness.com/products/premium-yoga-mat
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=22 AND NOT ((lower(trim(name)) = lower('yoga_mat') OR zh_name = '瑜伽垫') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 22,'yoga_mat','瑜伽垫','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=22);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/22-equipment.jpg"]')) WHERE id=22;
-- ID 23: 六角杠铃 ; https://repfitness.com/products/commercial-open-trap-bar
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=23 AND NOT ((lower(trim(name)) = lower('Hex Bar') OR zh_name = '六角杠铃') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 23,'Hex Bar','六角杠铃','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=23);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/23-equipment.jpg"]')) WHERE id=23;
-- ID 24: 高位下拉机 ; https://shop.lifefitness.com/products/axiom-series-pulldown
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=24 AND NOT ((lower(trim(name)) = lower('Lat Pulldown Machine') OR zh_name = '高位下拉机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 24,'Lat Pulldown Machine','高位下拉机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=24);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/24-equipment.jpg"]')) WHERE id=24;
-- ID 25: 髋外展/内收机 ; https://shop.lifefitness.com/products/axiom-series-hip-abductor-adductor
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=25 AND NOT ((lower(trim(name)) = lower('Hip Abduction/Adduction Machine') OR zh_name = '髋外展/内收机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 25,'Hip Abduction/Adduction Machine','髋外展/内收机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=25);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/25-equipment.jpg"]')) WHERE id=25;
-- ID 26: 平衡垫 ; https://shop.lifefitness.com/products/bosu®-pro-balance-trainer
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=26 AND NOT ((lower(trim(name)) = lower('Balance Pad') OR zh_name = '平衡垫') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 26,'Balance Pad','平衡垫','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=26);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/26-equipment.jpg"]')) WHERE id=26;
-- ID 27: 负重雪橇 ; https://shop.lifefitness.com/products/push-pull-weight-sled
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=27 AND NOT ((lower(trim(name)) = lower('Weight Sled') OR zh_name = '负重雪橇') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 27,'Weight Sled','负重雪橇','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=27);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/27-equipment.jpg"]')) WHERE id=27;
-- ID 28: 跳箱 ; https://shop.lifefitness.com/products/hammer-strength-3-in-1-soft-plyo-box
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=28 AND NOT ((lower(trim(name)) = lower('Plyo Box') OR zh_name = '跳箱') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 28,'Plyo Box','跳箱','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=28);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/28-equipment.jpg"]')) WHERE id=28;
-- ID 29: 沙袋 ; https://bellsofsteel.us/products/sandbag-sets
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=29 AND NOT ((lower(trim(name)) = lower('Sandbag') OR zh_name = '沙袋') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 29,'Sandbag','沙袋','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=29);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/29-equipment.jpg"]')) WHERE id=29;
-- ID 30: 攀爬绳 ; https://www.titaniumstrength.fr/corde-d-escalade-climbing-rope-8-m.html
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=30 AND NOT ((lower(trim(name)) = lower('Climbing Rope') OR zh_name = '攀爬绳') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 30,'Climbing Rope','攀爬绳','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=30);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/30-equipment.jpg"]')) WHERE id=30;
-- ID 31: 滑雪机 ; https://repfitness.com/products/concept2-ski-erg
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=31 AND NOT ((lower(trim(name)) = lower('Ski Erg') OR zh_name = '滑雪机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 31,'Ski Erg','滑雪机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=31);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/31-equipment.jpg"]')) WHERE id=31;
-- ID 32: 跳绳 ; https://bellsofsteel.us/products/jump-rope
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=32 AND NOT ((lower(trim(name)) = lower('Jump Rope') OR zh_name = '跳绳') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 32,'Jump Rope','跳绳','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=32);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/32-equipment.jpg"]')) WHERE id=32;
-- ID 33: 台阶机 ; https://repfitness.com/products/stepr-stair-climber
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=33 AND NOT ((lower(trim(name)) = lower('StairMill') OR zh_name = '台阶机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 33,'StairMill','台阶机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=33);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/33-equipment.jpg"]')) WHERE id=33;
-- ID 34: 攀爬机 ; https://precisionperformances.com/products/versaclimber-lx-model
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=34 AND NOT ((lower(trim(name)) = lower('VersaClimber') OR zh_name = '攀爬机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 34,'VersaClimber','攀爬机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=34);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/34-equipment.jpg"]')) WHERE id=34;
-- ID 35: 旋转核心机 ; https://shop.lifefitness.com/products/insignia-series-torso-rotation
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=35 AND NOT ((lower(trim(name)) = lower('Torso Rotation Machine') OR zh_name = '旋转核心机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 35,'Torso Rotation Machine','旋转核心机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=35);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/35-equipment.jpg"]')) WHERE id=35;
-- ID 36: 卷腹机 ; https://shop.lifefitness.com/products/hammer-strength-select-abdominal-crunch
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=36 AND NOT ((lower(trim(name)) = lower('Abdominal Crunch Machine') OR zh_name = '卷腹机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 36,'Abdominal Crunch Machine','卷腹机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=36);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/36-equipment.jpg"]')) WHERE id=36;
-- ID 37: 站姿提踵机 ; https://shop.lifefitness.com/products/hammer-strength-select-standing-calf
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=37 AND NOT ((lower(trim(name)) = lower('Standing Calf Raise Machine') OR zh_name = '站姿提踵机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 37,'Standing Calf Raise Machine','站姿提踵机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=37);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/37-equipment.jpg"]')) WHERE id=37;
-- ID 38: 坐姿提踵机 ; https://shop.lifefitness.com/products/hammer-strength-plate-loaded-seated-calf-raise
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=38 AND NOT ((lower(trim(name)) = lower('Seated Calf Raise Machine') OR zh_name = '坐姿提踵机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 38,'Seated Calf Raise Machine','坐姿提踵机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=38);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/38-equipment.jpg"]')) WHERE id=38;
-- ID 39: 哈克深蹲机 ; https://shop.lifefitness.com/products/hammer-strength-plate-loaded-hack-squat
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=39 AND NOT ((lower(trim(name)) = lower('Hack Squat Machine') OR zh_name = '哈克深蹲机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 39,'Hack Squat Machine','哈克深蹲机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=39);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/39-equipment.jpg"]')) WHERE id=39;
-- ID 40: 侧平举机 ; https://shop.lifefitness.com/products/insignia-series-lateral-raise
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=40 AND NOT ((lower(trim(name)) = lower('Lateral Raise Machine') OR zh_name = '侧平举机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 40,'Lateral Raise Machine','侧平举机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=40);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/40-equipment.jpg"]')) WHERE id=40;
-- ID 41: T 杠划船机 ; https://shop.lifefitness.com/products/hammer-strength-plate-loaded-t-bar-row
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=41 AND NOT ((lower(trim(name)) = lower('T-Bar Row Machine') OR zh_name = 'T 杠划船机') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 41,'T-Bar Row Machine','T 杠划船机','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=41);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/41-equipment.jpg"]')) WHERE id=41;
-- ID 42: 圆木 ; https://titan.fitness/products/rackable-strongman-log-bars
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=42 AND NOT ((lower(trim(name)) = lower('log') OR zh_name = '圆木') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 42,'log','圆木','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=42);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/42-equipment.jpg"]')) WHERE id=42;
-- ID 43: EZ曲杆 ; https://bellsofsteel.us/products/rackable-ez-curl-bar
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=43 AND NOT ((lower(trim(name)) = lower('ez_curl_bar') OR zh_name = 'EZ曲杆') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 43,'ez_curl_bar','EZ曲杆','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=43);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/43-equipment.jpg"]')) WHERE id=43;
-- ID 44: 通用健身器械 ; https://repfitness.com/products/arcadia-functional-trainer
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=44 AND NOT ((lower(trim(name)) = lower('machine') OR zh_name = '通用健身器械') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 44,'machine','通用健身器械','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=44);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/44-equipment.jpg"]')) WHERE id=44;
-- ID 45: 药球 ; https://repfitness.com/products/medicine-ball-2-0
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=45 AND NOT ((lower(trim(name)) = lower('medicine_ball') OR zh_name = '药球') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 45,'medicine_ball','药球','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=45);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/45-equipment.jpg"]')) WHERE id=45;
-- ID 46: 其他器械 ; https://bellsofsteel.us/products/weighted-vest
INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=46 AND NOT ((lower(trim(name)) = lower('other') OR zh_name = '其他器械') AND COALESCE(d,0)=0)) THEN 1 ELSE 0 END;
INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT 46,'other','其他器械','','','{}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id=46);
UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{}' END ELSE '{}' END,'$.pics',json('["/media/equipment/46-equipment.jpg"]')) WHERE id=46;
DROP TABLE fithub_equipment_photo_guard;
COMMIT;
