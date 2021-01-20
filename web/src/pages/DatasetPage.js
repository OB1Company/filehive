import React, {useEffect, useState} from 'react'
import Header from '../Header'
import Footer from '../Footer'
import { useParams } from 'react-router-dom'
import {getAxiosInstance} from "../components/Auth";
import Modal from "react-modal";
import { ModalProvider, ModalConsumer } from '../components/modals/ModalContext';
import ModalRoot from '../components/modals/ModalRoot';


const DatasetPurchaseModal = (props) => {
    console.log(props.datasetId);
    return (
        <div>
            <h2>Purchase</h2>
            You’re almost finished. Please confirm the order details below to purchase the dataset.
            {props.datasetId}
        </div>
    )
}

const Modal1 = ({ onRequestClose, ...otherProps }) => (
    <Modal isOpen onRequestClose={onRequestClose} className="dataset-purchase-modal" {...otherProps}>
        <DatasetPurchaseModal datasetId={otherProps.datasetId}/>
    </Modal>
);

export default function DatasetPage() {
    const [datasetImageUrl, setDatasetImageUrl] = useState("/api/v1/image/");
    let { id } = useParams();
    const [dataset, setDataset] = useState({});
    const [publisher, setPublisher] = useState({});

    Modal.setAppElement('#root');

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
        <ModalProvider>
            <ModalRoot />
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
                            <div className="dataset-metadata-container">
                                <div className="dataset-metadata-price">{dataset.price} FIL</div>
                                <div className="dataset-metadata-description">Your payment helps support the dataset creator and Filecoin miners.</div>
                                <div className="dataset-metadata-button">
                                <ModalConsumer>
                                    {({ showModal }) => (
                                        <button className="orange-button" onClick={() => showModal(Modal1, { datasetId: dataset.id })}>Buy Now</button>
                                    )}
                                </ModalConsumer>
                                </div>
                                <div className="dataset-metadata-warning">The price includes the miner fee.</div>
                            </div>
                        </div>
                    </div>
                </div>
                <Footer/>
            </div>

        </ModalProvider>
    )
}