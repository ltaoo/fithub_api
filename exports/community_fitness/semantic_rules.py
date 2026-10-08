"""Explicit movement-equivalence rules. Never fuzzy-match for automatic merging."""
import re

# Qualifiers omitted here are ONLY omitted when the reviewed names below express the same movement.
GROUPS = [
 ('barbell-flat-bench-press',['Barbell Bench Press','Barbell Flat Bench Press','Flat Barbell Bench Press','Barbell Bench Press - Medium Grip','杠铃卧推','杠铃平板卧推','平板卧推'], '平凳、杠铃、常规握距；上斜、下斜、宽窄握、史密斯、哑铃均单独保留'),
 ('pull-up',['Pull Up','Pull-ups','Pull - ups','Pullups','引体向上','单杠引体'], '正握引体向上；反握、辅助、负重、单臂、肩胛骨变式分别保留'),
 ('bodyweight-squat',['Squat','Squats','Bodyweight Squat','Air Squat','深蹲','自重深蹲'], '标准自重深蹲；杠铃、前蹲、相扑、单腿、窄站距、跳蹲分别保留'),
 ('forearm-plank',['Plank','Forearm Plank','平板支撑','前撑'], '双前臂支撑；直臂、侧撑、抬腿、行走变式单独保留'),
 ('seated-dumbbell-shoulder-press',['Dumbbell Shoulder Press','Seated Dumbbell Press','Seated Dumbbell Shoulder Press','哑铃肩推','坐姿哑铃推肩'], '已核对起始姿势为坐姿；Standing Dumbbell Press 不合并'),
 ('conventional-barbell-deadlift',['Deadlift','Barbell Deadlift','Conventional Deadlift','硬拉','传统硬拉','屈腿硬拉','杠铃硬拉'], '常规站距杠铃硬拉；相扑、六角杠、罗马尼亚、直腿、器械硬拉单独保留'),
 ('standard-push-up',['Push Up','Push-ups','Push - ups','Pushups','俯卧撑','标准俯卧撑'], '标准俯卧撑；钻石、宽窄距、抬脚、单臂、跪姿单独保留'),
 ('machine-leg-press',['Leg Press','腿推','腿举','倒蹬','坐姿腿推'], '同类普通双腿腿举；未确认的 Seated Leg Press、单腿、窄站距和史密斯腿举保留'),
 ('single-arm-dumbbell-row',['Dumbbell Row','One-arm Dumbbell Row','One Arm Dumbbell Row','Single Arm Dumbbell Row','哑铃划船','单臂划船','单臂哑铃划船'], '现有 Dumbbell Row 描述为单臂凳上支撑；双臂、反握、胸托和壶铃划船保留'),
 ('barbell-bent-over-row',['Barbell Row','Bent Over Barbell Row','杠铃划船','俯身杠铃划船'], '正握俯身杠铃划船；反握、潘德雷、史密斯、长杆地雷架划船保留'),
 ('incline-barbell-bench-press',['Incline Barbell Bench Press','Barbell Incline Bench Press - Medium Grip','Barbell Incline Bench Press','上斜杠铃卧推'], '常规握距上斜杠铃卧推'),
 ('incline-dumbbell-bench-press',['Incline Dumbbell Bench Press','Incline Dumbbell Press','上斜哑铃卧推'], '普通上斜哑铃卧推；锤式握法和单臂保留'),
 ('close-grip-barbell-bench-press',['Close Grip Bench Press','Close-Grip Barbell Bench Press','窄距卧推','窄握杠铃卧推'], '窄握平板杠铃卧推；哑铃、曲杠、史密斯、上斜下斜保留'),
 ('single-arm-dumbbell-bench-press',['Single Arm Bench Press','One Arm Dumbbell Bench Press','Single Arm Dumbbell Bench Press','单臂卧推','单臂哑铃卧推'], '普通单臂哑铃平板卧推'),
 ('dumbbell-biceps-curl',['Dumbbell Curl','Dumbbell Bicep Curl','Dumbbell Biceps Curl','哑铃弯举'], '普通哑铃弯举；交替、锤式、坐姿、上斜、单臂、牧师凳保留'),
 ('incline-seated-dumbbell-curl',['Inclined Seated Dumbbell Curl','Incline Dumbbell Curl','上斜坐姿哑铃弯举'], '双臂同时上斜凳弯举；交替、扭转和 Flexor 变式保留'),
 ('ez-bar-biceps-curl',['Barbell Bicep Curl','EZ-Bar Curl','曲杠二头弯举'], '现有动作明确使用曲杆；直杆 Barbell Curl 保留'),
 ('ez-bar-lying-triceps-extension',['Lying Triceps Extension','EZ-Bar Skullcrusher','仰卧臂屈伸'], '现有动作明确使用曲杆；哑铃、直杠和绳索版本保留'),
 ('lying-machine-leg-curl',['Prone Leg Curl','Lying Leg Curls','Lying Leg Curl','俯卧腿弯举'], '俯卧器械腿弯举；坐姿、站姿、健身球、弹力带和单腿保留'),
 ('seated-machine-leg-extension',['Seated Leg Extension','Leg Extensions','Leg Extension','坐姿腿屈伸'], '坐姿双腿器械伸膝；单腿版本保留'),
 ('barbell-upright-row',['Upright Row','Upright Barbell Row','直立提拉','直立杠铃划船'], '现有描述为杠铃直立提拉；绳索、哑铃和史密斯保留'),
 ('standing-dumbbell-lateral-raise',['Dumbbell Lateral Raise','Side Lateral Raise','站姿哑铃侧平举'], '站姿双臂哑铃侧平举；坐姿、单臂、绳索和器械保留'),
 ('dumbbell-step-up',['Dumbbell Step Up on Bench','Dumbbell Step Ups','哑铃健身椅登阶'], '持哑铃登阶；自重、杠铃、提膝保留'),
 ('barbell-floor-press',['Floor Press','Barbell Floor Press','地面卧推'], '现有动作使用杠铃；哑铃、单臂、壶铃和桥式保留'),
 ('farmer-walk',['Farmer Walk',"Farmer's Walk",'农夫行走'], '普通双手农夫行走；单侧和过头负重行走保留'),
 ('treadmill-jogging',['Jogging on Treadmill','Jogging, Treadmill','跑步机慢跑'], '慢跑；步行、跑步和坡度版本保留'),
 ('chair-side-lower-back-stretch',['Chair Lower Back Stretch','椅子上腰部拉伸'], '现有两个记录的起始姿势和步骤完全相同'),
 ('overhead-triceps-stretch',['Triceps Stretch','肱三头肌拉伸','三头肌拉伸'], '单手辅助的过头三头肌静态拉伸；PNF 同伴辅助和侧拉伸保留'),
 ('mountain-climber',['Mountain Climber Crawl','Mountain Climbers','登山式爬行'], '双手固定、双膝交替向胸部移动；蜘蛛爬、侧爬和攀爬机保留'),
 ('quad-foam-roll',['Quadriceps Foam Roller Massage','Quadriceps-SMR','股四头肌泡沫轴按摩'], '股四头肌泡沫轴自我放松'),
 ('calf-foam-roll',['Calf Foam Roller Massage','Calves-SMR','小腿泡沫轴按摩'], '小腿泡沫轴自我放松'),
 ('glute-foam-roll',['Gluteus Maximus Foam Roller Massage','Glute-SMR','臀大肌泡沫轴按摩'], '臀肌泡沫轴自我放松；梨状肌放松保留'),
 ('single-arm-landmine-row',['Bent Over One-Arm Long Bar Row','One-Arm Long Bar Row','单臂长杆划船','俯身单臂长杆划船'], '同一地雷架单臂划船'),
 ('standing-cable-glute-kickback',['Low Cable Back Kick','One-Legged Cable Kickback','低位绳索后踢','单腿绳索后踢'], '绑带低位绳索站姿后踢'),
 ('high-cable-biceps-curl',['Overhead Cable Curl','High Cable Curls','高位绳索弯举'], '双臂高位滑轮弯举；低位和仰卧保留'),
 ('upper-back-interlaced-hands-stretch',['上背部拉伸'], '双手手指相扣向前推；与上背环抱手臂拉伸不同'),
 ('upper-back-crossed-arms-stretch',['上背拉伸'], '身前环抱交叉手臂；与手指相扣拉伸不同'),
 ('straight-legged-bench-dip',['Triceps Dip on Bench','Bench Dips','凳上肱三头肌臂屈伸'], '双腿伸直凳上臂屈伸；现有 Bench Dip 为屈膝版，分别保留'),
]
SOURCE_OVERRIDES={'Upper_Back_Stretch':'upper-back-interlaced-hands-stretch'}
# Two records share the English name but have different techniques. Do not match that label alone.
AMBIGUOUS_NAMES={'Upper Back Stretch'}

def normalize(value):
 # Only spelling/format, never remove meaningful qualifiers or compare by similarity.
 value=value.casefold().replace('’',"'")
 return re.sub(r'[\s\-_,.\'"()/®™：:]+','',value)

def explicit_aliases():
 result={}
 for key,aliases,_ in GROUPS:
  for alias in aliases:
   n=normalize(alias)
   if n in result and result[n]!=key:raise ValueError('Conflicting reviewed alias: '+alias)
   result[n]=key
 return result

def key_for(record, source=False):
 if source and record.get('id') in SOURCE_OVERRIDES:return SOURCE_OVERRIDES[record['id']]
 if record.get('name') in AMBIGUOUS_NAMES:
  return explicit_aliases().get(normalize(record.get('zh_name','')), 'ambiguous:'+normalize(record.get('zh_name') or record['name']))
 return explicit_aliases().get(normalize(record['name']), 'movement:'+normalize(record['name']))
