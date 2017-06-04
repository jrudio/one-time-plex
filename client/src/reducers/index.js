import { combineReducers } from 'redux'
import Users from './users'
import Search from './search'

/*
    Example state

    {
        users: [ {plexUserID: '1954', plexusername: 'jrudio' } ],
        search: {
            errorMSG: '',
            isSearching: false,
            results: [ { title: '...', mediaID: '3146' } ]
        }
    }
*/

export const rootReducer = combineReducers({
    users: Users,
    search: Search
})