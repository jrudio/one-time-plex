import { ADD_USER } from '../constants/users'
// import fetch from 'isomorphic-fetch'

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