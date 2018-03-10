import React, { Component } from 'react'
import logo from './logo.svg'
// import Users from './containers/users'
import AddUser from './containers/adduser'
// import Search from './containers/search'
import AssignedMedia from './containers/assignedmedia'
import { Grid, Cell } from 'react-mdl'
import './App.css'

window.otp = {
  url: '//localhost:6969/api'
  // url: 'http://localhost:6969/api'
}

class App extends Component {
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Welcome to One Time Plex</h2>
        </div>

        <div style={{ width: '80%', margin: 'auto'}}>
          <Grid>
            <Cell col={4}>
              <h4>Plex Friends:</h4>
              
              <AddUser />
            </Cell>

            <Cell col={4}>
              {/* <h4>Search Plex:</h4>

              <Search /> */}
              <h4>Assigned Media</h4>

              <AssignedMedia />
            </Cell>
          </Grid>
        </div>
      </div>
    )
  }
}

export default App
