import React from 'react'
import './style/Header.css';
import {Link, useHistory} from 'react-router-dom';
import { useCookies } from "react-cookie";
import {getAxiosInstance} from "./components/Auth";

function Header() {

    const history = useHistory();
    const email = localStorage.getItem('email');
    const name = localStorage.getItem('name');
    const loggedIn = (!(email == null || email === ""));
    const [token, getToken, removeToken] = useCookies(['token']);

    const HandleLogout = (e) => {

        localStorage.removeItem("email");
        localStorage.removeItem("name");

        const instance = getAxiosInstance();
        instance.post('/api/v1/logout')
            .then((result) => {
                console.log(result);
            })

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
          { loggedIn ? <Link to='/dashboard'>{name}</Link> : ""}
          { loggedIn ? <Link onClick={HandleLogout}>Log out</Link> : ""}

        <Link to ='/create'><input type="button" value="Create dataset"/></Link>
      </div>
    </div>
  )
}

export default Header