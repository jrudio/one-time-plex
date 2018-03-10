import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import AssignedMedia from '../components/assignedmedia'
import { getMonitoredUsers } from '../actions/users'

const AssignedMediaContainer = props => {
    let {
        currentlySelected,
        users
    } = props

    let user = users[currentlySelected]

    return <AssignedMedia user={user} {...props} />
}

const mapStateToProps = (state) => {
    let { users } = state
    let {
        isLoading,
        list,
        currentlySelected
    } = users

    return {
        currentlySelected,
        isLoading,
        users: list,
    }
}

const mapDispatchToProps = dispatch => ({
    getAssignedMedia: bindActionCreators(getMonitoredUsers, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(AssignedMediaContainer)