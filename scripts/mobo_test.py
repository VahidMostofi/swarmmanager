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

core_count = 22
good_ones = []

import time
def objective(x):
    s = sum(x)
    
    g = core_count * (x[0] / s)
    a = core_count * (x[1] / s)
    b = core_count * (x[2] / s)
    if a + b + g > 22:
        print('INVALID')
        return np.array([100000,100000,a+b+g])
    
    f1 = (6350 / (2.4 * g + 0.4 * a * a))
    f2 = (55000 / (2.4 * g + 3.4 * b * b))
    res = [f1,f2,a+b+g]
    print(res)
    if (f1 < 400 and f2 < 400):
        good_ones.append([f1,f2,a,b,g,a+b+g])
    return np.array(res)


max_worker_count = 10
NParam = len(keys)
PB = np.asarray([
    [0.005, 1],
    [0.005, 1],
    [0.005, 1],
    [0.005, 1]
])

Optimizer = mo.MOBayesianOpt(target=objective,
                             NObj=3,
                             pbounds=PB,
                             verbose=False,
                             max_or_min='min',
                             RandomSeed=5)
Optimizer.initialize(init_points=5)
front, pop = Optimizer.maximize(n_iter=6,
                                prob=0.1,
                                ReduceProb=False,
                                q=0.5)

print('done')
for good in good_ones:
    print(np.round(good[0],2),np.round(good[1],2),np.round(good[-1]),":", np.round(good[2:5],2))