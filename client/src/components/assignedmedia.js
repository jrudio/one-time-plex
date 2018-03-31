import React, { Component } from 'react'
import Search from '../containers/search'

class AssignedMedia extends Component {
    handleUnassign (id) {
        let { unassignFriend } = this.props

        unassignFriend(id)
    }
    renderChooseMedia (username = '') {
        if (!username) {
            console.error('renderChooseMedia() username is empty')
        }

        return (
            <div>
                <p>{username} not assigned to anything</p>

                <Search />
            </div>
        )
    }
    renderSelectAUser () {
        return <p>Select a user</p>
    }
    render () {
        let { user, currentlySelected, currentlySelectedFriend } = this.props

        if (!currentlySelected) {
            return this.renderSelectAUser()
        }

        if (user === undefined) {
            return this.renderChooseMedia(currentlySelectedFriend)
        }

        return (<div>
            <p>{user.plexUsername} is assigned to {user.assignedMedia.title}</p>
            <p>status: {user.assignedMedia.status}</p>
            <p><a onClick={() => this.handleUnassign(user.plexUserID)}>unassign</a></p>
        </div>)
    }
}

export default AssignedMedia