#/usr/bin/env bash

_vc()
{
    COMREPLY=()
    local CUR=${COMP_WORDS[COMP_CWORD]}
    if [[ $COMP_CWORD -gt 1 ]]; then
	case "${COMP_WORDS[1]}" in
	    show|insert|edit|rm|ls|cp|mv)

		local DIR
		if [[ -n ${CUR} ]]; then
		    # Remove everything after the last slash
		    BASE=$(echo $CUR | sed 's/\(.*\/\).*$/\1/')
		else
		    BASE=""
		fi

		COMPREPLY=(
		    $(compgen -W "$(vc ls -a ${BASE})" ${CUR})
		)
		compopt -o nospace
	esac
    else
        COMPREPLY=($(compgen -W "login show insert cp mv rm edit ls" -- $CUR))
    fi
}
complete -F _vc vc
