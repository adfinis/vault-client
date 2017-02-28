#compdef vc
# ------------------------------------------------------------------------------
# vault-client.src
# Copyright (C) 2017  Adfinis SyGroup AG
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
# ------------------------------------------------------------------------------
# Description
# -----------
#
#  Completion script for vault-client
#
# ------------------------------------------------------------------------------
# Authors
# -------
#
#  * Adfinis SyGroup AG
#
# ------------------------------------------------------------------------------

local -a commands

commands=(
    'show:show secret'
    'insert:insert new secret'
    'edit:edit secret'
    'rm:remove secret'
    'ls:list secrets in directory'
    'cp:copy secret'
    'mv:move secret'
)

_arguments -C -s -S -n \
    '(- 1 *)'{-v,--version}"[Show program\'s version number and exit]: :->full" \
    '(- 1 *)'{-h,--help}'[Show help message and exit]: :->full' \
    '1:cmd:->cmds' \
    '*:: :->args' \

case "$state" in
    (cmds)
        _describe -t commands 'commands' commands
        ;;
    (args)
        dir=$words[2]
        _call_program vc vc ls $dir
        ;;
    (*)
        ;;
esac

# vim: set ft=zsh sw=4 ts=4 et wrap tw=76:
