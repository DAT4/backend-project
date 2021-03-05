import requests
import uuid
import re

url = 'https://tmp.mama.sh/api/login'

ip = requests.get('http://httpbin.org/ip').json()['origin'].strip()
print(ip)
mac = ':'.join(re.findall('..', '%012x' % uuid.getnode()))
data = {
        'username':'martin',
        'password': 'T3stpass!',
        }

r = requests.post(url, json=data)

print(r)
print(r.text)
