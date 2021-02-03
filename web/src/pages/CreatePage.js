import React from 'react'
import { useParams } from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import Create from '../components/Create.js'
import Edit from '../components/Edit.js'

export default function CreatePage() {

    let { id } = useParams();

    return (
    <div className="container">
      <Header/>
      <div className="maincontent">
          {id === undefined
              ? <Create/>
              : <Edit/>
          }
      </div>
      <Footer/>
    </div>
    )
}