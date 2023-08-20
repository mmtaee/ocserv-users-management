const { defineConfig } = require('@vue/cli-service')
const path = require("path");

module.exports = defineConfig({
  transpileDependencies: ['vuetify'],
  lintOnSave: false,
  publicPath: "/",
  devServer: {
    host: "127.0.0.1",
    allowedHosts: ["0.0.0.0", "127.0.0.1", "localhost"],
    port: 9000,
    compress: true,
    proxy: {
      "^/api": {
        target: "http://127.0.0.1:8000/api/",
        changeOrigin: true,
      },
    },
  },
  pages: {
    index: {
      entry: "src/main.ts",
      template: "./public/index.html",
      filename: "index.html",
      title: "Ocserv User panel",
      chunks: ["chunk-vendors", "chunk-common", "index"],
    },
  },
})
