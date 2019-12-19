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
