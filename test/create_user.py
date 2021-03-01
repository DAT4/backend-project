import re
import requests
import uuid

if __name__ == '__main__':
    url = 'https://tmp.mama.sh/api/register'

    ip = requests.get('http://httpbin.org/ip').json()['origin'].strip()
    print(ip)
    mac = ':'.join(re.findall('..', '%012x' % uuid.getnode()))
    data = {
        'username': 'martin',
        'email': 's231321@student.dtu.dk',
        'password': 'T3stpass!',
        'macs': [mac],
        'ips': [ip],
    }

    r = requests.post(url, json=data)

    print(r)
    print(r.text)
