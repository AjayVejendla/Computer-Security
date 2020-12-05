import json
import os.path as path
import sys

#Declare dictionaries
#Will be written if file does not already exists
#Will be loaded too
users = {}
domains = {}
types = {}

argv = sys.argv

#Pad array with empty strings to account for missing command line args
for x in range(len(argv),5):
    argv.append("")

#Remove portal.py from args
argv = argv[1:]



def AddUser(args):
    user = args[0]
    password = args[1]

    if not user:
        print("Error: username missing")
    
    elif user in users:
        print("Error: user exists")

    else:
        users[user] = password
        print('Success')

def Authenticate(args):
    user = args[0]
    password = args[1]

    if user not in users:
        print('Error: no such user')
    elif users[user] == password:
        authenticatedUser = user
        print('success')
    else:
        print('Error: bad password')

def SetDomain(args):
    user = args[0]
    domain = args[1]
    if user not in users:
        print("Error: no such user")
    elif not domain:
        print("Error: missing domain")
    elif domain not in domains:
        domains[domain] =[user]
        print('Success')
    else:
        domainEntries = domains[domain]

        if user not in domainEntries:
            domainEntries.append(user)
            domains[domain] = domainEntries

        print('Success')

def DomainInfo(args):
    domain = args[0]
    if (domain not in domains):
        print("Error: missing domain")
    else:

        domainEntries = domains[domain]
        for user in domainEntries:
            print(user)
           
def SetType(args):
    objectName = args[0]
    type = args[1]

    if not objectName:
        print("Error: null object")

    elif not type:
        print("Error: null type")
    elif type not in types:
        types[type] = [[objectName],{}]
        print('Success')
    else:
        typeEntries = types[type][0]

        if objectName not in typeEntries:
            typeEntries.append(objectName)
            types[type][0] = typeEntries

        print('Success')

def TypeInfo(args):
    type = args[0]
    if (type not in types):
        print("Error: missing type")
    else:

        typeEntries = types[type][0]
        for objects in typeEntries:
            print(objects)

def AddAccessHelper(args):
    print(AddAccess(args)) 
def AddAccess(args):
    operation = args[0]
    domainName = args[1]
    typeName = args[2]
    
    if not operation:
        return("Error: missing operation")
    elif not domainName:
        return("Error: missing domain")
    elif not typeName:
        return("Error: missing type")
    elif typeName not in types:
        types[typeNameu] = [[],{}]
    elif domainName not in domains:
        domains[domainName] =[]
    
    permissions = types[typeName][1]
    if operation not in permissions:
        permissions[operation] = [domainName]
    if domainName not in permissions[operation]:
        permissions[operation].append(domainName)

    return "Success"

def CanAccessHelper(args):
    print(CanAccess(args))
def CanAccess(args):
    operation = args[0]
    user = args[1]
    object = args[2]

    
    operations = findObject(object)

    if len(operations) == 0:
        return 'Error: access denied'
    
    domainList = []
    
    #Load domains with every domain that can perform the giving operation
    for dictionary in operations:
        for key in dictionary:
            if operation == key:
                for domain in dictionary[key]:
                    domainList.append(domain)
    
    #If no domains can perform operation return error

    if len(domains) == 0:
        return 'Error: access denied'
    
    #Check if user is in any of the domains

    for domain in domainList:
        if user in domains[domain]:
            return 'Success'
    
    return 'Error: access denied'

#Helper method to keep code cleaner
def findObject(object):
    operations = []
    print(object)
    for key in types:
        print(types[key][0])
        if object in types[key][0]:
            operations.append(types[key][1])
    
    return operations

def NoFunction(args):
    print('Error: no valid operation')


#Check if file exists, if it doesn't create the file with empty dictionaries users, domains, types
if not path.exists('tables.json'):
    with open('tables.json','w') as tableFile:
        json.dump({"users": users, "domains": domains, "types":types},tableFile)

#Open file with read and write permissions
with open('tables.json','r+') as tableFile:
    data = json.load(tableFile)
    users = data['users']
    domains = data['domains']
    types = data['types']

#Dictionary mapping command line string to appropriate function
functions = {

    "AddUser":AddUser,
    "Authenticate":Authenticate,
    "SetDomain":SetDomain,
    "DomainInfo":DomainInfo,
    "SetType":SetType,
    "TypeInfo":TypeInfo,
    "AddAccess":AddAccessHelper,
    "CanAccess":CanAccessHelper,
}

#Get a ppropriate function from dictionary
if argv[0] in functions:
    function = functions[argv[0]]
else:
    function = NoFunction

#remove function from argv
args = argv[1:]

#call function with args
function(args)

#Write data to file
with open('tables.json','w') as tableFile:
    json.dump({"users": users, "domains": domains, "types":types},tableFile,indent = 4, sort_keys=True)

