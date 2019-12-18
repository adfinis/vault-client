import yaml


def load_config():
    with open(".vaultrc", "r") as f:
        return yaml.load(f, Loader=yaml.FullLoader)
