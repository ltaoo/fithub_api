
--up主
CREATE TABLE IF NOT EXISTS INFLUENCER(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL DEFAULT '', --
  nickname TEXT NOT NULL DEFAULT '', --昵称
  alias TEXT NOT NULL DEFAULT '', --别名
  avatar_url TEXT NOT NULL DEFAULT '', --头像地址
  bio TEXT NOT NULL DEFAULT '', --简介
  status INTEGER NOT NULL DEFAULT 0, --状态 0
  d INTEGER,
  coach_id INTEGER NOT NULL DEFAULT 0, --教练id
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

--平台信息
CREATE TABLE IF NOT EXISTS INFLUENCER_PLATFORM(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL DEFAULT '', --平台名称
  logo_url TEXT NOT NULL DEFAULT '', --平台logo
  homepage_url TEXT NOT NULL DEFAULT '', --平台官网
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

--up主平台信息
CREATE TABLE IF NOT EXISTS INFLUENCER_ACCOUNT_IN_PLATFORM(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  nickname TEXT NOT NULL DEFAULT '', -- 在平台的名称
  avatar_url TEXT NOT NULL DEFAULT '', -- 在平台的头像
  handle TEXT NOT NULL DEFAULT '', -- 帐号标志如 @xxx
  account_url TEXT NOT NULL DEFAULT '', -- 在平台的个人主页链接
  followers_count INTEGER NOT NULL DEFAULT 0, -- 粉丝数
  status INTEGER NOT NULL DEFAULT 0,
  d INTEGER NOT NULL DEFAULT 0, --软删除
  coach_id INTEGER NOT NULL DEFAULT 0, --up主
  platform_id INTEGER NOT NULL DEFAULT 0, --平台
  updated_at DATETIME, --更新时间
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

--up主发布的内容
CREATE TABLE IF NOT EXISTS COACH_CONTENT(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  content_type INTEGER NOT NULL DEFAULT 0, --内容类型 1视频 2短视频 3文章 4图片 5纯文本
  status INTEGER NOT NULL DEFAULT 0, --1审核通过 2
  publish INTEGER NOT NULL DEFAULT 0, --1公开 2私有
  title TEXT NOT NULL DEFAULT '', --内容标题
  description TEXT NOT NULL DEFAULT '', --概述
  like_count INTEGER NOT NULL DEFAULT 0, --点赞数
  content_url TEXT NOT NULL DEFAULT '', --内容地址
  cover_image_url TEXT NOT NULL DEFAULT '', --封面图url
  video_key TEXT NOT NULL DEFAULT '', --视频被手动上传到资源库了，视频才会有
  video_overview TEXT NOT NULL DEFAULT '{}', --视频内容概述
  image_keys TEXT NOT NULL DEFAULT '', --图片被手动上传到资源库了，图片才会有
  image_overview TEXT NOT NULL DEFAULT '{}', --图片内容概述
  content TEXT NOT NULL DEFAULT '', --文本内容
  published_at DATETIME,
  d INTEGER NOT NULL DEFAULT 0, --软删除
  coach_id INTEGER NOT NULL DEFAULT 0, --up主
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

-- 健身动作和up主内容的多对多关联
CREATE TABLE IF NOT EXISTS COACH_CONTENT_WITH_WORKOUT_ACTION(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  d INTEGER NOT NULL DEFAULT 0, --软删除
  sort_idx INTEGER NOT NULL DEFAULT 0, --排序
  start_point INTEGER NOT NULL DEFAULT 0, --内容空降点 视频和短视频才有
  coach_content_id INTEGER NOT NULL DEFAULT 0, --内容id
  workout_action_id INTEGER NOT NULL DEFAULT 0, --动作id
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);
-- 训练计划和up主内容的多对多关联
CREATE TABLE IF NOT EXISTS COACH_CONTENT_WITH_WORKOUT_PLAN(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  d INTEGER NOT NULL DEFAULT 0, --软删除
  sort_idx INTEGER NOT NULL DEFAULT 0, --排序
  coach_content_id INTEGER NOT NULL DEFAULT 0, --内容id
  workout_plan_id INTEGER NOT NULL DEFAULT 0, --训练计划id
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

-- 健身周报
CREATE TABLE IF NOT EXISTS WORKOUT_OUT_WEEKLY_STATS(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  start_date_text TEXT NOT NULL DEFAULT '', --周一 日期
  end_date_text TEXT NOT NULL DEFAULT '', --周日 日期
  title TEXT NOT NULL DEFAULT '', --周报标题 六月第三周
  workout_day_count INTEGER NOT NULL DEFAULT 0, -- 周共运动多少天
  workout_times INTEGER NOT NULL DEFAULT 0, -- 周共运动多少次
  volume_total INTEGER NOT NULL DEFAULT 0, -- 周总容量 kg

  workout_day_count_diff INTEGER NOT NULL DEFAULT 0, -- 周共运动多少天 和上周的差异
  workout_times_diff INTEGER NOT NULL DEFAULT 0, -- 周共运动多少次 和上周的差异
  volume_total_diff INTEGER NOT NULL DEFAULT 0, -- 周总容量 kg 和上周的差异

  workout_action_max_weight INTEGER NOT NULL DEFAULT 0, -- 最大重量的动作
  workout_action_max_volume INTEGER NOT NULL DEFAULT 0, -- 最大容量的动作
  workout_action_max_reps INTEGER NOT NULL DEFAULT 0, -- 最大次数的动作

  workout_action_max_weight_added INTEGER NOT NULL DEFAULT 0, -- 重量进步最大的动作

  chest_set_count INTEGER NOT NULL DEFAULT 0, --胸 训练总组数
  legs_set_count INTEGER NOT NULL DEFAULT 0, --腿 训练总组数
  hip_set_count INTEGER NOT NULL DEFAULT 0, --臀 训练总组数
  back_set_count INTEGER NOT NULL DEFAULT 0, --臀 训练总组数

  coach_id INTEGER NOT NULL DEFAULT 0, --
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

CREATE TABLE IF NOT EXISTS NOTIFICATION(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  status INTEGER NOT NULL DEFAULT 0, --通知状态 1未读 2已读 
  d INTEGER, --软删除
  type INTEGER NOT NULL DEFAULT 0, --通知内容类型 1文本
  content TEXT NOT NULL DEFAULT '', --通知内容
  coach_id INTEGER NOT NULL DEFAULT 0, --通知人
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

-- 关注表
CREATE TABLE IF NOT EXISTS COACH_FOLLOW(
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  status INTEGER NOT NULL DEFAULT 0, -- 1关注 2取消关注
  follower_id INTEGER NOT NULL DEFAULT 0, --谁关注了我
  following_id INTEGER NOT NULL DEFAULT 0, --我关注了谁
  updated_at DATETIME, -- 更新时间
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);
