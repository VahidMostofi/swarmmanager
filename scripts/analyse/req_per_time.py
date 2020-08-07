import zipfile
import json
import io
import numpy as np
from tqdm import tqdm

with zipfile.ZipFile('/home/vahid/Dropbox/data/swarm-manager-data/jaegers/9b5d17d0-fbd4-4f90-7064-90628a6f50c9.zip') as zf:
    with io.TextIOWrapper(zf.open("jaeger-info.json"), encoding="utf-8") as f:
        data = json.load(f)

for d in data['data']:
    spans = d['spans']
    for span in spans:
        span['endTime'] = int(np.round((span['startTime'] + span['duration'])/10000,0))
        span['startTime'] = int(np.round(span['startTime']/10000,0))

def strip_zeros(arr):
    j = 0
    while arr[j] == 0 and j < len(arr) - 1:
        j += 1
    start = j
    j = len(arr) - 1
    while arr[j] == 0 and j >= 0:
        j -= 1
    end = j
    return arr[start: end]

for serviceName in ["auth", "books"]:
    print("working on", serviceName)
    traces = []
    for d in data['data']:
        spans = d['spans']
        if len(spans) < 10:
            continue
        flag = False
        spans_dict = {}
        for span in spans:
            if span['operationName'] == serviceName:
                flag = True
            spans_dict[span['operationName']] = span
        if flag:
            traces.append(spans_dict)
        
    print('len(traces)',len(traces))
    min_start = 10000000000000000000000000
    max_end = 0
    for trace in traces:
        for opName,span in trace.items():
            if span["startTime"] < min_start:
                min_start = span["startTime"]
            if span["endTime"] > max_end:
                max_end = span["startTime"]
    times = {
        "gateway": [0] * (max_end-min_start+50),
        "total": [0] * (max_end-min_start+50),
        "service": [0] * (max_end-min_start+50)
    }
    for trace in traces:
        for j in range(trace['backend']['startTime'], trace['backend']['endTime']):
            j = j - min_start
            if j >= len(times['total']):
                print(serviceName, 'total', j)
            times['total'][j] += 1
    print('total    ',np.round(np.mean(strip_zeros(times['total'])),2))

    for trace in traces:
        for j in range(trace[serviceName+'_connect']['endTime'], trace[serviceName]['endTime']):
            j = j - min_start
            if j >= len(times['service']):
                print(serviceName, 'service', j)
            times['service'][j] += 1
    print('service  ',np.round(np.mean(strip_zeros(times['service'])),2))

    for trace in traces:
        for j in range(trace['backend']['startTime'], trace['backend']['endTime']):
            if j < trace[serviceName+'_connect']['endTime'] or j > trace[serviceName]['endTime']:
                continue
            j = j - min_start
            if j >= len(times['gateway']):
                print(serviceName, 'gateway', j)
            times['gateway'][j] += 1
    print('gateway  ',np.round(np.mean(strip_zeros(times['gateway'])),2))
    print('==================================')