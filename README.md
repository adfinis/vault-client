Vault Client
------------
`vc` is command-line interface to vault inspired by `pass`. 

* Makes secrets from the `generic` backend easiable accessible from the command-line
* Features autocompletion for bash

Configuration
-------------
```
$ echo ~/.vaultrc <<EOF
host: 127.0.0.1 
port: 9200
cacert: /etc/ssl/cert.pem
user: admin
gpgid: 0x518417425D442A7A
encrypted_password: "password"
index_file: ~/dev/shm/vaultindex
EOF
$ chmod 600 ~/.vaultrc
```
