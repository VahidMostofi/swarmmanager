#!/bin/bash

CPUs=(0.5 1.5 2.5 3.5 4.5)
for e in "${CPUs[@]}"
do
    for b in "${CPUs[@]}"
    do
        for a in "${CPUs[@]}"
        do
            go run main.go autoconfig \
                --appname bookstore_nodejs \
                --config configurations/bookstore_nodejs.yaml \
                --workload  "150_110_0.33_0.33_0.34" \
                --testName "single_${a}_${b}_${e}" \
                single entry $e auth $a books $b
        done
    done
done