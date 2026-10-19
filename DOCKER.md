# Docker deployment

```sh
docker build -t datastar-tiptap .
docker run --rm -p 127.0.0.1:3000:3000 datastar-tiptap
```

The build regenerates Tailwind CSS and the bundled browser dependencies,
then compiles the Go server. The runtime runs as UID/GID 10001 on port 3000.
For a dedicated domain, leave `BASE_PATH` unset. For the portfolio demo, pass
`-e BASE_PATH=/demo/datastar-tiptap`, route that path to port 3000, and strip
it before forwarding. The reverse proxy must support unbuffered SSE responses.

In Coolify choose Dockerfile, repository-root build context, `/Dockerfile`,
and exposed port 3000. Content is held in memory and resets on restart.

Deployments are managed by central Coolify on Shared. Pushes to the configured
production branch deploy automatically through the Forgejo bridge. Coolify's Caddy
proxy serves the public domains over HTTPS; there is no manual deploy step.
