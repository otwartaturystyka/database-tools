import json
import sys
import os

import qrcode

try:
    region = sys.argv[1]
except:
    print("No region specified")
    sys.exit(1)

with open(f'generated/{region}/data.json') as f:
    data = json.load(f)

prefix = 'https://otwartaturystyka.pl/regions/rudnik/places'

gen_dir = f'generated_qrcodes/{region}'

if not os.path.exists(gen_dir):
    os.makedirs(gen_dir)

for section in data['sections']:
    for place in section['places']:
        id = place['id']
        url = f'{prefix}/{id}'
        print(url)
                    
        img = qrcode.make(url)
        img.save(f'{gen_dir}/{id}.png')

