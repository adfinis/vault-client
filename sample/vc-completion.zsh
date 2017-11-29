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
#

_vc() {
    local state subcmds

    subcmds=(
        'cp:Copy an existing secret to another location'
        'edit:Edit a secret at specified path'
        'insert:Insert an new secret'
        'login:Authenticate against Vault using your prefered method'
        'ls:List all secrets at specified path'
        'mv:Move an existing secret to another location'
        'rm:Remove a secret at specified path'
        'show:Show an existing secret'
    )

    _arguments -C -s -S -n \
        '(- 1 *)'{-v,--version}"[Show program\'s version number and exit]: :->full" \
        '(- 1 *)'{-h,--help}'[Show help message and exit]: :->full' \
        '1: :->action' \
        '*: :->secret'

    case "$state" in
        (action)
            _describe -t subcmds 'subcmds' subcmds
            ;;
        (*)
            compadd "$@" $(vc ls -r)
            ;;
    esac
}

_vc

# vim: set ft=zsh sw=4 ts=4 et wrap tw=76:
