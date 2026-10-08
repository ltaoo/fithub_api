#!/usr/bin/env python3
"""Build a reviewed, semantically deduplicated SQLite action catalog for a NEW project."""
import argparse, collections, datetime, hashlib, json
from pathlib import Path
from semantic_rules import GROUPS,key_for,normalize
from import_mapping import muscles,equipment,muscle_ids,equipment_ids,types,levels,quote,encoded

GUARD_TABLES=['WORKOUT_PLAN','WORKOUT_PLAN_ACTION','WORKOUT_DAY','WORKOUT_DAY_ACTION','WORKOUT_DAY_STEP','WORKOUT_ACTION_HISTORY','COACH_CONTENT_WITH_WORKOUT_ACTION','COACH_PHYSICAL_TEST','USER_FAVORITE','COACH_WORKOUT_SCHEDULE']
GROUP_LOOKUP={key:(aliases,note) for key,aliases,note in GROUPS}


def source_row(r,revision):
 eq=r.get('equipment');images=['https://raw.githubusercontent.com/yuhonas/free-exercise-db/'+revision+'/exercises/'+i for i in r.get('images',[])]
 ids=lambda values:','.join(str(i) for i in sorted({muscle_ids[x] for x in values}))
 return {'name':r['name'],'zh_name':r['name'],'alias':'','overview':'','type':types[r['category']],'level':levels[r['level']],'tags1':','.join(r['primaryMuscles']),'tags2':','.join(filter(None,['community',r.get('force'),r.get('mechanic')])),'details':json.dumps({'start_position':'','steps':r['instructions']},ensure_ascii=False),'points':'[]','problems':'[]','equipment_ids':str(equipment_ids[eq]) if eq in equipment else '', 'muscle_ids':ids(r['primaryMuscles']+r['secondaryMuscles']),'primary_muscle_ids':ids(r['primaryMuscles']),'secondary_muscle_ids':ids(r['secondaryMuscles']),'alternative_action_ids':'','advanced_action_ids':'','regressed_action_ids':'','score':0,'sort_idx':0,'pattern':'','cover_url':images[0] if images else ''}


def richness(r):
 return sum(len(str(r.get(k) or '')) for k in ['details','points','problems','overview'])


def make_catalog(source, revision):
 root=Path(__file__).parent
 community=json.loads(source.read_text());legacy=json.loads((root/'existing_action_catalog.json').read_text())
 groups={}
 for r in legacy:groups.setdefault(key_for(r),{'legacy':[],'community':[]})['legacy'].append(r)
 for r in community:groups.setdefault(key_for(r,True),{'legacy':[],'community':[]})['community'].append(r)
 rows=[];report=[]
 for index,(key,g) in enumerate(sorted(groups.items()),1):
  preferred=max(g['legacy'],key=richness).copy() if g['legacy'] else source_row(g['community'][0],revision)
  aliases=set()
  for r in g['legacy']:
   aliases.update(filter(None,[r['name'],r['zh_name']]))
   aliases.update(x.strip() for x in (r.get('alias') or '').replace('，',',').split(',') if x.strip())
  for r in g['community']:aliases.add(r['name'])
  note='同名或纯拼写格式差异，保留动作限定词'
  if key in GROUP_LOOKUP:
   reviewed,note=GROUP_LOOKUP[key];aliases.update(reviewed)
   english=[s for s in reviewed if s.isascii()]
   chinese=[s for s in reviewed if not s.isascii()]
   if english:preferred['name']=english[0]
   if chinese:preferred['zh_name']=chinese[0]
  aliases.discard(preferred['name']);aliases.discard(preferred['zh_name'])
  # The app splits aliases on commas, so comma-containing names receive a space instead.
  preferred['alias']=','.join(sorted({s.replace(',',' ') for s in aliases},key=str.casefold))
  extra={'canonical_key':key,'semantic_dedup_version':1,'source':'yuhonas/free-exercise-db' if g['community'] else 'project-action-catalog','source_ids':[r['id'] for r in g['community']],'source_revision':revision,'license':'project content + Unlicense community metadata' if g['legacy'] and g['community'] else 'Unlicense' if g['community'] else 'project-owned','original_project_ids':[r['id'] for r in g['legacy']]}
  if g['community']:
   extra['community_variants']=[{k:r.get(k) for k in ['id','name','equipment','category','level','force','mechanic','primaryMuscles','secondaryMuscles','instructions','images']} for r in g['community']]
   extra['images']=['https://raw.githubusercontent.com/yuhonas/free-exercise-db/'+revision+'/exercises/'+image for r in g['community'] for image in r.get('images',[])]
   if not preferred.get('cover_url'):preferred['cover_url']=extra['images'][0] if extra['images'] else ''
  preferred.update({'id':index,'status':1,'owner_id':0,'d':0,'extra_config':json.dumps(extra,ensure_ascii=False,separators=(',',':')),'created_at':preferred.get('created_at') or '2026-10-09 00:00:00','updated_at':None})
  rows.append(preferred)
  if len(g['legacy'])+len(g['community'])>1:
   report.append({'canonical_key':key,'canonical_id':index,'name':preferred['name'],'zh_name':preferred['zh_name'],'reason':note,'members':[{'source':'project','id':r['id'],'name':r['name'],'zh_name':r['zh_name']} for r in g['legacy']]+[{'source':'community','id':r['id'],'name':r['name']} for r in g['community']]})
 old_id_map={old['id']:r['id'] for r in rows for old in groups[json.loads(r['extra_config'])['canonical_key']]['legacy']}
 for r in rows:
  for field in ['alternative_action_ids','advanced_action_ids','regressed_action_ids']:
   r[field]=','.join(str(i) for i in sorted({old_id_map[int(v)] for v in (r.get(field) or '').split(',') if v.strip() and int(v) in old_id_map and old_id_map[int(v)]!=r['id']}))
 normalized_unique=len({normalize(r['name']) for r in legacy+community})
 summary={'project_source_records':len(legacy),'community_source_records':len(community),'input_records':len(legacy)+len(community),'unique_movements':len(rows),'collapsed_records':len(legacy)+len(community)-len(rows),'normalized_name_unique':normalized_unique,'reviewed_semantic_rules':len(GROUPS),'groups_with_multiple_records':len(report),'policy':'explicit reviewed aliases and instructions; no fuzzy or muscle-only automatic merges'}
 (root/'semantic_dedup_report.json').write_text(json.dumps({'summary':summary,'merged_groups':report,'preserved_examples':['上斜/下斜/窄握/宽握卧推','站姿/坐姿哑铃肩推','直杠/曲杠弯举','相扑/传统/罗马尼亚/直腿硬拉','屈膝/直腿凳上臂屈伸','上背环抱手臂/手指相扣拉伸']},ensure_ascii=False,indent=2)+'\n')
 (root/'deduplicated_actions.json').write_text(json.dumps(rows,ensure_ascii=False,indent=2)+'\n')
 return rows,summary


def generate(source, output, revision='main'):
 root=Path(__file__).parent;rows,summary=make_catalog(source,revision)
 lines=['-- FitHub semantic action catalog for a NEW project. SQLite only.', '-- Source: https://github.com/yuhonas/free-exercise-db ; community license: Unlicense.', '-- Project templates are user-owned data; see semantic_dedup_report.json.', '-- Revision: '+revision,'-- Source SHA256: '+hashlib.sha256(source.read_bytes()).hexdigest(), '-- IMPORTANT: this REPLACES WORKOUT_ACTION seed data and resets action IDs.', '-- Refuses to run if plans, training data, favorites, custom actions, etc. exist.', '-- Run sqlite3 -bail; errors must stop and roll back the transaction.', 'BEGIN IMMEDIATE;', 'CREATE TEMP TABLE fithub_fresh_catalog_guard (valid INTEGER NOT NULL CHECK(valid = 1));']
 for t in GUARD_TABLES:
  lines.append('-- Must be empty: '+t)
  lines.append(f'INSERT INTO fithub_fresh_catalog_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM {t}) THEN 1 ELSE 0 END;')
 lines.append('INSERT INTO fithub_fresh_catalog_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM WORKOUT_ACTION WHERE owner_id <> 0 OR status <> 1) THEN 1 ELSE 0 END;')
 lines.append("INSERT INTO fithub_fresh_catalog_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM WORKOUT_ACTION_MISTAKE WHERE trim(solution_action_ids) <> '') THEN 1 ELSE 0 END;")
 entities=json.loads((root/'fixed_entities.json').read_text())
 for key,(name,zh) in muscles.items():
  if not any(r['id']==muscle_ids[key] for r in entities['MUSCLE']):entities['MUSCLE'].append({'id':muscle_ids[key],'name':name,'zh_name':zh})
 for key,(name,zh) in equipment.items():
  if not any(r['id']==equipment_ids[key] for r in entities['EQUIPMENT']):entities['EQUIPMENT'].append({'id':equipment_ids[key],'name':name,'zh_name':zh})
 for table in ['MUSCLE','EQUIPMENT']:
  for r in sorted(entities[table],key=lambda v:v['id']):
   check=f"(lower(trim(name))=lower({quote(r['name'])}) OR zh_name={quote(r['zh_name'])}) AND COALESCE(d,0)=0"
   lines.append(f"INSERT INTO fithub_fresh_catalog_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM {table} WHERE id={r['id']} AND NOT ({check})) THEN 1 ELSE 0 END;")
   columns=['id','name','zh_name','overview','tags','d','features'] if table=='MUSCLE' else ['id','name','zh_name','alias','overview','medias','tags','d']
   values=[quote(r.get(k) or ('[]' if k=='features' else '{}' if k=='medias' else '')) if k not in ['id','d'] else str(r['id'] if k=='id' else 0) for k in columns]
   lines.append('INSERT INTO '+table+' ('+','.join(columns)+') SELECT '+','.join(values)+f' WHERE NOT EXISTS (SELECT 1 FROM {table} WHERE id={r["id"]});')
 lines+=['DELETE FROM WORKOUT_ACTION;',"DELETE FROM sqlite_sequence WHERE name='WORKOUT_ACTION';"]
 columns=['id','status','name','zh_name','alias','overview','type','level','tags1','tags2','details','points','problems','equipment_ids','muscle_ids','primary_muscle_ids','secondary_muscle_ids','alternative_action_ids','advanced_action_ids','regressed_action_ids','owner_id','score','extra_config','sort_idx','pattern','d','cover_url','created_at','updated_at']
 numeric={'id','status','level','owner_id','score','sort_idx','d'}
 for r in rows:
  vals=['NULL' if r.get(k) is None and k=='updated_at' else str(r.get(k) or 0) if k in numeric else quote(r.get(k) or '') for k in columns]
  lines.append('INSERT INTO WORKOUT_ACTION ('+','.join(columns)+') VALUES ('+','.join(vals)+');')
 lines+=['DROP TABLE fithub_fresh_catalog_guard;','COMMIT;','']
 output.write_text('\n'.join(lines))
 if output.name=='deduplicated_workout_actions.sql':
  (root/'community_workout_actions.sql').write_text(output.read_text())
  image_sql=(root/'equipment_photos.sql').read_text()
  action_sql=output.read_text()
  combined=action_sql.replace('BEGIN IMMEDIATE;\n','',1).replace('COMMIT;\n','',1)+image_sql.replace('BEGIN IMMEDIATE;\n','',1).replace('COMMIT;\n','',1)
  (root/'import_all.sql').write_text('-- NEW PROJECT ONLY: replace action seed catalog with unique movements, then attach equipment images.\nBEGIN IMMEDIATE;\n'+combined+'COMMIT;\n')
 print(json.dumps(summary,ensure_ascii=False,indent=2))
if __name__=='__main__':
 p=argparse.ArgumentParser();p.add_argument('source',type=Path);p.add_argument('output',type=Path);p.add_argument('--revision',default='main');a=p.parse_args();generate(a.source,a.output,a.revision)
