import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import Setup from '../components/setup'
import {
    checkPlexPin,
    clearServerSelection,
    fetchPlexPin,
    getPlexServers,
    selectConnection,
    selectServer,
    setPlexServer
} from '../actions/setup'

const mapStateToProps = (state) => {
    let { setup } = state

    return setup
}

const mapDispatchToProps = dispatch => ({
    clearServerSelection: bindActionCreators(clearServerSelection, dispatch),
    checkPlexPin: bindActionCreators(checkPlexPin, dispatch),
    getPlexServers: bindActionCreators(getPlexServers, dispatch),
    fetchPlexPin: bindActionCreators(fetchPlexPin, dispatch),
    selectConnection: bindActionCreators(selectConnection, dispatch),
    selectServer: bindActionCreators(selectServer, dispatch),
    setPlexServer: bindActionCreators(setPlexServer, dispatch),
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Setup)