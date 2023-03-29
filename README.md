# Preamble

It's a sort of web framework but don't use this, it's not for you.

# Intro

The idea was to make a kind-of-sort-of extension to (Gorilla Mux)[github.com/gorilla/mux] that allowed one to quickly
integrate with Postgresql (other DBs could be added to be supported)

But, as mentioned before, don't use this, it's *really* not for you!

# How-to

Ok, we're going here are we, I mean I told you not to use it but you've kept going.<br>
<br>
So, this was developed to make my life a little easier when knocking together websites.<br>

* It's probably not pretty
* It's probably not clever
* It's probably not intelligent
* It's probably not quicker or easier in the long run for you to learn it than say, Gorilla mux, but it works for me.
* It has some prebuilt "stuff" (models) for users, customers etc... Some need work.
* It has some not working stuff that I may, when I can/need the feature, get around to adding/fixing.
* It sort of integrates Stripe support.
* It supports a yaml config, so that's something I guess...

To be honest if you want any-more than that, just have a read through the code and comments that are there.
<br><br>
OR, as mentioned in previous sections
<br><br>
<b>just don't use it - it's not for you</b> :-D

# "Why is it not for me?"

I've made this public because I'm playing, maybe one day it will be for you, but that day isn't today :-)
<br>
If you really want to use it then have at it! I won't stop ya!

# Example

```go
package main

import (
	"context"
	"fmt"
	"gitea.viles.uk/dcp/web-framework/app"
	"gitea.viles.uk/dcp/web-framework/environment"
	"gitea.viles.uk/dcp/web-framework/responses"
	"gitea.viles.uk/dcp/web-framework/routing"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type YourServerConfig struct {
	*app.Server
}

func (s *YourServerConfig) Index(w http.ResponseWriter, _ *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to the site!")
}

func (s *YourServerConfig) RoutePathsDefinitions() *[]routing.Route {
	routeDefinitions := &[]routing.Route{
		{
			Name:            "Home",
			Description:     "Index of the site",
			Path:            "/index",
			HandlerFunc:     s.Index,
			RequestMethod:   http.MethodGet,
			HasJSONResponse: true,
		},
	}
	return routeDefinitions
}

func WithoutSSL() {
	var err error
	server := &YourServerConfig{
		&app.Server{},
	}
	server.Config, err = environment.Initialise()
	if err != nil {
		log.Fatalln(err)
	}
	server.Initialise(server.RoutePathsDefinitions())
	endpoint := fmt.Sprintf("%s:%d", server.Config.App.IP, server.Config.App.Port)
	log.Printf("Listening on address %s", server.Config.App.IP)
	log.Printf("Listening on port %d", server.Config.App.Port)
	log.Fatal(http.ListenAndServe(endpoint, server.Router))
}

func WithSSL() {
	var err error
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	server := &YourServerConfig{
		&app.Server{},
	}
	server.Config, err = environment.Initialise()
	if err != nil {
		log.Fatalln(err)
	}
	server.Initialise(server.RoutePathsDefinitions())
	httpsSrv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", server.Config.App.IP, server.Config.App.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      server.Router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Listening on address %s", server.Config.App.IP)
		log.Printf("Listening on port %d", server.Config.App.Port)
		if err := httpsSrv.ListenAndServeTLS("ssl/cert.pem", "ssl/key.pem"); err != nil {
			log.Println("server error:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	httpsSrv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

func main() {
	WithoutSSL()
	//OR
	WithSSL()
}
```

You can set up an `app.config` to configure the app like so:

```yaml
#app.config
web:
  fqdn: "https://example.com:8081"
  env: "DEV"
  ip: ""
  port: 8081
  domain_short: "example.com"
  ssl:
    private_key: ""
    public-key: ""
    ca_key: ""
api:
  api_secret: "PASSWORD"
  api_endpoint: "http://example.com:8082"
db:
  host: "IP_ADDR"
  port: 5432
  name: "DB_NAME"
  username: "USERNAME"
  password: "PASSWORD"
stripe:
  secret_key: "ENTER"
  public_key: "ENTER"
  webhook_secret: "ENTER"
  account_id: "ENTER"
```