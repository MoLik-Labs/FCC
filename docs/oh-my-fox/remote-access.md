# Remote access checkpoint

The FCC service defaults to `HOST=0.0.0.0`, and Admin and Code routes do not
have an inbound loopback-only dependency. Tests exercise direct non-loopback
client addresses and tunnel-style Cloudflare, ngrok and Tailscale Host/Origin
values, including Code session creation and readback. These are in-process ASGI
checks, not proof of a deployed tunnel, firewall configuration or browser
behavior on a real external network.

Claude Desktop gateway validation now permits remote HTTP as well as HTTPS
rather than requiring loopback for HTTP. Validation still rejects unsupported
schemes, URL credentials, queries, fragments and non-root paths. Desktop stores
a managed bearer token in its configuration: remote HTTP does not provide HTTP
transport encryption. Use HTTPS or a suitably protected private tunnel and
protect the endpoint. Admin and Code do not have a separate Admin authentication
boundary; Proxy Authentication must not be mistaken for one. A tunnel alone is
not an access policy; configure independent network/tunnel authorization.

This change does **not** remove loopback addresses that serve other purposes:
local provider defaults and same-machine launcher URLs, the OpenAI OAuth
callback, and the `web_fetch` egress restrictions that protect against SSRF.
Those are not inbound restrictions on Admin or Code, and disabling them would
change security or third-party protocol behavior rather than enable tunnels.
No newly deployed tunnel or public endpoint is configured by this branch.
