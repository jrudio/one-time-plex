import {
    FRIENDS_FETCH,
    FRIENDS_DONE_FETCHING,
    FRIENDS_ADD
 } from '../constants/users'

export default (state = { friendList: [], isFriendListLoading: true }, action) => {
    switch (action.type) {
        case FRIENDS_FETCH:
            return Object.assign({}, state, { isFriendListLoading: true })
        case FRIENDS_DONE_FETCHING:
            return Object.assign({}, state, { isFriendListLoading: false })
        case FRIENDS_ADD:
            let friendList = [].concat(action.friends)

            return Object.assign({}, state, { friendList })
        default:
            return state
    }
}
