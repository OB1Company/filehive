const { createProxyMiddleware } = require('http-proxy-middleware');

require('dotenv').config()

module.exports = function(app) {
    app.use(
        '/api',
        createProxyMiddleware({
            target: process.env.REACT_APP_TARGET_SERVER,
            changeOrigin: process.env.REACT_APP_CHANGE_ORIGIN,
            cookieDomainRewrite: process.env.REACT_APP_COOKIE_DOMAIN,
            cookiePathRewrite: process.env.REACT_APP_COOKIE_PATH,
            debug: true,
            onProxyReq: function (proxyReq, req, res) {
                // console.log(req.headers['cookies']);
            },
            onProxyRes: function (proxyRes, req, res) {
                // console.log(proxyRes.headers);
            },
        })
    );
};