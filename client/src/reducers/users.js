import { 
    ADD_USER,
    MONITORED_USERS_HAVE_BEEN_FETCHED,
    MONITORED_USERS_FETCH,
    SELECT_USER,
    REMOVE_USER
} from '../constants/users'

export default (state = {}, action) => {
    let {
        id
    } = action
    
    switch (action.type) {
        case ADD_USER:
        // console.log(action)
        
            delete action.type
            let newUser = action
            
            return Object.assign({}, state, {
                list: [].push(state.list, newUser)
            })
        case REMOVE_USER:
            // aka unassign
            let { list } = state

            list.filter(user => user.id !== id)

            return Object.assign({}, state, {
                list
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
        case SELECT_USER:
            return Object.assign({}, state, {
                currentlySelected: id
            })

        default:
            if (state.isLoading === undefined) {
                state.isLoading = true
            }

            if (state.list === undefined) {
                state.list = []
            }

            if (state.currentlySelected === undefined) {
                state.currentlySelected = ''
            }

            return state
    }
}