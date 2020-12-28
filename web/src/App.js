import React from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'

import HomePage from './pages/HomePage'
import UserPage from './pages/UserPage'
import LoginPage from './pages/LoginPage'
import SignupPage from './pages/SignupPage'
import CreatePage from './pages/CreatePage'

export default function App() {
  return (
    <Switch>
      <Route exact path="/">
          {true ? <Redirect to="/datasets/trending" /> : <HomePage />}
      </Route>
      <Route path="/user/:id" component={UserPage} />
      <Route path="/login" component={LoginPage} />
      <Route path="/signup" component={SignupPage} />
      <Route path="/create" component={CreatePage} />
      <Route path="/datasets/trending" component={HomePage} />
      <Route path="/datasets/latest" component={HomePage} />
    </Switch>
  )
}