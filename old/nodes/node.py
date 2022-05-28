#!/usr/bin/python3

import time
import sys
import random
import csv
import requests
import logging 

ID = int(sys.argv[1])

logfile = "../logs/" + str(ID) + ".log"
logging.basicConfig(filename = logfile, 
					format = '%(asctime)s %(message)s', 
					filemode = 'w')

logger = logging.getLogger() 

logger.setLevel(logging.INFO)
# logger.setLevel(logging.DEBUG)

idList = []
url = 'http://localhost:8080/'

miningTimeVariance = 2
avgMiningTime = 6

with open('../logs/IDs.csv', 'r') as f:
    csv_reader = csv.reader(f, delimiter=',')
    for row in csv_reader:
        if int(row[0]) == ID:
            idList.append(row[1])

# print(ID)

def broadcastBlock():
    idCount = 0
    while idList:
        currentID = idList[-1]
        
        # get the latest block index as old height
        oldHeight = requests.get(url).json()[-1]['Index']

        # sleep to simulate the pow mining time
        time.sleep(
            random.randint(-1000 * miningTimeVariance, 1000 * miningTimeVariance)/1000 
            + avgMiningTime)

        height = requests.get(url).json()[-1]['Index']

        # if no new block is added then post this block and select the next identity
        if oldHeight == height:
            requests.post(url, json = {"powID": currentID})
            idCount += 1
            logger.info("ID count " + str(idCount) + ", id: " + currentID)
            idList.pop()
            chainResponse = requests.get(url)
            height = chainResponse.json()[-1]['Index']
            oldHeight = height
        # else, start mining again
        else:
            oldHeight = height


broadcastBlock()
