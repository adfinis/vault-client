# coding: future_fstrings

from sys import argv

from vc.ansi import color
from vc.cli import cli, login
from vc.config import load_config
from vc.kv_client import KvClient

config = load_config()

if config.get("tls"):
    protocol = "https"
else:
    protocol = "http"

url = f"{protocol}://{config['host']}:{config['port']}"
client = KvClient(url, config.get("verify_tls"))

if len(argv) != 1 and argv[1] != 'login':
    token = config.get('token')
    if not token:
        print(f'{color.BOLD}You do not have a token set. Please login first (vc login).{color.END}\n')
        exit(0)
    else:
        client.set_token(token)

ctx = {"client": client, "config": config}
cli(obj=ctx)
