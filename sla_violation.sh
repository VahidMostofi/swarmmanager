#!/bin/bash

ti_config="serviceb 2.5 servicec 2.5 serviced 2.5 servicee 2.5 servicef 2.5"
ti_config_str=2.5_2.5_2.5_2.5_2.5

r_config="serviceb 0.8 servicec 2.5 serviced 2.5 servicee 2.5 servicef 2.5"
r_config_str=0.8_2.5_2.5_2.5_2.5

workload=10_1500_0.45_0.25_0.15_0.15

go run main.go autoconfig \
    --appname muck_star-small \
    --config configurations/muck_star.yaml \
    --workload  $workload \
    --testName "signle_${ti_config_str}" \
    single $ti_config


go run main.go autoconfig \
    --appname muck_star-small \
    --config configurations/muck_star.yaml \
    --workload  $workload \
    --testName "signle_${r_config_str}" \
    single $r_config
