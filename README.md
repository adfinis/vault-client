Vault Client
============
[![Travis](https://img.shields.io/travis/adfinis-sygroup/vault-client.svg?style=flat-square)](https://travis-ci.org/adfinis-sygroup/vault-client)
[![License](https://img.shields.io/github/license/adfinis-sygroup/vault-client.svg?style=flat-square)](LICENSE)

`vc` is a command-line interface to
[HashiCorp's Vault](https://www.vaultproject.io/) inspired by
[pass](https://www.passwordstore.org/).

* Makes secrets from `generic` backends easy accessible
* Features auto completion for `bash`

Demo
----
![gif](sample/demo.gif)

Installation
------------
1. Download the
[latest release](https://github.com/adfinis-sygroup/vault-client/releases).
2. Unzip and move `vc` into a directory of choice.

Build Instructions
------------------
To build vault-client you need a Go compiler and Git.
```
$ apt-get install git go
$ git clone https://github.com/adfinis-sygroup/vault-client.git
$ cd vault-client
$ make build
```
`make build` will install Go dependencies and build vault-client. After you
should have a binary `vc` in your working directory.

Configuration
-------------
The configuration happens through a simple yaml file.
```
$ echo "host: 127.0.0.1
port: 8200
token: password
backend_method: ldap
backend_auth: ldap
tls: true
verify_tls: true" >  ~/.vaultrc
$ chmod 600 ~/.vaultrc
```

Contributions
-------------
Contributions are more than welcome! Please feel free to open new issues or
pull requests.

License
-------
Copyright (C) 2017  Adfinis SyGroup AG

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

See the	[LICENSE](LICENSE) file.
