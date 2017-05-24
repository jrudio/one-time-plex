import { combineReducers } from 'redux'
import Users from './users'

export const rootReducer = combineReducers({
    users: Users
})