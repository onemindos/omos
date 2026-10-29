package tak.server.plugins;

import io.nats.client.Connection;
import io.nats.client.ConnectionListener;
import io.nats.client.Dispatcher;
import io.nats.client.Nats;
import io.nats.client.Options;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.function.Consumer;

/**
 * The NATS side of the bridge, without the TAK Server SDK: connection,
 * the inject subscription (queue group shared with tak-client), the
 * forward publish and the periodic status log line.
 */
public final class NatsLink implements AutoCloseable {

    private static final Logger logger = LoggerFactory.getLogger(NatsLink.class);

    private final BridgeCore.Config cfg;
    private final Consumer<String> injector;
    private final BridgeCore.Stats stats = new BridgeCore.Stats();

    private volatile Connection connection;
    private Dispatcher dispatcher;
    private ScheduledExecutorService statusTimer;
    private boolean closed;

    /**
     * @param injector called with each accepted inject payload (CoT XML); it
     *                 hands the event to TAK Server and throws on failure.
     */
    public NatsLink(BridgeCore.Config cfg, Consumer<String> injector) {
        this.cfg = cfg;
        this.injector = injector;
    }

    public BridgeCore.Stats stats() {
        return stats;
    }

    public String natsStatus() {
        Connection c = connection;
        return c == null ? "NOT_STARTED" : c.getStatus().name();
    }

    /**
     * Starts connecting in the background and returns at once, so TAK Server's
     * plugin manager is never blocked by a missing NATS server. The inject
     * subscription is made on the first CONNECTED (or RECONNECTED) event; the
     * client replays it after every later reconnect.
     */
    public void start() throws InterruptedException {
        Options options = new Options.Builder()
            .server(cfg.natsUrl)
            .maxReconnects(-1)
            .reconnectWait(Duration.ofSeconds(2))
            .connectionName("takserver-nats-bridge")
            .connectionListener(this::onConnectionEvent)
            .build();
        Nats.connectAsynchronously(options, true);
        if (cfg.statusIntervalSeconds > 0) {
            statusTimer = Executors.newSingleThreadScheduledExecutor(r -> {
                Thread t = new Thread(r, "nats-bridge-status");
                t.setDaemon(true);
                return t;
            });
            statusTimer.scheduleAtFixedRate(() -> logger.info(stats.statusLine(natsStatus())),
                cfg.statusIntervalSeconds, cfg.statusIntervalSeconds, TimeUnit.SECONDS);
        }
    }

    private void onConnectionEvent(Connection conn, ConnectionListener.Events event) {
        logger.info("NatsBridgePlugin NATS {}", event);
        // A first connect that succeeds only after retries is reported as RECONNECTED.
        if (event != ConnectionListener.Events.CONNECTED && event != ConnectionListener.Events.RECONNECTED) {
            return;
        }
        synchronized (this) {
            if (closed) {
                // stop() ran while the first connect was still retrying.
                new Thread(() -> closeQuietly(conn), "nats-bridge-late-close").start();
                return;
            }
            if (dispatcher != null) {
                return;
            }
            connection = conn;
            dispatcher = conn.createDispatcher(msg -> onInject(new String(msg.getData(), StandardCharsets.UTF_8)));
            if (cfg.queueGroup.isEmpty()) {
                dispatcher.subscribe(cfg.consumerSubject);
            } else {
                dispatcher.subscribe(cfg.consumerSubject, cfg.queueGroup);
            }
        }
        logger.info("NatsBridgePlugin subscribed to {}{}", cfg.consumerSubject,
            cfg.queueGroup.isEmpty() ? "" : " (queue " + cfg.queueGroup + ")");
    }

    private static void closeQuietly(Connection c) {
        try {
            c.close();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    void onInject(String xml) {
        String reason = BridgeCore.rejectReason(xml);
        if (reason != null) {
            stats.injectRejected.incrementAndGet();
            logger.warn("NatsBridgePlugin refused {} message: {}", cfg.consumerSubject, reason);
            return;
        }
        try {
            injector.accept(xml);
            stats.injected.incrementAndGet();
        } catch (RuntimeException e) {
            stats.injectFailed.incrementAndGet();
            logger.error("NatsBridgePlugin could not inject CoT into TAK", e);
        }
    }

    /** TAK → NATS. Drops (and counts) while disconnected; never blocks TAK. */
    public void forward(String cotXml) {
        if (!cfg.producerEnabled) {
            return;
        }
        Connection c = connection;
        if (c == null || c.getStatus() != Connection.Status.CONNECTED) {
            stats.forwardDropped.incrementAndGet();
            return;
        }
        c.publish(cfg.producerSubject, cotXml.getBytes(StandardCharsets.UTF_8));
        stats.forwarded.incrementAndGet();
    }

    @Override
    public void close() {
        Connection c;
        synchronized (this) {
            closed = true;
            c = connection;
        }
        if (statusTimer != null) {
            statusTimer.shutdownNow();
        }
        if (c != null) {
            closeQuietly(c); // also closes the dispatcher
        }
        logger.info(stats.statusLine("CLOSED"));
    }
}
