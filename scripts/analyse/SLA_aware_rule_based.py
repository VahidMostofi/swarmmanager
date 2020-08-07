import sys
import yaml
from itertools import groupby
import pandas as pd
import numpy as np

dir_path = "/home/vahid/Dropbox/data/swarm-manager-data/results/nodejs_bookstore/WORKLOAD/CPUUsageIncrease"
auth_sla = 500
books_sla = 500
key_name = "rt_ti_u_bound_c90_p95"
# key_name = "responseTimes95th"

def does_specs_meet_sla(config):
    return config['auth']['responseTimes']['total'][key_name] < auth_sla and config['books']['responseTimes']['total'][key_name] < books_sla

rules = sorted(['cpu90_50.yml','cpu90_60.yml','cpu90_70.yml','cpu90_80.yml','cpu90_90.yml','cpu95_50.yml','cpu95_60.yml','cpu95_70.yml','cpu95_80.yml','cpu95_90.yml','cpu_mean_50.yml','cpu_mean_60.yml','cpu_mean_70.yml','cpu_mean_80.yml','cpu_mean_90.yml'])
rule2index = {}
for i,rule in enumerate(rules):
    rule2index[rule] = i + 1

workloads = ["300_120_0.3_10","400_120_0.5_10","420_120_0.7_10","500_120_0.65_10"]
headers = ['workload']
headers.extend(rules)
results = [headers]

workload_to_valid_rules = {}
valid_rules = set()
for widx, workload in enumerate(workloads):
    widx += 1
    workload_to_valid_rules[workload] = []
    results.append([''] * (len(rules)+1))
    results[widx][0] = workload

    base_path = dir_path.replace("WORKLOAD", workload)
    
    for ruleIdx, rule_name in enumerate(rules):
        path = base_path + "/" + rule_name
        with open(path) as f:
            data = yaml.load(f, Loader=yaml.FullLoader)
            found_valid = False
            for config in data['steps']:
                if does_specs_meet_sla(config['info']):
                    cpu_count = config['specs']['auth']['replicaCount'] * config['specs']['auth']['CPULimits'] + config['specs']['books']['replicaCount'] * config['specs']['books']['CPULimits'] + config['specs']['gateway']['replicaCount'] * config['specs']['gateway']['CPULimits']
                    workload_to_valid_rules[workload].append((rule_name, cpu_count))
                    valid_rules.add(rule_name)
                    found_valid = True
                    break
            if found_valid:
                results[widx][ruleIdx+1] = str(np.round(cpu_count,2)) + " G"
            else:
                config = data['steps'][-1]
                cpu_count = config['specs']['auth']['replicaCount'] * config['specs']['auth']['CPULimits'] + config['specs']['books']['replicaCount'] * config['specs']['books']['CPULimits'] + config['specs']['gateway']['replicaCount'] * config['specs']['gateway']['CPULimits']
                results[widx][ruleIdx+1] = str(np.round(cpu_count,2)) + " R"

for row in results:
    rowStr = '{:^20} '.format(row[0].replace('_',' ',3).replace('.yml','',1))
    for v in row[1:]:
        rowStr += '{:^14} '.format(v.replace('_',' ',3).replace('.yml','',1))
    print(rowStr)