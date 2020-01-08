# coding: future_fstrings

import yaml

from os import getenv
from os.path import exists, expanduser

_DEFAULT_CONFIG = {"host": "localhost", "port": 8200, "tls": True, "verify_tls": True}


def load_config():
    path = config_path()
    try:
        with open(path, "r") as f:
            cfg = _DEFAULT_CONFIG.copy()
            cfg.update(
                yaml.load(f, Loader=yaml.FullLoader)
            )
            return cfg
    except FileNotFoundError:
        question = "Would you like to copy a sample config file to ~/.vaultrc? [y/N]: "
        reply = str(input(question)).lower().strip()
        if reply[0] == "y":
            create_default_config(path)
        exit(0)


def create_default_config(path):
    with open(path, "w") as f:
        yaml.dump(_DEFAULT_CONFIG, f)


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
