package tak.server.plugins;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertTrue;
import static org.junit.Assume.assumeTrue;

import io.nats.client.Connection;
import io.nats.client.Dispatcher;
import io.nats.client.Nats;

import java.io.File;
import java.net.ServerSocket;
import java.net.Socket;
import java.nio.charset.StandardCharsets;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.BooleanSupplier;

import org.junit.After;
import org.junit.Before;
import org.junit.Test;

/**
 * NatsLink against a disposable nats-server (binary from NATS_SERVER_BIN or
 * PATH; skipped without one). Covers the inject queue group shared with
 * tak-client, validation, forwarding and the connect-later path.
 */
public class NatsLinkTest {

    private Process server;
    private int port;

    private static String serverBin() {
        String env = System.getenv("NATS_SERVER_BIN");
        if (env != null && new File(env).canExecute()) {
            return env;
        }
        for (String dir : System.getenv().getOrDefault("PATH", "").split(File.pathSeparator)) {
            File f = new File(dir, "nats-server");
            if (f.canExecute()) {
                return f.getPath();
            }
        }
        return null;
    }

    private void startServer() throws Exception {
        String bin = serverBin();
        assumeTrue("nats-server not found", bin != null);
        server = new ProcessBuilder(bin, "-a", "127.0.0.1", "-p", Integer.toString(port)).inheritIO().start();
        waitFor(() -> {
            try (Socket ignored = new Socket("127.0.0.1", port)) {
                return ignored.isConnected();
            } catch (Exception e) {
                return false;
            }
        });
    }

    @Before
    public void setUp() throws Exception {
        try (ServerSocket s = new ServerSocket(0)) {
            port = s.getLocalPort();
        }
    }

    @After
    public void tearDown() throws Exception {
        if (server != null) {
            server.destroyForcibly().waitFor();
        }
    }

    private static void waitFor(BooleanSupplier cond) throws InterruptedException {
        long deadline = System.currentTimeMillis() + 10_000;
        while (System.currentTimeMillis() < deadline) {
            if (cond.getAsBoolean()) {
                return;
            }
            Thread.sleep(20);
        }
        throw new AssertionError("timed out");
    }

    private BridgeCore.Config cfg(boolean producer) {
        Map<String, Object> m = new HashMap<>();
        m.put("natsUrl", "nats://127.0.0.1:" + port);
        m.put("natsProducerEnabled", producer);
        m.put("statusIntervalSeconds", 1);
        return BridgeCore.parseConfig(m::get);
    }

    private static String pli(String uid) {
        return "<event version=\"2.0\" uid=\"" + uid + "\" type=\"a-n-G\" how=\"m-g\" time=\"2026-10-29T12:00:00Z\""
            + " start=\"2026-10-29T12:00:00Z\" stale=\"2026-10-29T12:01:00Z\"><point lat=\"1\" lon=\"2\" hae=\"0\" ce=\"9\" le=\"9\"/>"
            + "<detail/></event>";
    }

    @Test
    public void injectsSharesQueueAndRejects() throws Exception {
        startServer();
        List<String> toTak = new CopyOnWriteArrayList<>();
        try (NatsLink link = new NatsLink(cfg(true), toTak::add);
             Connection bus = Nats.connect("nats://127.0.0.1:" + port)) {
            link.start();
            waitFor(() -> "CONNECTED".equals(link.natsStatus()));

            // tak-client stand-in in the same queue group.
            AtomicInteger takClient = new AtomicInteger();
            Dispatcher d = bus.createDispatcher(m -> takClient.incrementAndGet());
            d.subscribe("tak.cot.inject", "tak-cot-inject");
            bus.flush(java.time.Duration.ofSeconds(2));
            Thread.sleep(100);

            for (int i = 0; i < 40; i++) {
                bus.publish("tak.cot.inject", pli("OM.track.q" + i).getBytes(StandardCharsets.UTF_8));
            }
            waitFor(() -> toTak.size() + takClient.get() == 40);
            Thread.sleep(200);
            assertEquals(40, toTak.size() + takClient.get());
            assertTrue("plugin got none of the 40", toTak.size() > 0);
            assertTrue("tak-client got none of the 40", takClient.get() > 0);
            assertEquals(toTak.size(), link.stats().injected.get());

            // Garbage and control types never reach TAK.
            d.unsubscribe("tak.cot.inject");
            bus.flush(java.time.Duration.ofSeconds(2));
            int before = toTak.size();
            bus.publish("tak.cot.inject", "nonsense".getBytes(StandardCharsets.UTF_8));
            bus.publish("tak.cot.inject", pli("x").replace("a-n-G", "t-x-takp-q").getBytes(StandardCharsets.UTF_8));
            bus.publish("tak.cot.inject", (pli("a") + pli("evil")).getBytes(StandardCharsets.UTF_8));
            waitFor(() -> link.stats().injectRejected.get() == 3);
            assertEquals(before, toTak.size());

            // TAK → NATS forward.
            List<String> out = new CopyOnWriteArrayList<>();
            Dispatcher o = bus.createDispatcher(m -> out.add(new String(m.getData(), StandardCharsets.UTF_8)));
            o.subscribe("tak.cot.out");
            bus.flush(java.time.Duration.ofSeconds(2));
            link.forward(pli("ANDROID-1"));
            waitFor(() -> out.size() == 1);
            assertTrue(out.get(0).contains("ANDROID-1"));
            assertEquals(1, link.stats().forwarded.get());
            assertTrue(link.stats().statusLine(link.natsStatus()).contains("injectRejected=3"));
        }
    }

    @Test
    public void injectorFailureIsCountedNotFatal() throws Exception {
        startServer();
        AtomicInteger calls = new AtomicInteger();
        try (NatsLink link = new NatsLink(cfg(false), xml -> {
            if (calls.incrementAndGet() == 1) {
                throw new IllegalStateException("TAK send failed");
            }
        }); Connection bus = Nats.connect("nats://127.0.0.1:" + port)) {
            link.start();
            waitFor(() -> "CONNECTED".equals(link.natsStatus()));
            bus.publish("tak.cot.inject", pli("a").getBytes(StandardCharsets.UTF_8));
            bus.publish("tak.cot.inject", pli("b").getBytes(StandardCharsets.UTF_8));
            waitFor(() -> link.stats().injected.get() == 1 && link.stats().injectFailed.get() == 1);
            link.forward(pli("x")); // producer disabled: no-op
            assertEquals(0, link.stats().forwarded.get());
        }
    }

    @Test
    public void startsBeforeNatsIsUpAndConnectsLater() throws Exception {
        assumeTrue("nats-server not found", serverBin() != null);
        List<String> toTak = new CopyOnWriteArrayList<>();
        try (NatsLink link = new NatsLink(cfg(true), toTak::add)) {
            link.start(); // no server yet: must not throw (the old plugin failed here for good)
            link.forward(pli("dropped"));
            assertEquals(1, link.stats().forwardDropped.get());
            startServer();
            waitFor(() -> "CONNECTED".equals(link.natsStatus()));
            try (Connection bus = Nats.connect("nats://127.0.0.1:" + port)) {
                Thread.sleep(200); // subscription replayed on connect
                bus.publish("tak.cot.inject", pli("late").getBytes(StandardCharsets.UTF_8));
                waitFor(() -> toTak.size() == 1);
            }
        }
    }
}
