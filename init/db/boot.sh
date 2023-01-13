#!/bin/bash

set -e

# 初始化数据库
DBNAME=growerlab
USERNAME=growerlab

createuser $USERNAME -P -e

# create database for growerlab
createdb $DBNAME --owner $USERNAME -e --encoding=UTF8

psql -d $DBNAME -U $USERNAME -c "GRANT ALL ON SCHEMA public TO ${USERNAME};"
psql -d $DBNAME -U $USERNAME -f ./growerlab.sql
psql -d $DBNAME -U $USERNAME -f ./seed.sql
