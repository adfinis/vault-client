import pytest
import re


@pytest.mark.parametrize(
    "backend", [pytest.lazy_fixture("v1_backend"), pytest.lazy_fixture("v2_backend")]
)
@pytest.mark.parametrize(
    "command,path,rc,output",
    [
        ("show", "secret1", 0, "key: value\n\n"),
        ("show", "nonexistent", 1, 'Path ".*" does not exist.\n'),
    ],
)
def test_commands(run_cmd, command, backend, path, rc, output):
    result = run_cmd([command, f"{backend}/{path}"])
    assert result.exit_code == rc
    assert re.match(output, result.output)


@pytest.mark.parametrize(
    "command,args",
    [
        ("delete", []),
        ("edit", []),
        ("list", []),
        ("show", []),
        ("cp", ["othersecret"]),
        ("mv", ["othersecret"]),
        ("insert", ["key=val"])
    ]
)
def test_nonexistent_mountpoint(run_cmd, command, args):
    path = "nonexistent/secret"
    result = run_cmd([command, path, *args])
    assert result.exit_code == 1
    assert result.output.endswith('is not under a valid mount point.\n')
