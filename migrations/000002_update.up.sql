ALTER TABLE WORKOUT_PLAN ADD COLUMN cover_url TEXT;
ALTER TABLE MUSCLE ADD COLUMN medias TEXT;
ALTER TABLE EQUIPMENT ADD COLUMN tags TEXT;

CREATE TABLE IF NOT EXISTS WORKOUT_PLAN_SET (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,  -- 自增主键
    title TEXT NOT NULL DEFAULT '',                   -- 标题
    overview TEXT NOT NULL DEFAULT '',                   -- 概要
    icon_url TEXT NOT NULL DEFAULT '',                   -- 额外标签
    idx INTEGER NOT NULL DEFAULT 0,                   -- 排序
    details TEXT NOT NULL DEFAULT '{}',                          -- 集合内容
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP                    -- 创建时间
);

CREATE TABLE IF NOT EXISTS SUBSCRIPTION_PLAN (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 订阅计划 id
  name TEXT NOT NULL DEFAULT '',             -- 订阅计划名称（如“年度会员”）
  details TEXT NOT NULL DEFAULT '{}',             -- 计划的详细信息
  unit_price INTEGER NOT NULL DEFAULT 0,               -- 单位价格
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP  -- 创建时间
);

CREATE TABLE IF NOT EXISTS COACH_PERMISSION (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 权限 id
  name TEXT NOT NULL DEFAULT '', -- 权限名称
  details TEXT NOT NULL DEFAULT '{}', -- 详细说明
  sort_idx INTEGER NOT NULL DEFAULT 0 -- 排序
);

-- 订阅和权限多对多关联表
CREATE TABLE IF NOT EXISTS SUBSCRIPTION_PLAN_COACH_PERMISSION (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, -- id
  subscription_plan_id INTEGER NOT NULL DEFAULT 0, -- 订阅计划 id
  permission_id INTEGER NOT NULL DEFAULT 0, -- 权限 id
  checked INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS DISCOUNT_POLICY (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 折扣id
  name TEXT NOT NULL DEFAULT '',             -- 折扣名称（如“？”）
  rate INTEGER NOT NULL DEFAULT 100,                 -- 折扣比例
  count_require INTEGER NOT NULL DEFAULT 0,               -- 
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- 创建时间
);

-- 订阅计划和计划折扣政策的多对多关联表
CREATE TABLE IF NOT EXISTS SUBSCRIPTION_PLAN_DISCOUNT_POLICY (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, -- id
  subscription_plan_id INTEGER NOT NULL DEFAULT 0, -- 订阅计划 id
  discount_policy_id INTEGER NOT NULL DEFAULT 0, -- 折扣 id
  enabled INTEGER NOT NULL DEFAULT 0 -- 该条折扣政策是否生效 0不生效 1生效
);

CREATE TABLE IF NOT EXISTS SUBSCRIPTION_ORDER (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 订阅订单 id
  amount INTEGER NOT NULL DEFAULT 0, -- 支付总价金额 使用货币的最小单位，如使用 230分，避免出现 2.30元
  discount INTEGER NOT NULL DEFAULT 0, -- 折扣金额
  discount_details TEXT NOT NULL DEFAULT '{}', -- 折扣详细说明
  subscription_plan_id INTEGER NOT NULL DEFAULT 0, -- 订阅计划
  invoice_id INTEGER NOT NULL DEFAULT 0, -- 账单id
  coach_id INTEGER NOT NULL DEFAULT 0 -- 购买者id
);

CREATE TABLE IF NOT EXISTS INVOICE (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 账单id
  status INTEGER NOT NULL DEFAULT 1, -- 1待支付 2支付完成 3取消支付
  amount INTEGER NOT NULL DEFAULT 0, -- 实际需要支付的金额（已减去折扣金额）
  currency_unit INTEGER NOT NULL DEFAULT 1, -- 货币类型 1人民币(fen) 2美元(cent) 3欧元(cent) 4英镑(pence) 5日元(yen) 6卢比(rupee)
  order_type INTEGER NOT NULL DEFAULT 0, -- 订单来源类型 1订阅 2实物
  order_id INTEGER NOT NULL DEFAULT 0, -- 关联的订单id
  coach_id INTEGER NOT NULL DEFAULT 0, -- 购买者id
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
  paid_at DATETIME,  -- 账单支付时间
  canceled_at DATETIME,  -- 账单取消时间
  cancel_reason TEXT NOT NULL DEFAULT '' -- 账单取消原因
);

CREATE TABLE IF NOT EXISTS SUBSCRIPTION (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 订阅id
  step INTEGER NOT NULL DEFAULT 1, -- 订阅状态 1待生效 2生效中 3已过期 4暂停中 5invalid
  count INTEGER NOT NULL DEFAULT 0, -- 订阅天数
  expect_expired_at DATETIME, -- 预期订阅到期时间
  coach_id INTEGER NOT NULL DEFAULT 0, -- 订阅所属人
  subscription_plan_id INTEGER NOT NULL DEFAULT 0, -- 关联订阅
  active_at DATETIME, -- 生效时间
  expired_at DATETIME, -- 过期时间
  paused_at DATETIME, -- 暂停时间
  pause_count INTEGER NOT NULL DEFAULT 0, -- 暂停次数
  invalid_at DATETIME, -- 作废时间
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

CREATE TABLE IF NOT EXISTS QUIZ (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 题目id
  content TEXT NOT NULL DEFAULT '', -- 题目内容
  overview TEXT NOT NULL DEFAULT '', -- 概述
  medias TEXT NOT NULL DEFAULT '', --题目附带的图片
  type INTEGER NOT NULL DEFAULT 1, -- 1单选 2多选 3判断 4填空 5简答
  difficulty INTEGER NOT NULL DEFAULT 1, -- 难度评分
  tags TEXT NOT NULL DEFAULT '', -- 标签
  analysis TEXT NOT NULL DEFAULT '', -- 题目解析
  choices TEXT NOT NULL DEFAULT '{}', -- 题目选项
  answer TEXT NOT NULL DEFAULT '{}', -- 答案
  creator_id INTEGER NOT NULL DEFAULT 0, -- 创建者id
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

CREATE TABLE IF NOT EXISTS PAPER (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 试卷id
  name TEXT NOT NULL DEFAULT '', --试卷名称
  overview TEXT NOT NULL DEFAULT '', -- 概述
  tags TEXT NOT NULL DEFAULT '', --试卷标签
  duration INTEGER NOT NULL DEFAULT 0, -- 试卷答题时长，单位 分钟
  pass_score INTEGER NOT NULL DEFAULT 0, --通过分
  quiz_count INTEGER NOT NULL DEFAULT 0, --总题数
  creator_id INTEGER NOT NULL DEFAULT 0, --创建人id
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 试卷创建时间
);

CREATE TABLE IF NOT EXISTS PAPER_QUIZ (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     --
  paper_id INTEGER NOT NULL DEFAULT 0, -- 试卷id
  quiz_id INTEGER NOT NULL DEFAULT 0, -- 题目id
  score INTEGER NOT NULL DEFAULT 1, -- 题目分数
  sort_idx INTEGER NOT NULL DEFAULT 0, -- 排序
  visible INTEGER NOT NULL DEFAULT 1 -- 是否可见
);

-- 开始考试就要创建一条该记录
CREATE TABLE IF NOT EXISTS EXAM (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 试卷结果id
  status INTEGER NOT NULL DEFAULT 1, --考试状态 1待开始 2进行中 3已完成 4手动放弃
  score INTEGER NOT NULL DEFAULT 0,  --试卷分数
  correct_rate INTEGER NOT NULL DEFAULT 0, --正确率
  pass INTEGER NOT NULL DEFAULT 0, --是否通过 0否 1是
  cur_quiz_id INTEGER NOT NULL DEFAULT 0, --如果进行中的考试 当前的进度
  paper_id INTEGER NOT NULL DEFAULT 0, --试卷id
  student_id INTEGER NOT NULL DEFAULT 0, -- 答题人id
  started_at DATETIME, --开始时间
  completed_at  DATETIME, --完成时间 包含交卷、超时自动提交 
  give_up_at  DATETIME, --放弃时间
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 创建时间
);

-- 答题记录 会有两种情况，一、题目答完直接给结果；二、试卷中答题，全部答完一起结算给结果
CREATE TABLE IF NOT EXISTS QUIZ_ANSWER (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     -- 试卷结果id
  quiz_id INTEGER NOT NULL DEFAULT 0,  --题目id
  exam_id INTEGER NOT NULL DEFAULT 0, --考试id 上面的情况二才会有该值
  paper_id INTEGER NOT NULL DEFAULT 0, --试卷id 上面的情况二才会有该值
  status INTEGER NOT NULL DEFAULT 0, --题目结果 0 1正确 2失败 3跳过
  answer TEXT NOT NULL DEFAULT '{}', -- 答题内容，JSON 根据题目类型有不同结构
  score INTEGER NOT NULL DEFAULT 0, --得分
  student_id INTEGER NOT NULL DEFAULT 0, --答题人id 用来获取答题列表
  updated_at DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 答题时间
);
