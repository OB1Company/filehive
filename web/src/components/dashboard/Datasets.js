import React, {useEffect, useState} from 'react'
import Header from "../../Header";
import TabbedLinks from "../TabbedLinks";
import DataSetsRows from "../DataSetsRows";
import Footer from "../../Footer";
import {Link, useLocation} from "react-router-dom";
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

export default function Datasets() {

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

        <div className="maincontent">
            <div className="dashboard-header">
                <h2>My datasets</h2>
            </div>
            <div>
                { datasets.length == 0 &&
                    <div className="margins-30">
                        <p className="mini-description dashboard-p">You do not have any datasets uploaded yet.</p>
                        <Link to='/create'><input type="button" className="orange-button" value="Create new dataset"/></Link>
                    </div>
                }
                { datasets.length > 0 &&
                <DataSetsRows sortby="trending" datasets={datasets} rowType="edit"/>
                }
            </div>

        </div>



    );
}