# Terraform Provider for Spotify

## Abstract

This document provides a comprehensive overview of the development and implementation of a Terraform provider for Spotify. The provider automates the creation of personalized playlists based on a user's listening history. The project integrates Terraform's infrastructure-as-code approach with Spotify's music service, offering an automated way to create playlists without manual input.

## 1. Introduction

This project demonstrates the creation of a custom Terraform provider for Spotify that automates the creation of personalized playlists based on a user’s listening history. The provider interacts with Spotify's Web API to fetch a user's most-listened tracks, based on defined time ranges (e.g., 4 weeks, 6 months, or 1 year). This integration enhances the automation capabilities of Terraform by allowing users to generate playlists without manual interaction. The primary goal is to deepen the understanding of API integration and provider development within the context of DevOps automation.

### Background

Spotify is a popular music streaming service that offers a vast library of songs, albums, and playlists. However, Spotify does not provide an easy way to automate playlist creation based on a user's listening history over defined periods. Users typically rely on manual selection or third-party applications (such as [stats.fm](https://stats.fm)) to create personalized playlists. 

This project attempts to fill this gap by creating a Terraform provider that automates this process using Spotify’s Web API.

### Objective

The primary objective of this project is to develop a Terraform provider that automates the creation of personalized Spotify playlists based on a user's listening history. This project aims to deepen the understanding of automation, API integration, and provider development within the context of a DevOps course. By leveraging the Spotify API and Terraform, users can automate playlist creation without manual input.


## 2. Methodology 

This section outlines the steps taken to complete the project, from setup to implementation, and the reasoning behind each choice.

### 2.1. Prerequisites - Tools and Technologies

To complete this project, the following tools and technologies are required:

- [**Go**](https://go.dev/): A programming language used to build the custom provider.
- [**Terraform**](https://www.terraform.io/): An infrastructure-as-code tool for automating resource creation and management.
- [**Spotify Developer Account**](https://developer.spotify.com/): Required to create an app and obtain API credentials for [Spotify's Web API](https://developer.spotify.com/documentation/web-api).

### 2.2. Environment Setup

#### 2.2.1. Directory Structure

Set up the project directory structure as follows:

```
/provider
    - main.go
    - resource_playlist.go
/spotify-auth
    - main.go
/terraform
    - main.tf
    - variables.tf
    - outputs.tf
    - terraform.tfvars
```

- The **provider** directory contains the code for the custom provider.
- The **spotify-auth** directory contains the authentication logic for communicating with Spotify’s API.
- The **terraform** directory contains the Terraform configuration files used to create and manage resources.

#### 2.2.2. Go Module Initialization

After setting up Go and Terraform, initialize a Go module for the `provider` directory. This allows Go to manage dependencies. Run the following command:

```
go mod init github.com/username/terraform-provider-spotify
```

Example:
```
go mod init github.com/KieuTrucTran/terraform-provider-spotify
```

#### 2.2.3. Dependencies

Install the required Go packages for the provider:

```
go get github.com/hashicorp/terraform-plugin-sdk/v2
go get github.com/zmb3/spotify
go get golang.org/x/oauth2
```
Note: `github.com/zmb3/spotify` is a Go wrapper for working with Spotify's [Web API](https://developer.spotify.com/documentation/web-api). It aims to support every task listed in the Web API Endpoint Reference, located [here](https://developer.spotify.com/documentation/web-api).

## 3. Implementation Details

This section explains the structure and functionality of the provider code. 

### 3.1. Custom Provider Structure

The provider code consists of two primary files:

- **main.go**: Initializes the provider, defines the configuration schema, and handles authentication.
- **resource_playlist.go**: Defines the "spotify_playlist" resource, which manages the creation, reading, and deletion of playlists.

### 3.2. Code Walkthrough

#### 3.2.1. `main.go`

This file sets up the provider schema and handles the communication with Spotify's API. Key components include:

- **Provider Schema**: Configuration parameters such as `auth_server`, `token_id`, `username`, and `api_key` allow users to authenticate and interact with the Spotify API.
- **Authentication**: OAuth2 authentication is implemented using the `spotify.New()` function to create a Spotify client with the necessary credentials.

Explanation:
- `auth_server`, `token_id`, and `username` define the authentication method and credentials used to communicate with the OAuth2 server. The API key is essential for authorization and access to Spotify's API.

#### 3.2.2. `resource_playlist.go`

This file defines the logic for managing Spotify playlists in Terraform. The resource includes:

- **Create**: Fetches the user's top tracks based on a time range (e.g., short-term, medium-term, long-term) and creates a new playlist.
- **Read**: Retrieves details about the playlist.
- **Delete**: Removes the playlist from the user’s Spotify account.

#### 3.2.3. Building the Executable

The final executable is built using Go’s `go build` command. This results in the `terraform-provider-spotify.exe` file. It should be copied into the appropriate directory to be used by Terraform:

```
go mod tidy
go build -o terraform-provider-spotify.exe
```

Place the executable into the plugin directory:
```
C:\Users\<your_user>\AppData\Roaming\terraform.d\plugins\<your_hostname>\spotifyprovider\spotify
```

For example:
```
C:\Users\kieut\AppData\Roaming\terraform.d\plugins\kateeruce\spotifyprovider\spotify
```

For Windows, the general path is `%APPDATA%\terraform.d\plugins\${host_name}/${namespace}/${type}/${version}/${target}`.


### 3.3. Authentication with Spotify

To authenticate with Spotify's API, users need to create a Spotify Developer application. The application provides the client ID and client secret used to authenticate with the API.

Follow this [tutorial](https://developer.hashicorp.com/terraform/tutorials/community-providers/spotify-playlist#create-spotify-developer-app):

1. Set up a [Spotify Developer](https://developer.spotify.com/) account and create an app to get API access.
2. If you plan to run this proxy locally, configure the redirect URI for the app as `http://localhost:27228/spotify_callback`, like described [here](https://github.com/conradludgate/terraform-provider-spotify/tree/main/spotify_auth_proxy).
3. Store the credentials (client ID and client secret) in an `.env` file:

```
SPOTIFY_CLIENT_ID=<your_client_id>
SPOTIFY_CLIENT_SECRET=<your_client_secret>
```


## 4. OAuth2 Authentication

The OAuth2 authentication flow is implemented to securely access Spotify's API. The `main.go` file in the `spotify-auth` directory handles the authentication process.

To create an instance of a Spotify authentication server, follow these steps:

1. Set up the `spotify-auth` directory and create a `main.go` file.
2. For a simple way to manage the spotify oauth2 tokens is to use this [source code](https://github.com/conradludgate/terraform-provider-spotify/blob/main/spotify_auth_proxy/main.go), which acts as an interface between a client and the Spotify oauth API.
3. Add the [`user-top-read`]((https://pkg.go.dev/github.com/zmb3/spotify#pkg-constants)) scope to seek read access to a user's top tracks and artists.

The authentication server will handle the OAuth2 flow and generate access tokens, allowing Terraform to interact with Spotify's API. For more information, read [this](https://github.com/conradludgate/terraform-provider-spotify/tree/main/spotify_auth_proxy).

### 4.1. Go Module Initialization

Initialize a Go module for the `spotify-auth` directory. This allows Go to manage dependencies. Run the following command:
```
go mod init github.com/conradludgate/terraform-provider-spotify/spotify_auth_proxy
```

### 4.2. Run Authorization

1. Navigate to the `spotify-auth` directory:
    - `go mod tidy`
    - `go build -o spotify-auth.exe`
2. Run the executable:

```
.\spotify-auth.exe
```

Example output:
```
APIKey:   xxxxxxx...
Auth URL: xxxxxxx...
Authorization successful
```

- Save the API key to use in the Terraform configuration (`terraform.tfvars`).
- Open a browser and navigate to the Auth URL. It should redirect you to Spotify to log in. After you log in, the auth server will redirect you back to the page where it should confirm that you've authorized correctly.
(Alternatively, simply click on the `Auth URL`, authenticate, and upon successful authentication, the server will display `Authorization successful`.)


## 5. Terraform Configuration

The Terraform configuration files define the necessary variables, outputs, and the main provider configuration.

### 5.1. `main.tf`

Defines the custom provider and the resource to be managed:

```hcl
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
```

### 5.2. `variables.tf`

Declares variables for the user to customize playlist properties and the API key:

```hcl
variable "spotify_api_key" {
  type = string
  description = "The API key for authenticating with the Spotify API"
}
```

### 5.3. `outputs.tf`

Defines outputs, such as the playlist ID:

```hcl
output "playlist_id" {
  value       = spotify_playlist.top_tracks.id
  description = "The ID of the created Spotify playlist"
}
```

### 5.4. `terraform.tfvars`

Stores variables like the Spotify API key, which should be added to `.gitignore` for security:

```hcl
spotify_api_key = "your_api_key_here"
```

---

### 5.5. Terraform Configuration Example

```hcl
provider "spotify" {
  api_key = "your_spotify_api_key"
}

resource "spotify_playlist" "top_tracks" {
  name        = "My Top Tracks"
  description = "A playlist of my top tracks from the last month"
  public      = true
  time_range  = "short_term"
  track_count = 50
}

output "playlist_id" {
  value = spotify_playlist.top_tracks.id
}
```

This configuration creates a public playlist named "My Top Tracks" with the user's top 50 tracks from the short-term time range.


## 6. Usage Guide

### 6.1. Initialize and Plan

Run the following command to initialize Terraform:

```
terraform init
```

Then, run `terraform plan` to view what changes Terraform will make.

### 6.2. Apply Changes - Create the Playlist

To create the playlist on Spotify, run:

```
terraform apply
```

### 6.3. Clean Up

To destroy the resources and remove unused configurations, run:

```
terraform destroy
```


## 7. Results

The Terraform provider successfully automates the creation of personalized Spotify playlists based on a user's listening history. The provider integrates seamlessly with Terraform, allowing users to define their playlist preferences declaratively.


## 8. References

- [Create a Spotify playlist with Terraform](https://developer.hashicorp.com/terraform/tutorials/community-providers/spotify-playlist)
- [zmb3/spotify Go Library - Go wrapper for working with Spotify's Web API](https://pkg.go.dev/github.com/zmb3/spotify)
- [terraform-provider-spotify](https://github.com/conradludgate/terraform-provider-spotify)
- [Dein eigener Provider](https://www.terraform-in-der-praxis.de/provider/eigener-provider.html)
- [Get User's Top Items](https://developer.spotify.com/documentation/web-api/reference/get-users-top-artists-and-tracks)
- [(Video) Devops project: Manage SPOTIFY using TERRAFORM!](https://www.youtube.com/watch?v=LjJLZRi_zGU&t=19s&ab_channel=CloudChamp)
- [Writing Custom Terraform Providers](https://www.hashicorp.com/blog/writing-custom-terraform-providers)
- [Beginner’s Guide to Creating a Terraform Provider](https://www.integralist.co.uk/posts/terraform-build-a-provider/)