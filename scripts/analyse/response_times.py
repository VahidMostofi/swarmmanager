
import sys
import yaml
dir_path = "/home/vahid/Dropbox/data/swarm-manager-data/results/java-mm-bookstore/300_80_0.3_10/CPUUsageIncrease"
auth_sla = 350
books_sla = 350
key_name = "rt_ti_u_bound_c90_p95"
path = sys.argv[1]


def print_something_with_cores(path):
    with open(path) as f:
        data = yaml.load(f)
    print('name: {:s}'.format(data['name']))
    print("based on:",key_name,'and total replica count meet')
    print('# {:^11} {:^11} {:^27} {:^11}'.format('auth','books','configs','count'))
    print('   {:^22} {:^11}'.format(' ',' ar, ac ,aw,br, bc ,bw,gr, gc ,gw'))
    for idx, steps in enumerate(data['steps']):
        print('{:2d} {:8.2f}({:s}) {:8.2f}({:s}) {:2d},{:2.2f},{:2d}:{:2d},{:2.2f},{:2d}:{:2d},{:2.2f},{:2d}:{:2.2f} {:5d} {:1s} {:8s}'.format(
            idx+1,
            steps['info']['auth']['responseTimes']['total'][key_name], 'Y' if steps['info']['auth']['responseTimes']['total'][key_name] <= auth_sla else 'N',
            steps['info']['books']['responseTimes']['total'][key_name], 'Y' if steps['info']['books']['responseTimes']['total'][key_name] <= books_sla else 'N', 
            steps['specs']['auth']['replicaCount'],steps['specs']['auth']['CPULimits'],int(steps['specs']['auth']['envs'][-1].split('=')[1]), 
            steps['specs']['books']['replicaCount'],steps['specs']['books']['CPULimits'],int(steps['specs']['books']['envs'][-1].split('=')[1]), 
            steps['specs']['gateway']['replicaCount'],steps['specs']['gateway']['CPULimits'],int(steps['specs']['gateway']['envs'][-1].split('=')[1]), 
            steps['specs']['auth']['replicaCount']    * steps['specs']['auth']['CPULimits']+ 
            steps['specs']['books']['replicaCount']   * steps['specs']['books']['CPULimits']+ 
            steps['specs']['gateway']['replicaCount'] * steps['specs']['gateway']['CPULimits'],
            steps['info']['gateway']['requestCount'], 
            'Y' if steps['info']['auth']['responseTimes']['total'][key_name] <= auth_sla and steps['info']['books']['responseTimes']['total'][key_name] <= books_sla else 'N',
            steps['hash'][:8]
            ))

if path == 'all_rule_based':
    paths = ['cpu90_50.yml','cpu90_60.yml','cpu90_70.yml','cpu90_80.yml','cpu90_90.yml','cpu95_50.yml','cpu95_60.yml','cpu95_70.yml','cpu95_80.yml','cpu95_90.yml','cpu_mean_50.yml','cpu_mean_60.yml','cpu_mean_70.yml','cpu_mean_80.yml','cpu_mean_90.yml']
    paths = [dir_path + "/" + p for p in paths]
    for p in paths:
        print_something_with_cores(p)
        print('=============================================')
else:
    print_something_with_cores(path)