import React from 'react'
import './style/Header.css';
import { Link } from 'react-router-dom';

function Header() {
  return (
    <div class="Header">
      <div>
        <Link to ='/'><h1>Filehive</h1></Link>
        <input type="text"/>
      </div>
      <div class="Header-Right">
        <Link to ='/login'>Log in</Link>
        <Link to ='/signup'>Sign up</Link>
        <Link to ='/create'><input type="button" value="Create dataset"/></Link>
      </div>
    </div>
  )
}

export default Header