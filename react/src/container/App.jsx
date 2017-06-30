import React from 'react'
import { connect } from 'react-redux'
import { hashHistory } from 'react-router'

import './main.scss'
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