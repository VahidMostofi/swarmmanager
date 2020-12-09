# split CPU cores, based on share of each service, create multiple containers. at most 1 CPU core per container

# https://github.com/ppgaluzio/MOBOpt
# @article{GALUZIO2020100520,
# title = "MOBOpt — multi-objective Bayesian optimization",
# journal = "SoftwareX",
# volume = "12",
# pages = "100520",
# year = "2020",
# issn = "2352-7110",
# doi = "https://doi.org/10.1016/j.softx.2020.100520",
# url = "http://www.sciencedirect.com/science/article/pii/S2352711020300911",
# author = "Paulo Paneque Galuzio and Emerson Hochsteiner [de Vasconcelos Segundo] and Leandro dos Santos Coelho and Viviana Cocco Mariani"
# }
import json
import mobopt as mo
import numpy as np
import sys

cache = {}
services_count, request_count,sla = 0,0,0
with open('/home/vahid/Desktop/values.txt') as f:
    services_count, request_count,sla = [int(a) for a in f.read().split(',')]
core_count = 100 * services_count

import time
def objective(x):
    
    s = sum(x)
    allocations = [0] * services_count
    for i in range(services_count):
        allocations[i] = core_count * (x[i] / s)
        if allocations[i] < 1:
            allocations[i] = 1
    not_used = core_count - sum(allocations)
    
    config = {}
    for i in range(services_count):
        config[str(i)] = {
            "cpu_count": np.round(allocations[i] / np.ceil(allocations[i]), 2),
            "container_count": int(np.ceil(allocations[i])),
            "worker_count": 1 #this is 1, because it's multi container 
        }

    key = json.dumps(config, sort_keys=True)
    if key in cache:
        return cache[key]

    print(json.dumps(config), flush=True)
    line = "default"
    for line in sys.stdin:
        data = json.loads(line.strip())
        break
    with open("/home/vahid/Desktop/log.python.mobo", "w+") as f:
        f.write(str(x))
    f.close()

    SLA_target = sla
    respones_times = [0] * request_count
    res = []
    
    for i in range(request_count):
        if data['feedbacks'][i] > SLA_target:
            respones_times[i] = data['feedbacks'][i] - SLA_target
        res.append(respones_times[i])
    res.append(not_used)
    cache[key] = np.array(res)

    return np.array(np.array(res))

st = 1 / core_count
temp = [[st + 0.01, 0.94]] * (services_count+1)
PB = np.asarray(temp)
NParam = PB.shape[0]

Optimizer = mo.MOBayesianOpt(target=objective,
                             NObj=request_count+1,
                             pbounds=PB,
                             verbose=False,
                             max_or_min='min', # whether the optimization problem is a maximization problem ('max'), or a minimization one ('min')
                             RandomSeed=10)
Optimizer.initialize(init_points=5) 
# there is no minimize function. maximize() starts optimization. Performs minimizing or maximizing based on max_or_min
front, pop = Optimizer.maximize(n_iter=200,
                                prob=0.1,
                                ReduceProb=False,
                                q=0.5)
print('done')