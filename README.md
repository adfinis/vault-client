Vault Client
------------
`vc` is a command-line interface to [HashiCorp's Vault](https://www.vaultproject.io/) inspired by [pass](https://www.passwordstore.org/).

* Makes secrets from `generic` backends easy accessible
* Features auto completion for `bash`

Demo
----

![gif](sample/demo.gif)

Installation
------------
1. Download the [latest release](https://github.com/adfinis-sygroup/vault-client/releases).
2. Unzip and move `vc` into a directory of choice.

Configuration
-------------
The configuration happens through a simple yaml file.
```
$ echo ~/.vaultrc <<EOF
host: 127.0.0.1 
port: 8200
token: password
EOF
$ chmod 600 ~/.vaultrc
```

Contributions
-------------
Contributions are more than welcome! Please feel free to open new issues or
pull requests.

License
-------
GNU GENERAL PUBLIC LICENSE Version 3

See the	[LICENSE](LICENSE) file.
