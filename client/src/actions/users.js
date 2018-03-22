import {
    ADD_USER,
    FRIENDS_ADD,
    FRIENDS_FETCH,
    FRIENDS_DONE_FETCHING,
    FRIENDS_ERROR,
    MONITORED_USERS_HAVE_BEEN_FETCHED,
    MONITORED_USERS_FETCH,
    SELECT_USER
 } from '../constants/users'
// import { format } from 'url';
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

const _selectUser = id => ({
    id,
    type: SELECT_USER
})

export const addUser = (user = {}) => {
    return dispatch => {
        // make POST request to server

        // get result
        // let res = {
        //     err: null,
        //     result: user
        // }

        // we are expecting an object with props of
        // - plexUserID
        // - mediaID

        // if (user && !user.mediaID || user && !user.plexUserID) {
        //     console.error('addUser() requires a user object with mediaID and a plexUserID')
        //     return
        // }
        // console.log(user)

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
        dispatch(fetchFriends)

        return fetch(window.otp.url + '/friends')
            .then(r => r.json())
            .then(r => {
                console.log(r)
                dispatch(friendsHaveBeenFetched)

                return dispatch({
                    type: FRIENDS_ADD,
                    friends: r.result
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

export const selectUser = (id = '') => {
    return dispatch => {
        dispatch(_selectUser(id))
    }
}