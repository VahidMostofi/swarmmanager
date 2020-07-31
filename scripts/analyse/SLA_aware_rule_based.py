import sys
import yaml
from itertools import groupby
import pandas as pd

dir_path = "/home/vahid/Dropbox/data/swarm-manager-data/results/WORKLOAD/cpu_util_rule_based"
auth_sla = 350
books_sla = 350
key_name = "rt_ti_u_bound_c90_p95"
# key_name = "responseTimes95th"

def does_specs_meet_sla(config):
    return config['auth'][key_name] < auth_sla and config['books'][key_name] < books_sla

rules = ['cpu90_50.yml','cpu90_60.yml','cpu90_70.yml','cpu90_80.yml','cpu90_90.yml','cpu95_50.yml','cpu95_60.yml','cpu95_70.yml','cpu95_80.yml','cpu95_90.yml','cpu_mean_50.yml','cpu_mean_60.yml','cpu_mean_70.yml','cpu_mean_80.yml','cpu_mean_90.yml']
workloads = ["300_80_0.3_10","400_80_0.5_10","420_80_0.7_10","500_80_0.65_10"]

workload_to_valid_rules = {}
valid_rules = set()
for workload in workloads:
    workload_to_valid_rules[workload] = []

    base_path = dir_path.replace("WORKLOAD", workload)
    
    for rule_name in rules:
        path = base_path + "/" + rule_name
        with open(path) as f:
            data = yaml.load(f, Loader=yaml.FullLoader)

            for config in data['configs']:
                if does_specs_meet_sla(config['info']):
                    cpu_count = config['specs']['auth']['replicaCount'] * config['specs']['auth']['CPULimits'] + config['specs']['books']['replicaCount'] * config['specs']['books']['CPULimits'] + config['specs']['gateway']['replicaCount'] * config['specs']['gateway']['CPULimits']
                    workload_to_valid_rules[workload].append((rule_name, cpu_count))
                    valid_rules.add(rule_name)
                    break

intersection = set(rules)
for workload in workloads:
    intersection = intersection.intersection(set([x[0] for x in workload_to_valid_rules[workload]]))

valid_rules = list(intersection)
valid_rules = sorted(valid_rules)

temp = {'workload': []}
for valid_rule in valid_rules:
    temp[valid_rule] = []

for workload,_ in workload_to_valid_rules.items():
    temp['workload'].append(workload)
    for (rule_name, cpu_count) in workload_to_valid_rules[workload]:
        if rule_name in temp:
            temp[rule_name].append(cpu_count)

df = pd.DataFrame(temp)
df.to_csv("~/Desktop/sla_aware_rule_based.csv",index=False)