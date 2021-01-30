import React, {useState, useEffect} from 'react'
import {useLocation} from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import TabbedLinks from "../components/TabbedLinks";
import DataSetsRows from "../components/DataSetsRows";
import axios from "axios";


const getDatasets = async (tabName) => {

    const csrftoken = localStorage.getItem('csrf_token');
    const instance = axios.create({
        baseURL: "",
        headers: { "x-csrf-token": csrftoken }
    })

    const loginUrl = "/api/v1/"+tabName;
    console.log(loginUrl);
    const apiReq = await instance.get(
        loginUrl
    );

    const datasets = (apiReq.data.hasOwnProperty("datasets")) ? apiReq.data.datasets : [];

    return datasets;
}



export default function HomePage() {

    const [datasets, setDatasets] = useState([]);

    const location = useLocation();
    const tabName = location.pathname.substring(location.pathname.lastIndexOf('/') + 1);
    console.log(tabName);

    useEffect(() => {
        const fetchData = async() => {
            const ds = await getDatasets(tabName);
            console.log(tabName);
            setDatasets(ds);
        };
        fetchData();
    }, [tabName]);

    const linkNames = [
        { name: 'Trending', link: '/datasets/trending' },
        { name: 'Latest', link: '/datasets/latest' }
    ];

    return (
    <div className="container">
      <Header/>
      <TabbedLinks linkNames={linkNames} activeLink={location.pathname}/>
      <div className="maincontent">
        <DataSetsRows sortby="trending" datasets={datasets}/>
      </div>
      <Footer/>
    </div>
  )
}