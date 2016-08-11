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
