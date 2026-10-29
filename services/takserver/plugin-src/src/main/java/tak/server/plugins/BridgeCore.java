package tak.server.plugins;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.LinkedHashSet;
import java.util.Set;
import java.util.concurrent.atomic.AtomicLong;
import java.util.function.Function;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * The parts of the NATS bridge that do not need the TAK Server plugin SDK:
 * configuration, inject validation and counters. Kept separate so they are
 * unit-tested without the SDK JAR (which comes from tak.gov and is not in Git).
 */
public final class BridgeCore {

    public static final String VERSION = "1.1.0";

    static final String DEFAULT_NATS_URL = "nats://localhost:4222";
    static final String DEFAULT_COT_OUT = "tak.cot.out";
    static final String DEFAULT_COT_INJECT = "tak.cot.inject";
    /** Shared with services/bridge/bridges/tak-client (INJECT_QUEUE): one injection per message. */
    static final String DEFAULT_QUEUE_GROUP = "tak-cot-inject";
    static final long DEFAULT_STATUS_INTERVAL_S = 60;

    private BridgeCore() {}

    /** Plugin settings from conf/plugins/tak.server.plugins.NatsBridgePlugin.yaml. */
    public static final class Config {
        public final String natsUrl;
        public final String consumerSubject;
        public final String producerSubject;
        /** Empty string = plain subscription (no queue group). */
        public final String queueGroup;
        public final boolean producerEnabled;
        public final long statusIntervalSeconds;
        public final Set<String> groups;

        Config(String natsUrl, String consumerSubject, String producerSubject, String queueGroup,
               boolean producerEnabled, long statusIntervalSeconds, Set<String> groups) {
            this.natsUrl = natsUrl;
            this.consumerSubject = consumerSubject;
            this.producerSubject = producerSubject;
            this.queueGroup = queueGroup;
            this.producerEnabled = producerEnabled;
            this.statusIntervalSeconds = statusIntervalSeconds;
            this.groups = groups;
        }

        @Override
        public String toString() {
            return "url=" + natsUrl + " consumer=" + consumerSubject
                + (queueGroup.isEmpty() ? "" : " queue=" + queueGroup)
                + " producer=" + (producerEnabled ? producerSubject : "off")
                + " groups=" + groups + " status=" + statusIntervalSeconds + "s";
        }
    }

    /**
     * Reads the config through a lookup (the plugin passes its PluginConfiguration).
     * Missing or empty values fall back to the defaults; wrong types are errors,
     * so a typo in the YAML fails the plugin start visibly instead of silently.
     */
    public static Config parseConfig(Function<String, Object> lookup) {
        String url = str(lookup, "natsUrl", DEFAULT_NATS_URL);
        String consumer = str(lookup, "natsConsumerSubject", DEFAULT_COT_INJECT);
        String producer = str(lookup, "natsProducerSubject", DEFAULT_COT_OUT);
        String queue = str(lookup, "natsQueueGroup", DEFAULT_QUEUE_GROUP);
        boolean producerEnabled = bool(lookup, "natsProducerEnabled", true);
        long status = lng(lookup, "statusIntervalSeconds", DEFAULT_STATUS_INTERVAL_S);
        if (status < 0) {
            throw new IllegalArgumentException("statusIntervalSeconds must be >= 0");
        }
        if (consumer.contains("*") || consumer.contains(">")) {
            throw new IllegalArgumentException("natsConsumerSubject must be a literal subject, got " + consumer);
        }
        Set<String> groups = new LinkedHashSet<>();
        Object g = lookup.apply("groups");
        if (g instanceof Collection<?>) {
            for (Object o : (Collection<?>) g) {
                if (o != null && !o.toString().isBlank()) {
                    groups.add(o.toString().trim());
                }
            }
        } else if (g != null) {
            throw new IllegalArgumentException("groups must be a list");
        }
        if (groups.isEmpty()) {
            groups.add("__ANON__");
        }
        return new Config(url, consumer, producer, queue.equals("none") ? "" : queue, producerEnabled, status,
            Collections.unmodifiableSet(groups));
    }

    private static String str(Function<String, Object> lookup, String key, String def) {
        Object v = lookup.apply(key);
        if (v == null || v.toString().isBlank()) {
            return def;
        }
        if (!(v instanceof String)) {
            throw new IllegalArgumentException(key + " must be a string");
        }
        return ((String) v).trim();
    }

    private static boolean bool(Function<String, Object> lookup, String key, boolean def) {
        Object v = lookup.apply(key);
        if (v == null) {
            return def;
        }
        if (v instanceof Boolean) {
            return (Boolean) v;
        }
        String s = v.toString().trim().toLowerCase();
        if (s.equals("true") || s.equals("false")) {
            return Boolean.parseBoolean(s);
        }
        throw new IllegalArgumentException(key + " must be true or false");
    }

    private static long lng(Function<String, Object> lookup, String key, long def) {
        Object v = lookup.apply(key);
        if (v == null) {
            return def;
        }
        if (v instanceof Number) {
            return ((Number) v).longValue();
        }
        try {
            return Long.parseLong(v.toString().trim());
        } catch (NumberFormatException e) {
            throw new IllegalArgumentException(key + " must be a number");
        }
    }

    // ─── Inject validation ───────────────────────────────────────────────────

    private static final Pattern EVENT_START = Pattern.compile("<event[\\s>]");
    private static final Pattern TYPE_ATTR = Pattern.compile("<event\\b[^>]*?\\stype\\s*=\\s*[\"']([^\"']*)[\"']");

    /** The CoT type of the first event, or null. */
    public static String cotType(String xml) {
        Matcher m = TYPE_ATTR.matcher(xml);
        return m.find() ? m.group(1) : null;
    }

    /**
     * Why an inject payload is refused, or null when it may go to TAK Server.
     * Refused: empty, not one CoT event, more than one event in one message
     * (a second, smuggled event), or stream plumbing (protocol negotiation,
     * pings), which only the TAK connection itself may send.
     */
    public static String rejectReason(String xml) {
        if (xml == null || xml.isBlank()) {
            return "empty message";
        }
        Matcher m = EVENT_START.matcher(xml);
        int events = 0;
        while (m.find()) {
            events++;
        }
        if (events == 0 || !xml.contains("</event>") && !xml.trim().endsWith("/>")) {
            return "not a CoT event";
        }
        if (events > 1) {
            return "more than one CoT event in one message";
        }
        String type = cotType(xml);
        if (type == null || type.isEmpty()) {
            return "event has no type";
        }
        if (type.startsWith("t-x-takp") || type.equals("t-x-c-t") || type.equals("t-x-c-t-r")) {
            return "stream control type " + type;
        }
        return null;
    }

    // ─── Counters ────────────────────────────────────────────────────────────

    /** Counters reported in the periodic status line. */
    public static final class Stats {
        public final AtomicLong injected = new AtomicLong();
        public final AtomicLong injectRejected = new AtomicLong();
        public final AtomicLong injectFailed = new AtomicLong();
        public final AtomicLong forwarded = new AtomicLong();
        public final AtomicLong forwardDropped = new AtomicLong();

        public String statusLine(String natsStatus) {
            ArrayList<String> parts = new ArrayList<>();
            parts.add("nats=" + natsStatus);
            parts.add("injected=" + injected.get());
            parts.add("injectRejected=" + injectRejected.get());
            parts.add("injectFailed=" + injectFailed.get());
            parts.add("forwarded=" + forwarded.get());
            parts.add("forwardDropped=" + forwardDropped.get());
            return "NatsBridgePlugin v" + VERSION + " status " + String.join(" ", parts);
        }
    }
}
