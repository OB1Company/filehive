import React, { useEffect } from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'

import HomePage from './pages/HomePage'
import UserPage from './pages/UserPage'
import LoginPage from './pages/LoginPage'
import SignupPage from './pages/SignupPage'
import CreatePage from './pages/CreatePage'
import DashboardPage from './pages/DashboardPage'
import axios from "axios";

const token = localStorage.getItem("username");

const VerifyCSRF = ({children}) => {
    const getCsrfToken = async () => {
        try {
            const token = localStorage.getItem("csrf_token");
            if(token == null) {
                const {data} = await axios.get('/api/v1/user', {withCredentials: true});
            }
        } catch(err) {
            console.log(err.response);
            localStorage.setItem("csrf_token", err.response.headers['x-csrf-token']);
            axios.defaults.headers.post['x-csrf-token'] = err.response.headers['x-csrf-token'];
        }
    };
    getCsrfToken();

    return children;
}

export const VerifyAuthenticated = ({children}) => {
    const token = localStorage.getItem("username");
    if(token == null || token === "") {
        return (
            <Redirect to="/login" />
        )
    }
    return children;
}

const PrivateRoute = ({ component: Component, ...rest }) => (
    <Route {...rest} render={(props) => (
        localStorage.getItem("username") == null || localStorage.getItem("username") === ""
            ? <Redirect to='/login' />
            : <Component {...props} />
    )
    } />
)

export default function App() {
  return (
      <VerifyCSRF>
          <Switch>
              <Route exact path="/">
                  {true ? <Redirect to="/datasets/trending" /> : <HomePage />}
              </Route>
              <Route exact path="/login" component={LoginPage} />
              <Route exact path="/signup" component={SignupPage} />
              <Route path="/datasets/trending" component={HomePage} />
              <Route path="/datasets/latest" component={HomePage} />
              <Route path="/user/:id" component={UserPage} />

                  <PrivateRoute path="/create" component={CreatePage} />
                  <Route exact path="/dashboard">
                      <Redirect to="/dashboard/datasets"/>
                  </Route>
                  <Route path="/dashboard/datasets" component={DashboardPage} />
                  <Route path="/dashboard/purchases" component={DashboardPage} />
                  <Route path="/dashboard/wallet" component={DashboardPage} />
                  <Route path="/dashboard/settings" component={DashboardPage} />

              <Route component={HomePage} />
          </Switch>
      </VerifyCSRF>
  )
}