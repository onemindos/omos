#!/usr/bin/env bash
# Custom Ignite startup using exec so the JVM is PID 1 (proper SIGTERM handling in K8s).

set -e

IGNITE_HOME="${IGNITE_HOME:-/opt/ignite}"
SCRIPTS_HOME="${IGNITE_HOME}/bin"

. "${SCRIPTS_HOME}/include/setenv.sh"

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
    --illegal-access=permit \
    -Djava.net.preferIPv4Stack=true \
    -DIGNITE_HOME="${IGNITE_HOME}" \
    -cp "${IGNITE_LIBS}" \
    org.apache.ignite.startup.cmdline.CommandLineStartup \
    "${CONFIG}"
