import re
from requests import post, get
import uuid

if __name__ == '__main__':
    url = 'https://backend.mama.sh/register'

    ip = get('http://httpbin.org/ip').json()['origin'].strip()
    print(ip)
    mac = ':'.join(re.findall('..', '%012x' % uuid.getnode()))
    data = {
        'username': 'martin',
        'email': 'hej@hej.hej',
        'password': 'T3stpass!',
    }

    r = post(url, json=data)

    print(r)
    print(r.text)
