import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import AssignedMedia from '../components/assignedmedia'
import { getMonitoredUsers, selectUser } from '../actions/users'
import { unassignFriend } from '../actions/assignedmedia'
import { reset } from '../actions/search'

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
    unassignFriend: bindActionCreators(unassignFriend, dispatch),
    resetSearch: bindActionCreators(reset, dispatch),
    selectUser: bindActionCreators(selectUser, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(AssignedMediaContainer)