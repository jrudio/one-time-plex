import fetch from 'isomorphic-fetch'

import {
    FETCH_PLEX_PIN,
    PLEX_PIN_ERROR,
    RECEIVED_PLEX_PIN,
    CHECK_PLEX_PIN,
    CHECK_PLEX_PIN_ERROR,
    CHECK_PLEX_PIN_SUCESS,
    ERROR_PIN_NOT_FOUND,
    PLEX_SERVER_ERROR,
    PLEX_SERVER_FETCHING,
    PLEX_SERVER_RECEIVED,
    PLEX_SERVER_SELECT,
    PLEX_SERVER_CONNECTION_SELECT,
    PLEX_SERVER_SET
} from '../constants/setup'

export const fetchPlexPin = () => {
    return dispatch => {
        dispatch({ type: FETCH_PLEX_PIN })

        fetch(window.otp.url + '/plex/pin')
            .then(r => {
                // unprocessable entity == already authorized
                if (r.status === 422) {
                    dispatch({
                        type: CHECK_PLEX_PIN_SUCESS
                    })
                }
                
                return r.json()
            })
            .then(r => {
                if (r.error) {
                    console.error(r.error)

                    dispatch({
                        type: PLEX_PIN_ERROR,
                        errMessage: r.error
                    })

                    return
                }

                const { result } = r

                dispatch({
                    pin: result.code,
                    type: RECEIVED_PLEX_PIN
                })
            })
            .catch(e => {
                console.error(e)

                dispatch({
                    type: PLEX_PIN_ERROR,
                    errMessage: e.message
                })
            })
    }
}

// checPlexPin checks if our plex pin was authorized
export const checkPlexPin = () => {
    return (dispatch, getState) => {
        // do not execute again until the previous checkPlexPin() has finished
        // const { setup } = getState()
        // const { isCheckingPin } = setup

        // if (isCheckingPin) {
        //     console.log('checking pin already in progess')
        //     return
        // }
        
        dispatch({ type: CHECK_PLEX_PIN })

        fetch(window.otp.url + '/plex/pin/check')
            .then(r => {
                if (r.status === 404) {
                    dispatch({ type: ERROR_PIN_NOT_FOUND })
                }

                return r.json()
            })
            .then(r => {
                if (r.error) {
                    // console.error(r.error)

                    dispatch({
                        type: CHECK_PLEX_PIN_ERROR,
                        errMessage: r.error
                    })

                    return
                }

                const { result } = r
                console.log(result)

                dispatch({
                    type: CHECK_PLEX_PIN_SUCESS
                })
            })
            .catch(e => {
                // console.error(e)

                dispatch({
                    type: CHECK_PLEX_PIN_ERROR,
                    errMessage: e.message
                })
            })
    }
}

export const getPlexServers = () => {
    return dispatch => {
        dispatch({ type: PLEX_SERVER_FETCHING })
        
        fetch(window.otp.url + '/plex/servers')
            .then(r => r.json())
            .then(r => {
                const { result } = r

                dispatch({
                    servers: result,
                    type: PLEX_SERVER_RECEIVED,
                })
            }).catch(e => {
                console.error(e.message)
                dispatch({
                    type: PLEX_SERVER_ERROR
                })
            })
    }
}

export const selectServer = (server = {}) => {
    return dispatch => {
        if (server.name === '') {
            console.error('server object is empty')
            return
        }

        dispatch({
            server,
            type: PLEX_SERVER_SELECT
        })
    }
}

export const selectConnection = (index = 0) => {
    return dispatch => {
        dispatch({
            index,
            type: PLEX_SERVER_CONNECTION_SELECT
        })
    }
}

export const clearServerSelection = (index = 0) => {
    return dispatch => {
        dispatch({
            type: PLEX_SERVER_SELECT,
            server: {
                name: '',
                connection: [],
                token: ''
            }
        })
    }
}

export const setPlexServer = (server = {}, serverConnectionIndex) => {
    return dispatch => {
        dispatch({
            type: PLEX_SERVER_SET
        })
        if (server.connection.length < 1) {
            console.error('invalid server when calling setPlexServer()')
            console.log(server)
            return
        }

        const connection = server.connection[serverConnectionIndex]
        const newServer = {
            name: server.name,
            url:  connection.uri
        }

        fetch(window.otp.url + '/plex/server', {
            body: JSON.stringify(newServer),
            headers: {
                'Accept': 'application/json'
            },
            method: 'PUT'
        }).then(r => {
            const result = r.json()

            if (r.status !== 200) {
                const errMessage = result.error

                throw new Error(`could not set plex server: ${r.status} status code - ${errMessage}`)
            }

            return result
        })
        .catch(e => {
            console.error(e.message)
        })
    }
}