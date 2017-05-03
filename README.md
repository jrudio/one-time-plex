# One Time Plex (OTP)

One Time Plex (OTP) allows a Plex user to access one movie or episode of a tv series from your Plex Media Server. OTP also features a REST api to allow you to interact with it programmatically.

### How It Works

* You (the Plex server owner) share your library to a Plex user (usually a family member or friend)
* Add that user’s Plex user id and the id of the desired media to OTP
* OTP will monitor what this user watches, and prevent the user from watching anything other than what you assigned to them
* Once the user is finished, OTP will automatically stop sharing the library with that Plex user

### Disclaimers

Currently, OTP uses Plex’s web hook feature to monitor users, which is a Plex Pass feature. An option that will not rely on this will be added in the future.

### Setup

- Add a plex webhook on your PMS and point it to one time plex
    The default address (of one time plex) should be `localhost:8080/webhook`

- Run `./otp -write` to create a default configuration file

- Edit it to reflect your settings


### Usage

- Grab the user id of the Plex user you are sharing your library with

- Grab the media id (called the rating key in Plex) of the movie or episode of that user is assigned to

- Make a `POST` request to `/api/users/add` with the following in the body:

  ```bash
    plexuserid: <plex-user-id>
    mediaID: <media-id>
  ```

- That's it!