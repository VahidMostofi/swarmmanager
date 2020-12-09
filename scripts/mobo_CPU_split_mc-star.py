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
core_count = 21
cache = {}

import time
def objective(x):

    s = sum(x)
    b = core_count * (x[0] / s)
    c = core_count * (x[1] / s)
    d = core_count * (x[2] / s)
    e = core_count * (x[3] / s)
    f = core_count * (x[4] / s)
    
    
    config = {
        'serviceb': {
            "cpu_count": np.round(b / np.ceil(b), 2),
            "container_count": int(np.ceil(b)),
            "worker_count": 1 #this is 1, because it's multi container 
        },
        'servicec': {
            "cpu_count": np.round(c / np.ceil(c), 2),
            "container_count": int(np.ceil(c)),
            "worker_count": 1
        },
        'serviced': {
            "cpu_count": np.round(d / np.ceil(d), 2),
            "container_count": int(np.ceil(d)),
            "worker_count": 1
        },
        'servicee': {
            "cpu_count": np.round(e / np.ceil(e), 2),
            "container_count": int(np.ceil(e)),
            "worker_count": 1
        },
        'servicef': {
            "cpu_count": np.round(f / np.ceil(f), 2),
            "container_count": int(np.ceil(f)),
            "worker_count": 1
        }
    }

    key = json.dumps(config, sort_keys=True)
    if key in cache:
        return cache[key]

    print(json.dumps(config), flush=True)
    line = "default"
    for line in sys.stdin:
        data = json.loads(line.strip())
        break
    with open("/home/vahid/Desktop/log.python.mobo", "w+") as ff:
        ff.write(str(x))
    ff.close()

    SLA_target = 300
    respones_times = [0] * 4
    for i in range(len(respones_times)):
        if data['feedbacks'][i] > SLA_target:
            respones_times[i] = data['feedbacks'][i] - SLA_target
    
    res = [respones_times[0],respones_times[1],respones_times[2],respones_times[3],b+c+d+e+f]
    cache[key] = np.array(res)

    return np.array(np.array(res))

PB = np.asarray([
    [0.06, 0.94],
    [0.06, 0.94],
    [0.06, 0.94],
    [0.06, 0.94],
    [0.06, 0.94],
    [0.06, 0.94]
])
NParam = PB.shape[0]

Optimizer = mo.MOBayesianOpt(target=objective,
                             NObj=5,
                             pbounds=PB,
                             verbose=False,
                             max_or_min='min', # whether the optimization problem is a maximization problem ('max'), or a minimization one ('min')
                             RandomSeed=10)
Optimizer.initialize(init_points=5) 
# there is no minimize function. maximize() starts optimization. Performs minimizing or maximizing based on max_or_min
front, pop = Optimizer.maximize(n_iter=20,
                                prob=0.1,
                                ReduceProb=False,
                                q=0.5)
print('done')