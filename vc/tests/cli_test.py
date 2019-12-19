import pytest
import re


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
        ("", 0, "secret0\nsecret1\nsecret2\nsecret3\nsecretdir/\n"),
        ("nonexistent", 1, 'Path ".*" does not exist.\n'),
    ],
)
def test_list(run_cmd, backend, path, rc, output):
    result = run_cmd(["list", f"{backend}/{path}"])
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
    ],
)
def test_insert(run_cmd, backend, path, data, rc, output):
    result = run_cmd(["insert", f"{backend}/{path}", data])
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
    result = run_cmd([command, path, *args])
    assert result.exit_code == 1
    assert result.output.endswith("is not under a valid mount point.\n")
