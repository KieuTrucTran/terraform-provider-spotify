/*
author: Truc Tran
date: 2025-01-04
description: Initializes the provider and registers the playlist resource.
AI Usage: For parts of this code, AI was used to improve structure and functionality.
*/

package main

import (
    "context"
    "encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
    
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
    "github.com/zmb3/spotify/v2"
    "golang.org/x/oauth2"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: func() *schema.Provider {
            return Provider()
        },
    })
}

// Provider defines the schema and configuration for the Spotify provider
func Provider() *schema.Provider {
    return &schema.Provider{
        Schema: map[string]*schema.Schema{
            "auth_server": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "http://localhost:27228",
				Description: "Oauth2 Proxy URL",
			},
			"token_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "terraform",
				Description: "Oauth2 Proxy token ID",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SpotifyAuthProxy",
				Description: "Oauth2 Proxy username",
			},
            "api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Oauth2 Proxy API Key",
            },
        },
        ResourcesMap: map[string]*schema.Resource{
            "spotify_playlist": resourceSpotifyPlaylist(),
        },
        ConfigureContextFunc: providerConfigure,
    }
}


// providerConfigure for spotify API access
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	server, err := url.Parse(d.Get("auth_server").(string))
	if err != nil {
		return nil, diag.Errorf("auth_server was not a valid url: %s", err.Error())
	}
	server.Path = path.Join(server.Path, "api/v1/token")
	server.Path = path.Join(server.Path, d.Get("token_id").(string))

	transport := &transport{
		Endpoint: server.String(),
		Username: d.Get("username").(string),
		APIKey:   d.Get("api_key").(string),
	}

	if err := transport.getToken(ctx); err != nil {
		return nil, diag.FromErr(err)
	}

	httpClient := &http.Client{Transport: transport}
	client := spotify.New(httpClient, spotify.WithRetry(true))

	return client, nil
}

type transport struct {
	Endpoint string
	Username string
	APIKey   string
	Base     http.RoundTripper
	token    *oauth2.Token
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if !t.token.Valid() {
		if err := t.getToken(req.Context()); err != nil {
			return nil, err
		}
	}

	t.token.SetAuthHeader(req)

	return t.base().RoundTrip(req)
}

func (t *transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

func (t *transport) getToken(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "POST", t.Endpoint, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(t.Username, t.APIKey)
	resp, err := t.base().RoundTrip(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(body))
	}

	t.token = &oauth2.Token{}
	if err := json.Unmarshal(body, t.token); err != nil {
		return err
	}

	if !t.token.Valid() {
		return errors.New("could not get a valid token")
	}

	return nil
}