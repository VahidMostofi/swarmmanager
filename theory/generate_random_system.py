import numpy as np
import time
from tqdm import tqdm
import json
import os
import sys

MIN_UTILIZATION = 0.3
MAX_UTILIZATION = 0.8

MIN_X = 5
MAX_X = 30
seed = int(time.time())
# seed = 1602304465
print('seed is', seed)
np.random.seed(seed)
def generate_single_system(K,C):
    out = {}
    bad_config = True
    while bad_config:
        print('generating config ...')
        throughput = np.random.uniform(MIN_X, MAX_X)
        out['throughput'] = throughput
        temp = np.random.uniform(0.1, 1, (C))
        PCs = temp / np.sum(temp)

        XCs = throughput * PCs
        
        class_probs = {}
        for c in range(C):
            class_probs[str(c)] = PCs[c]
        out['class_probs'] = class_probs
        demands = {}
        
        # print(throughput)
        classes = [str(c) for c in range(C)]
        resources = [str(k) for k in range(K)]
        out['classes'] = classes
        out['resources'] = resources
        zeros = ''
        for k in range(K):
            Uk = np.random.uniform(MIN_UTILIZATION, MAX_UTILIZATION)
            temp = np.random.uniform(0,1, (C))
            UCs = Uk * (temp / sum(temp)) # UCs[i] = Xc * Dc,k
            count = 0
            for c in range(C):
                demands[str(c)+'_'+str(k)] = UCs[c] / XCs[c]
                if demands[str(c)+'_'+str(k)] < 0.05:
                    demands[str(c)+'_'+str(k)] = 0
                if demands[str(c)+'_'+str(k)] < 1e-4:
                    count += 1
            zeros += str(count) + ", "
        print(zeros)
        out['demands'] = demands
        utilizations = {}
        responseTimes = {}

        for c in classes:
            res = ''
            for k in resources:
                f = '($N/(1 - ($D)))'
                flag = False
                if c+"_"+k in demands:
                    flag =True
                    f = f.replace('$N', '(' + str(demands[c+'_'+k]) + '/' + 'alphas["' + k + '"]' + ')')
                    p = ''
                    for c2 in classes:
                        if c2+'_'+k in demands:
                            p += '(('+str(class_probs[c2])+"*"+str(throughput) + "*" + str(demands[c2+'_'+k])+')/alphas["'+ k + '"]'+')'
                            p += '+'
                    utilizations[k] = p[:-1]
                    p = p[:-1]
                    f = f.replace('$D',p)
                if flag:
                    res += f + "+"
            res = res[:-1]
            responseTimes[c] = res
        
        bad_config = False
        alphas = {}
        for k in range(K):
            alphas[str(k)] = 1

        for k in range(K):
            assert eval(utilizations[str(k)]) < 1
            print(str(k), eval(utilizations[str(k)]))
            if eval(utilizations[str(k)]) < 1e-5:
                bad_config = True
                print('retrying... because of K')
                break
        for c in range(C):
            print(str(c), eval(responseTimes[str(c)]))
            if eval(responseTimes[str(c)]) < 1e-5:
                bad_config = True
                print('retrying... because of C')
                break

    python_code = """
def mean_response_timesF(alphas):
    A
    return B
    
def utilizationsF(alphas):
    return C
    """
    rts = ''
    ret = ''
    for key in range(C):
        key = str(key)
        value = responseTimes[key]
        rts += 'r' + key + ' = ' + value + '\n    '
        ret += 'r' + key + ' * 1000, '
    ret = ret[:-2]
    python_code = python_code.replace('A', rts)
    python_code = python_code.replace('B', ret)

    ut = ''
    for key in range(K):
        key = str(key)
        value = utilizations[key]
        ut += value + ', '
    ut = ut[:-2]
    python_code = python_code.replace('C', ut)

    exec(python_code, globals())
    
    def objective(_alphas):
        s = 0
        for key, value in _alphas.items():
            s += value
        return s
    start = 1
    stepSize = .1
    count = 30
    SLA = 150
    best_objective = 10000
    best = {}
    A = {}
    for i0 in tqdm(range(count)):
        for i1 in range(count):
            for i2 in range(count):
                for i3 in range(count):
                    for i4 in range(count):
                            A['0'] = start + i0*stepSize
                            A['1'] = start + i1*stepSize
                            A['2'] = start + i2*stepSize
                            A['3'] = start + i3*stepSize
                            A['4'] = start + i4*stepSize

                            mrts = mean_response_timesF(A)
                            meets = True
                            for v in mrts:
                                if v > SLA:
                                    meets = False
                                    break

                            if meets:
                                if objective(A) < best_objective:
                                    objective(A)
                                    best = {}
                                    best['0'] = A['0']
                                    best['1'] = A['1']
                                    best['2'] = A['2']
                                    best['3'] = A['3']
                                    best['4'] = A['4']
                                    best_objective = objective(best)
    print(best_objective, best)
    print(mean_response_timesF(best))
    out['best_objective'] = best_objective
    out['best'] = best
    out['best_mrt'] = mean_response_timesF(best)
    out['SLA'] = SLA
    # name = str(len(list(os.listdir('./systems/')))+1)
    name = str(time.time()).replace('.','')
    with open('./systems/' + name + '.json','w') as file:
        json.dump(out, file, indent=4, sort_keys=True)

generate_single_system(5, 12)

