import React from 'react'
import './style/Header.css';
import logo from './images/filecoin-icon-tiny.png';

function Footer() {
  return (
    <div class="Footer">
        <img src={logo} alt="Filecoin" width="18" height="18" />
        <a href='https://filecoin.io'>Powered by <strong>Filecoin</strong></a>
    </div>
  )
}

export default Footer