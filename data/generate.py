import random
import string
import json

class User:
    def __init__(self, firstName, lastName, UID, detachment):
        self.firstname = firstName
        self.lastname = lastName
        self.uid = UID
        self.tokilluid = "" 
        self.iskilled = False
        self.killed = []
        self.detachment = detachment
    def toJson(self):
        return json.dumps(self, default=lambda o: o.__dict__, ensure_ascii=False)

class UserEncoder(json.JSONEncoder):
    def default(self, o):
        return o.__dict__

def get_random_string(length):                                         
    letters = string.ascii_lowercase
    result_str = ''.join(random.choice(letters) for i in range(length))
    return result_str
f = open("parties.xml", "r", encoding='utf-8')

users = []
finalUsers = []
for line in f.readlines():         
    if len(line) < 5:
        continue
    firstName, lastName, detachment = line.split()[0], line.split()[2], line.split()[-1]
    uid = get_random_string(10)
    print(firstName, lastName, uid)
    if len(users) and detachment != users[-1].detachment:
        random.shuffle(users)
        for x in users:
            finalUsers.append(x)
        users = []
    users.append(User(firstName, lastName, uid, detachment))
random.shuffle(users)
for x in users:
    finalUsers.append(x)

fo = open("generated_user_data.json", "w", encoding='utf-8')
fo.write('{"Locker" : {},"knownusers" : [')
for i in range(len(users)):                           
    users[i].tokilluid = users[(i + 1) % len(users)].uid
    fo.write(users[i].toJson() + ',' * (i < len(users) - 1))                                                                               
fo.write("]}")
