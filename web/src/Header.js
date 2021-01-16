import React from 'react'
import './style/Header.css';
import {Link, useHistory} from 'react-router-dom';
import { useCookies } from "react-cookie";

function Header() {

    const history = useHistory();
    const username = localStorage.getItem('username');
    const loggedIn = (!(username == null || username === ""));
    const [token, getToken, removeToken] = useCookies(['token']);

    const HandleLogout = (e) => {

        localStorage.removeItem("username");
        localStorage.removeItem("email");
        removeToken("token");

        history.push("/login");
    }

  return (
    <div class="Header">
      <div>
        <Link to ='/'><h1>Filehive</h1></Link>
        <input type="text"/>
      </div>
      <div class="Header-Right">
          { !loggedIn ? <Link to='/login'>Log in</Link> : ""}
          { !loggedIn ? <Link to='/signup'>Sign up</Link> : ""}
          { loggedIn ? <Link to='/dashboard'>{username}</Link> : ""}
          { loggedIn ? <Link onClick={HandleLogout}>Log out</Link> : ""}

        <Link to ='/create'><input type="button" value="Create dataset"/></Link>
      </div>
    </div>
  )
}

export default Header