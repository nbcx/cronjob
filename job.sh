#!/usr/bin/env bash
SHELL=/bin/bash
PATH=/sbin:/bin:/usr/sbin:/usr/bin
MAILTO=root

# For details see man 4 crontabs

[script]
# Example of job definition:
# .---------------- id
# |  .------------- concurrent
# |  |  .---------- user
# |  |  |  .------- command
# |  |  |  |  .---- args
# |  |  |  |  |
# *  *  *  *  * user-name  command to be executed
0001        0 0 php php.php 1>>out.txt 2>&1



[crontab]