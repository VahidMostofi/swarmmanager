import json
import mobopt as mo
import numpy as np
import sys

gateway_cores = 2
auth_cores = 2
books_cores = 4

keys = [
    'gateway_container_count',
    'gateway_worker_per_core',
    'auth_container_count',
    'auth_worker_per_core',
    'books_container_count',
    'books_worker_per_core',    
]
import time
def objective(x):
    
    config = {
        'gateway': {
            "cpu_count": np.round(gateway_cores / np.round(x[0],0), 2),
            "container_count": int(np.round(x[0],0)),
            "worker_count": int(np.round(np.round(x[1],0) * np.round(gateway_cores / np.round(x[0],0), 2),0)) 
        },
        'auth': {
            "cpu_count": np.round(auth_cores / np.round(x[2],0), 2),
            "container_count": int(np.round(x[2],0)),
            "worker_count": int(np.round(np.round(x[3],0) * np.round(auth_cores / np.round(x[2],0), 2),0))
        },
        'books': {
            "cpu_count": np.round(books_cores / np.round(x[4],0), 2),
            "container_count": int(np.round(x[4],0)),
            "worker_count": int(np.round(np.round(x[5],0) * np.round(books_cores / np.round(x[4],0), 2),0))
        },
    }

    print(json.dumps(config), flush=True)
    for line in sys.stdin:
        data = json.loads(line.strip())
        break
    
    return np.array(data['feedbacks'])


max_worker_count = 10
NParam = len(keys)
PB = np.asarray([
    [0.5, gateway_cores+0.5],
    [0.5, max_worker_count+0.5],
    [0.5, gateway_cores+0.5],
    [0.5, max_worker_count+0.5],
    [0.5, gateway_cores+0.5],
    [0.5, max_worker_count+0.5],
])

Optimizer = mo.MOBayesianOpt(target=objective,
                             NObj=2,
                             pbounds=PB,
                             verbose=False,
                             max_or_min='min',
                             RandomSeed=10)
Optimizer.initialize(init_points=2)
front, pop = Optimizer.maximize(n_iter=10,
                                prob=0.1,
                                ReduceProb=False,
                                q=0.5)
print('done')