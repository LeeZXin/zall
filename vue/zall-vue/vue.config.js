const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    host: "0.0.0.0",
    port: 8080,
    https: false,
    open: false,
    client: {
      webSocketURL: "ws://0.0.0.0:8080/ws",
      overlay: false
    },
    proxy: {
      '/api': {
        target: "http://127.0.0.1:80",
        changeOrigin: true
      },
    }
  }
})
