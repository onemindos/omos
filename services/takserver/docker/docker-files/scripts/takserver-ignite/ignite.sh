#!/usr/bin/env bash
# Wrapper: converts OPTION_LIBS → USER_LIBS (the standard Ignite setenv.sh mechanism),
# then execs the stock ignite.sh so we get all the correct Java-17 --add-opens flags.
set -e

IGNITE_HOME="${IGNITE_HOME:-/opt/ignite}"

# Build USER_LIBS from OPTION_LIBS (comma-separated optional lib directory names).
# The stock setenv.sh prepends USER_LIBS to IGNITE_LIBS if set.
if [ -n "${OPTION_LIBS}" ]; then
    IFS=',' read -ra _OPTS <<< "${OPTION_LIBS}"
    for _lib in "${_OPTS[@]}"; do
        _dir="${IGNITE_HOME}/libs/optional/${_lib}"
        if [ -d "${_dir}" ]; then
            for _jar in "${_dir}"/*.jar; do
                [ -f "${_jar}" ] && USER_LIBS="${USER_LIBS:+${USER_LIBS}:}${_jar}"
            done
        fi
    done
    export USER_LIBS
fi

exec "${IGNITE_HOME}/bin/ignite-stock.sh" "$@"
