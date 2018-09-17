import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import Setup from '../components/setup'
import { fetchPlexPin } from '../actions/setup'

const mapStateToProps = (state) => {
    let { setup } = state

    return setup
}

const mapDispatchToProps = dispatch => ({
    fetchPlexPin: bindActionCreators(fetchPlexPin, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Setup)