# coding: future_fstrings

from functools import partial

import hvac
import pytest
from click.testing import CliRunner

from vc.cli import cli
from vc.config import load_config
from vc.kv_client import KvClient


@pytest.fixture()
def hvac_client():
    return hvac.Client("http://localhost:8200")


@pytest.fixture()
def v1_backend(hvac_client):
    path = "test/v1"
    hvac_client.sys.enable_secrets_engine(
        backend_type="kv", path=path, options={"version": 1}
    )

    kv_data = {"key": "value"}
    for i in range(4):
        hvac_client.secrets.kv.v1.create_or_update_secret(
            mount_point=path, path=f"secret{i}", secret=kv_data
        )

    for i in range(4):
        hvac_client.secrets.kv.v1.create_or_update_secret(
            mount_point=path, path=f"secretdir/subsecret{i}", secret=kv_data
        )

    yield path
    hvac_client.sys.disable_secrets_engine(path=path)


# TODO: There is probably a better way to share some code with v1_backend.
@pytest.fixture()
def v2_backend(hvac_client):
    path = "test/v2"
    hvac_client.sys.enable_secrets_engine(
        backend_type="kv", path=path, options={"version": 2}
    )

    kv_data = {"key": "value"}
    for i in range(4):
        hvac_client.secrets.kv.v2.create_or_update_secret(
            mount_point=path, path=f"secret{i}", secret=kv_data
        )

    for i in range(4):
        hvac_client.secrets.kv.v2.create_or_update_secret(
            mount_point=path, path=f"secretdir/subsecret{i}", secret=kv_data
        )

    yield path
    hvac_client.sys.disable_secrets_engine(path=path)


@pytest.fixture(scope="function", autouse=True)
def cleanup_kv_secret_engines(hvac_client):
    """Unmount test backends if last run did not terminate successfully."""
    path = "test/v1"
    hvac_client.sys.disable_secrets_engine(path=path)


@pytest.fixture()
def config():
    return {
        "host": "localhost",
        "port": 8200,
        "tls": False,
        "verify_tls": False,
        "token": "password",
    }


@pytest.fixture()
def kv_client(config):
    protocol = "https" if config.get("tls") else "http"
    url = f"{protocol}://{config['host']}:{config['port']}"
    client = KvClient(url, config.get("verify_tls"))
    client.set_token(config["token"])
    return client


@pytest.fixture()
def ctx(config, kv_client):
    return {"client": kv_client, "config": config}


@pytest.fixture()
def run_cmd(ctx):
    runner = CliRunner()
    return partial(runner.invoke, cli, obj=ctx)


@pytest.fixture()
def userpass_auth_backend(hvac_client):
    path = "vc_userpass"
    hvac_client.sys.enable_auth_method(method_type="userpass", path=path)
    yield path
    hvac_client.sys.disable_auth_method(path=path)


@pytest.fixture()
def userpass_credentials(hvac_client, userpass_auth_backend):
    path = userpass_auth_backend
    username = "user"
    password = "password"
    hvac_client.auth.userpass.create_or_update_user(
        username=username, password=password, mount_point=path
    )
    yield (username, password, path)
