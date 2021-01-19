import React, {useState, useEffect} from 'react'
import {useLocation} from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import TabbedLinks from "../components/TabbedLinks";
import DataSetsRows from "../components/DataSetsRows";
import axios from "axios";


const getDatasets = async () => {

    const csrftoken = localStorage.getItem('csrf_token');
    const instance = axios.create({
        baseURL: "",
        headers: { "x-csrf-token": csrftoken }
    })

    const loginUrl = "/api/v1/datasets";
    const apiReq = await instance.get(
        loginUrl
    );
    console.log(apiReq);

    return apiReq.data.datasets;
}



export default function HomePage() {

    const [datasets, setDatasets] = useState([]);

    useEffect(() => {
        const fetchData = async() => {
            const ds = await getDatasets();
            setDatasets(ds);
        };
        fetchData();
    }, []);

    const linkNames = [
        { name: 'Trending', link: '/datasets/trending' },
        { name: 'Latest', link: '/datasets/latest' }
    ];

    const location = useLocation();

  return (
    <div className="container">
      <Header/>
      <TabbedLinks linkNames={linkNames} activeLink={location.pathname}/>
      <div className="maincontent margins-30">
        <DataSetsRows sortby="trending" datasets={datasets}/>
      </div>
      <Footer/>
    </div>
  )
}