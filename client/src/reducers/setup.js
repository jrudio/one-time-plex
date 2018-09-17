import {
    FETCH_PLEX_PIN,
    PLEX_PIN_ERROR,
    RECEIVED_PLEX_PIN
} from '../constants/setup'

export default (state = {
    errMessage: '',
    isFetching: true,
    pin: ''
}, action) => {
     switch (action.type) {
        case FETCH_PLEX_PIN:
            return Object.assign({}, state, {
                isFetching: true
            })
        case RECEIVED_PLEX_PIN:
            return Object.assign({}, state, {
                isFetching: false,
                pin: action.pin || ''
            })
        case PLEX_PIN_ERROR:
            return Object.assign({}, state, {
                errMessage: action.errMessage
            })
        default:
            return state
     }
}