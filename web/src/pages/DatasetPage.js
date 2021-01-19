import React, {useEffect, useState} from 'react'
import Header from '../Header'
import Footer from '../Footer'
import { useParams } from 'react-router-dom'
import {getAxiosInstance} from "../components/Auth";

export default function DatasetPage() {
    const [datasetImageUrl, setDatasetImageUrl] = useState("/api/v1/image/");
    let { id } = useParams();
    const [dataset, setDataset] = useState({});
    const [publisher, setPublisher] = useState({});

    const pullDataset = async (datasetId) => {
        const instance = getAxiosInstance();
        const datasetUrl = "/api/v1/dataset/" + datasetId;
        await instance.get(datasetUrl, {withCredentials: true})
            .then((data)=>{
                const dataset = data.data;
                setDataset(data.data);
                setDatasetImageUrl(datasetImageUrl+dataset.imageFilename);

                const getPublisher = async () => {
                    instance.get("/api/v1/user/" + dataset.userID)
                        .then((publish) => {
                            setPublisher(publish.data);
                            console.log(publisher);
                        })

                }
                getPublisher();

            })
            //     // setDataset(data.data);
            //     // setDatasetImageUrl(datasetImageUrl+data.data.imageFilename);
            // })
            // .catch((err) => {
            //     //console.error(err);
            // })
    }

    useEffect(() => {
        const fetchData = async() => {
            const ds = await pullDataset(id);
        };
        fetchData();
    }, []);

    return (
        <div className="container">
            <Header/>
            <div className="maincontent">
                <div className="dataset-container-header">
                    <div>
                        <div className="dataset-header">
                            <div><h2>{dataset.title}</h2></div>
                            <div className="dataset-description">{dataset.shortDescription}</div>
                            <div className="mini-light-description">{dataset.fileType} 2.4GB 7d ago</div>
                        </div>
                    </div>
                    <div>
                        <div className="dataset-publisher">
                            <div>
                                <div className="dataset-publisher-name">{publisher.Name}</div>
                                <div className="dataset-publisher-location">{publisher.Country}</div>
                            </div>
                            <div className="dataset-publisher-avatar"><img src={publisher.Avatar}/></div>
                        </div>
                    </div>
                </div>
                <div className="dataset-container-body">
                    <div>
                        <div>
                            <img src={datasetImageUrl} alt="" className="dataset-hero-image"/>
                        </div>
                        <div className="dataset-maintext">{dataset.fullDescription}</div>
                    </div>
                    <div>
                        test
                    </div>
                </div>
            </div>
            <Footer/>
        </div>
    )
}