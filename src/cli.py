import click
import hvac

from kv_client import MountNotFound


@click.group()
def cli():
    pass

@cli.command()
@click.argument('query')
@click.pass_context
def search(ctx, query):
    client = ctx.obj['client']
    try:
        paths = client.traverse()
    except hvac.exceptions.InvalidPath:
        click.echo(f'Path "{path}" does not exist.', err=True)
    except MountNotFound:
        click.echo(f'Path "{path}" is not under a valid mount point.', err=True)

    results = []
    for path in paths:
        if query in path:
            results.append(path)

    if len(results) == 0:
        click.echo("No search results.")
    elif len(results) == 1:
        secret = client.get(path)
        click.echo(f"# {path}")
        click.echo(yaml.dump(secret))
    else:
        for path in results:
            click.echo(path)

@cli.command()
@click.argument('path')
@click.pass_context
def show(ctx, path):
    client = ctx.obj['client']
    try:
        secret = client.get(path)
        click.echo(yaml.dump(secret))
    except hvac.exceptions.InvalidPath:
        click.echo(f'Path "{path}" does not exist.', err=True)
    except MountNotFound:
        click.echo(f'Path "{path}" is not under a valid mount point.', err=True)

@cli.command()
@click.argument('src')
@click.argument('dest')
@click.pass_context
def mv(ctx, src, dest):
    client = ctx.obj['client']
    try:
        secret = client.get(src)
    except hvac.exceptions.InvalidPath:
        click.echo(f'Source path "{src}" does not exist.', err=True)
        return
    except MountNotFound:
        click.echo(f'Source path "{src}" is not under a valid mount point.', err=True)
        return

    try:
        secret = client.get(dest)
        click.echo('The destination secret already exists.')
        if not click.confirm('Do you want overwrite it?', abort=True):
            return

        client.delete(dest)
    except hvac.exceptions.InvalidPath:
        pass

    except MountNotFound:
        click.echo(f'Source path "{path}" is not under a valid mount point.', err=True)
        return

    client.put(dest, secret)
    client.delete(src)
    click.echo('Secret successfully moved!')

@cli.command()
@click.argument('src')
@click.argument('dest')
@click.pass_context
def cp(ctx, src, dest):
    client = ctx.obj['client']
    try:
        secret = client.get(src)
    except hvac.exceptions.InvalidPath:
        click.echo(f'Source path "{src}" does not exist.', err=True)
        return
    except MountNotFound:
        click.echo(f'Source path "{src}" is not under a valid mount point.', err=True)
        return

    try:
        secret = client.get(dest)
        click.echo('The destination secret already exists.')
        if not click.confirm('Do you want overwrite it?', abort=True):
            return

        client.delete(dest)
    except hvac.exceptions.InvalidPath:
        pass

    except MountNotFound:
        click.echo(f'Destination path "{path}" is not under a valid mount point.', err=True)
        return

    client.put(dest, secret)
    click.echo('Secret successfully copied!')

@cli.command()
@click.argument('path')
@click.pass_context
def edit(ctx, path):
    client = ctx.obj['client']
    secret = {}
    try:
        secret = client.get(path)
    except hvac.exceptions.InvalidPath:
        click.echo(f'Path "{path}" does not yet exist. Creating a new secret.')
    except MountNotFound:
        click.echo(f'Path "{path}" is not under a valid mount point.', err=True)
        return

    if secret:
        edited = click.edit(yaml.dump(secret))
    else:
        edited = click.edit()

    data = yaml.load(edited, Loader=yaml.FullLoader)
    client.put(path, data)
    click.echo("Secret successfully edited!")


@cli.command()
@click.argument('path')
@click.argument('data')
@click.pass_context
def insert(ctx, path, data):
    client = ctx.obj['client']
    kv_pair = data.split("=")
    if len(kv_pair) != 2:
        click.echo("Data is not a valid key/value pair.", err=True)
        return


    try:
        secret = client.put(path, {kv_pair[0]: kv_pair[1]})
        click.echo("Secret successfully inserted!")
    except hvac.exceptions.InvalidPath:
        click.echo(f'Path "{path}" does not exist.', err=True)
    except MountNotFound:
        click.echo(f'Path "{path}" is not under a valid mount point.', err=True)

@cli.command()
@click.argument('path', required=False)
@click.option('-r', '--recursive/--no-recursive', default=False)
@click.pass_context
def list(ctx, path, recursive):
    client = ctx.obj['client']
    try:
        if recursive:
            paths = client.traverse(path)
        else:
            paths = client.list(path)

        for p in paths:
            click.echo(p)
    except hvac.exceptions.InvalidPath:
        click.echo(f'Path "{path}" does not exist.', err=True)
    except MountNotFound:
        click.echo(f'Path "{path}" is not under a valid mount point.', err=True)

@cli.command()
@click.argument('path', required=False)
@click.pass_context
def delete(ctx, path):
    client = ctx.obj['client']
    try:
        client.get(path)
        client.delete(path)
        click.echo("Secret successfully deleted")
    except hvac.exceptions.InvalidPath:
        click.echo(f'Path "{path}" does not exist.', err=True)
    except MountNotFound:
        click.echo(f'Path "{path}" is not under a valid mount point.', err=True)
