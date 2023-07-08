# Python script that migrates our data from the old database to the new 
# append only logfile backed database.

import sqlite3
import sys
import os
import json

# Takes as argument the name of the old database
arg = sys.argv[1]

# Call out to the shell to backup the old database
if os.system('cp ' + arg + ' ' + arg + '.bak') != 0:
    print('Error backing up the old database')
    sys.exit(1)
# delete the old database
if os.system('rm ' + arg) != 0:
    print('Error deleting the old database')
    sys.exit(1)

# Open the old database
old_db = sqlite3.connect(arg + '.bak')

# Open the new logfile backed database
log = open(arg, 'w')

# Select all talks from the old database
talks = old_db.execute('SELECT * FROM talks')

# map talk type to string
talk_type_map = {
    0: 'forum topic',
    1: 'lightning talk',
    2: 'project update',
    3: 'announcement',
    4: 'after meeting slot'
}

# Write all talks to the new database
for talk in talks:
    print(talk)

    # Create the talk event
    talk_event = {
        'time': talk[7],
        'type': 'create',
        'create': {
            'id': talk[0],
            'name': talk[1],
            'type': talk_type_map[talk[2]],
            'description': talk[3],
            'week': talk[5]
        }
    }

    # Write the talk using the JSON format (without extra whitespace)
    text = json.dumps(talk_event, separators=(',', ':'))
    log.write(text + '\n')

    # If the talk is hidden, write the hide event
    if talk[4] == 1:
        hide_event = {
            'time': talk[7],
            'type': 'hide',
            'hide': {
                'id': talk[0]
            }
        }
        text = json.dumps(hide_event, separators=(',', ':'))
        log.write(text + '\n')
