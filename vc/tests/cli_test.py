import pytest

def test_show(v1_backend, run_cmd):
    result = run_cmd(['show', f'{v1_backend}/secret1'])
    assert result.exit_code == 0
    assert result.output == 'key: value\n\n'
