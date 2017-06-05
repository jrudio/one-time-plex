import { ADD_USER } from '../constants/users'

export default (state = [], action) => {
    switch (action.type) {
        case ADD_USER:
            let newState = []
            console.log(action)

            delete action.type

            let newUser = action

            newState.push(newUser)

            // console.log(newState)

            return [].concat(state, newState)
        default:
            return state
    }
}