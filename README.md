# One Time Plex Reloaded (otp-reloaded)

Limit your Plex users to one piece of media

### Core features

- Report user violation <sup>1</sup>
- Handle user violation appropriately
- Easily add user to be monitored
- Easily remove user from being monitored
- Able to add existing Plex accounts to PMS with appropriate access to media
- A terminal client to easily interact with the server

<sup>1</sup> A violation is a user that accesses media they were not assigned to

### Setup
                             
Add a plex webhook and point it to one time plex
The default address (of one time plex) should be `localhost:8080`

### Usage

- On first run it will create a `config.toml` file

- Edit it to your preference

- Make sure you supply a plex token in the config if you plan on inviting the user when using `/api/v1/request/add`

- run `otp`

##### Beware that the endpoints are *currently* NOT secured by an API key (Don't expose this app to the internet!)

`GET /api/v1/monitor/start`:

Will start monitoring your Plex sessions

`POST /api/v1/request/add`:

  - Add query `?plexpass=1` to restrict media with labels

  - By default an invite will be sent to the `plexUsername`. If the user is already invited to your server, then append the query `&invite=0`

  - required in post form:
    - `plexUsername: jrudio-guest`
    
    - `ratingKey: 6` (the id that Plex uses for media)
