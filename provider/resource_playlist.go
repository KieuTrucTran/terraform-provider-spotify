/*
author: Truc Tran
date: 2025-01-04
description: Defines the resource logic for creating and managing playlists based on top tracks.
AI Usage: For parts of this code, AI was used to improve structure and functionality.
*/

package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify/v2"
)

func resourceSpotifyPlaylist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpotifyPlaylistCreate,
		ReadContext:   resourceSpotifyPlaylistRead,
		DeleteContext: resourceSpotifyPlaylistDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Spotify playlist to be created",
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Spotify playlist",
				ForceNew:    true,
			},
			"public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether the playlist should be public",
				Default:     true,
				ForceNew:    true,
			},
			"time_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "short_term",
				Description: "The time range for top tracks (short_term, medium_term, long_term)",
				ForceNew:    true,
			},
			"track_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "Number of top tracks to include in the playlist",
				ForceNew:    true,
			},
		},
	}
}

func resourceSpotifyPlaylistCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*spotify.Client)

    // Read inputs from the Terraform configuration
  playlistName := d.Get("name").(string)
  playlistDescription := d.Get("description").(string)
  isPublic := d.Get("public").(bool)
  timeRange := d.Get("time_range").(string)
  trackCount := d.Get("track_count").(int)

  // Validate the time range
  validTimeRanges := map[string]bool{
      "short_term":  true,
      "medium_term": true,
      "long_term":   true,
  }
  if !validTimeRanges[timeRange] {
      return diag.FromErr(fmt.Errorf("invalid time range: %s", timeRange))
  }

  // Fetch the user's top tracks based on time range and limit
  var topTracks *spotify.FullTrackPage
  var err error

  switch timeRange {
  case "short_term":
      topTracks, err = client.CurrentUsersTopTracks(ctx, spotify.Limit(trackCount), spotify.Timerange("short_term"))
  case "medium_term":
      topTracks, err = client.CurrentUsersTopTracks(ctx, spotify.Limit(trackCount), spotify.Timerange("medium_term"))
  case "long_term":
      topTracks, err = client.CurrentUsersTopTracks(ctx, spotify.Limit(trackCount), spotify.Timerange("long_term"))
  }

  if err != nil {
      return diag.FromErr(fmt.Errorf("failed to fetch top tracks: %w", err))
  }

  // Extract track IDs
  trackIDs := make([]spotify.ID, len(topTracks.Tracks))
  for i, track := range topTracks.Tracks {
      trackIDs[i] = track.ID
  }

    // Fetch the current user
    user, err := client.CurrentUser(ctx)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to fetch current user: %w", err))
    }

    // Create a new playlist
    playlist, err := client.CreatePlaylistForUser(ctx, user.ID, playlistName, playlistDescription, isPublic, false)
    if err != nil {
        return diag.FromErr(fmt.Errorf("failed to create playlist: %w", err))
    }

  // Add tracks to the playlist
  _, err = client.AddTracksToPlaylist(ctx, playlist.ID, trackIDs...)
  if err != nil {
      return diag.FromErr(fmt.Errorf("failed to add tracks to playlist: %w", err))
  }

  // Set the playlist ID as the resource ID
  d.SetId(string(playlist.ID))

    return nil
}

func resourceSpotifyPlaylistRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	// Get the playlist ID from the resource data
	playlistID := spotify.ID(d.Id())

	// Fetch playlist details
	playlist, err := client.GetPlaylist(ctx, playlistID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch playlist: %w", err))
	}

	// Update resource data with playlist details
	d.Set("name", playlist.Name)
	d.Set("description", playlist.Description)
	d.Set("public", playlist.IsPublic)

	return nil
}

func resourceSpotifyPlaylistDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*spotify.Client)

	// Get the playlist ID from the resource data
	playlistID := spotify.ID(d.Id())

	// Unfollow (delete) the playlist
	err := client.UnfollowPlaylist(ctx, playlistID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete playlist: %w", err))
	}

	// Remove the resource ID
	d.SetId("")
	return nil
}