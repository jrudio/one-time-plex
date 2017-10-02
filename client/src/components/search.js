import React, { Component } from 'react'
import {
    Grid,
    Cell,
    Icon,
    List,
    ListItem,
    ListItemContent,
    ListItemAction,
    ProgressBar,
    Textfield
} from 'react-mdl'
import Proptypes from 'prop-types'

const styles = {
    errMsg: {
        color: 'red',
        fontSize: '1.2em'
    },
    navback: {
        float: 'left'
    },
    searchResults: {
        overflowY: 'overlay',
        height: '300px'
    },
    selected: {
        border: '#ff4081 0.2em solid'
    }
}

class Search extends Component {
    componentWillMount () {
        this.setState({
            selectedMediaTitle: '',
            selectedMediaID: ''
        })
    }
    handleAddMedia (title, mediaID) {
        this.setState({
            selectedMediaTitle: title,
            selectedMediaID: mediaID
        })
    }
    handleGetMetadata (mediaID) {
        if (mediaID === '') {
            console.error('media id is required to fetch metadata')
            return
        }

        let { getMetadata } = this.props

        getMetadata(mediaID)
    }
    handleSearch (e) {
        let { target, which } = e
        let { searchPlexForMedia } = this.props
        let searchText = target.value

        // user pressed enter
        if (which === 13 && searchText !== '') {
            searchPlexForMedia(searchText)
        }
    }
    handleSelectResult (e, mediaType, mediaID, title) {
        if (mediaType === 'show' || mediaType === 'season') {
            this.handleGetMetadata(mediaID)
            return
        }

        // it's either a movie or episode
        this.handleAddMedia(title, mediaID)
    }
    renderSearchResults (results) {
        let { selectedMediaID } = this.state

        return <List style={styles.searchResults} >  
            {results && results.map((r, i) => {
                let highlighted = {}

                if (selectedMediaID === r.mediaID) {
                    highlighted = styles.selected
                }
                
                /* default movie or episode */
                let iconType = 'add'

                if (r.type === 'show' || r.type === 'season') {
                    iconType = 'arrow_right'
                }
                
                return (
                    <ListItem key={i}>
                        <ListItemContent style={highlighted} onClick={(e) => this.handleSelectResult(e, r.type, r.mediaID, r.title)} >{r.title} {r.year && '(' + r.year + ')' }</ListItemContent>
                        <ListItemAction>
                            <a><Icon name={iconType} onClick={(e) => this.handleSelectResult(e, r.type, r.mediaID, r.title)} /></a>
                        </ListItemAction>
                    </ListItem>
                )
            })}
        </List>
    }
    render () {
        let {
            results,
            isSearching,
            errorMsg
        } = this.props

        let {
            selectedMediaTitle,
            selectedMediaID
        } = this.state
        
        return (
            <Grid>
                <Cell col={12}>
                    {/* Search input */}
                    <Textfield
                        label="Search..."
                        autoFocus
                        onKeyUp={e => this.handleSearch(e)}
                    />
                    {/*{results && results.length > 0 && (
                        <div>
                            <a><Icon name='arrow_back' style={styles.navback} /></a>
                        </div>
                    )}*/}
                    {selectedMediaID && (
                        <div>
                            <p>Media id for <i>{selectedMediaTitle}</i> is:</p>
                            <pre>{selectedMediaID}</pre>
                        </div>
                    )}
                </Cell>
                <Cell col={12}>
                    {/* Search results */}
                    {isSearching && (<ProgressBar indeterminate />)}
                    {!isSearching && !errorMsg && this.renderSearchResults(results)}
                    {!isSearching && errorMsg && (<p style={styles.errMsg}>{errorMsg}</p>)}
                </Cell>
            </Grid>
        )
    }
}   

Search.propTypes = {
    results: Proptypes.array.isRequired,
    searchPlexForMedia: Proptypes.func.isRequired,
    getMetadata: Proptypes.func.isRequired
}

export default Search