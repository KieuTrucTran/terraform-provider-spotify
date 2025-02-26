variable "spotify_api_key" {
  type = string
  description = "Set this as the APIKey that the authorization proxy server outputs"
}

variable "playlist_name" {
  type        = string
  description = "Name of the playlist to be created"
  default     = "My Top Tracks from the last month created by Terraform"
}

variable "playlist_description" {
  type        = string
  description = "Description of the playlist"
  default     = "A playlist of my favorite tracks from the last month"
}

variable "playlist_public" {
  type        = bool
  description = "Whether the playlist should be public"
  default     = true
}

variable "playlist_time_range" {
  type        = string
  description = "Time range for top tracks (short_term, medium_term, long_term)"
  default     = "long_term"
}

variable "playlist_track_count" {
  type        = number
  description = "Number of top tracks to include in the playlist"
  default     = 50
}
