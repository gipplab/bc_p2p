# This script helps me calculate certain time differences from the log created with rust-libp2p

import argparse
from datetime import datetime

parser = argparse.ArgumentParser(description='calculate time delta from input (delta.py "2020-05-24 17:49:40.018581" "2020-05-24 17:49:57.195546" )')
parser.add_argument('start', type=lambda s: datetime.strptime(s, '%Y-%m-%d %H:%M:%S.%f'), help='Start date') 
parser.add_argument('end', type=lambda s: datetime.strptime(s, '%Y-%m-%d %H:%M:%S.%f'), help='End date') 
args = parser.parse_args()

duration = args.end - args.start

print(duration)
