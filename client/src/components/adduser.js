import React, { Component } from 'react'
import {
    Grid,
    Cell,
    List,
    ListItem,
    ListItemContent,
    Spinner,
} from 'react-mdl'
import Proptypes from 'prop-types'

const styles = {
    friendList: {
        overflowY: 'overlay',
        height: '300px'
    },
    navback: {
        float: 'left',
        cursor: 'pointer'
    },
    pointer: {
        cursor: 'pointer'
    }
}

class AddUser extends Component {
    componentWillMount () {
        let { getFriends, getMonitoredUsers } = this.props

        getFriends()
        getMonitoredUsers()
    }
    handleSelectFriend (user) {
        let {
            id
        } = user

        let {
            selectUser
        } = this.props

        selectUser(id)
    }
    renderFriend (friend) {
        return (
            <div>
                {friend.username}
            </div>
        )
    }
    renderSelectedFriend (friend) {
        return (
            <b>
                {friend.username}
            </b>
        )
    }
    renderSelectFriend () {
        let {
            friends,
            isFriendListLoading,
            errorMsg,
            currentlySelected // selected plex user id
        } = this.props

        if (friends && friends.length === 0) {
            return (
                <div>
                    <h5>Friend List</h5>
                    
                    {isFriendListLoading && (<Spinner />)}
                    {!isFriendListLoading && errorMsg && (<pre>{errorMsg}</pre>)}
                </div>
            )
        }

        return (
            <Grid>
                <Cell col={12}>
                    <List style={styles.friendList}>
                        {friends.map((friend, i) => (<ListItem key={i} style={styles.pointer} onClick={(e) => {
                            this.handleSelectFriend({ id: friend.id, username: friend.username })
                            }} >
                            <ListItemContent>
                                {(() => {
                                    if (friend.id === currentlySelected) {
                                       return this.renderSelectedFriend(friend)
                                    }

                                    return this.renderFriend(friend)
                                })()}
                            </ListItemContent>
                        </ListItem>))
                        }
                    </List>
                </Cell>
            </Grid>
        )
    }
    render () {
        return (
            <Grid>
                <Cell col={11}>
                    {this.renderSelectFriend()}
                </Cell>
            </Grid>
        )
    }
}

AddUser.propTypes = {
    friends: Proptypes.array.isRequired,
    isFriendListLoading: Proptypes.bool.isRequired,
    errorMsg: Proptypes.string,
    getFriends: Proptypes.func.isRequired
}

export default AddUser
