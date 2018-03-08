import { 
    ADD_USER,
    MONITORED_USERS_HAVE_BEEN_FETCHED,
    MONITORED_USERS_FETCH
} from '../constants/users'

export default (state = {}, action) => {
    switch (action.type) {
        case ADD_USER:
        console.log(action)
        
        delete action.type
        let newUser = action
        
        return Object.assign({}, state, {
            list: [].push(state.list, newUser)
        })
        case MONITORED_USERS_FETCH:
        return Object.assign({}, state, {
            isLoading: true
        })
        case MONITORED_USERS_HAVE_BEEN_FETCHED:
            let { users } = action

            return Object.assign({}, state, {
                isLoading: false,
                list: users
            })
        default:
            if (state.isLoading === undefined) {
                state.isLoading = true
            }

            if (state.list === undefined) {
                state.list = []
            }

            return state
    }
}