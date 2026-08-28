#!/usr/bin/env bash
# Custom Ignite startup: Java 17 open flags, optional lib loading, exec so JVM is PID 1.
set -e

IGNITE_HOME="${IGNITE_HOME:-/opt/ignite}"
SCRIPTS_HOME="${IGNITE_HOME}/bin"

. "${SCRIPTS_HOME}/include/setenv.sh"

# Append optional libs specified in OPTION_LIBS (comma-separated directory names
# under libs/optional/). The stock setenv.sh doesn't process this — we do it here.
if [ -n "${OPTION_LIBS}" ]; then
    IFS=',' read -ra _OPTS <<< "${OPTION_LIBS}"
    for _lib in "${_OPTS[@]}"; do
        _dir="${IGNITE_HOME}/libs/optional/${_lib}"
        if [ -d "${_dir}" ]; then
            for _jar in "${_dir}"/*.jar; do
                [ -f "${_jar}" ] && IGNITE_LIBS="${IGNITE_LIBS}:${_jar}"
            done
        fi
    done
fi

# Also append USER_LIBS if set (allows passing extra JARs via env).
if [ -n "${USER_LIBS}" ]; then
    IGNITE_LIBS="${USER_LIBS}:${IGNITE_LIBS}"
fi

CONFIG="${1:-${IGNITE_HOME}/config/default-config.xml}"

exec java \
    ${JVM_OPTS:-} \
    --add-exports=java.base/jdk.internal.misc=ALL-UNNAMED \
    --add-exports=java.base/sun.nio.ch=ALL-UNNAMED \
    --add-exports=java.management/com.sun.jmx.mbeanserver=ALL-UNNAMED \
    --add-exports=jdk.internal.jvmstat/sun.jvmstat.monitor=ALL-UNNAMED \
    --add-exports=java.base/sun.reflect.generics.reflectiveObjects=ALL-UNNAMED \
    --add-opens=jdk.management/com.sun.management.internal=ALL-UNNAMED \
    --add-opens=java.base/jdk.internal.access=ALL-UNNAMED \
    --add-opens=java.base/java.nio=ALL-UNNAMED \
    --add-opens=java.base/jdk.internal.misc=ALL-UNNAMED \
    --add-opens=java.base/sun.net.util=ALL-UNNAMED \
    --add-opens=jdk.internal.jvmstat/sun.jvmstat.monitor=ALL-UNNAMED \
    --add-opens=java.management/com.sun.jmx.mbeanserver=ALL-UNNAMED \
    --add-opens=java.base/sun.reflect.generics.reflectiveObjects=ALL-UNNAMED \
    -Djava.net.preferIPv4Stack=true \
    -DIGNITE_HOME="${IGNITE_HOME}" \
    -cp "${IGNITE_LIBS}" \
    org.apache.ignite.startup.cmdline.CommandLineStartup \
    "${CONFIG}"
