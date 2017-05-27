import React, { Component } from 'react'
import { Grid, Cell, Button, Checkbox } from 'react-mdl'
// import Proptypes from 'prop-types'

class User extends Component {
    componentWillMount () {
        this.userInput = null
        this.setState({
            form: 'stepone'
        })
    }
    componentDidMount () {
        // this.userInput.focus()
    }
    showInviteForm () {
        this.setState({
            form: 'inviteform'
        })
    }
    showFinalForm () {
        this.setState({
            form: 'finalform'
        })
    }
    renderFinalForm () {
        return (
            <Grid>
                <Cell col={12}>
                    {/* Plex Username - optional */}
                    <label htmlFor="plexUsername">Plex Username:</label>
                    {' '}
                    <input
                        name="plexUsername"
                        type="text"
                        id="plexUsername"
                        ref={input => this.userInput = input}
                    />
                </Cell>
                <Cell col={12}>
                    {/* Plex User ID - optional */}
                    <label htmlFor="plexUserID">Plex User ID:</label>
                    {' '}
                    <input
                        name="plexUserID"
                        type="text"
                        id="plexUserID"
                    />
                </Cell>
                <Cell col={12}>
                    {/* Media ID - required */}
                    <label htmlFor="mediaID">Media ID:</label>
                    {' '}
                    <input
                        name="mediaID"
                        type="text"
                        id="mediaID"
                    />
                </Cell>
            </Grid>
        )
    }
    renderInviteForm () {
        return (
            <Grid>
                <Cell col={12}>
                    <label htmlFor="plexUsername">Plex Username:</label>
                    {' '}
                    <input
                        name="plexUsername"
                        type="text"
                        id="plexUsername"
                    />
                </Cell>
                <Cell col={12}>
                    <Checkbox label="Use labels (requires plex pass)" />
                </Cell>
                <Cell col={12}>
                    <Button>Invite</Button>
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
                    <Button onClick={() => { this.showFinalForm() }}>Yes</Button>
                </Cell>
                <Cell col={6}>
                    <Button onClick={() => { this.showInviteForm() }}>No, let's invite them</Button>
                </Cell>
            </Grid>
        )
    }
    render () {
        let { form } = this.state
        
        switch (form) {
            case 'inviteform':
                return this.renderInviteForm()
            default:
                return this.renderStepOne()
        }
    }
}

export default User
