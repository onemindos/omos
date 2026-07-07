import type { App } from 'vue';
import type { PluginAPI, PluginInstance } from '@tak-ps/cloudtak';
import AIPanel from './lib/AIPanel.vue';
import AIBottomBar from './lib/AIBottomBar.vue';
import IconAI from './lib/IconAI.vue';

const ROUTE_NAME = 'onemind-ai';

export default class OneMindAI implements PluginInstance {
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
            tooltip: 'Chat with Legacy AI',
        });

        this.api.bottomBar.add({
            key: 'onemind-ai-status',
            component: AIBottomBar,
        });
    }

    async disable(): Promise<void> {
        this.api.bottomBar.remove('onemind-ai-status');
        this.api.menu.remove(ROUTE_NAME);
        if (this.api.router.hasRoute(ROUTE_NAME)) {
            this.api.router.removeRoute(ROUTE_NAME);
        }
    }
}
