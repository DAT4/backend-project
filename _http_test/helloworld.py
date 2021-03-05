import requests,json, uuid, re

def createUser():
    url = 'https://api.backend.mama.sh/user'
    ip = requests.get('http://httpbin.org/ip').json()['origin']
    mac = ':'.join(re.findall('..', '%012x' % uuid.getnode()))
    user = {
            'email': 'newuser@emil.sh',
            'password': 'lalalalal',
            'name': 'Peter Pan',
            'ip': ip,
            'mac': mac,
            }

    res = requests.post(url,json=user)
    print(res.text)

def helloWorld():
    url = 'https://api.backend.mama.sh/'
    msg = requests.get(url)
    print(msg.text)

def listUsers():
    url = 'https://api.backend.mama.sh/user'
    msg = requests.get(url)
    print(msg.text)

createUser()
listUsers()
