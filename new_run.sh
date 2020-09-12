#!/bin/bash
go run main.go autoconfig \
    --appname bookstore_nodejs \
    --config configurations/bookstore_nodejs.yaml \
    --workload  1_100_0.5_5 \
    --testName adfc-350-1-mc \
    adfc \
    --property RTToleranceIntervalUBoundc90p95 \
    --value 350 \
    --stepsize 1 \
    --mc