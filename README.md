Vault Client
============
[![Travis](https://img.shields.io/travis/adfinis-sygroup/vault-client.svg?style=flat-square)](https://travis-ci.org/adfinis-sygroup/vault-client)
[![License](https://img.shields.io/github/license/adfinis-sygroup/vault-client.svg?style=flat-square)](LICENSE)

`vc` is a command-line interface to
[HashiCorp's Vault](https://www.vaultproject.io/) inspired by
[pass](https://www.passwordstore.org/).

* Makes secrets from `kv` backends easy accessible (`v1` and `v2`)
* Features auto completion for `bash` and `zsh`
* Supports `userpass` and `ldap` authentication backends

Demo
----
![gif](sample/demo.gif)


Configuration
-------------
The configuration happens through a simple yaml file.
```
$ echo "host: localhost
port: 8200
tls: false
verify_tls: false
authentication:
  type: userpass
  user: someuser" >  ~/.vaultrc

$ chmod 600 ~/.vaultrc
```

Development
-----------
If you would like to hack on `vc` you require:
- Python >= 3.4
- Docker and docker-compose

You can then start the Vault container using docker-compose, install all
required dependencies and run the test suite:
```
$ docker-compose up -d
$ pip install -r requirements.txt
$ pytest --isort --black --flake8
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
