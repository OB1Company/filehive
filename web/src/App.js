import React, { useEffect } from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'

import HomePage from './pages/HomePage'
import UserPage from './pages/UserPage'
import LoginPage from './pages/LoginPage'
import SignupPage from './pages/SignupPage'
import CreatePage from './pages/CreatePage'
import axios from "axios";

export default function App() {

    const getCsrfToken = async () => {
        try {
            const token = localStorage.getItem("token");
            console.log(token);
            if(token == null) {
                const {data} = await axios.get('/api/v1/user', {withCredentials: true});
            }
        } catch(err) {
            console.log(err.response);
            localStorage.setItem("token", err.response.headers['x-csrf-token']);
            axios.defaults.headers.post['x-csrf-token'] = err.response.headers['x-csrf-token'];
            console.log(localStorage);
        }
    };
    getCsrfToken();

  return (
      <Switch>
          <Route exact path="/">
              {true ? <Redirect to="/datasets/trending" /> : <HomePage />}
          </Route>
          <Route path="/login" component={LoginPage} />
          <Route path="/signup" component={SignupPage} />
          <Route path="/datasets/trending" component={HomePage} />
          <Route path="/datasets/latest" component={HomePage} />
          <Route path="/user/:id" component={UserPage} />
          {/*<VerifyAuthenticated>*/}
              <Route path="/create" component={CreatePage} />
          {/*</VerifyAuthenticated>*/}
      </Switch>
  )
}