import { SEARCH_PLEX } from '../constants/search'
import fetch from 'isomorphic-fetch'

export const reset = () => {
    return dispatch => {
        dispatch({
            type: SEARCH_PLEX,
            results: []
        })
    }
}

export const searchPlexForMedia = (title = '') => {
    return dispatch => {
        if (title === '') {
            console.error('title cannot be undefined when searching plex')
            return
        }

        title = encodeURIComponent(title)

        console.log(title)

        return fetch(window.otp.url + '/search?title=' + title, { mode: 'cors' })
            .then(r => r.json())
            .then(r => {
                let {
                    result
                } = r

                console.log(r)

                return dispatch({
                    type: SEARCH_PLEX,
                    results: result
                })
            })
            .catch(e => console.error(e))
    }
}

export const getMetadata = (mediaID = '') => {
    return dispatch => {
        if (mediaID === '') {
            console.error('mediaID required for getMetadata')
            return
        }

        return fetch(window.otp.url + '/metadata?mediaid=' + mediaID)
            .then(r => r.json())
            .then(r => {
                let {
                    result
                } = r

                console.log(r)

                return dispatch({
                    type: SEARCH_PLEX,
                    results: result
                })
            })
            .catch(e => console.error(e))
    }
}