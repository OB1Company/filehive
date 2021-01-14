import React from 'react'
import './style/Header.css';
import {Link} from 'react-router-dom';
import {VerifyAuthenticated} from './App'

const IsNotLoggedIn = ({children}) => {
    const token = localStorage.getItem('token');
    if(token == null) {
        return children;
    }
    return (null);
}

function Header() {



  return (
    <div class="Header">
      <div>
        <Link to ='/'><h1>Filehive</h1></Link>
        <input type="text"/>
      </div>
      <div class="Header-Right">
        <IsNotLoggedIn>
          <Link to='/login'>Log in</Link>
          <Link to ='/signup'>Sign up</Link>
        </IsNotLoggedIn>
        <VerifyAuthenticated>
            <Link to='/user/1'>Username</Link>
        </VerifyAuthenticated>
        <Link to ='/create'><input type="button" value="Create dataset"/></Link>
      </div>
    </div>
  )
}

export default Header