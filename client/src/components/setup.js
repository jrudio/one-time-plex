import React, {
    Component
} from 'react'
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
    selectButton: {
        marginLeft: '15px'
    },
    list: {
        listStyle: 'none'
    }
}

const availableServers = [
    {
        name: 'Server2',
        connection: [
            {
                local: true,
                url: 'https://192168120-local.plex.local/'
            },
            {
                local: false,
                url: 'https://1065311-remote.plex.tv/'
            }
        ]
    },
    {
        name: 'Justin\'s Server',
        connection: [
            {
                local: true,
                url: 'https://19216643120-local.plex.local/',
            },
            {
                local: false,
                url: 'https://15555311-remote.plex.tv/'
            }
        ]
    },
    {
        name: 'Justin\'s x1',
        connection: [
            {
                local: true,
                url: 'https://19436564620-local.plex.local/'
            },
            {
                local: false,
                url: 'https://8888311-remote.plex.tv/'
            }
        ]
    },
    {
        name: 'Random Server',
        connection: [
            {
                local: true,
                url: 'https://1946664360-local.plex.local/'
            },
            {
                local: false,
                url: 'https://103492841-remote.plex.tv/'
            }
        ]
    }
]

export default class Setup extends Component {
    componentWillMount () {
        const {
            fetchPlexPin
        } = this.props

        fetchPlexPin()
        // check if we have a plex token

        // if not just display the plex pin screen
        // and wait for 

        // 
        
        // available screens:
        // - plexPin
        // - serverList
        // - localOrRemoteIP
        // - final (show settings)
        this.setState({
            availableServers,
            isAuthorized: false,
            selectedServer: {
                connection: [
                    {
                        local: false,
                        url: ''
                    }
                ],
                name: '',
                url: '',
                token: 'xxxxxabc123'
            },
            screen: 'plexPin',
        })
    }
    showScreen (screen = '') {
        if (!screen) {
            console.error('showScreen() requires an argument')
            return
        }

        this.setState({
            screen
        })
    }
    handleServerSelect (index) {
        const {
            availableServers,
            selectedServer
        } = this.state

        if (index === null || index < 0 || index > availableServers.length - 1) {
            console.error('handleServerSelect() invalid server selection')
            return
        }

        let newSelectedServer = Object.assign({}, selectedServer, availableServers[index])


        this.setState({ selectedServer: newSelectedServer })

        this.showScreen('localOrRemoteIP')
    }
    handleIpSelect (index) {
        const { selectedServer } = this.state

        if (index === null || index < 0 || index > selectedServer.connection.length - 1) {
            console.error('handleIpSelect() invalid ip selection')
            return
        }

        console.log(selectedServer)

        console.log(index)

        let newSelectedServer = Object.assign({}, selectedServer, {
            url: selectedServer.connection[index].url
        })

        this.setState({
            selectedServer: newSelectedServer
        })

        this.showScreen('final')
    }
    renderPlexPin() {
        const {
            isFetching,
            pin
        } = this.props

        console.log(this.props)

        if (isFetching) {
            return (
                <div>
                    <h5>Fetching a Plex code...</h5>
                    <Spinner />
                </div>
            )
        }

        return (
            <div>
                <h5>Here is your Plex PIN: <pre>{pin || 'N/A'}</pre></h5>
                <p>Go to <a href="https://plex.tv/link" target="_blank">plex.tv/link</a> to grant One Time Plex access to your servers.</p>
            </div>
        )
    }
    renderServerList() {
        const { availableServers } = this.state

        return (
            <div>
                <h5>Please choose the plex server to monitor:</h5>
                <ul style={styles.list}>
                    {(() => {
                        let serverNodes = []

                        availableServers.forEach((server, i) => {
                            serverNodes.push(<li key={i}><Button colored onClick={() => this.handleServerSelect(i)}> {server.name} </Button></li>)
                        })

                        return serverNodes
                    })()}
                </ul>
            </div>
        )
    }
    renderLocalOrRemote () {
        const { selectedServer } = this.state
        
        return (
            <div>
                <h5>Choose a remote or local url:</h5>
                <ul style={styles.list}>
                    {(() => {
                        let ipNodes = []

                        selectedServer.connection.forEach((ip, i) => {
                            if (ip.local) {
                                ipNodes.push(<li key={i}><Button colored onClick={() => this.handleIpSelect(i)}>Local - {ip.url}</Button></li>)
                                return
                            }

                            ipNodes.push(<li key={i}><Button colored onClick={() => this.handleIpSelect(i)}>Remote - {ip.url}</Button></li>)
                        })
                        
                        return ipNodes
                    })()}
                </ul>
            </div>
        )
    }
    renderFinal () {
        const { selectedServer } = this.state

        return (
            <div>
                <h5>Here are your settings:</h5>

                <div>Server Name: <Textfield label='Server Name' readOnly value={selectedServer.name} /></div>
                <div>Plex Server URL: <Textfield label='Plex Server URL' readOnly value={selectedServer.url} /></div>
                <div>Plex Token: <Textfield label='Plex Token' readOnly value={'xxxxxxabc123'} /></div>

                <a onClick={() => this.showScreen('serverList')} href="#">Edit Settings</a>
            </div>
        )
    }
    render () {
        const { screen } = this.state

        return (
            <Grid>
                <Cell col={12}>
                    {(() => {
                        switch (screen) {
                            case 'plexPin':
                                return this.renderPlexPin()
                            case 'serverList':
                                return this.renderServerList()
                            case 'localOrRemoteIP':
                                return this.renderLocalOrRemote()
                            case 'final':
                                return this.renderFinal()
                            default:
                                return this.renderPlexPin()
                        }
                    })()}
                </Cell>
            </Grid>
        )
    }
}