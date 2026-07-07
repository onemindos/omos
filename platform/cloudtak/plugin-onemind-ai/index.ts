/**
 * OneMind AI Chat Plugin for CloudTAK
 *
 * Adds a "Legacy AI" sidebar panel that connects to Legacy (Hermes) via
 * the tak-agent NATS RPC bridge (svc.tak.chat). The plugin talks to a
 * lightweight HTTP bridge pod in the cluster that relays requests over NATS.
 *
 * Architecture:
 *   Browser → POST /api/onemind/chat (CloudTAK API proxy or sidecar)
 *     → tak-agent svc.tak.chat (NATS RPC)
 *       → Legacy AI (agents.prompt.hermes.zeus.legacy)
 *         → response back through chain → browser
 *
 * The BRIDGE_URL env var on the onemind-ai-bridge pod is injected into
 * the window.__ONEMIND_BRIDGE_URL__ global by the plugin's server-side
 * config endpoint (api/routes/config CLOUDTAK_Config_onemind__bridge_url).
 */

import { App } from 'vue';
import type { PluginAPI, PluginInstance } from '@tak-ps/cloudtak';
import AIPanel from './lib/AIPanel.vue';
import AIBottomBar from './lib/AIBottomBar.vue';
import IconAI from './lib/IconAI.vue';

const ROUTE_NAME = 'onemind-ai';

export default class OneMindAI {
    api: PluginAPI;

    constructor(api: PluginAPI) {
        this.api = api;
    }

    static async install(app: App, api: PluginAPI): Promise<PluginInstance> {
        return new OneMindAI(api);
    }

    async enable(): Promise<void> {
        this.api.routes.add(
            {
                path: '/onemind/ai',
                name: ROUTE_NAME,
                component: AIPanel,
            },
            'home-menu'
        );

        this.api.menu.add({
            key: ROUTE_NAME,
            label: 'Legacy AI',
            route: ROUTE_NAME,
            icon: IconAI,
            order: 10,
        });

        this.api.bottomBar.add({
            key: 'onemind-ai-status',
            component: AIBottomBar,
        });
    }

    async disable(): Promise<void> {
        this.api.bottomBar.remove('onemind-ai-status');
        this.api.menu.remove(ROUTE_NAME);
        this.api.routes.router.removeRoute(ROUTE_NAME);
    }
}
