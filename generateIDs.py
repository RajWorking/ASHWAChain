#!/usr/bin/python3

import aiohttp
import asyncio
import threading
import time
import os
import sys
import random
import binascii
import csv
import requests
# import dotenv
# from dotenv import load_dotenv

# load_dotenv()

# MY_ENV_VAR = os.getenv('NODES')

nodes = int(sys.argv[1])
count = int(sys.argv[2])

with open('./IDs.csv', mode='w') as csv_file:
   writer = csv.writer(csv_file)

   for i in range(count):
      for j in range(nodes):
         id = str(binascii.b2a_hex(os.urandom(3)))[2:-1]
         writer.writerow([j + 1, id])

