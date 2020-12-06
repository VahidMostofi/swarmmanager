#!/bin/bash
# VUS=100
# go run main.go autoconfig \
#         --appname bookstore_nodejs \
#         --config configurations/bookstore_nodejs-lt.yaml \
#         --workload  "${VUS}_180_0.33_0.33_0.34" \
#         --testName bnv2-250-2.0-mc-c-0.5 \
#         bnv2 \
#         --property ResponseTimes95Percentile  \
#         --value 250 \
#         --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands-2.yaml \
#         --mc \
#         --stepsize 2.0

# VUS=100
# go run main.go autoconfig \
#         --appname bookstore_nodejs \
#         --config configurations/bookstore_nodejs-lt.yaml \
#         --workload  "${VUS}_180_0.33_0.33_0.34" \
#         --testName bnv2-250-1.0-mc-c-0.5 \
#         bnv2 \
#         --property ResponseTimes95Percentile  \
#         --value 250 \
#         --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands-2.yaml \
#         --mc \
#         --stepsize 1.0

# VUS=150
# go run main.go autoconfig \
#         --appname bookstore_nodejs \
#         --config configurations/bookstore_nodejs-lt.yaml \
#         --workload  "${VUS}_180_0.33_0.33_0.34" \
#         --testName bnv2-250-1.0-mc-c-0.5 \
#         bnv2 \
#         --property ResponseTimes95Percentile  \
#         --value 250 \
#         --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands-2.yaml \
#         --mc \
#         --stepsize 1.0

VUS=75
go run main.go autoconfig \
        --appname bookstore_nodejs \
        --config configurations/bookstore_nodejs-lt.yaml \
        --workload  "${VUS}_180_0.33_0.33_0.34" \
        --testName bnv2-250-2.0-mc-c-0.5 \
        bnv2 \
        --property ResponseTimes95Percentile  \
        --value 250 \
        --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands-2.yaml \
        --mc \
        --stepsize 2.0

# VUS=125
# go run main.go autoconfig \
#         --appname bookstore_nodejs \
#         --config configurations/bookstore_nodejs-lt.yaml \
#         --workload  "${VUS}_180_0.33_0.33_0.34" \
#         --testName bnv2-250-1.0-mc-c-0.5 \
#         bnv2 \
#         --property ResponseTimes95Percentile  \
#         --value 250 \
#         --demands /home/vahid/Dropbox/data/swarm-manager-data/demands/bookstore_nodejs_demands-2.yaml \
#         --mc \
#         --stepsize 1.0