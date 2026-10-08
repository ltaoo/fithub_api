import concurrent.futures, hashlib, json, sqlite3, subprocess
from pathlib import Path
from urllib.parse import urlsplit,urlunsplit,parse_qsl,urlencode
from generate import equipment,equipment_ids,quote,encoded
out=Path(__file__).parent
photos=out/'media/equipment';photos.mkdir(parents=True,exist_ok=True)
manifest=json.loads((out/'equipment_sources.json').read_text())
assert len(manifest)==46

def download(r):
 u=urlsplit(r['image_url']);q=dict(parse_qsl(u.query));q['width']='800'
 url=urlunsplit((u.scheme,u.netloc,u.path,urlencode(q),u.fragment))
 tmp=photos/(str(r['id'])+'.download')
 subprocess.run(['curl','-fLsS','--retry','2','--max-time','45',url,'-o',str(tmp)],check=True)
 b=tmp.read_bytes()
 if b[:3]==b'\xff\xd8\xff':ext='.jpg'
 elif b[:8]==b'\x89PNG\r\n\x1a\n':ext='.png'
 elif b[:4]==b'RIFF' and b[8:12]==b'WEBP':ext='.webp'
 else:raise ValueError(f'Unexpected image format for {r["id"]}: {b[:30]!r}')
 filename=f'{r["id"]:02d}-equipment{ext}';tmp.rename(photos/filename)
 r.update({'local_path':'media/equipment/'+filename,'public_url':'/media/equipment/'+filename,'bytes':len(b),'sha256':hashlib.sha256(b).hexdigest(),'rights':'Product photo from retailer; no open-license grant verified.'})
 print('Downloaded',r['id'],r['zh_name'],len(b),flush=True)
 return r
with concurrent.futures.ThreadPoolExecutor(max_workers=4) as pool:
 manifest=list(pool.map(download,manifest))
manifest.sort(key=lambda r:r['id'])
(out/'equipment_sources.json').write_text(json.dumps(manifest,ensure_ascii=False,indent=2)+'\n')
sql=['-- SQLite fixed-ID equipment photo import. Copy media/equipment into dist/h5/media/equipment first.','-- Sources and checksums: equipment_sources.json. Images retain their original rights.','BEGIN IMMEDIATE;','CREATE TEMP TABLE fithub_equipment_photo_guard (valid INTEGER NOT NULL CHECK(valid = 1));']
for r in manifest:
 match=f"(lower(trim(name)) = lower({quote(r['name'])}) OR zh_name = {quote(r['zh_name'])}) AND COALESCE(d,0)=0"
 sql.append(f"-- ID {r['id']}: {r['zh_name']} ; {r['product_url']}")
 sql.append(f"INSERT INTO fithub_equipment_photo_guard SELECT CASE WHEN NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id={r['id']} AND NOT ({match})) THEN 1 ELSE 0 END;")
 sql.append(f"INSERT INTO EQUIPMENT (id,name,zh_name,alias,overview,medias,tags,d) SELECT {r['id']},{quote(r['name'])},{quote(r['zh_name'])},'','','{{}}','',0 WHERE NOT EXISTS (SELECT 1 FROM EQUIPMENT WHERE id={r['id']});")
 sql.append(f"UPDATE EQUIPMENT SET medias=json_set(CASE WHEN json_valid(medias) THEN CASE WHEN json_type(medias)='object' THEN medias ELSE '{{}}' END ELSE '{{}}' END,'$.pics',json({encoded([r['public_url']])})) WHERE id={r['id']};")
sql+=['DROP TABLE fithub_equipment_photo_guard;','COMMIT;','']
(out/'equipment_photos.sql').write_text('\n'.join(sql))
print('Generated equipment SQL and source manifest',flush=True)
