**NOTE:** I stopped working on this as Caddy has grown to be a more viable solution for me.

# 稲荷 (Inari)
A zero-config web server that gives you sensible defaults.

## Building from source
1. Install [Golang](https://golang.org/dl/)
2. Fetch dependencies and build (`go build`)

The binary `inari` is then able to be executed with environment variables. **You will need to manually export and set these variables, as 稲荷 will not read the .env file.** You can learn more about how to set environment variables on [Windows](https://docs.microsoft.com/powershell/module/microsoft.powershell.core/about/about_environment_variables), [macOS](https://support.apple.com/guide/terminal/apd382cc5fa-4f58-4449-b20a-41c53c006f8f), or [Linux](https://www.redhat.com/sysadmin/linux-environment-variables) with their respective links.

## Usage
稲荷 can be run in any directory you want to run a server. By default, it will use port `8442`, but you can change it with the `PORT` environment variable. In the event that you cannot have 稲荷 in the folder you want to deploy, you can use the `DIRECTORY` environment variable to decalre the absolute path to where your site's root is located.

稲荷 looks at the following variables:
- `CACHE_IMAGE_TIME`: Adjust how many minutes until the cache for images needs to be updated on the client (between sessions)
- `CACHE_FONT_TIME`: Adjust how many minutes until the cache for font files and stylesheets to be updated on the client (between sessions)
- `DIRECTORY`: Set the absolute path to where the web root will be
- `NOHSTS` (not recommended): Disabled HSTS and HSTS preloading
- `PORT`: Change the default port the server will deploy to
- `UNSAFE_FRAME` (not recommended): Control if framing of your site is allowed (say in an `<embed>` or `<iframe>`)

You can also disabled HTTP/2 by setting the `GODEBUG` variable to `http2server=0`. [Please report any issues before disabling HTTP/2,](https://github.com/doamatto/inari/issues/new) unless you know what you are doing.

Because HSTS is enabled by default, it is recommend to run 稲荷 behind a reverse proxy. You can use [Nginx](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/), [Caddy](https://caddyserver.com/docs/quick-starts/reverse-proxy), or any other web server solution (稲荷 was built to run painlessly on Uberspace meteors).

## Sensible Defaults
There are a few sensible defaults:
- **HSTS Preload enabled by default—** you can opt-out of HSTS with an environment variable (`NOHSTS`)
- **FLoC disabled by default—** this cannot be re-enabled and is left disabled for the privacy of the end user of your website
- **Sensible caching for images and font files—** these defaults can be edited with `CACHE_IMAGE_TIME` and `CACHE_FONT_TIME`, and can be disabled by setting both to `0`
- **Gzip compression out-of-the-box—** modern compression helps leagues with reducing data transfers, and is made readily available
- **Framing disabled by default—** this prevents click-through phising and can be disabled with `UNSAFE_FRAME` set to `true`
- **XSS protection out-of-the-box—** basic cross-site protection is provided out of the box, but further security measures should always be taken.

## Acknowledgements
稲荷 is licensed under the 3-Clause BSD license, which you can find in the root of this repository inside of the `LICENSE` file. Special thanks to the team behind Golang for creating amazing documentation to make creating 稲荷 a breeze.
