# go-pihole

A Golang Pi-hole client

Requires Pi-hole Web Interface >= 5.11

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