# Image Proxy

This package provides a container image that can proxy an HTTP request to a URL defined in the request path.

## Container Usage

This container is designed to be run behind a reverse proxy, such as Nginx or Cloudflare.

### Ports

- 80
- 443

The certificate presented by the TLS server on port 443 is a self-signed certificate generated when the application starts.

## Webhook Usage

To proxy an image, encode the full URL with the alternate unpadded variant of Base64. Provide that value as the path to the request to have the server proxy it.

If required, you may also provide a file extension to the end of the base64-encoded data. This extension will be ignored by the server.

Any headers will be passed along in the proxied request, with the exception of the Host header, which will be rewritten with the host from the provided URL.

Only HTTP HEAD and GET methods can be used. No more than 50MiB may be proxied per request.
