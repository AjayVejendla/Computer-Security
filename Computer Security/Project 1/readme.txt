Ajay Vejendla
178007104

Note that everything is case sensitive so that the calling program can decide whether or not they want to store data as case sensitive.
In other words, the calling program should convert to all lowercase or all upercase when storing data through portal if it doesn't want case sensitivity

The json file being used to store data will be created/looked for in same directory as the portal.py. Subdirectories will not be checked.

Program must use python3. Json parsing in python2 loads strings in arrays as unicode but python3 loads them as STR. Additional parsing would be required to use python2.

Program requires json, os.path, and sys, but these are already included in python3

Python treats null strings and '' identically for equality, so anything that needs to be checked for null values will also be considered as null when given an empty string. For example, you cannot have a user named ''. However a space is acceptable, so ' ' should be valid.

tables.json needs to be in a json format conforming to my file structures or it needs to not exist. Having an already existing tables.json that is incorrectly formatted will produce an error.
Json files are indented to make reading them easier.
My data is structured as follows:

users{'user':password}
domains{'domain:[user1,user2,...userN]}
types{'type':
            [
            [object1,object2....objectN],
            {
                'permission':
                    [domain1,domain2...domainN]    
            }
            ]


Users and domains should be easily readable. Types contains an array with an array and a dictionary. The inner array contains the objects in that type. The dictionary maps a permission or operation as the key to an array of domains that can perform that operation/has that permission.