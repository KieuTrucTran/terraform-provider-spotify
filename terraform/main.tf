terraform {
  required_providers {
    spotify = {
      source  = "kateeruce/spotifyprovider/spotify"
      version = "~> 1.0.0"
    }
  }
}

provider "spotify" {
  api_key = var.spotify_api_key
}

resource "spotify_playlist" "top_tracks" {
  name        = var.playlist_name
  description = var.playlist_description
  public      = var.playlist_public
  time_range  = var.playlist_time_range
  track_count = var.playlist_track_count
}
