# go-pihole

A Golang Pi-hole client

Requires Pi-hole Web Interface >= `5.11` (Docker tag >= `2022.02.1`)

## Usage

```go
import "github.com/ryanwholey/go-pihole"

client := pihole.New(pihole.Config{
	BaseURL:  "http://pi.hole"
	APIToken: "8c4e081d..."
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
docker-compose up -d --build
export PIHOLE_URL=http://localhost:8080
export PIHOLE_API_TOKEN=7b3d979ca8330a94fa7e9e1b466d8b99e0bcdea1ec90596c0dcc8d7ef6b4300c
make acceptance
```
