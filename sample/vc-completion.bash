_vc()
{
    local cur=${COMP_WORDS[COMP_CWORD]}
    if [[ $COMP_CWORD -gt 1 ]]; then
        COMPREPLY=( $(compgen -W "$(cat ~/.cache/vaultindex)" -- $cur) )
    else
        COMPREPLY=( $(compgen -W "ls insert rm cp mv edit index version help" -- $cur) )
    fi

}
complete -F _vc vc
