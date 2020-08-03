import sys
import yaml
from itertools import groupby

dir_path = "/home/vahid/Dropbox/data/swarm-manager-data/results/WORKLOAD/CPUUsageIncrease"
auth_sla = 350
books_sla = 350
# key_name = "rt_ti_u_bound_c90_p95"
key_name = "responseTimes95th"

def does_specs_meet_sla(config):
    return config['auth'][key_name] < auth_sla and config['books'][key_name] < books_sla

rules = ['cpu90_50.yml','cpu90_60.yml','cpu90_70.yml','cpu90_80.yml','cpu90_90.yml','cpu95_50.yml','cpu95_60.yml','cpu95_70.yml','cpu95_80.yml','cpu95_90.yml','cpu_mean_50.yml','cpu_mean_60.yml','cpu_mean_70.yml','cpu_mean_80.yml','cpu_mean_90.yml']
workloads = ["300_80_0.3_10","400_80_0.5_10","420_80_0.7_10","500_80_0.65_10"]

workload_to_valid_rules = {}

for workload in workloads:
    workload_to_valid_rules[workload] = []

    base_path = dir_path.replace("WORKLOAD", workload)
    
    for rule_name in rules:
        path = base_path + "/" + rule_name
        with open(path) as f:
            data = yaml.load(f, Loader=yaml.FullLoader)

            for config in data['configs']:
                if does_specs_meet_sla(config['info']):
                    workload_to_valid_rules[workload].append(rule_name)

intersection = set(rules)
for workload in workloads:
    intersection = intersection.intersection(set(workload_to_valid_rules[workload]))
    # print(workload, set(workload_to_valid_rules[workload]))
print('there are',len(intersection),'rules that reach to a configuration which meets the SLA:')
print(intersection)

intersection = set(rules)
for workload in workloads:
    valid_configs = workload_to_valid_rules[workload]
    no_waste = []
    for rule, group in groupby(valid_configs):
        count = len(list(group))
        if count == 1:
            no_waste.append(rule)
    print(len(no_waste))
    intersection = intersection.intersection(set(no_waste))
print(intersection)