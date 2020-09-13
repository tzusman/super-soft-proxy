## Super Soft Proxy

Each macOS computer get's a friendly `.local` host which is accessible to other computers on the network. This proxy automatically provisions an SSL certificate for that host. By downloading the generated certificate-authority, other devices will accept the SSL certificate.

This is especially helpful when using browser services like camera access which are only available on secure connections.

### How to run it

`go run main.go`

This will create the directory `~/.config/super-soft-proxy`.  On the first run, it will generate a certificate authority as well as the key and certificate for your computer's hostname. Any subsequent runs will use the existing certificates.

### Access your computer with SSL from a mobile device

When you run the proxy, it will output something like this (my hostname is `tz.local`):

```
Download ca cert to your phone at http://tz.local:7001
Proxy is at https://tz.local:7000
```

Visit `http://tz.local:7001` on your device or computer and install the certificate authority certificate which will automatically download.  This is easiest on iOS devices by visiting the url in Safari.

Once the certificate authority is installed, visit `https://tz.local:7000` or any other request on that origin and it will be proxied to port `8080` on the host computer by default.