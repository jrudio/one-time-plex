// import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import AddUser from '../components/adduser'
import { addUser, getFriends, selectUser, getMonitoredUsers } from '../actions/users'
// import { getAssignedMedia } from '../actions/assignedmedia'

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
    selectUser: bindActionCreators(selectUser, dispatch),
    addUser: bindActionCreators(addUser, dispatch),
    getFriends: bindActionCreators(getFriends, dispatch),
    getMonitoredUsers: bindActionCreators(getMonitoredUsers, dispatch),
    // getAssignedMedia: bindActionCreators(getAssignedMedia, dispatch)

})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(AddUser)