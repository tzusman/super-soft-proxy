## Super Soft Proxy

Zero configuration macOS SSL proxy for your network

Each macOS computer get's a friendly `.local` host which is accessible to other computers on the network. This proxy automatically provisions an SSL certificate for that host. By downloading the generated certificate-authority, other devices will accept the SSL certificate.

This is especially helpful when using browser services like camera access which are only available on secure connections.

### How to run it

To create the config file at ` ~/.config/super-soft-proxy/config` first run:

```
go build -o super-soft-proxy main.go
./super-soft-proxy init
```

The first run after a configuration file has been generated, `super-soft-proxy` will generate a certificate authority as well as the key and certificate for your computer's hostname. Any subsequent runs will use the existing certificates.

```
./super-soft-proxy
```

No configuration is needed, but if you want to adjust the settings:

```
Usage of ./super-soft-proxy:
  -ca-cert-server int
      port where the CA certificate will be served (default 7001)
  -proxy-port int
      proxy port (default 7000)
  -upstream int
      upstream service port (default 8080)
```

### Access your computer with SSL from a mobile device

When you run the proxy, it will output something like this (my hostname is `tz.local`):

```
https://tz.local:7000 is the proxy address
 http://tz.local:7001 serves the CA cert for easy device installation
```

Visit `http://tz.local:7001` on your device or computer and install the certificate authority certificate which will automatically download.  This is easiest on iOS devices by visiting the url in Safari.

Once the certificate authority is installed, visit `https://tz.local:7000` or any other request on that origin and it will be proxied to port `8080` on the host computer.