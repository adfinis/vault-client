from cli import cli
from config import load_config
from kv_client import KvClient

config = load_config()
client = KvClient(f"http://{config['host']}:{config['port']}")
cli(
    obj={
        'client': client,
        'config': config
    }
)
