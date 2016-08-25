# One Time Plex Reloaded (otp-reloaded)

A server with a simple API to limit Plex users to one piece of media

##### This is a leaner version of [one-time-plex](https://github.com/jrudio/one-time-plex)

### Goals

Here are some goals I hope to achieve with this version:

1. Accept an API key for authentication
2. Endpoints need to be clear and concise
3. Stick to core features and expand later
4. Don't focus on frontend code until core is complete

### Core features

- Report user violation <sup>1</sup>
- Handle user violation appropriately
- Easily add user to be monitored
- Easily remove user from being monitored
- Able to add existing Plex accounts to PMS with appropriate access to media
- A terminal client to easily interact with the server
- That's all I can think of for now... lol

<sup>1</sup> A violation is a user that accesses media they were not assigned to

### Build

To compile this source you need:

  - Go
  - govendor

Then you need to:

1. cd into `cmd/otp`
2. `govendor sync` (which gathers required dependencies)
3. `go install`

### Usage

- On first run it will create a `config.toml` file

- Edit it to your preference

- Make sure you supply a plex token in the config if you plan on inviting the user when using `/api/v1/request/add`

- run `otp`


### [Command-line tool](cmd/tool) (Use the command-line tool to manually execute certain tasks)

### API (Using the api will help automate some tasks)

##### Beware that the endpoints are *currently* NOT secured by an API key (Don't expose this app to the internet!)

`GET /api/v1/monitor/start`:

Will start monitoring your Plex sessions

`POST /api/v1/request/add`:

  - Add query `?plexpass=1` to restrict media with labels

  - ~~By default an invite will be sent to the `plexUsername`. If the user is already invited to your server, then append the query `&invite=0`~~ NOT IMPLEMENTED

  - required in post form:
    - `plexUsername: jrudio-guest`
    
    - `ratingKey: 6` (the id that Plex uses for media)


##### The following endpoints require OTP to be on the Plex Server and must have proper permission to copy files/folders:

~~`POST /api/v1/library/shared/new`~~

  - ~~Post form:~~
    - ~~`plexUsername: jrudio-guest`~~
    
    - ~~`ratingKey: 6` (the id that Plex uses for media~~ NOT IMPLEMENTED
  
