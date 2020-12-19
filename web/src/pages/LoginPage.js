import React from 'react'
import Header from '../Header'
import Footer from '../Footer'
import Login from '../Login'

export default function LoginPage() {
  return (
    <div className="container">
      <Header/>
      <div className="subBody">
        <Login/>
      </div>
      <Footer/>
    </div>
  )
}