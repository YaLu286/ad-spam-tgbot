#!/Users/yalu/Developer/GO_projects/ad-spam-tgbot/src/py_sender/bin/python3

import configparser
import json

from telethon.sync import TelegramClient
from telethon import connection

from datetime import date, datetime

# from telethon.tl.functions.channels import GetParticipantsRequest
# from telethon.tl.types import ChannelParticipantsSearch

# from telethon.tl.functions.messages import GetHistoryRequest

import psycopg2

config = configparser.ConfigParser()
config.read("/Users/yalu/Developer/GO_projects/ad-spam-tgbot/src/py_sender/config.ini")

api_id   = config['Telegram']['api_id']
api_hash = config['Telegram']['api_hash']
username = config['Telegram']['username']

client = TelegramClient(username, api_id, api_hash)

client.start()

try:
    conn = psycopg2.connect(dbname='postgres', user='postgres', password='123', host='localhost')
    cursor = conn.cursor()
    cursor.execute("SELECT * from messages")
    msg = cursor.fetchall()
    # print(msg[0][0])

    cursor.execute("SELECT * FROM recievers;")
    recievers = cursor.fetchall()

    # print(recievers)

    for i in range(len(recievers)) :
        id, name, _ = recievers[i]
        # print(msg)
        # client.send_message(name, msg[0][0])
        # print(id, name)
        
    cursor.close()
    conn.close()

    # dialogs = client.get_dialogs()

    # for j in range(len(dialogs)) :
    #     client.send_message(dialogs[j], str(j))
    
    
except:
    print('Can`t establish connection to database')



