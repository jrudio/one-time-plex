import { ADD_USER } from '../constants/users'

export default (state = [], action) => {
    switch (action.type) {
        case ADD_USER:
            let newState = []

            let newUser = {
                name: action.name,
                plexUserID: action.plexUserID,
                assignedMediaID: action.assignedMediaID
            }

            newState.push(newUser)

            console.log(newState)

            return [].concat(state, newState)
        default:
            return state
    }
}