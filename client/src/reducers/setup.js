import {
    ERROR_PIN_NOT_FOUND,
    // ERROR_PIN_UNAUTHORIZED,
    FETCH_PLEX_PIN,
    PLEX_PIN_ERROR,
    RECEIVED_PLEX_PIN,
    CHECK_PLEX_PIN,
    CHECK_PLEX_PIN_ERROR,
    CHECK_PLEX_PIN_SUCESS,
    PLEX_SERVER_ERROR,
    PLEX_SERVER_FETCHING,
    PLEX_SERVER_RECEIVED,
    PLEX_SERVER_CONNECTION_SELECT,
    PLEX_SERVER_SELECT,
} from '../constants/setup'

export default (state = {
    connectionIndex: 0,
    errMessage: '',
    isAuthorized: false,
    isExpired: false,
    isCheckingPin: false,
    isFetching: true,
    pin: '',
    selectedServer: {
        connection: [
            { local: 0, uri: '' }
        ],
        name: '',
        token: ''
    },
    servers: []
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
        case CHECK_PLEX_PIN:
            return Object.assign({}, state, {
                isCheckingPin: true
            })
        case CHECK_PLEX_PIN_ERROR:
            return Object.assign({}, state, {
                errMessage: action.errMessage,
                isAuthorized: false,
                isCheckingPin: false
            })
        case ERROR_PIN_NOT_FOUND:
            return Object.assign({}, state, {
                isExpired: true
            })
        case CHECK_PLEX_PIN_SUCESS:
            return Object.assign({}, state, {
                isAuthorized: true,
                isCheckingPin: false
            })
        case PLEX_SERVER_FETCHING:
            return Object.assign({}, state, {
                isFetching: true
            })
        case PLEX_SERVER_RECEIVED:
            return Object.assign({}, state, {
                isFetching: false,
                servers: action.servers
            })
        case PLEX_SERVER_ERROR:
            return Object.assign({}, state, {
                errMessage: action.errMessage,
                isFetching: false
            })
        case PLEX_SERVER_SELECT:
            return Object.assign({}, state, {
                selectedServer: {
                    connection: action.server.connection || [{
                        local: 0,
                        uri: ''
                    }],
                    name: action.server.name,
                    token: action.server.accessToken
                }
            })
        case PLEX_SERVER_CONNECTION_SELECT:
            return Object.assign({}, state, {
                connectionIndex: action.index
            })
        default:
            return state
     }
}