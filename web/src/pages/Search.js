import React, {useState, useEffect} from 'react'
import {useLocation} from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import DataSetsRows from "../components/DataSetsRows";
import axios from "axios";
import Pluralize from 'react-pluralize'


const useQuery = () => {
    return new URLSearchParams(useLocation().search);
}

const getDatasets = async (query) => {

    const csrftoken = localStorage.getItem('csrf_token');
    const instance = axios.create({
        baseURL: "",
        headers: { "x-csrf-token": csrftoken }
    })

    const loginUrl = "/api/v1/search?query="+query;
    const apiReq = await instance.get(
        loginUrl
    );
    console.debug(apiReq);

    const datasets = (apiReq.data.hasOwnProperty("datasets")) ? apiReq.data.datasets : [];

    return datasets;
}



export default function SearchPage() {

    const query = useQuery();

    const [datasets, setDatasets] = useState([]);
    const [searchQuery] = useState(query.get("q"));

    useEffect(() => {
        const fetchData = async() => {
            const ds = await getDatasets(searchQuery);
            setDatasets(ds);
        };
        fetchData();
    }, []);

  return (
    <div className="container">
      <Header/>
        <div className="search-results-header"><strong>{datasets.length}</strong> <Pluralize singular={'result'} showCount={false} count={datasets.length} /> matching "{searchQuery}"</div>
      <div className="maincontent margins-30">
        <DataSetsRows sortby="trending" datasets={datasets}/>
      </div>
      <Footer/>
    </div>
  )
}