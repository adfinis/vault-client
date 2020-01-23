# coding: future_fstrings

import hvac

from contextlib import contextmanager
from collections import namedtuple


VaultPath = namedtuple("VaultPath", ["mount_path", "secret_path", "kv_version"])


class MountNotFound(Exception):
    pass


class LoginFailed(Exception):
    pass


class KvClient:
    """hvac client wrapper that transparently can deal with different version of
    of vault"""

    def __init__(self, url, verify_tls=True):
        self.client = hvac.Client(url, verify=verify_tls)

    @contextmanager
    def vault_path(self, path):
        kv_mounts = self._get_kv_mounts()
        match = ""
        for m in kv_mounts.keys():
            if path.lstrip("/").startswith(m) and len(m) > len(match):
                match = m

        if not match:
            raise MountNotFound(path)

        try:
            version = kv_mounts[match]["options"]["version"]
        except TypeError:
            # Old secret engines do not use options.
            version = "1"

        yield VaultPath(
            mount_path=match, secret_path=path[len(match) :], kv_version=version
        )

    def _get_kv_mounts(self):
        mounts = self.client.sys.list_mounted_secrets_engines()["data"]
        return {k: v for k, v in mounts.items() if v["type"] == "kv"}

    def set_token(self, token):
        self.client.token = token
        if not self.client.is_authenticated():
            raise LoginFailed()

    def login(self, user, password, mount_point, auth_type):
        try:
            if auth_type == 'ldap':
                resp = self.client.auth.ldap.login(
                    username=user, password=password, mount_point=mount_point
                )
            elif auth_type == 'userpass':
                resp = self.client.auth.userpass.login(
                    username=user, password=password
                )
        except hvac.exceptions.InvalidRequest as e:
            raise LoginFailed(e)

        if not self.client.is_authenticated():
            raise LoginFailed()
        token = resp["auth"]["client_token"]
        self.set_token(token)
        return token

    def get(self, path):

        with self.vault_path(path) as vpath:
            if vpath.kv_version == "1":
                secret = self.client.secrets.kv.v1.read_secret(
                    mount_point=vpath.mount_path, path=vpath.secret_path
                )
            elif vpath.kv_version == "2":
                secret = self.client.secrets.kv.v2.read_secret_version(
                    mount_point=vpath.mount_path, path=vpath.secret_path
                )["data"]
            else:
                raise NotImplementedError

            return secret["data"]

    def list(self, path):

        if path in ["", "/", None]:
            return self._get_kv_mounts().keys()

        with self.vault_path(path) as vpath:
            try:
                if vpath.kv_version == "1":
                    secrets = self.client.secrets.kv.v1.list_secrets(
                        mount_point=vpath.mount_path, path=vpath.secret_path
                    )
                elif vpath.kv_version == "2":
                    secrets = self.client.secrets.kv.v2.list_secrets(
                        mount_point=vpath.mount_path, path=vpath.secret_path
                    )
                else:
                    raise NotImplementedError

                return secrets["data"]["keys"]

            except hvac.exceptions.Forbidden:
                return []

            except hvac.exceptions.InvalidPath as exc:
                # Listing an empty backend will result in an InvalidPath
                # exception. I think an empty list is more appropriate.
                if vpath.secret_path == "/":
                    return []
                raise exc


    def traverse(self, path):
        paths = []
        childs = self.list(path)
        for child in childs:

            if path:
                full_path = f"{path}{child}"
            else:
                full_path = child

            if child.endswith("/"):
                paths += self.traverse(full_path)
            else:
                paths.append(full_path)

        return paths

    def put(self, path, data):

        with self.vault_path(path) as vpath:
            if vpath.kv_version == "1":
                return self.client.secrets.kv.v1.create_or_update_secret(
                    mount_point=vpath.mount_path, path=vpath.secret_path, secret=data
                )
            elif vpath.kv_version == "2":
                return self.client.secrets.kv.v2.create_or_update_secret(
                    mount_point=vpath.mount_path, path=vpath.secret_path, secret=data
                )
            else:
                raise NotImplementedError

    def delete(self, path):

        with self.vault_path(path) as vpath:
            if vpath.kv_version == "1":
                return self.client.secrets.kv.v1.delete_secret(
                    mount_point=vpath.mount_path, path=vpath.secret_path
                )
            elif vpath.kv_version == "2":
                self.client.secrets.kv.v2.delete_secret_versions(
                    mount_point=vpath.mount_path, path=vpath.secret_path, versions=[1]
                )
            else:
                raise NotImplementedError
