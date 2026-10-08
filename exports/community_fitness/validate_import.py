#!/usr/bin/env python3
"""Validate NEW-project SQL in memory against a read-only reference schema."""
import argparse,hashlib,json,sqlite3
from pathlib import Path


def validate(source_path):
 root=Path(__file__).parent;sql=(root/'import_all.sql').read_text();catalog=json.loads((root/'deduplicated_actions.json').read_text())
 source=sqlite3.connect(source_path.resolve().as_uri()+'?mode=ro',uri=True)
 schemas=[r[0] for r in source.execute("SELECT sql FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND sql IS NOT NULL")]
 def empty_schema():
  c=sqlite3.connect(':memory:')
  for schema in schemas:c.execute(schema)
  return c
 def assert_catalog(c):
  c.row_factory=sqlite3.Row
  actual=[dict(r) for r in c.execute('SELECT * FROM WORKOUT_ACTION ORDER BY id')]
  assert len(actual)==len(catalog)
  by_key={}
  for r in actual:
   for field in ['details','points','problems','extra_config']:json.loads(r[field])
   key=json.loads(r['extra_config'])['canonical_key'];assert key not in by_key;by_key[key]=r
   for field,table in [('equipment_ids','EQUIPMENT'),('muscle_ids','MUSCLE'),('primary_muscle_ids','MUSCLE'),('secondary_muscle_ids','MUSCLE'),('alternative_action_ids','WORKOUT_ACTION'),('advanced_action_ids','WORKOUT_ACTION'),('regressed_action_ids','WORKOUT_ACTION')]:
    for i in filter(None,r[field].split(',')):assert c.execute('SELECT 1 FROM '+table+' WHERE id=? AND COALESCE(d,0)=0',(int(i),)).fetchone()
  bench=by_key['barbell-flat-bench-press']
  for name in ['杠铃卧推','杠铃平板卧推','平板卧推']:
   assert name in [bench['name'],bench['zh_name']]+bench['alias'].split(',')
  # Important counterexamples must remain independent movements.
  source_lookup={i:r['id'] for r in actual for i in json.loads(r['extra_config'])['source_ids']}
  pairs=[('Barbell_Bench_Press_-_Medium_Grip','Barbell_Incline_Bench_Press_-_Medium_Grip'),('Barbell_Bench_Press_-_Medium_Grip','Close-Grip_Barbell_Bench_Press'),('Barbell_Bench_Press_-_Medium_Grip','Decline_Barbell_Bench_Press'),('Barbell_Bench_Press_-_Medium_Grip','Wide-Grip_Barbell_Bench_Press'),('Dumbbell_Shoulder_Press','Standing_Dumbbell_Press'),('Barbell_Curl','EZ-Bar_Curl')]
  for a,b in pairs:assert source_lookup[a]!=source_lookup[b],(a,b)
  assert source_lookup['Dumbbell_Shoulder_Press']==source_lookup['Seated_Dumbbell_Press']
  assert source_lookup['Bent_Over_One-Arm_Long_Bar_Row']==source_lookup['One-Arm_Long_Bar_Row']
  assert by_key['upper-back-interlaced-hands-stretch']['id']!=by_key['upper-back-crossed-arms-stretch']['id']
  for r in json.loads((root/'equipment_sources.json').read_text()):
   assert hashlib.sha256((root/r['local_path']).read_bytes()).hexdigest()==r['sha256']
   assert json.loads(c.execute('SELECT medias FROM EQUIPMENT WHERE id=?',(r['id'],)).fetchone()[0])['pics']==[r['public_url']]
  return actual
 c=empty_schema();c.executescript(sql);actual=assert_catalog(c)
 snapshot={t:[tuple(r) for r in c.execute('SELECT * FROM '+t+' ORDER BY id')] for t in ['WORKOUT_ACTION','MUSCLE','EQUIPMENT']}
 c.executescript(sql)
 for t in snapshot:assert [tuple(r) for r in c.execute('SELECT * FROM '+t+' ORDER BY id')]==snapshot[t],t
 # Stock project action templates must be replaced, rather than left next to canonical movements.
 seeded=empty_schema()
 for r in json.loads((root/'existing_action_catalog.json').read_text()):
  fields=list(r);seeded.execute('INSERT INTO WORKOUT_ACTION ('+','.join(fields)+') VALUES ('+','.join('?' for _ in fields)+')',[r[k] for k in fields])
 seeded.commit();seeded.executescript(sql);assert_catalog(seeded)
 # Refuse any environment with existing plans/history/custom actions; no partial writes.
 for setup in ["INSERT INTO WORKOUT_PLAN (title) VALUES ('existing-plan')", "UPDATE WORKOUT_ACTION SET owner_id=123 WHERE id=1", "UPDATE EQUIPMENT SET name='other',zh_name='other' WHERE id=1"]:
  blocked=empty_schema();blocked.executescript(sql);blocked.execute(setup);blocked.commit()
  before={t:blocked.execute('SELECT * FROM '+t+' ORDER BY id').fetchall() for t in ['WORKOUT_ACTION','MUSCLE','EQUIPMENT','WORKOUT_PLAN']}
  try:blocked.executescript(sql)
  except sqlite3.IntegrityError:blocked.rollback()
  else:raise AssertionError('Guard should reject: '+setup)
  for t in before:assert blocked.execute('SELECT * FROM '+t+' ORDER BY id').fetchall()==before[t],t
  blocked.close()
 # Existing deployed schema with actual plans is used only to verify refusal in a backup.
 reference=sqlite3.connect(':memory:');source.backup(reference)
 if reference.execute('SELECT count(*) FROM WORKOUT_PLAN').fetchone()[0]:
  before=reference.execute('SELECT * FROM WORKOUT_ACTION ORDER BY id').fetchall()
  try:reference.executescript(sql)
  except sqlite3.IntegrityError:reference.rollback()
  else:raise AssertionError('Existing application should be rejected')
  assert reference.execute('SELECT * FROM WORKOUT_ACTION ORDER BY id').fetchall()==before
 report={'canonical_actions':len(actual),'semantic_examples':'bench synonyms merged; seated shoulder synonyms merged; landmine row synonyms merged','distinct_variants':'incline/decline/close/wide-grip/standing/EZ vs straight preserved','same_name_different_movement':'upper-back stretch variants preserved','seed_catalog_replaced':'passed','repeat_import':'same records and IDs','nonempty_project_guard':'passed, rollback verified','fixed_id_conflict':'passed, rollback verified','images_and_json_and_references':'passed','actual_database':'read-only, never imported'}
 print(json.dumps(report,ensure_ascii=False,indent=2))
 source.close();reference.close();c.close();seeded.close();return report
if __name__=='__main__':
 p=argparse.ArgumentParser();p.add_argument('database',type=Path);a=p.parse_args();validate(a.database)
