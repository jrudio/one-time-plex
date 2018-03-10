import React, { Component } from 'react'

class AssignedMedia extends Component {
    render () {
        let {
            user,
            currentlySelected
        } = this.props

        if (user === undefined) {
            return <p>No media assigned</p>
        }

        return (<div>
            <p>{user.plexUsername} is assigned to {user.assignedMedia.title}</p>
            <p>status: {user.assignedMedia.status}</p>
        </div>)
    }
}

export default AssignedMedia