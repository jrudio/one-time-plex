import { SEARCH_PLEX } from '../constants/search'

/*
    {
        errorMsg: '',
        isSearching: <bool>
        results: []
    }
*/
export default (state = { errorMsg: '', isSearching: false, results: [] }, action) => {
    switch (action.type) {
        case SEARCH_PLEX:
            let { results } = action

            return Object.assign({}, state, { results })
        default:
            return state
    }
}