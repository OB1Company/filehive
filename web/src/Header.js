import React, {useEffect, useState} from 'react'
import './style/Header.css';
import {Link, useHistory} from 'react-router-dom';
import {getAxiosInstance} from "./components/Auth";
import defaultAvatar from './images/avatar-placeholder.png';

function Header() {

    const history = useHistory();
    const email = localStorage.getItem('email');
    const name = localStorage.getItem('name');
    const admin = localStorage.getItem('admin');
    const loggedIn = (!(email == null || email === ""));
    const [avatar, setAvatar] = useState("");

    const HandleLogout = () => {

        localStorage.removeItem("email");
        localStorage.removeItem("name");
        localStorage.removeItem("admin");

        const instance = getAxiosInstance();
        instance.post('/api/v1/logout')
            .then((result) => {
                history.push("/login");
            })
            .catch((err) => {
                console.error(err);
                localStorage.removeItem("_gorilla_csrf");
                history.push("/login");
            })
    }

    const HandleSearchSubmit = (e)=>{
        e.preventDefault();
        if(e.target.q.value === "") {
            return false;
        }
        e.target.submit();
    }

    useEffect(() => {
        if(localStorage.getItem("email")) {
            const instance = getAxiosInstance();
            instance.get("/api/v1/user")
                .then((data) => {
                    if (data.data.Avatar !== "") {
                        setAvatar("/api/v1/image/" + data.data.Avatar);
                    } else {
                        setAvatar(defaultAvatar);
                    }
                })
                .catch((error) => {
                    console.log(error.data);
                })
        }

    });

  return (
    <div className="Header">
      <div>
        <Link to='/'><h1>Filehive</h1></Link>
        <form action="/search" className="filehive-search-form" onSubmit={HandleSearchSubmit}>
            <input type="text" name="q" placeholder="Search Filehive"/>
        </form>
      </div>
      <div className="Header-Right">
          { !loggedIn ? <Link to='/login'>Log in</Link> : ""}
          { !loggedIn ? <Link to='/signup'>Sign up</Link> : ""}
          { loggedIn ? <img src={avatar} className="header-avatar"/> : ""}
          { loggedIn ? <Link to='/dashboard'>{name}</Link> : ""}
          { admin == "true" ? <a href="/admin">Admin</a> : ""}
          { loggedIn ? <Link onClick={HandleLogout}>Log out</Link> : ""}

        <Link to='/create'><input type="button" value="Create dataset" className="raise"/></Link>
      </div>
    </div>
  )
}

export default Header
