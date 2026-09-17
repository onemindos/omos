/**
 * Node-RED Settings for Legacy AI Infrastructure
 * Production configuration with security enabled
 */

module.exports = {
    // Flow file settings
    flowFile: 'flows.json',
    flowFilePretty: true,

    // User directory (where flows, credentials are stored)
    userDir: '/data',

    // Node-RED UI settings
    uiHost: "0.0.0.0",
    uiPort: process.env.PORT || 1880,

    // Security
    adminAuth: {
        type: "credentials",
        users: [{
            username: "admin",
            password: "$2a$08$lrBPueiIJASas2ZubtU4Uu5C4uB69TAJtM009JpOBrr5ivYrAO.Be",
            permissions: "*"
        }]
    },

    // HTTPS settings (handled by Caddy reverse proxy)
    https: false,

    // Enable projects feature
    editorTheme: {
        projects: {
            enabled: true,
            workflow: {
                mode: "auto"
            }
        },
        header: {
            title: "Legacy AI - Node-RED",
            image: null
        },
        palette: {
            editable: true,
            catalogues: [
                'https://catalogue.nodered.org/catalogue.json'
            ]
        },
        codeEditor: {
            lib: "monaco",
            options: {
                theme: "vs-dark"
            }
        }
    },

    // Logging
    logging: {
        console: {
            level: "info",
            metrics: false,
            audit: false
        }
    },

    // Context storage
    contextStorage: {
        default: {
            module: "localfilesystem"
        }
    },

    // Function node settings
    functionGlobalContext: {
        // Global context available in function nodes
    },

    // Export settings
    exportGlobalContextKeys: false,

    // Node settings
    debugMaxLength: 1000,
    mqttReconnectTime: 15000,
    serialReconnectTime: 15000,

    // Editor settings
    runtimeState: {
        enabled: true,
        ui: true
    }
}
