import React, { useEffect } from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'

import HomePage from './pages/HomePage'
import UserPage from './pages/UserPage'
import LoginPage from './pages/LoginPage'
import Logout from './Logout'
import SignupPage from './pages/SignupPage'
import CreatePage from './pages/CreatePage'
import DashboardPage from './pages/DashboardPage'
import DatasetPage from './pages/DatasetPage'
import SearchPage from './pages/Search'
import axios from "axios";


const PrivateRoute = ({ component: Component, ...rest }) => (
    <Route {...rest} render={(props) => (
        localStorage.getItem("email") == null || localStorage.getItem("username") === ""
            ? <Redirect to='/login' />
            : <Component {...props} />
    )
    } />
)

export default function App() {

    const checkCsrfToken = async () => {
        let token = localStorage.getItem("csrf_token");
        await axios.get('/api/v1/user', {withCredentials: true})
            .then((data)=>{

            })
            .catch((err) => {
                console.log('GET User call failed', err.response);

                if(err.response.status === 400) {
                    localStorage.removeItem("email");
                    localStorage.removeItem("name");
                }

                token = err.response.headers['x-csrf-token'];

                localStorage.setItem("csrf_token", token);
                axios.defaults.headers.post['x-csrf-token'] = token;
            });
    };
    checkCsrfToken();

    return (
      <Switch>
          <Route exact path="/">
              {true ? <Redirect to="/datasets/trending" /> : <HomePage />}
          </Route>
          <Route exact path="/login" component={LoginPage} />
          <Route exact path="/logout" component={Logout} />
          <Route exact path="/signup" component={SignupPage} />
          <Route path="/datasets/trending" component={HomePage} />
          <Route path="/datasets/latest" component={HomePage} />
          <Route path="/dataset/:id" component={DatasetPage} />
          <Route path="/user/:id" component={UserPage} />
          <Route path="/search" component={SearchPage} />

              <PrivateRoute path="/create" component={CreatePage} />
              <PrivateRoute exact path="/dashboard">
                  <Redirect to="/dashboard/datasets"/>
              </PrivateRoute>
              <PrivateRoute path="/dashboard/datasets/:id" component={CreatePage} />
              <PrivateRoute path="/dashboard/datasets" component={CreatePage} />
              <PrivateRoute path="/dashboard/purchases" component={DashboardPage} />
              <PrivateRoute path="/dashboard/wallet" component={DashboardPage} />
              <PrivateRoute path="/dashboard/settings" component={DashboardPage} />

          <Route component={HomePage} />
      </Switch>
  )
}