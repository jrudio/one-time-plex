import { SEARCH_PLEX } from '../constants/search'

export const searchPlexForMedia = (title = '') => {
    return dispatch => {
        let mrRobotEpisodes = [
            {
                type: 'episode',
                year: '2015',
                title: 'hacknation.mov',
                mediaID: '4913'
            }
        ]

        let mrRobotSeasons = [
            {
                type: 'season',
                year: '2015',
                title: 'Season 1',
                mediaID: '4920'
            },
            {
                type: 'season',
                year: '2016',
                title: 'Season 2',
                mediaID: '4921'
            },
            {
                type: 'season',
                year: '2017',
                title: 'Season 3',
                mediaID: '4922'
            }
        ]

        let results = [
            { type: 'movie', year: '2015', title: 'John Wick', mediaID: '2912' },
            { type: 'movie', year: '2016', title: 'John Wick 2', mediaID: '3912' },
            { type: 'series', year: '2015', title: 'Mr: Robot', mediaID: '4912', episodes: mrRobotEpisodes }
        ]
        
        return dispatch({
            type: SEARCH_PLEX,
            results: mrRobotSeasons
        })
    }
}