import fetch from 'isomorphic-fetch'

import {
    FETCH_PLEX_PIN,
    PLEX_PIN_ERROR,
    RECEIVED_PLEX_PIN
} from '../constants/setup'

export const fetchPlexPin = () => {
    return dispatch => {
        dispatch({ type: FETCH_PLEX_PIN })

        fetch(window.otp.url + '/plex/pin')
            .then(r => r.json())
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
                    pin: result.pin,
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