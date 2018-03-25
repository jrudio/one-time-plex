import {
    ISFETCHING,
    ERRMESSAGE,
    FETCHED,
    RECEIVEDINFO,
    TEST_CONNECTION
} from '../constants/server'

import fetch from 'isomorphic-fetch'

export const fetchServerInfo = () => {
    return dispatch => {
        dispatch({
            type: ISFETCHING
        })

        return fetch(window.otp.url + '/plex/server')
            .then(r => r.json())
            .then(r => {
                dispatch({ type: FETCHED })

                if (r.error) {
                    dispatch({
                        type: ERRMESSAGE,
                        message: r.error
                    })

                    return
                }

                let {
                    url,
                    token
                } = r.result
                
                dispatch({
                    token,
                    url,
                    type: RECEIVEDINFO
                })
            })
    }
}

export const saveServerInfo = (info = {}) => {
    return dispatch => {
        if (!info.url || !info.token) {
            console.error('need url and token before sending request')
            return
        }
        
        return fetch(window.otp.url + '/plex/server', {
            body: JSON.stringify(info),
            headers: {
                'Accept': 'application/json'
            },
            method: 'PUT',
        }).then(r => r.json())
        .then(r => {
            if (r.error) {
                dispatch({
                    type: ERRMESSAGE,
                    message: r.error
                })

                return
            }

            if (r.result === true) {
                dispatch({
                    type: ERRMESSAGE,
                    message: 'sucessfully saved'
                })
            }
        })
    }
}

export const testServer = (info = {}) => {
    return dispatch => {
        if (!info.url || !info.token) {
            console.error('testServer() needs a token and url')
            return
        }

        let endpoint = window.otp.url + '/plex/test'
        endpoint += '?url=' + encodeURIComponent(info.url)
        endpoint += '&token=' + encodeURIComponent(info.token)
        
        return fetch(endpoint)
            .then(r => r.json())
            .then(r => {
                if (r.error) {
                    console.error('testing connection to plex server failed: ' + r.error)
                    dispatch({
                        type: ERRMESSAGE,
                        errMessage: r.error
                    })
                    return
                }

                let status = "connection to plex server failed"

                if (r.result === true) {
                    status = "connection to plex server sucessful"
                }

                dispatch({
                    status,
                    type: TEST_CONNECTION
                })
            })
    }
}