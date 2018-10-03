import React, {
    Component
} from 'react'
import {
    Button,
    Cell,
    Grid,
    Icon,
    Spinner,
    Textfield
} from 'react-mdl'
// import { getPlexServers } from '../actions/setup';

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

export default class Setup extends Component {
    constructor () {
        super()

        this.hasCalledGetPlexServers = false
    }
    stopPlexPinCheck () {
        console.log('stopping plex pin check')
        clearInterval(this.pinTimer)
    }
    componentWillReceiveProps (newProps) {
        const { getPlexServers } = this.props
        const {
            isAuthorized,
            selectedServer
        } = newProps

        if (!!isAuthorized && selectedServer.name === '') {
            if (!this.hasCalledGetPlexServers) {
                getPlexServers()
                this.hasCalledGetPlexServers = true
            }

            this.setState({
                screen: 'serverList'
            })
        } else if (!!isAuthorized && selectedServer.name !== '') {
            this.setState({
                screen: 'localOrRemoteIP'
            })
        }
    }
    componentWillMount () {
        const {
            fetchPlexPin
        } = this.props

        fetchPlexPin()

        this.pinTimer = setInterval(() => {
            const {
                checkPlexPin,
                isAuthorized
            } = this.props

            if (isAuthorized) {
                console.log('we are authorized')
                this.stopPlexPinCheck()
                return
            }

            // check if we have a plex token
            checkPlexPin()
        }, 3 * 1000)

        // available screens:
        // - plexPin
        // - serverList
        // - localOrRemoteIP
        // - final (show settings)
        this.setState({
            screen: 'plexPin',
        })
    }
    componentWillUnmount () {
        this.stopPlexPinCheck()
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
            servers,
            selectedServer,
            selectServer
        } = this.props


        if (index === null || index < 0 || index > servers.length - 1) {
            console.error('handleServerSelect() invalid server selection')
            return
        }

        let newSelectedServer = Object.assign({}, selectedServer, servers[index])

        // this.setState({ selectedServer: newSelectedServer })
        selectServer(newSelectedServer)

        this.showScreen('localOrRemoteIP')
    }
    handleIpSelect (index) {
        const {
            selectConnection,
            selectedServer,
            setPlexServer
        } = this.props

        if (index === null || index < 0 || index > selectedServer.connection.length - 1) {
            console.error('handleIpSelect() invalid ip selection')
            return
        }

        console.log(selectedServer)

        // let newSelectedServer = Object.assign({}, selectedServer, {
        //     url: selectedServer.connection[index].url
        // })

        // this.setState({
        //     selectedServer: newSelectedServer
        // })
        selectConnection(index)

        if (selectedServer.name === '' || selectedServer.connection.length < 1) {
            console.error('selected server is invalid; rejecting setPlexServer()')
            console.log(selectedServer)
            return
        }

        setPlexServer(selectedServer, index)

        this.showScreen('final')
    }
    handleEditSettings () {
        const { clearServerSelection } = this.props

        clearServerSelection()

        this.setState({
            screen: 'serverList'
        })
    }
    renderPlexPin() {
        const {
            errMessage,
            isFetching,
            isExpired,
            pin
        } = this.props

        console.log(this.props)

        // if (isFetching) {
        //     return (
        //         <div>
        //             <h5>Fetching a Plex code...</h5>
        //             <Spinner />
        //         </div>
        //     )
        // }

        return (
            <div>
                <h5>Here is your Plex PIN: {(() => isExpired && <a onClick={() => console.log('refresh')}><Icon name="refresh" /></a>)()}</h5>
                <p>{errMessage && errMessage}</p>
                {(() => {
                    if (isFetching) {
                        return <pre>{pin || 'N/A'}</pre>
                    }
                    
                    return <Spinner />
                })()}
                
                <p>Go to <a href="https://plex.tv/link" target="_blank">plex.tv/link</a> to grant One Time Plex access to your servers.</p>
            </div>
        )
    }
    renderServerList() {
        const { servers } = this.props

        // console.log(this.props.servers)

        return (
            <div>
                <h5>Please choose the plex server to monitor:</h5>
                <ul style={styles.list}>
                    {(() => {
                        if (servers.length === 0) {
                            return (
                                <div>
                                    Fetching available servers...
                                    <Spinner />
                                </div>
                            )
                        }
                        
                        let serverNodes = []

                        servers.length > 0 && servers.forEach((server, i) => {
                            serverNodes.push(<li key={i}><Button colored onClick={() => this.handleServerSelect(i)}> {server.name} </Button></li>)
                        })

                        return serverNodes
                    })()}
                </ul>
            </div>
        )
    }
    renderLocalOrRemote () {
        const { selectedServer } = this.props

        return (
            <div>
                <h5>Choose a remote or local url:</h5>
                <ul style={styles.list}>
                    {(() => {
                        let ipNodes = []

                        selectedServer && selectedServer.connection.forEach((ip, i) => {
                            let addressType = 'Remote'

                            if (ip.local === 1) {
                                addressType = 'Local'
                            }

                            ipNodes.push(<li key={i}>{addressType} - {ip.uri}<Button colored onClick={() => this.handleIpSelect(i)}>Select</Button></li>)
                        })
                        
                        return ipNodes
                    })()}
                </ul>
            </div>
        )
    }
    renderFinal () {
        const {
            connectionIndex,
            selectedServer
        } = this.props

        return (
            <div>
                <h5>Here are your settings:</h5>

                <div>Server Name: <Textfield label='Server Name' readOnly value={selectedServer.name} /></div>
                <div>Plex Server URL: <Textfield label='Plex Server URL' readOnly value={selectedServer.connection[connectionIndex].uri} /></div>
                <div>Plex Token: <Textfield label='Plex Token' readOnly value={selectedServer.token} /></div>

                <a onClick={() => this.handleEditSettings()} href="#">Edit Settings</a>
            </div>
        )
    }
    render () {
        const { screen } = this.state

        console.log(this.props)
        console.log('screen', screen)

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