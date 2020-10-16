import numpy as np
import time
from tqdm import tqdm
import json
import os
import sys
import tempfile
from amplpy import AMPL, Environment

MIN_UTILIZATION = 0.3
MAX_UTILIZATION = 0.99

MIN_X = 5
MAX_X = 10
seed = int(time.time())
# seed = 1602304465
print('seed is', seed)
np.random.seed(seed)
def generate_single_system(K,C):
    out = {}
    bad_config = True
    SLA = 200
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
        de = []
        for k in range(K):
            Uk = np.random.uniform(MIN_UTILIZATION, MAX_UTILIZATION)
            temp = np.random.uniform(0,1, (C))
            UCs = Uk * (temp / sum(temp)) # UCs[i] = Xc * Dc,k
            count = 0
            for c in range(C):
                demands[str(c)+'_'+str(k)] = np.round(UCs[c] / XCs[c], 5)
                if demands[str(c)+'_'+str(k)] < 0.05:
                    demands[str(c)+'_'+str(k)] = 0
                
                if demands[str(c)+'_'+str(k)] < 1e-4:
                    count += 1
                else:
                    de.append(demands[str(c)+'_'+str(k)] * 1000)
            zeros += str(count) + ", "
        
        print('mean_demands', np.mean(de))
        factor = 1.2 * (K+C) / (K * C)
        for k in range(K):
            for c in range(C):
                demands[str(c)+'_'+str(k)] *= factor
        print(factor)
        # print(zeros)
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

        bad_config = False
        alphas = {}
        for k in range(K):
            alphas[str(k)] = 1

        max_response_time = 0
        min_response_time = 10000000000000000000000
        for r in mean_response_timesF(alphas):
            if r > max_response_time:
                max_response_time = r
            if r < min_response_time:
                min_response_time = r
            if r < 1e-5:
                bad_config = True
                print('retrying... because of C')
                break
        print('max_response_time', max_response_time)
        print('min_response_time', min_response_time)
        out['min_response_time'] = min_response_time
        out['max_response_time'] = max_response_time
        out['mean_response_time'] = np.mean(mean_response_timesF(alphas))
        print('mean_response_time', out['mean_response_time'])
    
    sts = []
    obj = 'minimize objective:'
    for r in resources:
        sts.append('var A'+r +';')
        obj += ' A' + r + ' + '
    obj = obj[:-3] + ';'
    sts.append('')
    sts.append(obj)
    sts.append('')

    for r in resources:
        sts.append('subject to count_' + r + ": A" + r + " >= 1;")
    sts.append('')

    for c in classes:
        res = ''
        for k in resources:
            f = '($N/(1 - ($D)))'
            flag = False
            if c+"_"+k in demands:
                flag =True
                f = f.replace('$N', '(' + str(demands[c+'_'+k]) + '/' + 'A' + k + ')')
                p = ''
                for c2 in classes:
                    if c2+'_'+k in demands:
                        p += '(('+str(class_probs[c2])+"*"+str(throughput) + "*" + str(demands[c2+'_'+k])+') / A'+ k + ''+')'
                        p += '+'
                p = p[:-1]
                f = f.replace('$D',p)
            if flag:
                res += f + "+"
        res = res[:-1]
        sts.append('subject to R' + c + ' : ' +res + ' <= '+str(SLA/1000.0)+' ;')
        sts.append('')
    
    ampl = AMPL(Environment('/home/vahid/apps/amplide.linux64/'))
    ampl.setOption('solver', 'conopt')

    _, path = tempfile.mkstemp()
    
    with open(path, 'w') as tempf:
        for s in sts:
            tempf.write(s + '\n')
    
    ampl.read(path)
    ampl.solve()
    best_objective = ampl.getObjective('objective').value()
    best = {}
    for k in range(K):
        best[str(k)] = ampl.getVariable('A' + str(k)).value()
    print(best)
    print(best_objective)

    print(best_objective, best)
    print(mean_response_timesF(best))
    out['best_objective'] = best_objective
    out['best'] = best
    out['best_mrt'] = list(mean_response_timesF(best))
    out['SLA'] = SLA
    # out['python_code'] = python_code
    name = str(len(list(os.listdir('./systems/')))+1)
    # name = str(time.time()).replace('.','')
    with open('./systems/' + name + '.json','w') as file:
        json.dump(out, file, indent=4)

K_COUNT = np.random.randint(5,50)
C_COUNT = np.random.randint(K_COUNT, K_COUNT * 2)
print(K_COUNT, C_COUNT)
generate_single_system(K_COUNT, C_COUNT)

import generate_demand_files