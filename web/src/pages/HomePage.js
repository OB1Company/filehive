import React from 'react'
import { Link } from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'

export default function HomePage() {
  return (
    <div className="container">
      <Header/>

      <div class="subBody">

        <Link to='/datasets/trending'>Trending</Link>
        <Link to='/datasets/recent'>Recent</Link>

      </div>

      <Footer/>
    </div>
  )
}