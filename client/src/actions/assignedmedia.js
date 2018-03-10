import fetch from 'isomorphic-fetch'

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