#!/bin/bash

VAULT_INDEX="/home/patrick/.cache/vaultindex"

path=$(cat $VAULT_INDEX | dmenu)
secretkey=$(vc show $path | awk '{ print $1 }' | dmenu)
