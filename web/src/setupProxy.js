const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
    app.use(
        '/api',
        createProxyMiddleware({
            target: 'http://localhost:8080',
            changeOrigin: true,
            cookieDomainRewrite: "localhost",
            cookiePathRewrite: "/",
            debug: true,
            onProxyRes: function (proxyRes, req, res) {
            },
        })
    );
};