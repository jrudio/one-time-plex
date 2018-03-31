import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import AssignedMedia from '../components/assignedmedia'
import { getMonitoredUsers } from '../actions/users'
import { unassignFriend } from '../actions/assignedmedia'

const AssignedMediaContainer = props => {
    let {
        currentlySelected,
        users
    } = props

    let user = users[currentlySelected]

    return <AssignedMedia user={user} {...props} />
}

const mapStateToProps = (state) => {
    let { users, friends } = state
    let {
        isLoading,
        list,
        currentlySelected
    } = users

    let currentlySelectedFriend = null

    if (currentlySelected) {
        // search through friends array to find selected friend
        currentlySelectedFriend = friends.friendList.find(friend => friend.id === currentlySelected).username
    }

    return {
        currentlySelected,
        currentlySelectedFriend,
        isLoading,
        users: list
    }
}

const mapDispatchToProps = dispatch => ({
    getAssignedMedia: bindActionCreators(getMonitoredUsers, dispatch),
    unassignFriend: bindActionCreators(unassignFriend, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(AssignedMediaContainer)