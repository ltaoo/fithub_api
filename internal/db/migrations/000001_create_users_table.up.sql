
    DROP TABLE IF EXISTS LLM_AGENT;
    CREATE TABLE LLM_AGENT(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        name TEXT(255) NOT NULL, --名称
        desc TEXT(255), --描述信息
        avatar_uri TEXT(255), --头像地址
        prompt TEXT(255) NOT NULL, --提示词
        tags TEXT(255), --标签
        agent_type INTEGER NOT NULL, --类型
        llm_config TEXT(1000) NOT NULL DEFAULT '{}', --LLM 配置
        llm_provider_id TEXT(255) NOT NULL, --使用厂商ID
        llm_model_id TEXT(255) NOT NULL, --使用厂商指定模型id
        builtin INTEGER NOT NULL, --是否系统内置，不能删除
        config TEXT(1000) NOT NULL DEFAULT '{}', --agent 的配置
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP --
    
    ); --Agent
    
  
    DROP TABLE IF EXISTS LLM_PROVIDER_MODEL_PROFILE;
    CREATE TABLE LLM_PROVIDER_MODEL_PROFILE(
    
        llm_provider_model_id TEXT(255) NOT NULL, --关联 model id
        desc TEXT(255) NOT NULL DEFAULT '', --描述
        tags TEXT(255) NOT NULL DEFAULT '', --标签
        model_type INTEGER NOT NULL DEFAULT 1, --模型类型
        version TEXT(255) NOT NULL DEFAULT '', --模型版本，也算一个额外标签吧
        architecture_type TEXT(255) NOT NULL DEFAULT '', --架构类型
        parameter_count TEXT(255) NOT NULL DEFAULT '', --参数量
        training_data TEXT(255) NOT NULL DEFAULT '', --训练数据
        training_method TEXT(255) NOT NULL DEFAULT '', --训练方式
        evaluation_metrics TEXT(255) NOT NULL DEFAULT '', --评估指标
        usage_restrictions TEXT(255) NOT NULL DEFAULT '', --使用限制
        cost_information TEXT(255) NOT NULL DEFAULT '' --费用信息
    
    ); --Model详情信息
    
  
    DROP TABLE IF EXISTS LLM_PROVIDER_MODEL;
    CREATE TABLE LLM_PROVIDER_MODEL(
    
        id TEXT(255) NOT NULL PRIMARY KEY, --id
        name TEXT(255) NOT NULL, --模型名称
        enabled INTEGER NOT NULL, --是否启用
        llm_provider_id TEXT(255) NOT NULL, --所属厂商id
        builtin INTEGER NOT NULL --是否厂商自带
    
    ); --LLM厂商模型
    
  
    DROP TABLE IF EXISTS LLM_PROVIDER;
    CREATE TABLE LLM_PROVIDER(
    
        id TEXT(255) NOT NULL PRIMARY KEY, --id
        name TEXT(90) NOT NULL, --厂商名称
        logo_uri TEXT(255) NOT NULL, --厂商logo
        api_address TEXT(255) NOT NULL, --请求地址
        configure TEXT(1000) NOT NULL DEFAULT '{}', --允许的配置项
        api_proxy_address TEXT(255), --用户输入的配置项
        api_key TEXT(255), --用户输入的api key
        enabled INTEGER NOT NULL DEFAULT 0 --该厂商是否启用
    
    ); --LLM厂商
    
  
    DROP TABLE IF EXISTS LLM_RESPONSE;
    CREATE TABLE LLM_RESPONSE(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        llm_provider_id TEXT(255) NOT NULL, --该次调用使用的厂商id
        llm_provider_model_id TEXT(255) NOT NULL, --该次调用使用的厂商模型id
        llm_agent_id INTEGER NOT NULL, --该次调用的发起者
        body TEXT(1000) NOT NULL DEFAULT '{}', --该次调用的请求体
        response TEXT(1000), --该次调用的响应体
        error TEXT(1000), --该次调用是否发生错误
        prompt_tokens INTEGER NOT NULL, --输入提示所消耗的 token
        completion_tokens INTEGER NOT NULL, --生成回复内容所消耗的 token
        total_tokens INTEGER NOT NULL, --次对话请求和响应总共消耗
        response_id TEXT(255), --厂商返回的对话完成的唯一标识符
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP --响应时间
    
    ); --调用LLM记录
    
  
    DROP TABLE IF EXISTS CONFIG;
    CREATE TABLE CONFIG(
    
        file_rootpath TEXT(255) NOT NULL DEFAULT '' --笔记保存根路径
    
    ); --应用全局配置
    
  
    DROP TABLE IF EXISTS CHAT_BOX;
    CREATE TABLE CHAT_BOX(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        sender_id INTEGER NOT NULL, --所属 agent id
        chat_session_id INTEGER NOT NULL, --所属对话 id
        payload TEXT(1000) NOT NULL DEFAULT '{}', --具体内容
        box_type INTEGER NOT NULL, --类型
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP --创建时间
    
    ); --对话内容
    
  
    DROP TABLE IF EXISTS CHAT_SESSION;
    CREATE TABLE CHAT_SESSION(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        title TEXT(255) NOT NULL, --对话概述
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP --创建时间
    
    ); --对话
    
  
    DROP TABLE IF EXISTS CHAT_MEMBER;
    CREATE TABLE CHAT_MEMBER(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        llm_agent_id INTEGER NOT NULL, --对应 agent
        chat_session_id INTEGER NOT NULL --所属 session
    
    ); --对话成员
  
    DROP TABLE IF EXISTS MEDIA_RESOURCE;
    CREATE TABLE MEDIA_RESOURCE(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        key TEXT(255) NOT NULL, --资源 key
        media_type INTEGER NOT NULL, --资源类型
        width INTEGER NOT NULL DEFAULT 0, --宽度
        height INTEGER NOT NULL DEFAULT 0, --高度
        size INTEGER NOT NULL DEFAULT 0, --资源大小
        resolution TEXT(255) NOT NULL DEFAULT '', --分辨率
        duration TEXT NOT NULL DEFAULT '', --视频时长
        title TEXT(255) NOT NULL DEFAULT '', --
        attachments TEXT(255) NOT NULL DEFAULT '', --
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP --
    
    ); --媒体资源
    
  
    DROP TABLE IF EXISTS FITNESS_MEMBER;
    CREATE TABLE FITNESS_MEMBER(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        name TEXT NOT NULL DEFAULT '', --名称
        age INTEGER NOT NULL DEFAULT 0, --年龄
        gender INTEGER NOT NULL DEFAULT 1, --性别
        body_type INTEGER NOT NULL DEFAULT 2, --体型;1偏瘦 2中等 3偏胖 4肌肉 5匀称
        height INTEGER NOT NULL DEFAULT 0, --身高，单位 cm
        weight REAL NOT NULL DEFAULT 0, --体重，单位 kg
        body_fat_percent INTEGER NOT NULL DEFAULT 0, --体脂率，12 就是 12%
        risk_screenings TEXT(1000) NOT NULL DEFAULT '{}', --健康风险评估
        training_goals TEXT(1000) NOT NULL DEFAULT '{}', --训练目标
        training_frequency INTEGER NOT NULL DEFAULT 1, --训练频率;一周几次
        training_preferences TEXT(1000) NOT NULL DEFAULT '{}', --训练偏好
        diet_preferences TEXT(1000) NOT NULL DEFAULT '{}', --饮食偏好
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP --
    
    ); --学员
    
  
    DROP TABLE IF EXISTS PHYSICAL_TEST;
    CREATE TABLE PHYSICAL_TEST(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        member_id INTEGER NOT NULL DEFAULT 0, --学员id
        action_ids TEXT(255) NOT NULL DEFAULT '', --包含的动作
        result TEXT(255) NOT NULL DEFAULT '', --评估结果
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP --创建时间
    
    ); --体能评估
    
  
    DROP TABLE IF EXISTS WORKOUT_ACTION;
    CREATE TABLE WORKOUT_ACTION(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        name TEXT(255) NOT NULL DEFAULT '', --动作名称
        zh_name TEXT(255) NOT NULL DEFAULT '', --中文名称
        alias TEXT(255) NOT NULL DEFAULT '', --别名
        overview TEXT NOT NULL DEFAULT '', --简要说明
        type TEXT NOT NULL DEFAULT "resistance", --动作类型;resistance、cardio、balance、flexibility、strength
        level INTEGER NOT NULL DEFAULT 1, --难度等级;1-10
        tags1 TEXT(255) NOT NULL DEFAULT '', --标签;逗号分割
        tags2 TEXT(255) NOT NULL DEFAULT '', --标签;逗号分割
        details TEXT(900) NOT NULL DEFAULT '{}', -- 详情 JSON
        points TEXT(255) NOT NULL DEFAULT '{}', --动作要点;逗号分割
        problems TEXT(255) NOT NULL DEFAULT '{}', --常见错误;逗号分割
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, --创建时间
        updated_at DATETIME, --更新时间
        equipment_ids TEXT(255) NOT NULL DEFAULT '', --器械;逗号分割
        muscle_ids TEXT(255) NOT NULL DEFAULT '', --肌肉;逗号分割
        alternative_action_ids TEXT(255) NOT NULL DEFAULT '', --替代动作;逗号分割
        advanced_action_ids TEXT(255) NOT NULL DEFAULT '', --进阶动作;逗号分割
        regressed_action_ids TEXT(255) NOT NULL DEFAULT '' --退阶动作;逗号分割
    
    ); --健身动作
    
  
    DROP TABLE IF EXISTS MUSCLE;
    CREATE TABLE MUSCLE(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        name TEXT NOT NULL DEFAULT '', --肌肉名称
        zh_name TEXT NOT NULL DEFAULT '', --肌肉名称
        overview TEXT NOT NULL DEFAULT '', --简要说明
        tags TEXT NOT NULL DEFAULT '', --标签;逗号分割
        features TEXT NOT NULL DEFAULT '{}' --功能
    
    ); --肌肉


    INSERT INTO MUSCLE (name, zh_name, overview, tags, features)
    VALUES 
    ('Brachialis', '肱肌', '位于上臂前侧肱二头肌的深层', '手臂', '[{"title":"肘关节屈曲","details":"当前臂向靠近上臂移动时，肱肌强烈收缩以完成这一动作"}]'),
    ('Biceps brachii', '肱二头肌', '位于上臂前侧，是上臂的主要屈肌', '手臂', '[{"title":"肘关节屈曲","details":"当前臂向靠近上臂移动时，肱二头肌强烈收缩以完成这一动作"}]'),
    ('Triceps brachii', '肱三头肌', '位于上臂后侧，是上臂的主要伸肌', '手臂', '[{"title":"肘关节伸直","details":"当前臂向远离上臂移动时，肱三头肌强烈收缩以完成这一动作"}]'),
    ('Deltoid', '三角肌', '位于肩部，是肩部的主要肌肉', '肩膀', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Trapezius', '斜方肌', '位于背部，是背部的主要肌肉', '背部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Pectoralis major', '胸大肌', '位于胸部，是胸部的主要肌肉', '胸部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Latissimus dorsi', '背阔肌', '位于背部，是背部的主要肌肉', '背部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Gluteus maximus', '臀大肌', '位于臀部，是臀部的主要肌肉', '臀部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Quadriceps', '股四头肌', '位于大腿前侧，是腿部的主要肌肉', '腿部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Hamstrings', '股二头肌', '位于大腿后侧，是腿部的主要肌肉', '腿部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Gluteus maximus', '臀大肌', '位于臀部，是臀部的主要肌肉', '臀部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Abdominals', '腹直肌', '位于腹部，是腹部的主要肌肉', '腹部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Obliques', '腹斜肌', '位于腹部，是腹部的主要肌肉', '腹部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Serratus anterior', '前锯肌', '位于胸部，是胸部的主要肌肉', '胸部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Rhomboids', '菱形肌', '位于背部，是背部的主要肌肉', '背部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Erector spinae', '竖脊肌', '位于背部，是背部的主要肌肉', '背部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]'),
    ('Sartorius', '梨状肌', '位于臀部，是臀部的主要肌肉', '臀部', '[{"title":"肩关节外展","details":"当前臂向远离上臂移动时，三角肌强烈收缩以完成这一动作"}]');

    
    DROP TABLE IF EXISTS EQUIPMENT;
    CREATE TABLE EQUIPMENT(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        name TEXT NOT NULL DEFAULT '', --器械名称
        zh_name TEXT NOT NULL DEFAULT '', --器械名称
        alias TEXT NOT NULL DEFAULT '', --别名
        overview TEXT NOT NULL DEFAULT '', --简要说明
        medias TEXT NOT NULL DEFAULT '{}' --图片、视频等
    ); --器械

    INSERT INTO EQUIPMENT (name, zh_name, overview)
    VALUES
    ('dumbbell', '哑铃', '哑铃是一种用于增强肌肉力量训练的简单器材，可进行多种部位的力量练习。'),
    ('barbell', '杠铃', '杠铃是举重及健身练习的常见器械，可进行全身性的力量训练。'),
    ('treadmill', '跑步机', '跑步机是家庭及健身房常见的健身器材，是有氧运动的有效工具。'),
    ('elliptical', '椭圆机', '椭圆机是一种全身性的有氧运动器械，对膝关节压力较小。'),
    ('ab_machine', '仰卧板', '仰卧板主要用于腹部等部位的力量训练，辅助进行仰卧起坐等动作。'),
    ('chest_press_machine', '坐姿推胸器', '这是一种固定健身器械，主要用于锻炼胸大肌，通过坐姿推动把手来完成动作，能有效增强胸部力量和维度。'),
    ('butterfly_machine', '蝴蝶机', '专门针对胸肌内侧进行孤立训练的固定器械，模拟蝴蝶展翅动作，可塑造胸部线条。'),
    ('row_machine', '坐姿划船器', '帮助锻炼背部肌群的固定器械，尤其对背阔肌、斜方肌中下束等有很好的刺激效果，增强背部的宽度和厚度。'),
    ('leg_extension_machine', '腿部推蹬机', '主要用于锻炼腿部的大型肌群，如股四头肌、股二头肌和臀大肌等，提供稳定的发力环境来增强腿部力量。'),
    ('leg_curl_machine', '腿部伸展机', '专门针对股四头肌进行锻炼的器械，可有效孤立训练股四头肌，提升其力量和线条。'),
    ('leg_curl_machine', '腿部弯举机', '主要用于锻炼股二头肌，通过弯曲小腿的动作来刺激股二头肌的收缩和生长。'),
    ('smith_machine', '史密斯机', '一种多功能的固定训练器械，具有安全保护装置，可进行深蹲、卧推、硬拉等多种训练动作。'),
    ('cable_machine', '龙门架', '多功能训练设备，可搭配各种配件进行全身各部位的训练，如绳索夹胸、高位下拉、坐姿腿屈伸等动作。'),
    ('pull_up_machine', '引体向上器', '经典的背部训练器械，主要锻炼背阔肌和肱二头肌等，也可辅助进行腹肌训练。'),
    ('roman_chair', '罗马椅', '常用于训练腰背肌和腹肌的器械，可进行罗马椅挺身、仰卧抬腿等动作，增强核心力量。'),
    ('kettlebell', '壶铃', '一种古老的小工具，可用于全身力量、爆发力和协调性训练，常见动作有壶铃摇摆、高脚杯深蹲等。'),
    ('resistance_band', '弹力带', '轻便且多功能的小工具，具有不同的阻力级别，可用于热身、康复训练以及力量辅助训练等，增加训练的多样性。'),
    ('foam_roller', '泡沫轴', '主要用于放松肌肉、缓解肌肉紧张和改善柔韧性的小工具，通过自身重量在泡沫轴上滚动来按摩肌肉。'),
    ('yoga_ball', '瑜伽球', '可用于平衡、核心力量和柔韧性训练的小工具，常见动作有球上平板支撑、仰卧抬腿等。'),
    ('battle_rope', '健身战绳', '一种高强度的训练小工具，通过快速甩动战绳来锻炼手臂、肩部、背部和核心等部位的力量和耐力。'),
    ('trx_suspension_trainer', 'TRX悬挂训练带', '利用自身重量进行全身力量、平衡、柔韧性和核心稳定性训练的工具，训练动作丰富多样。');
  
    DROP TABLE IF EXISTS WORKOUT_ACTION_MISTAKE;
    CREATE TABLE WORKOUT_ACTION_MISTAKE(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        name TEXT(255) NOT NULL DEFAULT '', --错误名称
        reason TEXT(255) NOT NULL DEFAULT '', --常见错误原因
        solution_direction TEXT(255) NOT NULL DEFAULT '', --解决方案
        solution_action_ids TEXT(255) NOT NULL DEFAULT '', --解决动作
        solution_action_text TEXT(255) NOT NULL DEFAULT '' --解决动作描述
    
    ); --动作错误点
    
  
    DROP TABLE IF EXISTS WORKOUT_PLAN;
    CREATE TABLE WORKOUT_PLAN(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        title TEXT(255) NOT NULL DEFAULT '', --计划名称
        overview TEXT(255) NOT NULL DEFAULT '', --训练计划描述
        level INTEGER NOT NULL DEFAULT 1, --适宜什么水平的人群;1-10
        tags TEXT(255) NOT NULL DEFAULT '', --部位标签
        estimated_duration INTEGER NOT NULL DEFAULT 0, --预计耗时;单位 min
        equipment_ids TEXT(255) NOT NULL DEFAULT '', --所需器械 id 列表
        muscle_ids TEXT(255) NOT NULL DEFAULT '', --该计划练到的肌肉 id 列表
        details TEXT(1000) NOT NULL DEFAULT '{}', --内容详情 JSON;在创建、更新 plan 时，根据内容统计出的动作列表 JSON
        points TEXT(255) NOT NULL DEFAULT '', --注意事项
        suggestions TEXT(255) NOT NULL DEFAULT '', --训练计划建议
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, --创建时间
        updated_at DATETIME, --更新时间
        owner_id INTEGER NOT NULL DEFAULT 0 --创建人id
    
    ); --训练计划
    
  
    DROP TABLE IF EXISTS WORKOUT_PLAN_STEP;
    CREATE TABLE WORKOUT_PLAN_STEP(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        title TEXT(255) NOT NULL DEFAULT '', --动作(阶段)名称
        step_type TEXT(255) NOT NULL DEFAULT 'strength', --动作(阶段)类型;warmup、strength、stretch、cool_down、cardio、heart、performance
        idx INTEGER NOT NULL DEFAULT 1, --第几个动作(阶段)
        set_count INTEGER NOT NULL DEFAULT 1, --组数
        sets TEXT(1000) NOT NULL DEFAULT '{}', --组详情;JSON 数组
        note TEXT(255) NOT NULL DEFAULT '', --该组额外信息说明
        workout_plan_id INTEGER NOT NULL DEFAULT 0 --所属训练计划id
    
    ); --训练计划中的一组
    
  
    DROP TABLE IF EXISTS WORKOUT_PLAN_ACTION;
    CREATE TABLE WORKOUT_PLAN_ACTION(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        idx INTEGER NOT NULL DEFAULT 1, --动作顺序;超级组、循环组会需要该字段
        action_id INTEGER NOT NULL DEFAULT 0, --要做的动作id
        set_idx INTEGER NOT NULL DEFAULT 1, --第几组
        reps INTEGER NOT NULL DEFAULT 0, --数量
        unit TEXT(255) NOT NULL DEFAULT '个', --单位
        weight TEXT(255) NOT NULL DEFAULT '', --阻力;字符串，直接写成 60%1RM 或 12RM 或 自重、无负重等等？
        tempo TEXT(255) NOT NULL DEFAULT '4/1/2', --节奏;4/1/2 表示离心4s，停顿1s，向心2s
        rest_interval TEXT(32) NOT NULL DEFAULT '0', --间歇时间
        note TEXT(255) NOT NULL DEFAULT '', --该组备注信息
        workout_plan_step_id INTEGER NOT NULL DEFAULT 0 --关联的训练计划中组id
    
    ); --训练计划中的动作要求
    
  
    DROP TABLE IF EXISTS WORKOUT_DAY_ACTION;
    CREATE TABLE WORKOUT_DAY_ACTION(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        action_id INTEGER NOT NULL DEFAULT 0, --做的动作id
        reps INTEGER NOT NULL DEFAULT 0, --实际数量
        unit TEXT(255) NOT NULL DEFAULT '个', --单位
        weight TEXT(255) NOT NULL DEFAULT '', --实际阻力;字符串，直接写成 60%1RM 或 12RM 或 自重、无负重等等？
        note TEXT(255) NOT NULL DEFAULT '', --教练评价、备注。比如动作不标准之类
        feedback TEXT(255) NOT NULL DEFAULT '', --动作反馈
        extra_medias TEXT(1000) NOT NULL DEFAULT '{}', --拍照、录像记录;存 字符串列表
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, --创建时间;CURRENT_TIMESTAMP
        workout_day_id INTEGER NOT NULL DEFAULT 0, --训练日 id
        workout_day_set_id INTEGER NOT NULL DEFAULT 0, --训练日哪个组
        workout_plan_action_id INTEGER NOT NULL DEFAULT 0 --
    
    ); --动作执行记录
    
  
    DROP TABLE IF EXISTS COACH;
    CREATE TABLE COACH(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        nickname TEXT(255) NOT NULL DEFAULT '', --昵称
        avatar_url TEXT(255) NOT NULL DEFAULT '', --头像链接
        config TEXT(1000) NOT NULL DEFAULT '{}', --配置信息
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, --创建时间
        updated_at DATETIME --更新时间
    
    ); --用户
  
    DROP TABLE IF EXISTS COACH_ACCOUNT;
    CREATE TABLE COACH_ACCOUNT(
    
        provider_type TEXT(255) NOT NULL, --帐号授权方式;email、phone、wxapp 等等
        provider_id TEXT(255) NOT NULL, --帐号唯一标志;如果 email，这里就是 email 帐号，可以发送验证码来验证
        provider_arg1 TEXT(255) NOT NULL DEFAULT '', --帐号授权参数
        provider_arg2 TEXT(255) NOT NULL DEFAULT '', --授权参数2
        provider_arg3 TEXT(255) NOT NULL DEFAULT '', --授权参数3;可以放很多额外信息，比如 expires
        coach_id INTEGER NOT NULL DEFAULT 0, --帐号关联的用户id
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, --帐号创建时间

        PRIMARY KEY (provider_type, provider_id)
    ); --用户帐号

    DROP TABLE IF EXISTS WORKOUT_DAY;
    CREATE TABLE WORKOUT_DAY(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        time TEXT NOT NULL DEFAULT '', --训练时间;年月日时分秒
        status INTEGER NOT NULL DEFAULT 1, --训练日状态;1等待进行 2进行中 3已完成 4已过期 5手动作废
        estimated_duration INTEGER NOT NULL DEFAULT 0, --预计时间
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, --创建时间
        updated_at DATETIME, --更新时间
        coach_id INTEGER NOT NULL DEFAULT 0, --教练id
        member_id INTEGER NOT NULL DEFAULT 0, --学员id
        workout_plan_id INTEGER NOT NULL DEFAULT 0 --关联的训练计划id
    
    ); --训练日
    
  
    DROP TABLE IF EXISTS WORKOUT_DAY_SET;
    CREATE TABLE WORKOUT_DAY_SET(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        completed INTEGER NOT NULL, --是否完成;1完成 2未完成
        step_order INTEGER NOT NULL, --第几阶段
        set_order INTEGER NOT NULL, --第几组
        feedback TEXT(255) NOT NULL DEFAULT '', --反馈，自感疲劳度？
        workout_day_id INTEGER NOT NULL DEFAULT 0, --训练日id
        workout_plan_set_id INTEGER NOT NULL DEFAULT 0 --训练计划组id
    
    ); --训练日组记录
    
  
    DROP TABLE IF EXISTS BODY_MEASUREMENT;
    CREATE TABLE BODY_MEASUREMENT(
    
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, --id
        height INTEGER NOT NULL DEFAULT 0, --身高;单位厘米
        weight REAL NOT NULL DEFAULT 0, --体重;单位kg
        body_fat_percentage REAL NOT NULL DEFAULT 0, --体脂率
        heart_rate INTEGER NOT NULL DEFAULT 0, --静息心率;次每分钟
        chest REAL NOT NULL DEFAULT 0, --胸围
        waist REAL NOT NULL DEFAULT 0, --腰围
        hip REAL NOT NULL DEFAULT 0, --臀围
        arm REAL NOT NULL DEFAULT 0, --臂围
        thigh REAL NOT NULL DEFAULT 0, --大腿围度
        notes TEXT(255) NOT NULL DEFAULT '', --额外备注信息
        created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP --体测时间
    
    ); --体测记录
    
