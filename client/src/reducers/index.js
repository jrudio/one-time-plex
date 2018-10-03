import { combineReducers } from 'redux'
import Users from './users'
import Friends from './friends'
import Search from './search'
import Server from './server'
import Setup from './setup'

/*
    Example state

    {
        users: {
            list: { 
                '1234': {
                    plexUserID: '1234',
                    plexUsername: 'jrudio-guest',
                    assignedMediaID: {
                        id: 123,
                        title: 'atlanta - s1e2,
                        status: 'finished'
                    }
                }
            },
            currentlySelected: '1234',
            isLoading: true
        },
        friends: {
            friendList: [ { id: '1234', username: 'jrudio-guest' } ],
            isFriendListLoading: true,
            errMsg: ''
        },
        search: {
            errMessage: '',
            message: '',
            isSearching: false,
            results: [ { title: '...', mediaID: '3146' } ]
        },
        server: {
            url: 'http://192.168.1.200:32400',
            token: 'abc123',
            errMessage: '',
            message: '',
            isFetching: false
        },
        setup: {
            connectionIndex: 0,
            errMessage: '',
            isAuthorized: false,
            isExpired: false,
            isCheckingPin: false,
            isFetching: true,
            pin: 'abc1',
            selectedServer: {
                connection: [{
                    local: 0,
                    uri: ''
                }],
                name: '',
                token: ''
            },
            servers: []
        }
    }
*/

export const rootReducer = combineReducers({
    users: Users,
    friends: Friends,
    search: Search,
    server: Server,
    setup: Setup
})