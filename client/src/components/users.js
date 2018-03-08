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
        
        return <div>{(() => {
            let userBoxes = []
            
            for (const key of Object.keys(users)) {
                userBoxes.push(<p key={key} data-user-id={users[key].plexUserID}>{users[key].plexUserID} - <i>{users[key].title}</i> {'(' + users[key].mediaID + ')'} </p>)
            }

            return userBoxes
        })()}
        </div>
    }
}

Users.propTypes = {
    users: Proptypes.object.isRequired,
    addUser: Proptypes.func.isRequired
}

export default Users