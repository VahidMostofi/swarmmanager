#!/bin/bash

workload=75_1500_0.33_0.33_0.34
# 75_110_0.33_0.33_0.34/bnv2/bnv2-250-2.0-mc-c-0.5.yml
#    edit_book-rt  get_book-rt  login-rt  meets_rt  edit_book-ti  get_book-ti  \
# 7         203.5        193.0     189.0      True           263          236   

#    login-ti  meets_ti  auth_a  books_a  entry_a  sum(a)  
# 7       198     False     4.5      0.6      2.5     7.6  
rt_config="auth 4.5 books 0.6 entry 2.5"
rt_config_str=4.5_0.6_2.5

go run main.go autoconfig \
    --appname bookstore_nodejs \
    --config configurations/bookstore_nodejs.yaml \
    --workload  $workload \
    --testName "single_${rt_config_str}" \
    single $rt_config

go run main.go violations $rt_config_str $workload > bad1.txt
# #################################################################################

workload=75_1500_0.33_0.33_0.34
# 75_110_0.33_0.33_0.34/bnv2/bnv2-250-0.5-mc-c-0.5.yml
#    edit_book-rt  get_book-rt  login-rt  meets_rt  edit_book-ti  get_book-ti  \
# 9         138.0        132.5     244.0      True           152          151   

#    login-ti  meets_ti  auth_a  books_a  entry_a  sum(a)  
# 9       260     False     4.0      0.8      1.5     6.3  
rt_config="auth 4.0 books 0.8 entry 1.5"
rt_config_str=4.0_0.8_1.5

go run main.go autoconfig \
    --appname bookstore_nodejs \
    --config configurations/bookstore_nodejs.yaml \
    --workload  $workload \
    --testName "single_${rt_config_str}" \
    single $rt_config

go run main.go violations $rt_config_str $workload > bad2.txt
# #################################################################################

workload=75_1500_0.33_0.33_0.34
# 75_110_0.33_0.33_0.34/bnv2/bnv2-250-0.2-mc-c-0.5.yml
#     edit_book-rt  get_book-rt  login-rt  meets_rt  edit_book-ti  get_book-ti  \
# 19         249.0        231.5     227.5      True           287          279   

#     login-ti  meets_ti  auth_a  books_a  entry_a  sum(a)  
# 19       234     False     3.9      0.6      1.5     6.0  
rt_config="auth 3.9 books 0.6 entry 1.5"
rt_config_str=3.9_0.6_1.5

go run main.go autoconfig \
    --appname bookstore_nodejs \
    --config configurations/bookstore_nodejs.yaml \
    --workload  $workload \
    --testName "single_${rt_config_str}" \
    single $rt_config

go run main.go violations $rt_config_str $workload > bad3.txt
# #################################################################################

workload=100_1500_0.33_0.33_0.34
# 100_110_0.33_0.33_0.34/bnv2/bnv2-250-2.0-mc-c-0.5.yml
#    edit_book-rt  get_book-rt  login-rt  meets_rt  edit_book-ti  get_book-ti  \
# 2          86.5         82.0     244.5      True            94           89   

#    login-ti  meets_ti  auth_a  books_a  entry_a  sum(a)  
# 2       252     False     4.5      2.5      2.5     9.5  
rt_config="auth 4.5 books 2.5 entry 2.5"
rt_config_str=4.5_2.5_2.5

go run main.go autoconfig \
    --appname bookstore_nodejs \
    --config configurations/bookstore_nodejs.yaml \
    --workload  $workload \
    --testName "single_${rt_config_str}" \
    single $rt_config

go run main.go violations $rt_config_str $workload > bad4.txt
# #################################################################################

workload=100_1500_0.33_0.33_0.34
# 100_110_0.33_0.33_0.34/bnv2/bnv2-250-0.2-mc-c-0.5.yml
#     edit_book-rt  get_book-rt  login-rt  meets_rt  edit_book-ti  get_book-ti  \
# 20         141.5        125.5     245.5      True           151          133   

#     login-ti  meets_ti  auth_a  books_a  entry_a  sum(a)  
# 20       255     False     4.5      0.9      1.7     7.1  
rt_config="auth 4.5 books 0.9 entry 1.7"
rt_config_str=4.5_0.9_1.7

go run main.go autoconfig \
    --appname bookstore_nodejs \
    --config configurations/bookstore_nodejs.yaml \
    --workload  $workload \
    --testName "single_${rt_config_str}" \
    single $rt_config

go run main.go violations $rt_config_str $workload > bad5.txt
# #################################################################################

workload=125_1500_0.33_0.33_0.34
# 125_110_0.33_0.33_0.34/bnv2/bnv2-250-0.2-mc-c-0.5.yml
#     edit_book-rt  get_book-rt  login-rt  meets_rt  edit_book-ti  get_book-ti  \
# 31         219.0        206.0     247.0      True           239          222   

#     login-ti  meets_ti  auth_a  books_a  entry_a  sum(a)  
# 31       257     False     6.5      1.2      2.3    10.0  

rt_config="auth 6.5 books 1.2 entry 2.3"
rt_config_str=6.5_1.2_2.3

go run main.go autoconfig \
    --appname bookstore_nodejs \
    --config configurations/bookstore_nodejs.yaml \
    --workload  $workload \
    --testName "single_${rt_config_str}" \
    single $rt_config

go run main.go violations $rt_config_str $workload > bad6.txt
