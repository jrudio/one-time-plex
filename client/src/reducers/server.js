import {
    ISFETCHING,
    ERRMESSAGE,
    FETCHED,
    RECEIVEDINFO,
    TEST_CONNECTION
} from '../constants/server'

export default (state = { url: '', token: '', message: '', isFetching: false }, action) => {
    switch (action.type) {
        case FETCHED:
            return Object.assign({}, state, {
                isFetching: false
            })
        case ISFETCHING:
            return Object.assign({}, state, {
                isFetching: true
            })
        case ERRMESSAGE:
            let { errMessage } = action
        
            return Object.assign({}, state, {
                errMessage,
                message: ''
            })
        case RECEIVEDINFO:
            let { token, url } = action
        
            return Object.assign({}, state, {
                token,
                url
            })
        case TEST_CONNECTION:
            let { status } = action
        
            return Object.assign({}, state, {
                errMessage: '',
                message: status
            })
        default:
            return state
    }
}