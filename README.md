# MateBot Web

MateBot Web is a web frontend for [MateBot](https://github.com/hopfenspace/MateBot). 
It is served as single-side application to allow embedding as Matrix Custom Integration.

## Install from source

In order to install MateBot Web, the following packages are required (based on Debian 11):

- git
- build-essential

**Go 1.18**:

As Debian 11 provides an old version of golang, use the installation instructions from
[the official site](https://go.dev/doc/install)

**NodeJS 16 LTS**:
```bash
curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
apt-get install -y nodejs
```

Compile and install the project itself:
```bash
git clone https://github.com/hopfenspace/matebot-web.git
cd matebot-web
make
make install
```

## Configuration

```bash
cp /etc/matebot-web/example.config.toml /etc/matebot-web/config.toml
```

### Reverse Proxy
It is recommended to deploy matebot-web behind a reverse proxy. Nginx is used as reference.

Install nginx:
```bash
apt-get install nginx
```

Copy nginx configuration file:
```bash
cp matebot-web.nginx /etc/nginx/sites-available/
ln -s /etc/nginx/sites-available/matebot-web.nginx /etc/nginx/sites-enabled/
```

Set `SERVER_NAME` to the vhost, your server can be accessed by:
```bash
sed -i -E 's/SERVER_NAME/your-server-name/' /etc/nginx/sites-available/matebot-web.nginx
```

Set `SSL_CERT_PATH` and `SSL_CERT_KEY_PATH` to a certificate, that is valid for `SERVER_NAME`:
```bash
sed -i -E 's/SSL_CERT_PATH/\/path\/to\/cert\.pem/' /etc/nginx/sites-available/matebot-web.nginx
sed -i -E 's/SSL_CERT_KEY_PATH/\/path\/to\/key\.pem/' /etc/nginx/sites-available/matebot-web.nginx
```

Finally, start nginx:
```bash
systemctl start nginx
```