import React, { Component } from 'react'
import Search from '../containers/search'

class AssignedMedia extends Component {
    handleUnassign (id) {
        let { unassignFriend } = this.props

        unassignFriend(id)
    }
    renderChooseMedia () {
        return (
            <div>
                <p>Not assigned</p>

                <Search />
            </div>
        )
    }
    render () {
        let {
            user,
            currentlySelected
        } = this.props

        if (user === undefined) {
            return this.renderChooseMedia()
        }

        return (<div>
            <p>{user.plexUsername} is assigned to {user.assignedMedia.title}</p>
            <p>status: {user.assignedMedia.status}</p>
            <p><a onClick={() => this.handleUnassign(user.plexUserID)}>unassign</a></p>
        </div>)
    }
}

export default AssignedMedia