import requests
import uuid
import re

url = 'http://localhost:8055/register'

ip = requests.get('http://httpbin.org/ip').json()['origin'].strip()
print(ip)
mac = ':'.join(re.findall('..', '%012x' % uuid.getnode()))
data = {
        'username':'martinimain',
        'email':'s195469@student.dtu.dk',
        'password': 'HeJ123!',
        'macs': [mac],
        'ips': [ip],
        }

r = requests.post(url, json=data)

print(r)
print(r.text)
