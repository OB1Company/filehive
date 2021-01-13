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
                if (proxyRes.headers['set-cookie'] != undefined) {
                    console.log(proxyRes);
                    console.log(res, req);
                    // res.headers['set-cookie'] = proxyRes.headers['set-cookie'];
                    // req.session['cookie'] = proxyRes.headers['set-cookie'];  // must be or you will get new session for each call
                    // req.session['proxy-cookie'] = proxyRes.headers['set-cookie'];  // add to other key because cookie will be lost
                }
                // console.log("response: " + req.session.id);
                // console.log(req.session);
            },
        })
    );
};