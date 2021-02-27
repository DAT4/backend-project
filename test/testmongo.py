import pymongo
from os import environ as e

cli = pymongo.MongoClient('mongodb://localhost:27017')
col = cli['backend']['users']


def insert(data):
    global col
    col.insert_one(data)

data = {
        'name':'martin',
        'no':123
        }

colll = (col.find())

for x in colll:
    print(x)


