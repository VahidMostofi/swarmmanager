
import sys
import yaml
dir_path = "/home/vahid/Dropbox/data/swarm-manager-data/results/400_80_0.5_10/cpu_util_rule_based"
auth_sla = 350
books_sla = 350
key_name = "rt_ti_u_bound_c90_p95"
path = sys.argv[1]


def print_something_with_cores(path):
    with open(path) as f:
        data = yaml.load(f)
    print('name: {:s}'.format(data['name']))
    print("based on:",key_name,'and total replica count meet')
    print('# {:^10} {:^10} {:^27} {:^11}'.format('auth','books','configs','count'))
    print('   {:^21} {:^8}'.format('','ar,ac,aw,br,bc,bw,gr,gc,gw'))
    for idx, config in enumerate(data['configs']):
        print('{:2d} {:7.2f}({:s}) {:7.2f}({:s}) {:2d},{:2d},{:2d}:{:2d},{:2d},{:2d}:{:2d},{:2d},{:2d}:{:2d} {:5d} {:1s} {:8s}'.format(
            idx+1,
            config['info']['auth'][key_name], 'Y' if config['info']['auth'][key_name] <= auth_sla else 'N',
            config['info']['books'][key_name], 'Y' if config['info']['books'][key_name] <= books_sla else 'N', 
            config['specs']['auth']['replicaCount'],config['specs']['auth']['CPULimits'],int(config['specs']['auth']['envs'][-1].split('=')[1]), 
            config['specs']['books']['replicaCount'],config['specs']['books']['CPULimits'],int(config['specs']['books']['envs'][-1].split('=')[1]), 
            config['specs']['gateway']['replicaCount'],config['specs']['gateway']['CPULimits'],int(config['specs']['gateway']['envs'][-1].split('=')[1]), 
            config['specs']['auth']['replicaCount']    * config['specs']['auth']['CPULimits']+ 
            config['specs']['books']['replicaCount']   * config['specs']['books']['CPULimits']+ 
            config['specs']['gateway']['replicaCount'] * config['specs']['gateway']['CPULimits'],
            config['info']['gateway']['requestCount'], 
            'Y' if config['info']['auth'][key_name] <= auth_sla and config['info']['books'][key_name] <= books_sla else 'N',
            config['hash'][:8]
            ))

if path == 'all_rule_based':
    paths = ['cpu90_50.yml','cpu90_60.yml','cpu90_70.yml','cpu90_80.yml','cpu90_90.yml','cpu95_50.yml','cpu95_60.yml','cpu95_70.yml','cpu95_80.yml','cpu95_90.yml','cpu_mean_50.yml','cpu_mean_60.yml','cpu_mean_70.yml','cpu_mean_80.yml','cpu_mean_90.yml']
    paths = [dir_path + "/" + p for p in paths]
    for p in paths:
        print_something_with_cores(p)
        print('=============================================')
else:
    print_something_with_cores(path)