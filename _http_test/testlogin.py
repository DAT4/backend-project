import requests
import uuid
import re

url = 'http://localhost:8056/login'
#ip = requests.get('http://httpbin.org/ip').json()['origin'].strip()
#mac = ':'.join(re.findall('..', '%012x' % uuid.getnode()))
data = {
        'username':'martin',
        'password': 'T3stpass!',
}

r = requests.post(url, json=data)

print(r)
print(r.text)
