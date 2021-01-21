import React, {useEffect, useState} from 'react'
import Header from '../Header'
import Footer from '../Footer'
import { useParams } from 'react-router-dom'
import {getAxiosInstance} from "../components/Auth";
import Modal from "react-modal";
import { ModalProvider, ModalConsumer } from '../components/modals/ModalContext';
import ModalRoot from '../components/modals/ModalRoot';
import { Countries } from '../constants/Countries'

const instance = getAxiosInstance();

const DatasetPurchaseModal = (props) => {
    console.log(props);



    return (
        <div className="modal-container">
            <div className="modal-title">Purchase</div>
            <div>You’re almost finished. Please confirm the order details below to purchase the dataset.</div>
            <div className="modal-center-text-bold">Pay {props.price} FIL</div>
            {/*<div className="modal-button-container"><button className="normal-button">Top up wallet</button></div>*/}
            <div className="modal-button-container"><button className="orange-button">Confirm Order</button></div>
            {/*<div className="mini-light-description text-center top-32">You don’t have enough funds in your wallet. Please add at least 5.1834 FIL to your wallet.</div>*/}
            <div className="mini-light-description text-center top-32">The funds will automatically be deducted from your wallet once you proceed.</div>
        </div>
    )
}

const Modal1 = ({ onRequestClose, ...otherProps }) => (
    <Modal isOpen onRequestClose={onRequestClose} className="dataset-purchase-modal" {...otherProps}>
        <DatasetPurchaseModal datasetId={otherProps.datasetId} price={otherProps.price}/>
    </Modal>
);

export default function DatasetPage() {
    const [datasetImageUrl, setDatasetImageUrl] = useState("/api/v1/image/");
    let { id } = useParams();
    const [dataset, setDataset] = useState({});
    const [publisher, setPublisher] = useState({});

    Modal.setAppElement('#root');

    const pullDataset = async (datasetId) => {
        const datasetUrl = "/api/v1/dataset/" + datasetId;
        await instance.get(datasetUrl, {withCredentials: true})
            .then((data)=>{
                const dataset = data.data;
                setDataset(data.data);
                setDatasetImageUrl(datasetImageUrl+dataset.imageFilename);

                const getPublisher = async () => {
                    instance.get("/api/v1/user/" + dataset.userID)
                        .then((publish) => {
                            publish.data.avatarFilename = "/api/v1/image/"+publish.data.Avatar;

                            // Convert country code to name
                            const countryObject = Countries.find(c => c.value === publish.data.Country);
                            publish.data.countryName = countryObject.label;
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
                                    <div className="dataset-publisher-location">{publisher.countryName}</div>
                                </div>
                                <div className="dataset-publisher-avatar"><img src={publisher.avatarFilename} className="dataset-metadata-avatar"/></div>
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
                                        <button className="orange-button" onClick={() => showModal(Modal1, { datasetId: dataset.id, price: dataset.price })}>Buy Now</button>
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