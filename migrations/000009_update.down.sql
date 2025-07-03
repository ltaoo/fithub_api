-- 1. 禁用外键检查（避免删除旧表时触发外键约束错误）
PRAGMA foreign_keys = OFF;

-- 3. 创建新表（复制原表结构，修改`weight`字段类型）
CREATE TABLE WORKOUT_ACTION_HISTORY_OLD (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
    action_id INTEGER NOT NULL DEFAULT 0, --做的动作id
    d INTEGER NOT NULL DEFAULT 0, --软删除
    reps INTEGER NOT NULL DEFAULT 0, --实际数量
    reps_unit TEXT(255) NOT NULL DEFAULT '次', --单位
    weight INTEGER NOT NULL DEFAULT 0, --实际阻力;12kg 24lbs
    weight_unit TEXT(255) NOT NULL DEFAULT '', --实际阻力单位;kg 24lbs
    remark TEXT(255) NOT NULL DEFAULT '', --教练评价、备注。比如动作不标准之类
    feedback TEXT(255) NOT NULL DEFAULT '', --学员反馈，自感疲劳度，动作感受，等等
    extra_medias TEXT(1000) NOT NULL DEFAULT '{}', --拍照、录像记录;存 字符串列表
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, --创建时间;CURRENT_TIMESTAMP
    workout_day_id INTEGER NOT NULL DEFAULT 0, --训练日 id
    student_id INTEGER NOT NULL DEFAULT 0
);

-- 4. 复制原表数据到新表（需指定字段名，避免顺序问题）
INSERT INTO WORKOUT_ACTION_HISTORY_OLD (id, d, action_id, reps, reps_unit, weight, weight_unit, remark, feedback, extra_medias, created_at, workout_day_id, student_id)
SELECT id, d, action_id, reps, reps_unit, weight, weight_unit, remark, feedback, extra_medias, created_at, workout_day_id, student_id
FROM WORKOUT_ACTION_HISTORY;

-- 5. 删除旧表
DROP TABLE WORKOUT_ACTION_HISTORY;

-- 6. 重命名新表为原表名（保持应用兼容性）
ALTER TABLE WORKOUT_ACTION_HISTORY_OLD RENAME TO WORKOUT_ACTION_HISTORY;

-- 8. 恢复外键检查（若原表有外键，必须开启）
PRAGMA foreign_keys = ON;

DROP TABLE IF EXISTS USER_FAVORITE;
DROP TABLE IF EXISTS USER_FAVORITE_FOLDER;


ALTER TABLE WORKOUT_PLAN DROP COLUMN type;
ALTER TABLE WORKOUT_DAY DROP COLUMN duration;
ALTER TABLE WORKOUT_DAY DROP COLUMN total_volume;
ALTER TABLE WORKOUT_DAY DROP COLUMN type;
ALTER TABLE WORKOUT_DAY DROP COLUMN title;
ALTER TABLE WORKOUT_ACTION_HISTORY DROP COLUMN step_uid;
ALTER TABLE WORKOUT_ACTION_HISTORY DROP COLUMN set_uid;
ALTER TABLE WORKOUT_ACTION_HISTORY DROP COLUMN act_uid;
