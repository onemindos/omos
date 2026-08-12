package tak.server.plugins;

import io.nats.client.Connection;
import io.nats.client.Dispatcher;
import io.nats.client.Nats;
import io.nats.client.Options;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tak.server.cot.CotEventContainer;
import tak.server.proto.StreamingProtoBufHelper;
import atakmap.commoncommo.protobuf.v1.MessageOuterClass.Message;

import java.lang.invoke.MethodHandles;
import java.nio.charset.StandardCharsets;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

/**
 * TAK Server ↔ NATS Bridge Plugin
 *
 * Bidirectional bridge between TAK Server CoT events and NATS subjects.
 *
 * TAK → NATS:  Every CoT message arriving at TAK Server is published to
 *               natsProducerSubject (default: tak.cot.out) as UTF-8 XML.
 *
 * NATS → TAK:  Messages published to natsConsumerSubject
 *               (default: tak.cot.inject) are parsed as CoT XML and
 *               injected into TAK Server, appearing on all connected maps.
 *
 * Configuration (plugin-config.yaml next to JAR in tak/lib/):
 *   natsUrl:             nats://localhost:4222
 *   natsConsumerSubject: tak.cot.inject
 *   natsProducerSubject: tak.cot.out
 *   groups:
 *     - __ANON__
 */
@TakServerPlugin(
    name = "TAK Server NATS Bridge Plugin",
    description = "Bidirectional CoT bridge between TAK Server and NATS JetStream"
)
@SuppressWarnings("unused")
public class NatsBridgePlugin extends MessageSenderReceiverBase {

    private static final Logger logger = LoggerFactory.getLogger(MethodHandles.lookup().lookupClass());

    private static final String DEFAULT_NATS_URL     = "nats://localhost:4222";
    private static final String DEFAULT_COT_OUT      = "tak.cot.out";
    private static final String DEFAULT_COT_INJECT   = "tak.cot.inject";

    private final String natsUrl;
    private final String producerSubject;
    private final String consumerSubject;
    private final Set<String> groups;

    private Connection natsConnection;
    private Dispatcher natsDispatcher;

    @SuppressWarnings("unchecked")
    public NatsBridgePlugin() {
        natsUrl         = getStringProp("natsUrl",             DEFAULT_NATS_URL);
        producerSubject = getStringProp("natsProducerSubject", DEFAULT_COT_OUT);
        consumerSubject = getStringProp("natsConsumerSubject", DEFAULT_COT_INJECT);

        Set<String> grps = new HashSet<>();
        if (config.containsProperty("groups")) {
            grps.addAll((List<String>) config.getProperty("groups"));
        }
        if (grps.isEmpty()) {
            grps.add("__ANON__");
        }
        groups = grps;

        logger.info("NatsBridgePlugin config — url={} producer={} consumer={} groups={}",
            natsUrl, producerSubject, consumerSubject, groups);
    }

    // ─── TAK → NATS ───────────────────────────────────────────────────────────

    @Override
    public void onMessage(Message message) {
        if (natsConnection == null || natsConnection.getStatus() != Connection.Status.CONNECTED) {
            logger.warn("NATS not connected — dropping outbound CoT");
            return;
        }
        try {
            CotEventContainer cot = StreamingProtoBufHelper.proto2cot(message.getPayload());
            String cotXml = cot.asXml();
            natsConnection.publish(producerSubject, cotXml.getBytes(StandardCharsets.UTF_8));
            logger.debug("TAK→NATS [{}]: {}", producerSubject, cotXml.substring(0, Math.min(120, cotXml.length())));
        } catch (Exception e) {
            logger.error("Error forwarding CoT to NATS", e);
        }
    }

    // ─── NATS → TAK ───────────────────────────────────────────────────────────

    private void onNatsMessage(io.nats.client.Message natsMsg) {
        try {
            String cotXml = new String(natsMsg.getData(), StandardCharsets.UTF_8);
            logger.debug("NATS→TAK [{}]: {}", consumerSubject, cotXml.substring(0, Math.min(120, cotXml.length())));

            Message takMessage = getConverter().cotStringToDataMessage(
                cotXml,
                groups,
                Integer.toString(System.identityHashCode(this))
            );
            send(takMessage);
        } catch (Exception e) {
            logger.error("Error injecting NATS message into TAK", e);
        }
    }

    // ─── Lifecycle ────────────────────────────────────────────────────────────

    @Override
    public void start() {
        logger.info("NatsBridgePlugin starting — connecting to {}", natsUrl);
        try {
            Options options = new Options.Builder()
                .server(natsUrl)
                .maxReconnects(-1)          // reconnect forever
                .reconnectWait(java.time.Duration.ofSeconds(2))
                .connectionName("takserver-nats-bridge")
                .build();

            natsConnection = Nats.connect(options);
            logger.info("NatsBridgePlugin connected to NATS");

            // Subscribe to inbound CoT subject
            natsDispatcher = natsConnection.createDispatcher(this::onNatsMessage);
            natsDispatcher.subscribe(consumerSubject);
            logger.info("NatsBridgePlugin subscribed to NATS subject: {}", consumerSubject);

        } catch (Exception e) {
            logger.error("NatsBridgePlugin failed to connect to NATS at {}", natsUrl, e);
        }
    }

    @Override
    public void stop() {
        logger.info("NatsBridgePlugin stopping");
        try {
            if (natsDispatcher != null) {
                natsDispatcher.unsubscribe(consumerSubject);
            }
            if (natsConnection != null) {
                natsConnection.close();
            }
        } catch (Exception e) {
            logger.error("Error stopping NatsBridgePlugin", e);
        }
    }

    // ─── Helpers ──────────────────────────────────────────────────────────────

    private String getStringProp(String key, String defaultVal) {
        if (config.containsProperty(key)) {
            return (String) config.getProperty(key);
        }
        return defaultVal;
    }
}
