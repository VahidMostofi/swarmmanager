import os
import yaml
import json

for f in os.listdir('./systems'):
    with open('./systems/' + f) as file:
        data = json.load(file)

        demands = {}
        for k in data['resources']:
            demands[k] = {}
            for c in data['classes']:
                if str(c) + '_' + str(k) in data['demands']:
                    d = data['demands'][str(c) + '_' + str(k)]
                else:
                    d = 0
                demands[k][c] = int(d * 1000)
        with open('./demands/' + f[:-5] +'.yml', 'w') as pyf:
            yaml.dump(demands, pyf, default_flow_style=False)
    