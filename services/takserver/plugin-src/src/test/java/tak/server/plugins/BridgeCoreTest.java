package tak.server.plugins;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertThrows;
import static org.junit.Assert.assertTrue;

import java.util.Arrays;
import java.util.HashMap;
import java.util.Map;

import org.junit.Test;

public class BridgeCoreTest {

    private static BridgeCore.Config parse(Map<String, Object> m) {
        return BridgeCore.parseConfig(m::get);
    }

    @Test
    public void defaultsWhenEmpty() {
        BridgeCore.Config c = parse(new HashMap<>());
        assertEquals("nats://localhost:4222", c.natsUrl);
        assertEquals("tak.cot.inject", c.consumerSubject);
        assertEquals("tak.cot.out", c.producerSubject);
        assertEquals("tak-cot-inject", c.queueGroup);
        assertTrue(c.producerEnabled);
        assertEquals(60, c.statusIntervalSeconds);
        assertEquals(1, c.groups.size());
        assertTrue(c.groups.contains("__ANON__"));
    }

    @Test
    public void deployedConfigParses() {
        // services/takserver/plugin/tak.server.plugins.NatsBridgePlugin.yaml as SnakeYAML returns it.
        Map<String, Object> m = new HashMap<>();
        m.put("natsUrl", "nats://tak-nats-leaf.takserver.svc.cluster.local:4222");
        m.put("natsConsumerSubject", "tak.cot.inject");
        m.put("natsProducerSubject", "tak.cot.out");
        m.put("natsQueueGroup", "tak-cot-inject");
        m.put("statusIntervalSeconds", 60);
        m.put("groups", Arrays.asList("__ANON__", "Red"));
        BridgeCore.Config c = parse(m);
        assertEquals("nats://tak-nats-leaf.takserver.svc.cluster.local:4222", c.natsUrl);
        assertEquals(Arrays.asList("__ANON__", "Red"), c.groups.stream().toList());
        assertTrue(c.toString().contains("queue=tak-cot-inject"));
    }

    @Test
    public void queueGroupCanBeTurnedOff() {
        Map<String, Object> m = new HashMap<>();
        m.put("natsQueueGroup", "none");
        assertEquals("", parse(m).queueGroup);
    }

    @Test
    public void producerCanBeTurnedOff() {
        Map<String, Object> m = new HashMap<>();
        m.put("natsProducerEnabled", "false");
        BridgeCore.Config c = parse(m);
        assertFalse(c.producerEnabled);
        assertTrue(c.toString().contains("producer=off"));
    }

    @Test
    public void badConfigFailsLoudly() {
        Map<String, Object> wildcard = new HashMap<>();
        wildcard.put("natsConsumerSubject", "tak.cot.>");
        assertThrows(IllegalArgumentException.class, () -> parse(wildcard));

        Map<String, Object> groups = new HashMap<>();
        groups.put("groups", "Red");
        assertThrows(IllegalArgumentException.class, () -> parse(groups));

        Map<String, Object> url = new HashMap<>();
        url.put("natsUrl", 4222);
        assertThrows(IllegalArgumentException.class, () -> parse(url));

        Map<String, Object> bool = new HashMap<>();
        bool.put("natsProducerEnabled", "yes");
        assertThrows(IllegalArgumentException.class, () -> parse(bool));

        Map<String, Object> neg = new HashMap<>();
        neg.put("statusIntervalSeconds", -1);
        assertThrows(IllegalArgumentException.class, () -> parse(neg));
    }

    private static final String PLI = "<?xml version=\"1.0\"?>\n<event version=\"2.0\" uid=\"OM.drone.sim-x1\" type=\"a-f-A-M-H-Q\""
        + " how=\"m-g\" time=\"2026-10-29T12:00:00Z\" start=\"2026-10-29T12:00:00Z\" stale=\"2026-10-29T12:01:00Z\">"
        + "<point lat=\"30.2\" lon=\"-97.7\" hae=\"150\" ce=\"10\" le=\"10\"/><detail><contact callsign=\"SIM-X1\"/></detail></event>";

    @Test
    public void acceptsOneEvent() {
        assertNull(BridgeCore.rejectReason(PLI));
        assertEquals("a-f-A-M-H-Q", BridgeCore.cotType(PLI));
    }

    @Test
    public void rejectsBadInjects() {
        assertNotNull(BridgeCore.rejectReason(""));
        assertNotNull(BridgeCore.rejectReason("hello"));
        assertEquals("more than one CoT event in one message", BridgeCore.rejectReason(PLI + PLI.replace("OM.drone", "evil")));
        assertEquals("stream control type t-x-takp-q", BridgeCore.rejectReason(PLI.replace("a-f-A-M-H-Q", "t-x-takp-q")));
        assertEquals("stream control type t-x-c-t", BridgeCore.rejectReason(PLI.replace("a-f-A-M-H-Q", "t-x-c-t")));
        assertEquals("event has no type", BridgeCore.rejectReason(PLI.replace(" type=\"a-f-A-M-H-Q\"", "")));
    }

    @Test
    public void typeAttributeNotConfusedWithOtherAttributes() {
        String x = "<event version=\"2.0\" subtype=\"zzz\" type='b-m-p-s-m' uid=\"u\"></event>";
        assertEquals("b-m-p-s-m", BridgeCore.cotType(x));
        assertNull(BridgeCore.rejectReason(x));
    }

    @Test
    public void statusLineIsGreppable() {
        BridgeCore.Stats s = new BridgeCore.Stats();
        s.injected.addAndGet(3);
        s.injectRejected.incrementAndGet();
        String line = s.statusLine("CONNECTED");
        assertTrue(line.startsWith("NatsBridgePlugin v" + BridgeCore.VERSION + " status nats=CONNECTED injected=3 injectRejected=1"));
    }
}
