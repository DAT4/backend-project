import requests
import uuid
import re
import json

url = 'http://localhost:1001/login'
#url = 'https://api.backend.mama.sh/login'
#ip = requests.get('http://httpbin.org/ip').json()['origin'].strip()
#mac = ':'.join(re.findall('..', '%012x' % uuid.getnode()))
data = {
        'username':'martin',
        'password': 'T3stpass!',
}

r = requests.post(url, json=data)

print(r)

if r.status_code == 200:
    with open('localstorage.json', 'w') as f:
        json.dump(f, r.json())

print(r.text)
