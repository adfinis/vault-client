from cli import cli
from config import load_config
from kv_client import KvClient


config = load_config()

if config.get('tls'):
    protocol = 'https'
else:
    protocol = 'http'

url = f"{protocol}://{config['host']}:{config['port']}"

client = KvClient(
    url,
    config.get('verify_tls')
)
client.set_token(config['token'])

cli(
    obj={
        'client': client,
        'config': config
    }
)
