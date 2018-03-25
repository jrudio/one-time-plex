// import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import Server from '../components/server'
import { fetchServerInfo, saveServerInfo, testServer } from '../actions/server'

const mapStateToProps = (state) => {
    let { server } = state

    return server
}

const mapDispatchToProps = dispatch => ({
    fetchServerInfo: bindActionCreators(fetchServerInfo, dispatch),
    saveServerInfo: bindActionCreators(saveServerInfo, dispatch),
    testServer: bindActionCreators(testServer, dispatch)
})

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Server)