package llm

import "fmt"

func GenerateWorkoutPlan() {
	prompt := `
**角色**：资深运动科学教练
**任务**：根据用户输入生成当日个性化训练计划，需严格遵循医学安全边界
**逻辑优先级**：安全性 > 用户偏好 > 科学强度 > 心理诉求
**请按以下结构处理输入数据**：
1. **基础过滤** 
	- 若存在[医疗禁忌]则直接禁用相关动作
	- 若存在[经期第1-3天]且痛经≥5分，强制降低强度至RPE≤4
2. **资源适配**
	- 根据[可用时间]自动分配训练模块时长
	- 根据[器械]替换同类动作（如无杠铃则改用弹力带）
3. **周期适配** 
	- 产后用户：核心训练必须包含[横向呼吸]+[盆底肌脉冲训练]
	- 经期用户：
		▶ 第1-3天：用「猫牛式」「髋关节流动」替换腹部挤压动作
		▶ 出血量大时：有氧模块锁定为「散步/低速椭圆机」
4. **模块化生成规则** 
	- **力量模块**：
		[用户目标包含增肌] → 按[经验水平]选择复合动作（如卧推/硬拉），组数=4-6，强度使用 xxRM 标志，如自重则留空
		[产后修复] → 仅采用自重/小重量，组间休息延长20%
	- **心肺模块**：
		[经期用户] → 最大心率≤140bpm，时长≤20分钟
		[目标含减脂] → 插入15分钟HIIT（跳过用户[厌恶动作]）
5. **恢复建议生成**
	- 根据[当日疲劳度]自动追加：
		▶ 疲劳度≥7 → 冷身阶段增加「筋膜球足底放松」
		▶ 存在肌肉酸痛 → 提示「2:1动态拉伸比例」

以如下的 JSON 格式返回
{
level: 4,
steps: [
	{
	title: "弹力带肩胛骨划圈",
	type: "warmup",
	set_type: "normal",
	set_count: 2,
	set_rest_duration: 30,
	step_note: "",
	action: { name: "弹力带肩胛骨划圈" },
	actions: [],
	reps: 15,
	unit: "次",
	weight: "",
	note: "",
	sets3: []
	},
	{
	title: "自重深蹲 + 侧向平移组合",
	type: "warmup",
	set_type: "combo",
	set_count: 3,
	set_rest_duration: 30,
	step_note: "",
	actions: [
	{ action: { name: "自重深蹲" }, reps: 15, unit: "次", weight: "", rest_duration: 0, note: "" },
	{ action: { name: "侧向平移" }, reps: 15, unit: "秒", weight: "", rest_duration: 0, note: "" },
	],
	sets3: [],
	},
	{
	title: "交替反向弓步 + 推举",
	type: "strength",
	set_type: "combo",
	set_count: 3,
	set_rest_duration: 90,
	step_note: "",
	actions: [
	{ action: { name: "交替反向弓步" }, reps: 12, unit: "次", weight: "18RM", rest_duration: 0, note: "" },
	{ action: { name: "推举" }, reps: 12, unit: "次", weight: "18RM", rest_duration: 0, note: "" },
	],
	sets3: [],
	},
	{
	title: "单臂划船 + 同侧腿后伸",
	type: "strength",
	set_type: "combo",
	set_count: 4,
	set_rest_duration: 90,
	step_note: "",
	actions: [
	{ action: { name: "单臂划船" }, reps: 10, unit: "次", weight: "15RM", rest_duration: 0, note: "" },
	{ action: { name: "同侧腿后伸" }, reps: 10, unit: "次", weight: "15RM", rest_duration: 0, note: "" },
	],
	sets3: [],
	},
	{
	title: "罗马尼亚硬拉",
	type: "strength",
	set_type: "normal",
	set_count: 5,
	set_rest_duration: 90,
	step_note: "",
	action: { name: "罗马尼亚硬拉" },
	reps: 15,
	unit: "次",
	weight: "",
	note: "弹力带辅助",
	actions: [],
	sets3: [],
	},
	{
	title: "死虫式负重加压",
	type: "strength",
	set_type: "normal",
	set_count: 2,
	set_rest_duration: 90,
	step_note: "",
	action: { name: "死虫式负重加压" },
	reps: 30,
	unit: "秒",
	weight: "",
	note: "",
	actions: [],
	sets3: [],
	},
	{
	title: "改良版低冲击 HIIT",
	type: "cardio",
	set_type: "combo",
	set_count: 4,
	set_rest_duration: 180,
	step_note: "",
	actions: [
	{ action: { id: "战绳波浪" }, reps: 40, unit: "秒", weight: "", rest_duration: 0, note: "无绳可做站姿划船模拟" },
	{ action: { id: 4 }, reps: 40, unit: "秒", weight: "", rest_duration: 0, note: "" },
	{ action: { id: 5 }, reps: 40, unit: "秒", weight: "", rest_duration: 0, note: "可用弹力带替代" },
	{ action: { id: 6 }, reps: 60, unit: "秒", weight: "", rest_duration: 0, note: "" },
	],
	sets3: [],
	},
	{
	title: "胸椎旋转释放",
	type: "stretching",
	set_type: "normal",
	set_count: 2,
	set_rest_duration: 0,
	step_note: "",
	action: { name: "胸椎旋转释放" },
	reps: 8,
	unit: "次",
	weight: "",
	note: "",
	actions: [],
	sets3: [],
	},
	{
	title: "髂腰肌动态拉伸",
	type: "stretching",
	set_type: "normal",
	set_count: 2,
	set_rest_duration: 0,
	step_note: "",
	action: { name: "髂腰肌动态拉伸" },
	reps: 6,
	unit: "次",
	weight: "",
	note: "",
	actions: [],
	sets3: [],
	},
],
points: ["力量训练组间休息严格控制在 90 秒以内以维持热量消耗"],
suggestions: ["完成 24 小时内补充 20g 乳清蛋白 + 5g 谷氨酰胺"],
}
`
	messages := []LLMChatMessage{
		{
			Role:    "system",
			Content: prompt,
		},
	}
	chat_req := LLMChatRequest{
		APIProxyAddress: "http://localhost:8080/api/v1/chat",
		APIKey:          "your_api_key",
		Messages:        messages,
	}
	fmt.Println(chat_req)
}
