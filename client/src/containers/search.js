import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import Search from '../components/search'
import { searchPlexForMedia, getMetadata } from '../actions/search'

const SearchContainer = props => {
    return <Search {...props} />
}

const mapStateToProps = (state) => {
    let {
        search,
    } = state

    console.log(search)

    return search
}

const mapDispatchToProps = dispatch => ({
    searchPlexForMedia: bindActionCreators(searchPlexForMedia, dispatch),
    getMetadata: bindActionCreators(getMetadata, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SearchContainer)