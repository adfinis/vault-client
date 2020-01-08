# coding: future_fstrings

from vc.cli import cli
from vc.config import load_config
from vc.kv_client import KvClient


config = load_config()

if config.get("tls"):
    protocol = "https"
else:
    protocol = "http"

url = f"{protocol}://{config['host']}:{config['port']}"

client = KvClient(url, config.get("verify_tls"))
client.set_token(config["token"])

cli(obj={"client": client, "config": config})
