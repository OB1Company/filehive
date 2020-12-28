import React from 'react'
import { Link } from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import TabbedLinks from "../components/TabbedLinks";
import DataSetsRows from "../components/DataSetsRows";

export default function HomePage() {

    const linkNames = [
        { name: 'Trending', link: '/datasets/trending' },
        { name: 'Latest', link: '/datasets/latest' }
    ];

  return (
    <div className="container">
      <Header/>
      <TabbedLinks linkNames={linkNames} />
      <DataSetsRows sortby="trending"/>
      <Footer/>
    </div>
  )
}