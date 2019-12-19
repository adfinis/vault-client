import yaml

from os import getenv
from os.path import exists

_DEFAULT_CONFIG = {"host": "localhost", "port": 8200, "tls": True, "verify_tls": True}


def load_config():
    path = config_path()
    with open(path, "r") as f:
        return {**_DEFAULT_CONFIG, **yaml.load(f, Loader=yaml.FullLoader)}


def config_path():
    envvar = getenv("VC_CONFIG")
    if envvar:
        return envvar
    elif exists(".vaultrc"):
        return ".vaultrc"
    return expanduser("~/.vaultrc")


def update_config_token(token):
    path = config_path()
    with open(path, "r") as f:
        config = f.readlines()

    updated = []
    for line in config:
        if line.startswith("token:"):
            updated += f"token: {token}\n"
            continue
        updated += line

    with open(path, "w") as f:
        f.writelines(updated)
