package tak.server.plugins;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tak.server.cot.CotEventContainer;
import tak.server.proto.StreamingProtoBufHelper;
import atakmap.commoncommo.protobuf.v1.MessageOuterClass.Message;

import java.lang.invoke.MethodHandles;

/**
 * TAK Server ↔ NATS Bridge Plugin
 *
 * NATS → TAK:  messages on natsConsumerSubject (tak.cot.inject) are checked
 *              (one CoT event, no stream-control types), converted and sent
 *              into TAK Server. The subscription uses the queue group
 *              natsQueueGroup (tak-cot-inject), shared with the tak-client
 *              bridge, so each message is injected once.
 *
 * TAK → NATS:  CoT delivered to onMessage is published to natsProducerSubject
 *              (tak.cot.out). In TAK Server cluster mode the plugin manager
 *              delivers nothing to onMessage (services/takserver/cot-tap/README.md),
 *              so the tak-client bridge (services/bridge/bridges/tak-client) is
 *              the TAK → bus leg there; set natsProducerEnabled: false once it runs.
 *
 * Visibility: the constructor logs "NatsBridgePlugin v… loaded", start() logs
 * the subscription, and a status line with counters is logged every
 * statusIntervalSeconds (default 60). If none of these lines appears in the
 * takserver-plugins log, the plugin JAR is not on the plugin manager's classpath
 * (see services/takserver/plugin-src/README.md, "Is it loaded?").
 *
 * The SDK-free parts (config, validation, NATS) are BridgeCore and NatsLink;
 * this class is only the TAK Server glue.
 */
@TakServerPlugin(
    name = "TAK Server NATS Bridge Plugin",
    description = "Bidirectional CoT bridge between TAK Server and NATS JetStream"
)
@SuppressWarnings("unused")
public class NatsBridgePlugin extends MessageSenderReceiverBase {

    private static final Logger logger = LoggerFactory.getLogger(MethodHandles.lookup().lookupClass());

    private final BridgeCore.Config cfg;
    private volatile NatsLink link;

    public NatsBridgePlugin() {
        cfg = BridgeCore.parseConfig(key -> config.containsProperty(key) ? config.getProperty(key) : null);
        logger.info("NatsBridgePlugin v{} loaded: {}", BridgeCore.VERSION, cfg);
    }

    // ─── TAK → NATS ───────────────────────────────────────────────────────────

    @Override
    public void onMessage(Message message) {
        NatsLink l = link;
        if (l == null) {
            return;
        }
        try {
            CotEventContainer cot = StreamingProtoBufHelper.proto2cot(message.getPayload());
            l.forward(cot.asXml());
        } catch (Exception e) {
            logger.error("NatsBridgePlugin could not forward CoT to NATS", e);
        }
    }

    // ─── NATS → TAK ───────────────────────────────────────────────────────────

    private void inject(String cotXml) {
        try {
            Message takMessage = getConverter().cotStringToDataMessage(
                cotXml, cfg.groups, Integer.toString(System.identityHashCode(this)));
            send(takMessage);
        } catch (Exception e) {
            throw new IllegalStateException("convert or send failed", e);
        }
    }

    // ─── Lifecycle ────────────────────────────────────────────────────────────

    @Override
    public void start() {
        logger.info("NatsBridgePlugin starting, connecting to {}", cfg.natsUrl);
        try {
            NatsLink l = new NatsLink(cfg, this::inject);
            l.start();
            link = l;
        } catch (Exception e) {
            logger.error("NatsBridgePlugin failed to start NATS link to {}", cfg.natsUrl, e);
        }
    }

    @Override
    public void stop() {
        logger.info("NatsBridgePlugin stopping");
        NatsLink l = link;
        link = null;
        if (l != null) {
            l.close();
        }
    }
}
