import React, { Component } from 'react'
import { Button } from 'react-mdl'
import Search from '../containers/search'

class AssignedMedia extends Component {
    handleUnassign (id) {
        let { unassignFriend, resetSearch } = this.props

        unassignFriend(id)
        resetSearch()
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
            <p><Button colored onClick={() => this.handleUnassign(user.plexUserID)}>unassign</Button></p>
        </div>)
    }
}

export default AssignedMedia