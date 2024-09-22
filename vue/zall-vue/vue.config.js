const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
    transpileDependencies: true,
    devServer: {
        host: "0.0.0.0",
        port: 7501,
        https: true,
        open: false,
        allowedHosts: "all",
        client: {
            webSocketURL: "wss://0.0.0.0:7501/ws",
            overlay: false
        },
        proxy: {
            '^/api': {
                target: "http://127.0.0.1:80",
                changeOrigin: true
            },
        },
    },
    pwa: {
        iconPaths: {
            favicon32: "logo.png",
            favicon16: "logo.png"
        }
    }
})