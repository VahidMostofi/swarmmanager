import numpy as np
import time
from tqdm import tqdm
import json
import os
import sys
import tempfile
from amplpy import AMPL, Environment

seed = int(time.time())
# seed = 1602304465
print('seed is', seed)
np.random.seed(seed)

# each path is a list of resources
def generate_paths(count, resources, min_path_length, max_path_length):
    while True:
        usedResources = set()
        paths = []
        for _ in range(count):
            path = []
            path_length = np.random.randint(min_path_length, max_path_length)
            for _ in range(path_length):
                option = resources[np.random.randint(0,len(resources))]
                if option not in path:
                    path.append(option)
                    usedResources.add(option)
            paths.append(path)
        # check
        if len(usedResources) == len(resources):
            break
    return paths

def get_class_probabilities(classCount):
    temp = np.random.uniform(0.1, 1, (classCount))
    return temp / np.sum(temp)

def get_demands(paths, resourceCount, XCs, MIN_UTILIZATION, MAX_UTILIZATION):
    demands = {}
    resource2path = {}
    for k in range(resourceCount):
        resource2path[k] = []
        for i, path in enumerate(paths):
            if k in path:
                resource2path[k].append(i)

    for k in range(resourceCount):
        Uk = np.random.normal()
        Uk = np.random.uniform(0.3,0.9)

        temp = np.random.uniform(0,1, (len(resource2path[k])))
        UCs = Uk * (temp / sum(temp)) # UCs[i] = Xc * Dc,k
        
        for i,c in enumerate(resource2path[k]):
            demands[str(c)+"_"+str(k)] = np.round(UCs[i] / XCs[c], 5)
    return demands

def generate_functions(resources, paths, demands, throughput, class_probabilities):
    utilizations = {}
    responseTimes = {}

    for pathIdx, path in enumerate(paths):        
        res = ''
        for k in path:
            k = str(k)
            f = '($N/(1 - ($D)))'
            if str(pathIdx) + "_" + k in demands:
                f = f.replace('$N', '(' + str(demands[str(pathIdx)+'_'+k]) + '/' + 'alphas["' + k + '"]' + ')')
                p = ''
                for c2, _ in enumerate(paths):
                    c2 = str(c2)
                    if c2 + '_' + k in demands:
                        p += '(('+str(class_probabilities[int(c2)])+"*"+str(throughput) + "*" + str(demands[c2+'_'+k])+')/alphas["'+ k + '"]'+')'
                        p += '+'
                utilizations[k] = p[:-1]
                p = p[:-1]
                f = f.replace('$D',p)
            else:
                assert False
            res += f + "+"
        res = res[:-1]
        responseTimes[str(pathIdx)] = res

    python_code = """
def mean_response_timesF(alphas):
    A
    return B
    
def utilizationsF(alphas):
    return C
        """
    rts = ''
    ret = ''
    for key, path in enumerate(paths):
        key = str(key)
        value = responseTimes[key]
        rts += 'r' + key + ' = ' + value + '\n    '
        ret += 'r' + key + ' * 1000, '
    ret = ret[:-2]
    python_code = python_code.replace('A', rts)
    python_code = python_code.replace('B', ret)

    ut = ''
    
    for key in resources:
        key = str(key)
        value = utilizations[key]
        ut += value + ', '
    ut = ut[:-2]
    python_code = python_code.replace('C', ut)

    exec(python_code, globals())
    return python_code

def solve_ampl(resources, classes, class_probs, throughput, SLA):
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
                        p += '(('+str(class_probs[int(c2)])+"*"+str(throughput) + "*" + str(demands[c2+'_'+k])+') / A'+ k + ''+')'
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


    return best_objective, best

MIN_X = 5
MAX_X = 10

K = 10
C = 25
SLA = 450

paths = generate_paths(C, list(range(K)), 2, 8)
throughput = np.round(np.random.uniform(MIN_X, MAX_X))
class_probabilities = get_class_probabilities(len(paths))
demands = get_demands(paths, K, throughput * class_probabilities, 0.3, 0.99)
generate_functions(list(range(K)), paths, demands, throughput, class_probabilities)
alpha_tests = {}
for k in range(K):
    alpha_tests[str(k)] = 1
# print(mean_response_timesF(alpha_tests))
# print(utilizationsF(alpha_tests))

best_objective, best = solve_ampl([str(s) for s in list(range(K))], [str(s) for s in list(range(len(paths)))], class_probabilities, throughput, SLA)
out = {}
out['throughput'] = throughput
out['class_probs'] = {}
for i in range(class_probabilities.shape[0]):
    out['class_probs'][str(i)] = class_probabilities[i]
out['best_objective'] = best_objective
out['classes'] = [str(s) for s in list(range(len(paths)))]
out['resources'] = [str(s) for s in list(range(K))]
out['demands'] = demands
out['best'] = best
out['best_mrt'] = list(mean_response_timesF(best))
out['SLA'] = SLA
name = str(len(list(os.listdir('./systems/')))+1)
with open('./systems/' + name + '.json','w') as file:
    json.dump(out, file, indent=4)


import generate_demand_files