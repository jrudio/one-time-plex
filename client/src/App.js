import React, { Component } from 'react'
import logo from './logo.svg'
// import Users from './containers/users'
import AddUser from './components/adduser'
import Search from './containers/search'
import { Grid, Cell } from 'react-mdl'
import './App.css'

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
              <Search />
              {/*Restricted Users:
              <Users />*/}
            </Cell>

            <Cell col={8}>
              <h4>Add User:</h4>
              
              <AddUser />
            </Cell>
          </Grid>
        </div>
      </div>
    )
  }
}

export default App
