import { combineReducers } from 'redux'
import Users from './users'
import Friends from './friends'
import Search from './search'

/*
    Example state

    {
        users: {
            list: [ 
                { 
                    plexUserID: '1234',
                    plexUsername: 'jrudio-guest',
                    assignedMediaID: {
                        id: 123,
                        title: 'atlanta - s1e2,
                        status: 'finished'
                    }
                }
            ],
            currentlySelected: '1234',
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