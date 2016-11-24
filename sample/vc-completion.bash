_vc1()
{
    local cur=${COMP_WORDS[COMP_CWORD]}
    if [[ $COMP_CWORD -gt 1 ]]; then
        COMPREPLY=( $(compgen -W "$(cat ~/.cache/vaultindex)" -- $cur) )
    else
        COMPREPLY=( $(compgen -W "rm cp mv edit index version help" -- $cur) )
    fi

}
_vc()
{
    local cur=${COMP_WORDS[COMP_CWORD]}
    COMPREPLY=( $(compgen -W "$(cat ~/.cache/vaultindex) rm cp mv edit index version help" -- $cur) )
    vc index > /dev/null
}
complete -F _vc vc
