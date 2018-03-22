# One Time Plex (OTP)

One Time Plex (OTP) allows a Plex user to access one movie or episode of a tv series from your Plex Media Server.

### How It Works

* 

### How to setup

If you would like, you can download the appropriate binary for your system [here](https://github.com/jrudio/one-time-plex/releases).

If not, here's how you can build this repo from scratch:

1. clone this repo
2. make sure [Go is installed](https://golang.org/dl/)
3. then install [Go dep](https://github.com/golang/dep)
4. `cd server/`
3. run `dep ensure`
4. run `go build -o one-time-plex`


### Notes

*The following instructions are for mainly for building a binary when changes to the front end occur*

Here is how to bake the front end files into the server, so we get a single binary:

Make sure these are installed:

- Go
- Go dep (dependency tool)
- npm or yarn
- [go-bindata](github.com/jteeuwen/go-bindata/)
- [go-bindata-assetfs](github.com/elazarl/go-bindata-assetfs/)

Front end

1. `yarn run build` or `npm run build` to create a production build of the front end
2. copy `client/build` to `server/build`


Back end

1. `cd server/`
2. run `go-bindata-assetfs ./build/...` this generates `bindata.go`
3. `go build -o one-time-plex`





