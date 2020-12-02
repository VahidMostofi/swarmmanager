#!/bin/bash
VUS=20
go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs-lt.yaml \
        --workload  "${VUS}_110_0.33_0.33_0.34" \
        --testName bnv2-250-2.0-mc-c-0.5 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /Users/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands.yaml \
        --mc \
        --stepsize 2.0