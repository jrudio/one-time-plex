import React, { Component } from 'react'
import logo from './logo.svg'
import Users from './containers/users'
import './App.css'

class App extends Component {
  render() {
    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Welcome to React</h2>
        </div>

        Restricted Users:
        <Users />
      </div>
    )
  }
}

export default App
