# split CPU cores, based on share of each service, create multiple containers. at most 1 CPU core per container
import json
import mobopt as mo
import numpy as np
import sys
core_count = 20
cache = {}

import time
def objective(x):

    s = sum(x)
    g = core_count * (x[0] / s)
    a = core_count * (x[1] / s)
    b = core_count * (x[2] / s)
    # this configuration would use a + b + g cores. x[3] is the amount which is not being used
    
    config = {
        'gateway': {
            "cpu_count": np.round(g / np.ceil(g), 2),
            "container_count": int(np.ceil(g)),
            "worker_count": 1 #this is 1, because it's multi container 
        },
        'auth': {
            "cpu_count": np.round(a / np.ceil(a), 2),
            "container_count": int(np.ceil(a)),
            "worker_count": 1
        },
        'books': {
            "cpu_count": np.round(b / np.ceil(b), 2),
            "container_count": int(np.ceil(b)),
            "worker_count": 1
        },
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
    res = [data['feedbacks'][0],data['feedbacks'][1],a+b+g]
    cache[key] = np.array(res)

    return np.array(np.array(res))

PB = np.asarray([
    [0.06, 0.94],
    [0.06, 0.94],
    [0.06, 0.94],
    [0.06, 0.94]
])
NParam = PB.shape[0]

Optimizer = mo.MOBayesianOpt(target=objective,
                             NObj=3,
                             pbounds=PB,
                             verbose=False,
                             max_or_min='min',
                             RandomSeed=10)
Optimizer.initialize(init_points=5)
front, pop = Optimizer.maximize(n_iter=35,
                                prob=0.1,
                                ReduceProb=False,
                                q=0.5)
print('done')