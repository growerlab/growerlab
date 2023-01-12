#!/bin/bash

set -e

# 初始化数据库

# default password growerlab
createuser growerlab -P -e

# create database for growerlab
createdb growerlab --owner growerlab -e --encoding=UTF8

# GRANT ALL ON SCHEMA public TO growerlab;