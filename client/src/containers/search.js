import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import Search from '../components/search'
import { searchPlexForMedia, getMetadata } from '../actions/search'
import { addUser, getMonitoredUsers } from '../actions/users'


const mapStateToProps = (state) => {
    let {
        search,
        users
    } = state

    let { currentlySelected } = users

    return {
        ...search,
        currentlySelected
    }
}

const mapDispatchToProps = dispatch => ({
    searchPlexForMedia: bindActionCreators(searchPlexForMedia, dispatch),
    getMetadata: bindActionCreators(getMetadata, dispatch),
    getMonitoredUsers: bindActionCreators(getMonitoredUsers, dispatch),
    addUser: bindActionCreators(addUser, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Search)