#!/usr/bin/python3

import aiohttp
import asyncio
import threading
import time
import os
import random
import binascii
import csv
import requests

class Node (threading.Thread):
   def __init__(self, nodeID, idList):
      threading.Thread.__init__(self)
      self.nodeID = nodeID
      self.idList = idList
      self.url = 'http://localhost:8080/'
      self.oldHeight = 0
      self.height = 0
      self.avgMiningTime = 3.5
      self.miningTimeVariance = 1.5


   def run(self):
      # loop = asyncio.get_event_loop()
      print("Starting: ", self.nodeID)
      self.loop = asyncio.new_event_loop()
      asyncio.set_event_loop(self.loop)
      self.loop.run_until_complete(self.broadcastPowBlock())

   
   async def broadcastPowBlock(self):
      # Get lock to synchronize threads
      async with aiohttp.ClientSession() as session:
         for id in self.idList:
            # requests.get(self.url).json()
            # print(self.powChain[-1]['Index'])


            # sleep to simulate the pow mining time. AVG time = 3.5s, variance = 1.5s
            time.sleep(
               random.randint(-1000 * self.miningTimeVariance, 1000 * self.miningTimeVariance)/1000 
               + self.avgMiningTime)

            chainResponse = requests.get(self.url)
            self.height = chainResponse.json()[-1]['Index']

            threadLock.acquire()
            if self.oldHeight == self.height:
               await session.post(self.url, json = {"powID": id})
               print("NodeID: " + str(self.nodeID) + ", id: " + id)
               # print_time(self.name, self.counter, 3)
               # Free lock to release next thread
            else:
               self.oldHeight = self.height
            threadLock.release()


# if __name__ == "__main__":

totalNodes = 3

d = [(i + 1, []) for i in range(totalNodes)]
nodeIDs = dict(d)

with open('./IDs.csv', mode='w') as csv_file:
   #  fieldnames = ['emp_name', 'dept', 'birth_month']
   writer = csv.writer(csv_file)

   # [random.randrange(0, totalNodes) + 1, binascii.b2a_hex(os.urandom(3))]
   # writer.writeheader()
   for i in range(1000):
      nodeID = random.randrange(0, totalNodes) + 1
      id = str(binascii.b2a_hex(os.urandom(3)))[2:-1]
      # nodeIDs[nodeID].append(id)
      row = [nodeID, id]
      writer.writerow(row)


threadLock = threading.Lock()
threads = []

for node in range(totalNodes):
   threads.append(Node(node + 1, nodeIDs[node + 1]))

for node in range(totalNodes):
   threads[node].start()
   # loop = asyncio.get_event_loop()
   # loop.run_until_complete(threads[node].start())

for t in threads:
   t.join()

# print("Exiting Main Thread")
# res = requests.post('http://localhost:8080/', json = {'powID':'1222'}, timeout=10)

# exit()

