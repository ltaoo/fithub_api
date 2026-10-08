#!/usr/bin/env python3
"""Fixed FitHub entity mappings and SQL quoting helpers."""
import argparse, hashlib, json
from pathlib import Path

muscles = {
 'abdominals': ('Abdominals','腹直肌'), 'abductors': ('Gluteus medius','臀中肌'),
 'adductors': ('Hip adductors','髋内收肌群'), 'biceps': ('Biceps brachii','肱二头肌'),
 'calves': ('Calf muscles','小腿肌群'), 'chest': ('Pectoralis major','胸大肌'),
 'forearms': ('Forearm muscles','前臂肌群'), 'glutes': ('Gluteus maximus','臀大肌'),
 'hamstrings': ('Hamstrings','股二头肌'), 'lats': ('Latissimus dorsi','背阔肌'),
 'lower back': ('Erector spinae','竖脊肌'), 'middle back': ('Rhomboids','菱形肌'),
 'neck': ('Neck muscles','颈部肌群'), 'quadriceps': ('Quadriceps','股四头肌'),
 'shoulders': ('Deltoid','三角肌'), 'traps': ('Trapezius','斜方肌'),
 'triceps': ('Triceps brachii','肱三头肌'),
}
equipment = {
 'bands': ('resistance_band','弹力带'), 'barbell': ('barbell','杠铃'),
 'cable': ('cable_machine','龙门架'), 'dumbbell': ('dumbbell','哑铃'),
 'e-z curl bar': ('ez_curl_bar','EZ曲杆'), 'exercise ball': ('yoga_ball','瑜伽球'),
 'foam roll': ('foam_roller','泡沫轴'), 'kettlebells': ('kettlebell','壶铃'),
 'machine': ('machine','通用健身器械'), 'medicine ball': ('medicine_ball','药球'),
 'other': ('other','其他器械'),
}
types = {'cardio':'cardio','stretching':'static_stretch','plyometrics':'performance',
 'strength':'resistance','powerlifting':'resistance','olympic weightlifting':'resistance','strongman':'comprehensive'}
levels = {'beginner':3,'intermediate':5,'expert':8}
# Fixed IDs from the project's seeded database; supplemental IDs are explicitly reserved.
muscle_ids = {'abdominals':12,'abductors':19,'adductors':20,'biceps':2,'calves':21,'chest':6,'forearms':22,'glutes':8,'hamstrings':10,'lats':7,'lower back':16,'middle back':15,'neck':23,'quadriceps':9,'shoulders':4,'traps':5,'triceps':3}
equipment_ids = {'bands':17,'barbell':2,'cable':13,'dumbbell':1,'e-z curl bar':43,'exercise ball':19,'foam roll':18,'kettlebells':16,'machine':44,'medicine ball':45,'other':46}

def quote(value):
 return "'" + str(value).replace("'", "''") + "'"

def encoded(value):
 return quote(json.dumps(value,ensure_ascii=False,separators=(',',':')))
