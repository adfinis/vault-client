from setuptools import setup, find_packages

setup(
    name="vc",
    version="2.0.0",
    url="https://github.com/adfinis-sygroup/vault-client.git",
    author="Patrick Winter",
    author_email="info@adfinis-sygroup.ch",
    description="A command-line interface to HashiCorp's Vault",
    packages=find_packages(),
    install_requires=[
        "PyYAML >= 5.2",
        "Click >= 7.0",
        "hvac >= 0.9.6",
        "future-fstrings",
    ],
    entry_points={"console_scripts": ["vc = vc.__main__:main"]},
)
