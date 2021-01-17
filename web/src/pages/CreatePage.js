import React from 'react'
import Header from '../Header'
import Footer from '../Footer'
import Create from '../components/Create.js'

export default function LoginPage() {
  return (
    <div className="container">
      <Header/>
      <div className="maincontent">
        <Create/>
      </div>
      <Footer/>
    </div>
  )
}