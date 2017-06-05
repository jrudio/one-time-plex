import {
    ADD_USER,
    FRIENDS_ADD,
    FRIENDS_FETCH,
    FRIENDS_DONE_FETCHING
 } from '../constants/users'
// import fetch from 'isomorphic-fetch'

const fetchFriends = {
    type: FRIENDS_FETCH
}

const friendsHaveBeenFetched = {
    type: FRIENDS_DONE_FETCHING
}

export const addUser = user => {
    return dispatch => {
        // make POST request to server

        // get result
        let res = {
            err: null,
            result: {
                name: 'jrudio',
                plexUserID: 1,
                assignedMediaID: '3214',
                title: 'Better Call Saul: Lawyer'
            }
        }

        return dispatch({
            type: ADD_USER,
            ...res.result
        })
    }
} 

export const getFriends = () => {
    return dispatch => {
        let friends = [
            { id: '1234', username: 'Guest' },
            { id: '2145', username: 'siirclutch-guest' }
        ]

        dispatch(fetchFriends)

        return fetch(window.otp.url + '/friends')
            .then(r => r.json())
            .then(r => {
                console.log(r)
                dispatch(friendsHaveBeenFetched)

                return dispatch({
                    type: FRIENDS_ADD,
                    friends
                })
            })
            .catch(e => console.error(e))
    }
}