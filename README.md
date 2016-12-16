Vault Client
------------
`vc` is a command-line interface to (HashiCorp's Vault)[https://www.vaultproject.io/] inspired by (pass)[https://www.passwordstore.org/]. 

* Makes secrets from `generic` backends easy accessible
* Features auto completion for `bash`

Demo
----

![gif](sample/demo.gif)

Configuration
-------------
```
$ echo ~/.vaultrc <<EOF
host: 127.0.0.1 
port: 9200
user: admin
EOF
$ chmod 600 ~/.vaultrc
```
