import React, { Component } from 'react'
import {
    Button,
    Grid,
    Cell,
    Icon,
    List,
    ListItem,
    ListItemContent,
    // ListItemAction,
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
            selectedMedia: {
                id: '',
                title: '',
                year: ''
            }
        })
    }
    handleAddMedia (title, id, year) {
        let { selectedMedia } = this.state

        let newSelectedMedia = Object.assign({}, selectedMedia, {
            id,
            title,
            year
        })

        this.setState({
            selectedMedia: newSelectedMedia
        })
    }
    handleAddUser() {
        if (!this.props || !this.state) {
            console.error()
            return
        }

        let { selectedMedia } = this.state
        let { addUser, currentlySelected } = this.props

        console.log('assigning media to user ' + currentlySelected + ' ' + selectedMedia.id + ' ' + selectedMedia.title)

        addUser({
            mediaID: selectedMedia.id,
            plexUserID: currentlySelected
        })
    }
    handleNo () {
        console.log('no')

        this.setState({
            selectedMedia: {
                id: '',
                title: '',
                year: ''
            }
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
    handleSelectResult (e, mediaType, mediaID, title, year) {
        // show tv show's seasons or a season's episodes
        if (mediaType === 'show' || mediaType === 'season') {
            this.handleGetMetadata(mediaID)
            return
        }

        // it's either a movie or episode
        this.handleAddMedia(title, mediaID, year)
    }
    renderAddSuccessful () {
        return <div><Icon name="check" className="text-green" /></div>
    }
    renderSearchResults (results) {
        let { selectedMedia } = this.state
        let { id } = selectedMedia

        return <List style={styles.searchResults} >
            {results && results.map((r, i) => {
                let highlighted = {}

                if (id === r.mediaID) {
                    highlighted = styles.selected
                }
                
                /* default movie or episode */
                // let iconType = 'add'

                // if (r.type === 'show' || r.type === 'season') {
                //     iconType = 'arrow_right'
                // }
                
                return (
                    <ListItem key={i}>
                        <ListItemContent style={highlighted} onClick={(e) => this.handleSelectResult(e, r.type, r.mediaID, r.title, r.year)} >{r.title} {r.year && '(' + r.year + ')' }</ListItemContent>
                        {/* <ListItemAction>
                            <a><Icon name={iconType} onClick={(e) => this.handleSelectResult(e, r.type, r.mediaID, r.title, r.year)} /></a>
                        </ListItemAction> */}
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
            selectedMedia
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
                </Cell>
                <Cell col={12}>
                    {/* Search results */}
                    {(() => {
                        if (selectedMedia.id !== '') {
                            return (
                                <div>
                                    <p>Are you sure you want to assign {selectedMedia.title} to this user?</p> 
                                    <Button onClick={() => this.handleAddUser()}>yes</Button>
                                    <Button onClick={() => this.handleNo()}>no</Button>
                                </div>
                            )
                        }


                    })()}
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