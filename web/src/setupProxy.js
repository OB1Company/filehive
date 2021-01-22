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
            onProxyReq: function (proxyReq, req, res) {
                // console.log(req.headers.cookie);
            },
            onProxyRes: function (proxyRes, req, res) {
                // console.log(req.headers.cookie);
            },
        })
    );
};