#compdef vc

# vc completion for zsh

_vc() {
    local state

    _arguments \
        '1: :->action'\
        '*: :->secret'

    case $state in
        (action) _arguments '1:action:(cp edit insert login ls mv rm show)';;
        (*) compadd "$@" $(vc ls -r) ;;
    esac

}

_vc
