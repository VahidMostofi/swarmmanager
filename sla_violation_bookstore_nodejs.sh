#!/bin/bash

#didn't work
# ti_config="auth 4.5 books 0.8 entry 1.5"
# ti_config_str=4.5_0.8_1.5

# ti_config="auth 3.9 books 0.9 entry 1.5"
# ti_config_str=3.9_0.9_1.5

ti_config="auth 4.5 books 0.7 entry 2.5"
ti_config_str=4.5_0.7_2.5


workload=75_1500_0.33_0.33_0.34

# go run main.go autoconfig \
#     --appname bookstore_nodejs \
#     --config configurations/bookstore_nodejs.yaml \
#     --workload  $workload \
#     --testName "single_${ti_config_str}" \
#     single $ti_config

# go run main.go autoconfig \
#     --appname muck_star-small \
#     --config configurations/muck_star.yaml \
#     --workload  $workload \
#     --testName "single_${r_config_str}" \
#     single $r_config

# go run main.go violations $r_config_str $workload
go run main.go violations $ti_config_str $workload
