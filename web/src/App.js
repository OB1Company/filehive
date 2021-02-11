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
import ConfirmPage from './pages/ConfirmPage'
import PasswordResetPage from './pages/PasswordResetPage'
import ChangePasswordPage from './pages/ChangePasswordPage'
import AdminPage from "./pages/AdminPage";
import axios from "axios";


const PrivateRoute = ({ component: Component, ...rest }) => (
    <Route {...rest} render={(props) => (
        localStorage.getItem("email") == null || localStorage.getItem("name") === ""
            ? <Redirect to='/login' />
            : <Component {...props} />
    )
    } />
)

const AdminRoute = ({ component: Component, ...rest }) => (
    <Route {...rest} render={(props) => (
        localStorage.getItem("email") != null
            && localStorage.getItem("name") !== ""
            && localStorage.getItem("admin")
            ? <Component {...props} />
            : <Redirect to='/dashboard' />
    )
    } />
)

export default function App() {

    const checkCsrfToken = async () => {
        let token = localStorage.getItem("csrf_token");
        await axios.get('/api/v1/user', {withCredentials: true})
            .then((data)=>{
                console.debug(data.headers);
                localStorage.setItem("csrf_token", data.headers['x-csrf-token']);
            })
            .catch((err) => {
                console.debug('GET User call failed', err.response);

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
          <Route path="/confirm_email" component={ConfirmPage} />
          <Route path="/password_reset" component={PasswordResetPage} />
          <Route path="/change_password" component={ChangePasswordPage} />

          <PrivateRoute path="/create" component={CreatePage} />

          <PrivateRoute exact path="/dashboard">
              <Redirect to="/dashboard/datasets"/>
          </PrivateRoute>
          <PrivateRoute exact path="/dashboard/datasets" component={DashboardPage} />
          <PrivateRoute path="/dashboard/purchases" component={DashboardPage} />
          <PrivateRoute path="/dashboard/sales" component={DashboardPage} />
          <PrivateRoute path="/dashboard/wallet" component={DashboardPage} />
          <PrivateRoute path="/dashboard/settings" component={DashboardPage} />
          <PrivateRoute exact path="/dashboard/datasets/:id">
              <CreatePage/>
          </PrivateRoute>

          <AdminRoute exact path="/admin">
              <Redirect to="/admin/users"/>
          </AdminRoute>
          <AdminRoute exact path="/admin/users" component={AdminPage} />
          <AdminRoute path="/admin/datasets" component={AdminPage} />
          <AdminRoute path="/admin/sales" component={AdminPage} />

          <Route component={HomePage} />
      </Switch>
  )
}