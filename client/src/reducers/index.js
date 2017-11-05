import { combineReducers } from 'redux'
import Users from './users'
import Friends from './friends'
import Search from './search'

/*
    Example state

    {
        users: {
            list: [ {plexUserID: '1954', plexusername: 'jrudio' } ],
            isLoading: true
        },
        friends: {
            friendList: [ { id: '1234', username: 'jrudio-guest' } ],
            isFriendListLoading: true,
            errMsg: ''
        },
        search: {
            errorMSG: '',
            isSearching: false,
            results: [ { title: '...', mediaID: '3146' } ]
        }
    }
*/

export const rootReducer = combineReducers({
    users: Users,
    friends: Friends,
    search: Search
})