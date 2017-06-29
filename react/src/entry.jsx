import React from 'react'
import ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import { Router, Route, IndexRoute, hashHistory } from 'react-router'
    
import './theme.scss'
import './main.scss'
import App from './container/App.jsx'
import store from './store.js'

ReactDOM.render(
  <div>
    <Provider store={store}>
      <Router history={hashHistory}>
        <Route path="/">
          <IndexRoute component={App}/>
        </Route>
      </Router>
    </Provider>
  </div>,
  document.getElementById('App')
);
