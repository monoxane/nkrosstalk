# nkrosstalk
Control an NK-IPS with RossTalk

## Development
It's golang, it's probably bad code, it's my first real golang thing.

There's also a Dockerfile becuase I run things in docker.

## Deployment

```
docker run --name nkrosstalk -d -p 7788:7788 -e NK_HOST=10.101.41.2 -e NK_SIZE=72 monoxane/nkrosstalk:latest
```

Send it RossTalk with line feed like: 
```
XPT <level>:<destination>:<source>
XPT 1:65:50
```