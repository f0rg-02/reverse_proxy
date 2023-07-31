# reverse_proxy
A stupid simple reverse proxy that takes a simple yaml file for options. Wrote this in response to Caddy, Traefik, Nginx, etc either being too complicated or just weren't behaving the way I wanted them to for simple use cases.

This was tested with self signed certificates so there were some minor issues since most things don't like self signed certificates all that much which includes browsers, the proxy itself, etc. Idk why testing with self signed is such a bitch, but here we are.

To get self signed to mostly behave,

**NOTE: FOR TESTING!!! DO NOT USE SELF SIGNED IN PRODUCTION!!!**

To generate we use openssl on Linux:

```ruby
openssl req -newkey rsa:4096 -nodes -keyout bind.key -x509 -days 365 -out bind.crt -addext 'subjectAltName = IP:127.0.0.1' -subj '/C=US/ST=CA/L=SanFrancisco/O=MyCompany/OU=RND/CN=localhost/'
```

Something like above^^^

------

On debian install `ca-certificates` but most likely already installed:

```ruby
sudo apt-get install ca-certificates
```

Copy the certificate with `.crt` extension:

```ruby
cp bind.crt /usr/share/ca-certificates
```

And then reconfigure `ca-certificates`:

```ruby
sudo dpkg-reconfigure ca-certificates
```

When prompted, select and press **ENTER** to activate the cert. This will install the cert locally on the Linux system. 

------

Link to [original so](https://unix.stackexchange.com/a/90607) and link to [archive so](https://web.archive.org/web/20230729053542/https://unix.stackexchange.com/questions/90450/adding-a-self-signed-certificate-to-the-trusted-list/90607)

------

To install this tool:

```ruby
git clone https://github.com/f0rg-02/reverse_proxy
cd reverse_proxy && go build
```

Keeping with my KISS theory mentality, the only required argument is a YAML file that contains all the necessary options.

Example config:

```ruby
listen: "0.0.0.0:443"
server: "127.0.0.1:4443"
cert_file: "bind.crt"
key_file: "bind.key"
paths: [ "/path1", "/path2" ]
default_domain: "url_to_route_everything_else"
```

`listen` is for the port the proxy server to listen on and server is the server it should reroute to.

`cert_file` and `key_file` is your x509 or whatever certificates that you either generated yourself or you obtained some other means.

`paths` is the paths you want to handle. This is so you can send only specific uri requests with path to proxy server while everything else gets rerouted to the `default_domain`.

I should note, this isn't 100% full proof so please use something more professional and well done in productin like [Caddy](https://caddyserver.com/) or [Traefik](https://traefik.io/traefik/).


