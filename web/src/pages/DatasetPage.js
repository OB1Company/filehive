import React, {useEffect, useState} from 'react'
import Header from '../Header'
import Footer from '../Footer'
import {useParams, Link, useHistory} from 'react-router-dom'
import {getAxiosInstance} from "../components/Auth";
import Modal from "react-modal";
import { Countries } from '../constants/Countries'
import {FiatPrice, HumanFileSize} from "../components/utilities/images";
import TimeAgo from 'javascript-time-ago'
import { GetWalletBalance } from "../components/dashboard/Wallet";
import defaultAvatar from '../images/avatar-placeholder.png';
import { decode } from 'html-entities';
import ReactMarkdown from 'react-markdown'

const instance = getAxiosInstance();

export default function DatasetPage() {

    const timeAgo = new TimeAgo('en-US');
    const history = useHistory();

    const [datasetImageUrl, setDatasetImageUrl] = useState("/api/v1/image/");
    let { id } = useParams();
    const [dataset, setDataset] = useState({});
    const [price, setPrice] = useState("");
    const [fiatPrice, setFiatPrice] = useState("");
    const [fileType, setFileType] = useState("");
    const [fileSize, setFileSize] = useState("");
    const [username, setUsername] = useState("");
    const [timestamp, setTimestamp] = useState("");
    const [publisher, setPublisher] = useState({});
    const [balance, setBalance] = useState("");
    const [modalIsOpen, setModalIsOpen] = useState(false);
    const [purchaseOpen, setPurchaseOpen] = useState(false);
    const [successOpen, setSuccessOpen] = useState(false);
    const [openModal, setOpenModal] = useState("purchase");
    const [disableBuy, setDisableBuy] = useState("orange-button-disable");

    Modal.setAppElement('#root');

    const handleCloseModal = () => {
        setOpenModal("purchase");
        setPurchaseOpen(false);
        setSuccessOpen(false);
    }

    const DatasetSuccessModal = (props) => {

        const HandleClickDownload = async (e) => {

            // 1. Try to retrieve via IPFS
            const getDatasetFile = async () => {
                const datasetUrl = "/api/v1/download/"+id;
                window.open(datasetUrl);
            }
            await getDatasetFile();

            // 2. Warn the user that it's being retrieved via Filecoin
            // 3. Server will attempt to retrieve it from Filecoin

            setOpenModal("purchase");
            setModalIsOpen(false);
        }

        return (
                <div className="modal-container-success">
                    <div className="modal-title">✅ Success</div>
                    <div>It may take a few hours to retrieve your dataset from the Filecoin network.</div>
                    <div className="modal-button-container"><button className="orange-button" onClick={HandleClickDownload}>Download</button></div>
                    <div className="mini-light-description text-center top-32">If download is not available yet, please try again in <Link to="/dashboard/datasets">your account</Link>.</div>
                </div>
        )
    }

    const DatasetPurchaseModal = (props) => {
        const HandleClickPurchase = async () => {

            // Send payment
            const sendPayment = async () => {
                const updateUserUrl = "/api/v1/purchase/"+id;
                await instance.post(
                    updateUserUrl
                ).then((data)=>{
                    console.log(data);
                }).catch((error)=>{
                    console.log(error);
                })
            }
            await sendPayment();

            // nextModal();
            console.log(props);
            props.nextModal();

        }

        const handleTopUp = () => {
            history.push("/dashboard/wallet");
        }

        const showWarning = (price > balance);

        return (
            <div className="modal-container">
                <div className="modal-title">Purchase</div>
                <div>You’re almost finished. Please confirm the order details below to purchase the dataset.</div>
                <div className="modal-center-text-bold">Pay {price} FIL</div>
                {showWarning && <div>
                    <div className="modal-button-container"><button className="normal-button" onClick={handleTopUp}>Top up wallet</button></div>
                    <div className="mini-light-description text-center top-32">You don’t have enough funds in your wallet.
                        Please add at least {price} FIL to your wallet.</div>
                </div>
                }
                {!showWarning &&
                    <div>
                        <div className="modal-button-container">
                            <button className="orange-button" onClick={HandleClickPurchase}>Confirm Order</button>
                        </div>
                        <div className="mini-light-description text-center top-32">The funds will automatically be deducted from
                        your wallet once you proceed.</div>
                    </div>
                }
            </div>

        )
    }

    const NextModal = (e)=>{
        console.log('Next Modal', this);
        setPurchaseOpen(false);
        setSuccessOpen(true);
    };

    const Modal1 = ({ onRequestClose, ...otherProps }) => (
        <Modal shouldCloseOnOverlayClick="true" isOpen={purchaseOpen} onRequestClose={handleCloseModal} nextModal={NextModal} className="dataset-purchase-modal" {...otherProps}>
            <DatasetPurchaseModal datasetId={otherProps.datasetId} nextModal={NextModal} price={otherProps.price}/>
        </Modal>
    );

    const Modal2 = ({ onRequestClose, ...otherProps }) => (
        <Modal shouldCloseOnOverlayClick="true" isOpen={successOpen} onRequestClose={handleCloseModal} className="dataset-purchase-modal" {...otherProps}>
            <DatasetSuccessModal datasetId={otherProps.datasetId} price={otherProps.price}/>
        </Modal>
    );

    useEffect(() => {
        const pullDataset = async (datasetId) => {
            const datasetUrl = "/api/v1/dataset/" + datasetId;
            await instance.get(datasetUrl, {withCredentials: true})
                .then((data)=>{
                    const dataset = data.data;
                    setDataset(data.data);
                    setDatasetImageUrl(datasetImageUrl+dataset.imageFilename);

                    setPrice(Number.parseFloat(dataset.price).toFixed(8).toString().replace(/\.?0+$/,""));
                    const getFiatPrice = async ()=>{
                        console.log(dataset.price);
                        setFiatPrice(await FiatPrice(dataset.price));
                    }
                    getFiatPrice();

                    setFileType(dataset.fileType);
                    setFileSize(HumanFileSize(dataset.fileSize, true));
                    setUsername(dataset.username);
                    setTimestamp(timeAgo.format(Date.parse(dataset.createdAt)));

                    const getPublisher = async () => {
                        instance.get("/api/v1/user/" + dataset.userID)
                            .then((publish) => {
                                const avatar = publish.data.Avatar;
                                publish.data.avatarFilename = (avatar === "") ? defaultAvatar : "/api/v1/image/"+publish.data.Avatar;

                                // Convert country code to name
                                const countryObject = Countries.find(c => c.value === publish.data.Country);
                                publish.data.countryName = countryObject.label;
                                setPublisher(publish.data);

                                const checkPublisher = (publish.data.Email === localStorage.getItem("email")) ? "orange-button-disable" : "";
                                setDisableBuy(checkPublisher);
                            })

                    }
                    getPublisher();

                    const grabBalance = async () => {
                        GetWalletBalance().then((balance)=>{setBalance(balance)});
                    }
                    grabBalance();


                })
            //     // setDataset(data.data);
            //     // setDatasetImageUrl(datasetImageUrl+data.data.imageFilename);
            // })
            // .catch((err) => {
            //     //console.error(err);
            // })
        }

        const fetchData = async() => {
            const ds = await pullDataset(id);
        };
        fetchData();
    }, []);

    const HandleBuyButton = () => {
        if(disableBuy === "") {
            setPurchaseOpen(true);
        }
    }

    return (
            <div className="container">
                <Modal1 onRequestClose={handleCloseModal}/>
                <Modal2 onRequestClose={handleCloseModal}/>
                <Header/>
                <div className="maincontent">
                    <div className="dataset-container-header">
                        <div>
                            <div className="dataset-header">
                                <div><h2>{decode(dataset.title)}</h2></div>
                                <div className="dataset-description">{decode(dataset.shortDescription)}</div>
                                <div className="mini-light-description tag-container">
                                    <div>{fileType}</div>
                                    <div>{fileSize}</div>
                                    <div>{timestamp}</div>
                                    <div>{username}</div>
                                </div>
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
                            <div className="dataset-maintext"><ReactMarkdown>{decode(dataset.fullDescription)}</ReactMarkdown></div>
                        </div>
                        <div>
                            <div className="dataset-metadata-container">
                                <div className="dataset-metadata-price">{dataset.price} FIL <span className="tiny-price">({fiatPrice})</span></div>
                                <div className="dataset-metadata-description">Your payment helps support the dataset creator and Filecoin miners.</div>
                                <div className="dataset-metadata-button">
                                        <button className={"orange-button "+disableBuy} onClick={HandleBuyButton}>Buy Now</button>
                                </div>
                                <div className="dataset-metadata-warning">The price includes the miner fee.</div>
                            </div>
                        </div>
                    </div>
                </div>
                <Footer/>
            </div>
    )
}