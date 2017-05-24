import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import Users from '../components/users'
import { addUser } from '../actions/users'

const UsersContainer = props => {
    console.log(props)
    return <Users users={props.users} addUser={props.addUser} />
}

const mapStateToProps = (state) => {
    let { users } = state
    console.log(state)
    return { users }
}

const mapDispatchToProps = dispatch => ({
    addUser: bindActionCreators(addUser, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(UsersContainer)