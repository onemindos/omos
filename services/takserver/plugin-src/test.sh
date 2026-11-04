#!/usr/bin/env bash
# Compiles and tests the SDK-free part of NatsBridgePlugin (BridgeCore, NatsLink) with
# plain javac and JUnit 4, so it runs without the TAK Server SDK JAR and without Gradle.
# Needs: JDK 17+ (javac/java on PATH or JAVA_HOME), curl, and nats-server on PATH or
# NATS_SERVER_BIN for the NatsLink tests (they are skipped without one).
# NatsBridgePlugin.java itself needs the SDK and is compiled only by Gradle (README.md).
set -euo pipefail
cd "$(dirname "$0")"
JAVA_BIN="${JAVA_HOME:+$JAVA_HOME/bin/}"
LIB="${PLUGIN_TEST_LIB:-build/test-lib}"
OUT=build/test-classes-sh
mkdir -p "$LIB"
fetch() { # path sha256
  local f
  f="$LIB/$(basename "$1")"
  if [ ! -f "$f" ]; then curl -sSfL -o "$f.tmp" "https://repo1.maven.org/maven2/$1" && mv "$f.tmp" "$f"; fi
  echo "$2  $f" | shasum -a 256 -c --quiet - || { echo "checksum mismatch: $f" >&2; exit 1; }
}
fetch io/nats/jnats/2.20.2/jnats-2.20.2.jar 421d51c2a341ae7e109aafd6291c6ccd3fa624242a6dbbc61dca2a598f55cf7d
fetch org/slf4j/slf4j-api/2.0.13/slf4j-api-2.0.13.jar e7c2a48e8515ba1f49fa637d57b4e2f590b3f5bd97407ac699c3aa5efb1204a9
fetch org/slf4j/slf4j-simple/2.0.13/slf4j-simple-2.0.13.jar 3153fe1d689cffb94f1530b58470c306685ba68844de8857116e3b6ebb81d9f7
fetch junit/junit/4.13.2/junit-4.13.2.jar 8e495b634469d64fb8acfa3495a065cbacc8a0fff55ce1e31007be4c16dc57d3
fetch org/hamcrest/hamcrest-core/1.3/hamcrest-core-1.3.jar 66fdef91e9739348df7a096aa384a5685f4e875584cce89386a7a47251c4d8e9
rm -rf "$OUT" && mkdir -p "$OUT/main" "$OUT/test"
CP="$LIB/jnats-2.20.2.jar:$LIB/slf4j-api-2.0.13.jar"
"${JAVA_BIN}javac" -Xlint:all -Werror --release 17 -d "$OUT/main" -cp "$CP" \
  src/main/java/tak/server/plugins/BridgeCore.java src/main/java/tak/server/plugins/NatsLink.java
"${JAVA_BIN}javac" -Xlint:all,-try --release 17 -d "$OUT/test" -cp "$OUT/main:$CP:$LIB/junit-4.13.2.jar" \
  src/test/java/tak/server/plugins/*.java
"${JAVA_BIN}java" -cp "$OUT/main:$OUT/test:$CP:$LIB/slf4j-simple-2.0.13.jar:$LIB/junit-4.13.2.jar:$LIB/hamcrest-core-1.3.jar" \
  org.junit.runner.JUnitCore tak.server.plugins.BridgeCoreTest tak.server.plugins.NatsLinkTest
