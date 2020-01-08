# coding: future_fstrings

import pytest
import re


# TODO: Test search command
# TODO: Test login command


@pytest.mark.parametrize(
    "backend", [pytest.lazy_fixture("v1_backend"), pytest.lazy_fixture("v2_backend")]
)
@pytest.mark.parametrize(
    "path,rc,output",
    [
        ("secret1", 0, "key: value\n\n"),
        ("nonexistent", 1, 'Path ".*" does not exist.\n'),
    ],
)
def test_show(run_cmd, backend, path, rc, output):
    result = run_cmd(["show", f"{backend}/{path}"])
    assert re.match(output, result.output)
    assert result.exit_code == rc


@pytest.mark.parametrize(
    "backend", [pytest.lazy_fixture("v1_backend"), pytest.lazy_fixture("v2_backend")]
)
@pytest.mark.parametrize(
    "path,rc,output",
    [
        ("secret1", 0, "Secret successfully deleted\n"),
        ("nonexistent", 1, 'Path ".*" does not exist.\n'),
    ],
)
def test_delete(run_cmd, backend, path, rc, output):
    result = run_cmd(["delete", f"{backend}/{path}"])
    assert re.match(output, result.output)
    assert result.exit_code == rc


@pytest.mark.parametrize(
    "backend", [pytest.lazy_fixture("v1_backend"), pytest.lazy_fixture("v2_backend")]
)
@pytest.mark.parametrize(
    "path,rc,output",
    [
        # TODO: Should we mock click.edit()?
        ("nonexistent", 1, "Path .* does not yet exist. Creating a new secret.\n")
    ],
)
def test_edit(run_cmd, backend, path, rc, output):
    result = run_cmd(["edit", f"{backend}/{path}"])
    assert re.match(output, result.output)
    assert result.exit_code == rc


@pytest.mark.parametrize(
    "backend", [pytest.lazy_fixture("v1_backend"), pytest.lazy_fixture("v2_backend")]
)
@pytest.mark.parametrize(
    "path,opts,rc,output",
    [
        ("", [], 0, "secret0\nsecret1\nsecret2\nsecret3\nsecretdir/\n"),
        ("nonexistent", [], 1, 'Path ".*" does not exist.\n'),
        (
            "",
            ["-r"],
            0,
            "test/v[1,2]/secret0\ntest/v[1,2]/secret[1,2]\ntest/v[1,2]/secret2\ntest/v[1,2]/secret3\ntest/v[1,2]/secretdir/subsecret0\ntest/v[1,2]/secretdir/subsecret[1,2]\ntest/v[1,2]/secretdir/subsecret2\ntest/v[1,2]/secretdir/subsecret3\n",
        ),
    ],
)
def test_list(run_cmd, backend, path, opts, rc, output):
    cmd = ["list"]
    cmd.extend(opts)
    cmd.append(f"{backend}/{path}")
    result = run_cmd(cmd)
    assert re.match(output, result.output)
    assert result.exit_code == rc


@pytest.mark.parametrize(
    "backend", [pytest.lazy_fixture("v1_backend"), pytest.lazy_fixture("v2_backend")]
)
@pytest.mark.parametrize(
    "path,data,rc,output",
    [
        ("secret1", "key=value", 0, "Secret successfully inserted!\n"),
        ("secret1", "key=key=value", 1, "Data .* is not a valid key/value pair.\n"),
        ("secret1", "key", 1, "Data .* is not a valid key/value pair.\n"),
    ],
)
def test_insert(run_cmd, backend, path, data, rc, output):
    result = run_cmd(["insert", f"{backend}/{path}", data])
    assert re.match(output, result.output)
    assert result.exit_code == rc


@pytest.mark.parametrize(
    "backend", [pytest.lazy_fixture("v1_backend"), pytest.lazy_fixture("v2_backend")]
)
@pytest.mark.parametrize(
    "src,dest,rc,output",
    [
        ("secret1", "newsecret", 0, "Secret successfully copied!\n"),
        ("nonexistent", "newsecret", 1, "Source path .* does not exist.\n"),
        ("secret1", "secret2", 1, "The destination secret already exists..*"),
    ],
)
def test_copy(run_cmd, backend, src, dest, rc, output):
    result = run_cmd(["cp", f"{backend}/{src}", f"{backend}/{dest}"])
    assert re.match(output, result.output)
    assert result.exit_code == rc


@pytest.mark.parametrize(
    "backend", [pytest.lazy_fixture("v1_backend"), pytest.lazy_fixture("v2_backend")]
)
@pytest.mark.parametrize(
    "src,dest,rc,output",
    [
        ("secret1", "newsecret", 0, "Secret successfully copied!\n"),
        ("nonexistent", "newsecret", 1, "Source path .* does not exist.\n"),
        ("secret1", "secret2", 1, "The destination secret already exists..*"),
    ],
)
def test_move(run_cmd, backend, src, dest, rc, output):
    result = run_cmd(["cp", f"{backend}/{src}", f"{backend}/{dest}"])
    assert re.match(output, result.output)
    assert result.exit_code == rc


@pytest.mark.parametrize(
    "command,args",
    [
        ("delete", []),
        ("edit", []),
        ("list", []),
        ("show", []),
        ("cp", ["othersecret"]),
        ("mv", ["othersecret"]),
        ("insert", ["key=val"]),
    ],
)
def test_nonexistent_mountpoint(run_cmd, command, args):
    path = "nonexistent/secret"
    cmd = [command, path]
    cmd.extend(args)
    result = run_cmd(cmd)
    assert result.exit_code == 1
    assert result.output.endswith("is not under a valid mount point.\n")
