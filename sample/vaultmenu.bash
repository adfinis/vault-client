#!/bin/bash

SECRET_PATH=$(vc index | dmenu)
SECRET_KEY=$(vc show $SECRET_PATH | awk '{ print $1 }' | dmenu)
echo $(vc show $SECRET_PATH | grep $SECRET_KEY | awk '{ print $3 }')
