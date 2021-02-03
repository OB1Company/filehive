import React from 'react'
import './style/Header.css';
import {Link, useHistory} from 'react-router-dom';
import {getAxiosInstance} from "./components/Auth";

function Header() {

    const history = useHistory();
    const email = localStorage.getItem('email');
    const name = localStorage.getItem('name');
    const loggedIn = (!(email == null || email === ""));

    const HandleLogout = () => {

        localStorage.removeItem("email");
        localStorage.removeItem("name");

        const instance = getAxiosInstance();
        instance.post('/api/v1/logout')
            .then((result) => {
                console.log(result);
                history.push("/login");
            })
            .catch((err) => {
                console.error(err);
                localStorage.removeItem("_gorilla_csrf");
                history.push("/login");
            })
    }

    const HandleSearchSubmit = (e)=>{
        console.log(e.target.q.value);
        e.preventDefault();
        if(e.target.q.value === "") {
            return false;
        }
        e.target.submit();
    }

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
          { loggedIn ? <Link to='/dashboard'>{name}</Link> : ""}
          { loggedIn ? <Link onClick={HandleLogout}>Log out</Link> : ""}

        <Link to='/create'><input type="button" value="Create dataset"/></Link>
      </div>
    </div>
  )
}

export default Header