import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import AddUser from '../components/adduser'
import { addUser, getFriends } from '../actions/users'

const AddUserContainer = props => {
    console.log(props)
    return <AddUser {...props} />
}

const mapStateToProps = (state) => {
    let { friends } = state
    let {
        isFriendListLoading,
        errorMsg
    } = friends
    
    console.log(state)
    
    return {
        friends: friends.friendList,
        isFriendListLoading,
        errorMsg
    }
}

const mapDispatchToProps = dispatch => ({
    addUser: bindActionCreators(addUser, dispatch),
    getFriends: bindActionCreators(getFriends, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(AddUserContainer)