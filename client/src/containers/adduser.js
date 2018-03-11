// import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import AddUser from '../components/adduser'
import { addUser, getFriends, selectUser, getMonitoredUsers } from '../actions/users'
// import { getAssignedMedia } from '../actions/assignedmedia'

const mapStateToProps = (state) => {
    let { friends, users } = state
    let {
        isFriendListLoading,
        errorMsg
    } = friends
    let { currentlySelected } = users
    
    // console.log(state)
    
    return {
        currentlySelected,
        errorMsg,
        friends: friends.friendList,
        isFriendListLoading
    }
}

const mapDispatchToProps = dispatch => ({
    selectUser: bindActionCreators(selectUser, dispatch),
    addUser: bindActionCreators(addUser, dispatch),
    getFriends: bindActionCreators(getFriends, dispatch),
    getMonitoredUsers: bindActionCreators(getMonitoredUsers, dispatch)

})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(AddUser)