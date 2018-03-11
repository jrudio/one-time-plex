import fetch from 'isomorphic-fetch'

import REMOVE_USER from '../constants/users'

// getAssignedMedia returns all plex friends with any assigned media
export const getAssignedMedia = (mediaID = '') => {
    return dispatch => {
        if (mediaID === '') {
            console.error('getAssignedMedia requires media id')
            return
        }

        return fetch(window.otp.url + '/users')
            .then(r => r.json())
            .then(r => {
                let {
                    result
                } = r

                console.log(result)
            })
            .catch(e => console.error(e))
    }
}

export const unassignFriend = (plexUserID = '') => {
    return dispatch => {
        if (plexUserID === '') {
            console.error("plex user id is required to unassign")
            return
        }

        return fetch(window.otp.url + '/user/' + plexUserID, {
            method: 'DELETE'
        })
            .then(r => r.json())
            .then(r => {
                if (!r.result) {
                    console.error('failed to unassign user: ' + r.error)
                    return
                }

                dispatch({
                    type: REMOVE_USER,
                    id: plexUserID
                })

                return true
            })
    }
}