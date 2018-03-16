import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import Search from '../components/search'
import { searchPlexForMedia, getMetadata } from '../actions/search'
import { addUser } from '../actions/users'

const SearchContainer = props => {
    return <Search {...props} />
}

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
    addUser: bindActionCreators(addUser, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SearchContainer)