import React, { Component } from 'react'
import {
    Grid,
    Cell,
    // Button,
    // Checkbox,
    Icon,
    List,
    ListItem,
    ListItemContent,
    ListItemAction,
    ProgressBar,
    Textfield
} from 'react-mdl'
// import Proptypes from 'prop-types'

const styles = {
    errMsg: {
        color: 'red',
        fontSize: '1.2em'
    },
    navback: {
        float: 'left'
    }
}

class Search extends Component {
    handleAddMedia (mediaID) {
        console.log('Adding:', mediaID)
    }
    handleGetEpisodes (seriesID) {
        console.log('Getting more episodes of:', seriesID)
    }
    handleSearch (e) {
        let { which } = e
        let { searchPlexForMedia } = this.props

        // user pressed enter
        if (which === 13) {
            searchPlexForMedia('test')
        }
    }
    handleSelectResult (e, mediaType, mediaID) {
        if (mediaType === 'series') {
            this.handleGetSeasons(mediaID)
            return
        } else if (mediaType === 'season') {
            this.handleGetEpisodes(mediaID)
            return
        }

        // it's either a movie or episode
        this.handleAddMedia(mediaID)
    }
    renderSearchResults (results) {
        return <List>  
            {results && results.map((r, i) => {
                {/* default movie or episode */}
                let iconType = 'add'

                if (r.type === 'series' || r.type === 'season') {
                    iconType = 'arrow_right'
                }
                
                return (
                    <ListItem key={i}>
                        <ListItemContent>{r.title} ({r.year})</ListItemContent>
                        <ListItemAction>
                            <a><Icon name={iconType} onClick={(e) => this.handleSelectResult(e, r.type, r.mediaID)} /></a>
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
        
        return (
            <Grid>
                <Cell col={12}>
                    {/* Search input */}
                    <Textfield
                        label="Search..."
                        autoFocus
                        onKeyUp={e => this.handleSearch(e)}
                    />
                    {results && results.length > 0 && (
                        <div>
                            <a><Icon name='arrow_back' style={styles.navback} /></a>
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

export default Search