import {
    ADD_USER,
    FRIENDS_ADD,
    FRIENDS_FETCH,
    FRIENDS_DONE_FETCHING,
    FRIENDS_ERROR,
    MONITORED_USERS_HAVE_BEEN_FETCHED,
    MONITORED_USERS_FETCH
 } from '../constants/users'
// import fetch from 'isomorphic-fetch'

const fetchFriends = {
    type: FRIENDS_FETCH
}

const friendsHaveBeenFetched = {
    type: FRIENDS_DONE_FETCHING
}

const fetchMonitoredUsers = {
    type: MONITORED_USERS_FETCH
}

const monitoredUsersFetched = users => ({
    users,
    type: MONITORED_USERS_HAVE_BEEN_FETCHED
})

const fetchedErr = err => ({
    type: FRIENDS_ERROR,
    errorMsg: err
})

export const addUser = (user = {}) => {
    return dispatch => {
        // make POST request to server

        // get result
        // let res = {
        //     err: null,
        //     result: user
        // }

        return fetch(window.otp.url + '/users/add', {
            method: 'POST',
            body: JSON.stringify(user)
        })
            .then(r => r.json())
            .then(r => {
                dispatch({
                    type: ADD_USER,
                    ...r.result
                })
            })
            .catch(e => {
                console.error(e.message)
                dispatch(fetchedErr(e))
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
            .catch(e => {
                dispatch(friendsHaveBeenFetched)

                console.error(e.message)
                dispatch(fetchedErr(e.message))
            })
    }
}

export const getMonitoredUsers = () => {
    return dispatch => {
        dispatch(fetchMonitoredUsers)
        
        return fetch(window.otp.url + '/users')
            .then(r => r.json())
            .then(r => {
                // console.log(r)

                return dispatch(monitoredUsersFetched(r.result))
            })
            // .catch(e => {
                // console.log(e)
                // dispatch(monitoredUsersFetched())

                // console.error(e.message)
                // dispatch(fetchedErr(e.message))
            // })

    }
}

window.ayy = getMonitoredUsers