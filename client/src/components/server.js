import React, { Component } from 'react'
import {
    Button,
    Cell,
    Grid,
    Spinner,
    Textfield
} from 'react-mdl'

const styles = {
    success: {
        color: '#109d36',
        fontSize: '1.14em'
    },
    failed: {
        color: '#e11a1a',
        fontSize: '1.14em'
    },
    saveButton: {
        marginLeft: '15px'
    }
}

export default class Server extends Component {
    componentWillMount () {
        let {
            fetchServerInfo
        } = this.props

        this.setState({
            plexURL: '',
            plexToken: '',
            init: false
        })

        if (fetchServerInfo) {
            fetchServerInfo()
        } else {
            console.error('fetchServerInfo is not a function')
        }
    }
    componentWillReceiveProps (props = {}) {
        let {
            token,
            url
        } = props

        let {
            init,
            plexURL,
            plexToken
        } = this.state

        // prevent re-renders if already initialized
        if (init && plexURL !== '') {
            return
        }

        if (init && plexToken !== '') {
            return
        }

        this.setState({
            plexURL: url,
            plexToken: token,
            init: true
        })
    }
    handlePlexURL (text) {
        // TODO: client-side validation
        this.setState({
            plexURL: text
        })
    }
    handlePlexToken (text) {
        // TODO: client-side validation
        this.setState({
            plexToken: text
        })
    }
    handleTestConnection () {
        let {
            testServer
        } = this.props

        let {
            plexToken,
            plexURL
        } = this.state

        testServer({
            url: plexURL,
            token: plexToken
        })
    }
    handleSave () {
        let {
            plexURL,
            plexToken
        } = this.state

        let {
            saveServerInfo
        } = this.props
        
        saveServerInfo({
            url: plexURL,
            token: plexToken
        })
    }
    renderStatus () {
        let {
            errMessage,
            message
        } = this.props

        if (message && !errMessage) {
            return <p style={styles.success}>{message}</p>
        }

        if (!message && errMessage) {
            return <p style={styles.failed}>{errMessage}</p>
        }
    }
    render () {
        let {
            plexToken,
            plexURL,
        } = this.state

        let {
            isFetching
        } = this.props

        if (isFetching) {
            return <Spinner singleColor />
        }

        return (
            <Grid>
                <Cell col={12}>
                    <div>
                        <h4>Server</h4>
                        {this.renderStatus()}
                    </div>

                    <div>
                        <Textfield
                            onChange={(e) => this.handlePlexURL(e.target.value)}
                            label="plex server url"
                            floatingLabel
                            value={plexURL}
                            />
                    </div>
                    <div>
                        <Textfield
                            onChange={(e) => this.handlePlexToken(e.target.value)}
                            label="plex token"
                            floatingLabel
                            value={plexToken}
                        />
                    </div>

                    <div>
                        <Button
                            raised
                            ripple
                            onClick={() => this.handleTestConnection()}
                            >Test</Button>
                        <Button
                            style={styles.saveButton}
                            raised
                            colored
                            ripple
                            onClick={() => this.handleSave()}
                            >Save</Button>
                    </div>
                </Cell>
            </Grid>
        )
    }
}