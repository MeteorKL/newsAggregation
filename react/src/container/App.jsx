import React from 'react'
import { connect } from 'react-redux'
import { hashHistory } from 'react-router'

import './main.scss'
import GameTab from './GameTab/GameTab'
import UserTab from './UserTab/UserTab'
import ws from '../utils/websocket'
import Game from './GameTab/Game.jsx'
import SettingTab from './SettingTab/SettingTab.jsx'
import Header from './Header.jsx'
import Body from './Body.jsx'

class App extends React.Component {

  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div >
        <Header/>
        <Body/>
      </div>
    )
  }
}

function mapStateToProps(state) {
    return {
      user: state.user,
      tab: state.tab,
    }
}

export default connect(mapStateToProps, null)(App);