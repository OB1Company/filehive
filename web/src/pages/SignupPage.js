import React from 'react'
import Header from '../Header'
import Footer from '../Footer'
import Signup from "../Signup";

export default function LoginPage() {
  return (
    <div className="container">
      <Header/>
        <div className="subBody">
            <Signup/>
        </div>
      <Footer/>
    </div>
  )
}