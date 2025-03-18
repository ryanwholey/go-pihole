# go-pihole

A Golang Pi-hole client

Requires Pi-hole Web Interface >= `6`. For <6, use tag <= v0.0.4

## Usage

```go
import "github.com/ryanwholey/go-pihole"

client := pihole.New(pihole.Config{
	BaseURL:  "http://pi.hole"
	Password: "token"
})

record, err := client.LocalDNS.Create(context.Background(), "my-domain.com", "127.0.0.1")
if err != nil {
	log.Fatal(err)
}
```

## Test

```sh
make test
```

### Acceptance

```sh
docker compose up -d
export PIHOLE_URL=http://localhost:8080
export PIHOLE_PASSWORD=test
make acceptance
```
