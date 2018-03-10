import React, { Component } from 'react'

class AssignedMedia extends Component {
    render () {
        let {
            user,
            currentlySelected
        } = this.props

        console.log(user)
        console.log(currentlySelected)

        return <p>{user && user.plexUsername}</p>
    }
}

export default AssignedMedia