import React, { Component } from 'react'
import {
    Grid,
    Cell,
    Button,
    List,
    ListItem,
    ListItemContent,
    ListItemAction,
    Spinner,
    Icon,
    // Checkbox,
    Textfield
} from 'react-mdl'
import Proptypes from 'prop-types'

const styles = {
    friendList: {
        overflowY: 'overlay',
        height: '300px'
    },
    navback: {
        float: 'left'
    }
}

class AddUser extends Component {
    componentWillMount () {
        this.setState({
            // selectfriend, inviteform, finalform, and stepone
            form: 'stepone',
            previousForm: [],
            selectedFriend: {
                plexUserID: '',
                plexUsername: '',
                serverName: '',
                assignedMediaID: '',
                serverMachineID: ''
            },
        })

        let { getFriends } = this.props

        getFriends()
    }
    show (newForm = '') {
        let { form, previousForm } = this.state

        previousForm.push(form)

        this.setState({
            form: newForm,
            previousForm
        })
    }
    handleGoBack () {
        let { previousForm } = this.state
        let form = previousForm.pop()

        if (form === 'stepone') {
            return
        }
        
        this.setState({
            form,
            previousForm
        })
    }
    handleSelectFriend (user) {
        let {
            id,
            username
        } = user

        console.log('selecting ' + username + ' (' + id + ')')

        let {
            selectedFriend
        } = this.state

        selectedFriend.plexUserID = id
        selectedFriend.plexUsername = username

        this.setState({
            selectedFriend
        })

        this.show('finalform')
    }
    handleMediaID (e) {
        let { target, which } = e

        if (which === 13) {
            this.handleAddUser()
            return
        }

        let { selectedFriend } = this.state

        selectedFriend.assignedMediaID = target.value

        this.setState({
            selectedFriend
        })
    }
    handleAddUser () {
        let {
            selectedFriend
        } = this.state

        let {
            plexUserID,
            plexUsername,
            assignedMediaID
        } = selectedFriend

        if (plexUserID === '' || plexUserID === 0 || plexUsername === '' || !assignedMediaID) {
            console.error('missing required friend info')
            return
        }

        let {
            addUser
        } = this.props

        addUser(selectedFriend)
    }
    renderFinalForm () {
        let {
            selectedFriend
        } = this.state

        return (
            <div>
                <pre>{selectedFriend.plexUsername} - (Plex user id: {selectedFriend.plexUserID})</pre>

                <Textfield
                    label="Media id"
                    onKeyUp={e => this.handleMediaID(e)}
                    floatingLabel
                />

                <Button onClick={()=> this.handleAddUser()} >Add user</Button>
            </div>
        )
    }
    renderSelectFriend () {
        let { isFriendListLoading } = this.props

        if (isFriendListLoading) {
            return (
                <div>
                    <h5>Friend List</h5>
                    <Spinner />
                </div>
            )
        }

        return (
            <Grid>
                <Cell col={12}>
                    <h5>Friend List</h5>
                    
                    <List style={styles.friendList}>
                        <ListItem onClick={() => this.handleSelectFriend({ id: '9873', username: 'siirclutch'})} >
                            <ListItemContent>siirclutch</ListItemContent>
                            <ListItemAction><a><Icon name="arrow_forward" /></a></ListItemAction>
                        </ListItem>
                        <ListItem onClick={() => this.handleSelectFriend({ id: '9876', username: 'siirclutch-guest'})} >
                            <ListItemContent>siirclutch-guest</ListItemContent>
                            <ListItemAction><a><Icon  name="arrow_forward" /></a></ListItemAction>
                        </ListItem>
                    </List>
                </Cell>
            </Grid>
        )
    }
    renderInviteForm () {
        return (
            <Grid>
                <Cell col={12}>
                    <Textfield
                        label="Plex Username"
                        floatingLabel
                        disabled
                    />
                </Cell>
                <Cell col={12}>
                    {/*<Checkbox label="Use labels (requires plex pass)" />*/}
                    <Textfield
                        label="Media id"
                        floatingLabel
                        disabled
                    />
                </Cell>
                <Cell col={12}>
                    <Button disabled>Invite</Button>
                </Cell>
            </Grid>
        )
    }
    renderStepOne () {
        return (
            <Grid>
                <Cell col={12}>
                    Does the user already have access to your Plex library?
                </Cell>
                <Cell col={6}>
                    <Button onClick={() => { this.show('selectfriend') }}>Yes</Button>
                </Cell>
                <Cell col={6}>
                    <Button onClick={() => { this.show('inviteform') }}>No, invite them</Button>
                </Cell>
            </Grid>
        )
    }
    render () {
        let { form } = this.state

        return (
            <Grid>
                <Cell col={1}>
                    <a><Icon onClick={() => this.handleGoBack()} name="arrow_back" /></a>
                </Cell>
                <Cell col={11}>
                    {(() => {
                        switch (form) {
                            case 'inviteform':
                                return this.renderInviteForm()
                            case 'selectfriend':
                                return this.renderSelectFriend()
                            case 'finalform':
                                return this.renderFinalForm()
                            default:
                                return this.renderStepOne()
                        }
                    })()}
                </Cell>
            </Grid>
        )
    }
}

AddUser.propTypes = {
    friends: Proptypes.array.isRequired,
    isFriendListLoading: Proptypes.bool.isRequired,
    getFriends: Proptypes.func.isRequired
}

export default AddUser
