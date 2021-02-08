import React from 'react'
import { useParams } from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import Create from '../components/Create.js'
import Edit from '../components/Edit.js'
import {Helmet} from "react-helmet";

export default function CreatePage() {

    let { id } = useParams();

    return (
    <div className="container">
        <Helmet>
            <title>Filehive | Create Dataset</title>
        </Helmet>
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