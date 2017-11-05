import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import Users from '../components/users'
import { addUser, getMonitoredUsers } from '../actions/users'
import { Spinner } from 'react-mdl'

const UsersContainer = props => {
    let {
        getMonitoredUsers,
        isLoading,
        users,
        addUser
    } = props

    if (isLoading) {
        getMonitoredUsers()
        return <Spinner />
    }

    console.log(users)
    
    return <Users users={users} addUser={addUser} />
}

const mapStateToProps = (state) => {
    let { users } = state

    // console.log(state)

    return { users: users.list, isLoading: users.isLoading }
}

const mapDispatchToProps = dispatch => ({
    addUser: bindActionCreators(addUser, dispatch),
    getMonitoredUsers: bindActionCreators(getMonitoredUsers, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(UsersContainer)