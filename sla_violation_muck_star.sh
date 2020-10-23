#!/bin/bash

#didn't work
# ti_config="serviceb 2.5 servicec 2.5 serviced 2.5 servicee 2.5 servicef 2.5"
# ti_config_str=2.5_2.5_2.5_2.5_2.5

# r_config="serviceb 0.8 servicec 2.5 serviced 2.5 servicee 2.5 servicef 2.5"
# r_config_str=0.8_2.5_2.5_2.5_2.5

ti_config="serviceb 1.4 servicec 2 serviced 2 servicee 1.3 servicef 1.3"
ti_config_str=2.5_2.5_2.5_2.5_2.5

r_config="serviceb 0.5 servicec 0.5 serviced 0.5 servicee 0.5 servicef 0.5"
r_config_str=0.8_2.5_2.5_2.5_2.5


workload=10_1500_0.45_0.25_0.15_0.15

# go run main.go autoconfig \
#     --appname muck_star-small \
#     --config configurations/muck_star.yaml \
#     --workload  $workload \
#     --testName "single_${ti_config_str}" \
#     single $ti_config


# go run main.go autoconfig \
#     --appname muck_star-small \
#     --config configurations/muck_star.yaml \
#     --workload  $workload \
#     --testName "single_${r_config_str}" \
#     single $r_config

go run main.go violations $r_config_str $workload
go run main.go violations $ti_config_str $workload
