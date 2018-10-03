import React, { Component } from 'react'
import logo from './logo.svg'
import AddUser from './containers/adduser'
import Server from './containers/server'
import Setup from './containers/setup'
import AssignedMedia from './containers/assignedmedia'
import { Grid, Button, Cell } from 'react-mdl'
import './App.css'

window.otp = {
  url: '//localhost:6969/api'
}

class App extends Component {
  componentWillMount () {
    // page can be one of:
    // - home
    // - settings
    // - setup
    this.setState({
      // page: 'setup'
      page: 'home'
    })
  }
  showPage(page = '') {
    if (page === '') {
      console.log('showPage() please be explicit when calling showPage')
      page = 'home'
    }

    this.setState({
      page
    })
  }
  renderSetup () {
    return (
      <Grid>
        <Cell col={12}>
          <Button onClick={() => this.showPage('home')}>Go Home</Button>
        </Cell>

        <Setup />
      </Grid>
    )
  }
  renderSettings () {
    return (
      <Grid>
        <Cell col={12}>
          <Button onClick={() => this.showPage('home')}>Go Home</Button>
        </Cell>

        <Server />
      </Grid>
    )
  }
  renderUnknownPage () {
    return <h3>Page not defined</h3>
  }
  renderHome () {
    return (
      <Grid>
        <Cell col={12}>
          <Button onClick={() => this.showPage('setup')}>Server Settings</Button>
        </Cell>
        <Cell col={4}>
          <h4>Plex Friends:</h4>

          <AddUser />
        </Cell>

        <Cell col={4}>
          <h4>Assigned Media</h4>

          <AssignedMedia />
        </Cell>
      </Grid>
    )
  }
  render() {
    let {
      page
    } = this.state

    return (
      <div className="App">
        <div className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h2>Welcome to One Time Plex</h2>
        </div>

        <div style={{ width: '80%', margin: 'auto' }}>
          {(() => {
            switch (page) {
              case 'home':
                return this.renderHome()
              case 'settings':
                return this.renderSettings()
              case 'setup':
                return this.renderSetup()
              default:
                return this.renderUnknownPage()
            }
          })()}
        </div>
      </div>
    )
  }
}

export default App
