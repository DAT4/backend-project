import re
from requests import post, get
import uuid

if __name__ == '__main__':
    #url = 'https://tmp.mama.sh/api/register'
    url = 'http://localhost:8056/register'

    ip = get('http://httpbin.org/ip').json()['origin'].strip()
    print(ip)
    mac = ':'.join(re.findall('..', '%012x' % uuid.getnode()))
    data = {
        'username': 'martin',
        'email': 's111111@student.dtu.dk',
        'password': 'T3stpass!',
        'macs': [mac],
        'ips': [ip],
    }

    r = post(url, json=data)

    print(r)
    print(r.text)
