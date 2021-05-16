# SpaceRouter Authentication Server

[![](https://goreportcard.com//badge/github.com/SpaceRouter/authentication_server)](https://goreportcard.com/report/github.com/SpaceRouter/authentication_server)

## Run
```bash
docker build . -t spacerouter/auth_server
docker run --rm -p 8080:8080 -v /etc/:/etc/ --privileged spacerouter/auth_server
```