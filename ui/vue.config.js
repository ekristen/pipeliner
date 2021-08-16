module.exports = {
  "transpileDependencies": [
    "vuetify"
  ],
  configureWebpack: {
    resolve: {
      symlinks: false
    }
  },
  chainWebpack: config => {
    config
      .plugin('html')
      .tap(args => {
        args[0].title = "Pipeliner";
        return args;
      })
  },
  publicPath: process.env.NODE_ENV === 'production'
    ? '/ui/'
    : '/',
  devServer: {
    port: 4445,
    proxy: {
      '^/v1|/ws|/whoami|/logout|/version': {
        target: 'http://localhost:4444',
        ws: true,
        changeOrigin: true
      },
    }
  }
}