import React, { Component } from 'react'
import Proptypes from 'prop-types'

class Users extends Component {
    componentWillMount () {
        // let { addUser } = this.props

        // addUser()
        // addUser()
    }
    render () {
        let {
            users
        } = this.props
        
        // console.log(this.props)
        
        if (!users) {
            return <div></div>
        }

        return <div>{this.props.users.map((user, i) => <p key={i}>{user.plexUsername} - <i>{user.title}</i> {'(' + user.assignedMediaID + ')'} </p>)}</div>
    }
}

Users.propTypes = {
    users: Proptypes.array.isRequired,
    addUser: Proptypes.func.isRequired
}

export default Users