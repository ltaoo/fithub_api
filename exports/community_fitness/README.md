# FitHub 全新项目：语义去重动作库与器械图片

适用于完成当前项目迁移的 **SQLite 全新项目**。本版本将项目已有的 420 个公开动作模板与社区的 876 条数据按动作语义整理为 **962 个动作**，合并 334 条来源重复记录。

## 重要：这是全新项目初始化 SQL

`deduplicated_workout_actions.sql` 与 `import_all.sql` 会 **清空 WORKOUT_ACTION 中的默认动作数据并重建动作 ID**，并非此前的增量追加脚本。
SQL 先检查训练计划、训练日、训练动作记录、收藏、关联内容、体测、排期、自定义动作等数据；不符合全新项目条件时会失败并回滚，不会处理历史记录。
器械、肌肉继续使用项目固定编号。不能把该版本用于有业务数据的项目。
本次没有导入正在运行的本地数据库。

## 文件

- `deduplicated_workout_actions.sql`：仅初始化去重后的 962 个动作，并补齐固定 ID 肌肉、器械。
- `community_workout_actions.sql`：上述文件的兼容文件名，内容相同，同样只适用于全新项目。
- `equipment_photos.sql`：仅更新 46 个固定 ID 的器械图片，可以单独重复执行。
- `import_all.sql`：单一事务执行动作库初始化与器械图片关联。
- `media/equipment/`：已下载的 46 张图片。
- `semantic_dedup_report.json`：每个合并组、规范名称、来源名称、规则说明和数量统计。
- `deduplicated_actions.json`：SQL 对应的完整去重动作数据。
- `semantic_rules.py`：38 组明确的动作同义规则。没有使用模糊相似度或仅凭目标肌肉自动合并。
- `existing_action_catalog.json`：项目公开动作模板快照，含已整理的中文动作说明，不含用户或训练记录。
- `fixed_entities.json`、`import_mapping.py`：项目固定肌肉和器械编号、转换配置。
- `equipment_sources.json`：商品页、原图地址、文件大小、SHA256。
- `exercises.json`、`SOURCE_REVISION`、`LICENSE.free-exercise-db.md`：原始社区数据、上游精确版本、Unlicense。
- `generate.py`：重新生成去重动作 SQL、JSON 和报告。
- `validate_import.py`：只读获取参考库结构，在内存数据库验证导入。

## 部署与一键导入

将 `fithub-community-fitness.tar.gz` 上传到服务器。在 API 项目根目录执行，数据库名按实际 `DB_PATH` 修改：

```sh
mkdir -p imports/community_fitness
tar -xzf fithub-community-fitness.tar.gz -C imports/community_fitness

mkdir -p dist/h5/media/equipment h5/public/media/equipment
cp imports/community_fitness/media/equipment/* dist/h5/media/equipment/
cp imports/community_fitness/media/equipment/* h5/public/media/equipment/

# 路径必须是已迁移的现有数据库，避免意外创建空库
sqlite3 myapi.db ".backup 'myapi.before-initial-catalog.db'"

# 遇错立即停止；事务会回滚，不能在数据库工具中设置忽略错误继续执行
sqlite3 -bail -cmd '.timeout 10000' myapi.db < imports/community_fitness/import_all.sql
```

只导入动作时执行 `deduplicated_workout_actions.sql`，只导入器械图片时执行 `equipment_photos.sql`。

如果 `STATIC_DIR` 不是 `./dist`，将图片放到 `$STATIC_DIR/h5/media/equipment/`。
`h5/public/media/equipment/` 保留源码副本，使后续 Vite 构建继续携带图片。Go 无需重启。
SQL 上传工具可运行对应文件，但须遇错停止并回滚。图片文件必须随 SQL 部署；SQL 不存储图片二进制。

## 动作语义去重

已确认的同一动作使用一个条目，其他名称进入逗号分隔的 `alias`，供现有搜索接口检索。例如：

- 杠铃卧推、杠铃平板卧推、平板卧推、Barbell Bench Press - Medium Grip：同一普通平板杠铃卧推。
- Pull Up、Pull - ups、Pullups、引体向上：同一普通正握引体向上。
- Dumbbell Shoulder Press、Seated Dumbbell Shoulder Press、Seated Dumbbell Press：核对步骤后归为坐姿哑铃肩推；站姿版独立保留。
- Bent Over One-Arm Long Bar Row、One-Arm Long Bar Row：同一地雷架单臂划船。

保留器械、角度、握距、握法、姿势、单侧/双侧等有意义的变式。特别保留：上斜/下斜/宽窄握卧推，站姿/坐姿肩推，直杠/曲杠弯举，传统/相扑/罗马尼亚/直腿硬拉，屈膝/直腿凳上臂屈伸。

名称相同也未必合并。现有两个 Upper Back Stretch 分别是身前环抱手臂和手指相扣向前推，因此保留为两个动作。
未知或描述不足的相似项暂时保留，不能保证其他同义动作不存在；可在规则文件中明确确认后再次生成。报告中的名字相似度不作为自动删除依据。

项目模板和社区动作重叠时，优先采用项目中说明较完整的记录，保留中文内容；社区名称、步骤、图片、器械、肌肉和来源 ID 一并保存在 `extra_config`，不丢弃来源信息。未有中文模板的动作保留英文名称和说明。
动作 ID 根据规范动作键排序分配；相同输入生成相同 ID。同一版本重复执行得到相同数据。
动作图片仍使用上游版本固定的 GitHub 图片 URL；本包实际下载到本地的是器械图片。

## 固定肌肉与器械 ID

项目原有肌肉和器械编号保持不变。新增肌肉：19 臀中肌、20 髋内收肌群、21 小腿肌群、22 前臂肌群、23 颈部肌群。
新增器械：43 EZ曲杆、44 通用健身器械、45 药球、46 其他器械。
固定编号不存在时创建；编号已经被其他名称占用或已删除时停止并回滚。

器械 9 按中文定义为腿部推蹬机，10 为腿部伸展机，11 为腿部弯举机，兼容原英文 name 中的历史错配。
通用健身器械、其他器械图片分别用综合训练器和负重背心作代表。
难度 beginner=3、intermediate=5、expert=8 为导入规则。上游肌肉分类与项目解剖肌肉并非严格等价，转换映射及源分类均保留。

## 来源与图片

社区动作：https://github.com/yuhonas/free-exercise-db ，采用 Unlicense。
项目原有中文动作模板属于项目内容，不因合并而变成社区数据的许可证内容。
器械图片来自 Life Fitness、REP Fitness、Bells of Steel、Titan Fitness、Precision Performances 与 Titanium Strength 商城。具体来源与校验值见 `equipment_sources.json`。
商城图片未确认开放授权，公开或商业使用前应确认使用权或换成自有、已授权图片。

## 验证与再生成

```sh
python3 generate.py exercises.json deduplicated_workout_actions.sql --revision "$(cat SOURCE_REVISION)"
python3 validate_import.py /path/to/reference/myapi.db
```

验证程序只读参考库，在内存中测试：962 个规范动作、同义别名、独立变式、同名不同动作、固定引用、JSON、46 张图片校验、默认模板替换、重复执行、有业务数据和 ID 冲突时完整回滚。
以 `deduplicated_workout_actions.sql` 为输出名运行生成器时，会自动同步兼容文件名与组合导入 SQL；规则变更后重新运行即可。
未操作浏览器，请部署后手动检查动作与器械展示。
