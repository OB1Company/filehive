import React, {useEffect, useState} from 'react'
import Header from '../Header'
import Footer from '../Footer'
import {useParams, Link, useHistory} from 'react-router-dom'
import {getAxiosInstance} from "../components/Auth";
import Modal from "react-modal";
import { Countries } from '../constants/Countries'
import {FiatPrice, FilecoinPrice, HumanFileSize} from "../components/utilities/images";
import TimeAgo from 'javascript-time-ago'
import { GetWalletBalance } from "../components/dashboard/Wallet";
import defaultAvatar from '../images/avatar-placeholder.png';
import { decode } from 'html-entities';
import ReactMarkdown from 'react-markdown'
import useSWR from "swr";
import {Helmet} from "react-helmet";
import spinner from "../images/spinner.gif";
import QRCode from 'qrcode.react';

const instance = getAxiosInstance();
const qrSettings = {
    width: 188,
    height: 188
}

export default function DatasetPage() {

    const timeAgo = new TimeAgo('en-US');
    const history = useHistory();

    const [datasetImageUrl, setDatasetImageUrl] = useState("/api/v1/image/");
    let { id } = useParams();
    const [dataset, setDataset] = useState({});
    const [price, setPrice] = useState("");
    const [contentID, setContentID] = useState("");
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
    const [publisherAddress, setPublisherAddress] = useState("");

    Modal.setAppElement('#root');

    const handleCloseModal = () => {
        setOpenModal("purchase");
        setPurchaseOpen(false);
        setSuccessOpen(false);
    }

    const DatasetSuccessModal = (props) => {

        const HandleClickDownload = async (e) => {
            setOpenModal("purchase");
            setModalIsOpen(false);
            history.push('/dashboard/purchases');
        }

        const downloadLink = "/api/v1/download/"+id;

        return (
                <div className="modal-container-success">
                    <div className="modal-title">✅ Success</div>
                    <div>
                        <p>Click the button below to download your purchase.</p>
                    </div>
                    <div className="modal-button-container"><a href={downloadLink} download><button className="orange-button" onClick={HandleClickDownload}>Download</button></a></div>
                    <div className="mini-light-description text-center">
                        <p>If the file has to be retrieved from Filecoin, you will receive an email notifying you once it has be delivered to your account.</p>
                        If download is not available yet, please try again in <Link to="/dashboard/datasets" className="orange-link">your account</Link>.</div>
                </div>
        )
    }

    const DatasetPurchaseModal = (props) => {
        const [isPurchasing, setIsPurchasing] = useState(false);

        const HandleClickPurchase = async () => {
            setIsPurchasing(true);

            // Send payment
            const sendPayment = async () => {
                const updateUserUrl = "/api/v1/purchase/"+id;
                await instance.post(
                    updateUserUrl
                ).then((data)=>{
                    setIsPurchasing(false);
                    console.log(data);
                    props.nextModal();
                }).catch((error)=>{
                    setIsPurchasing(false);
                    console.log(error);
                })
            }
            await sendPayment();

        }

        const handleTopUp = () => {
            history.push("/dashboard/wallet");
        }

        const showWarning = (price > balance);

        const PurchaseButton = () => {
            if (!isPurchasing) {
                return  <button className="orange-button" onClick={HandleClickPurchase}>Confirm Order</button>
            } else {
                return <span className="spinner-span">
                    <img src={spinner} width="20" height="20" alt="spinner" className="noblock"/> Purchasing dataset...
                </span>
            }
        }

        return (
            <div className="modal-container">
                <div className="modal-title">Purchase</div>
                <div>You’re almost finished. Please confirm the order details below to purchase the dataset.</div>
                <div className="modal-center-text-bold">Pay {price} FIL<br/>
                    <span className="mini-light-description">Wallet balance: {balance} FIL</span>
                </div>
                {showWarning && <div>
                    <div className="modal-button-container"><button className="normal-button" onClick={handleTopUp}>Top up wallet</button></div>
                    <div className="mini-light-description text-center top-32">You don’t have enough funds in your wallet.
                        Please add at least {price} FIL to your wallet.</div>
                </div>
                }
                {!showWarning &&
                    <div>
                        <div className="modal-button-container">
                            <PurchaseButton/>
                        </div>
                        <div className="mini-light-description text-center top-32">The funds will automatically be deducted from
                        your wallet once you proceed.</div>
                    </div>
                }
            </div>

        )
    }

    const NextModal = (e)=>{
        setPurchaseOpen(false);
        setSuccessOpen(true);
    };

    const Modal1 = ({ onRequestClose, ...otherProps }) => (
        <Modal  isOpen={purchaseOpen} onRequestClose={handleCloseModal} nextModal={NextModal} className="dataset-purchase-modal" styles={{ modal: {}, overlay: { background: "rgba(156, 156, 156, 0.75)" } }}  {...otherProps}>
            <DatasetPurchaseModal datasetId={otherProps.datasetId} nextModal={NextModal} price={otherProps.price}/>
        </Modal>
    );

    const Modal2 = ({ onRequestClose, ...otherProps }) => (
        <Modal  isOpen={successOpen} onRequestClose={handleCloseModal} className="dataset-download-modal" styles={{ modal: {}, overlay: { background: "rgba(156, 156, 156, 0.75)" } }}  {...otherProps}>
            <DatasetSuccessModal datasetId={otherProps.datasetId} price={otherProps.price}/>
        </Modal>
    );

    const filecoinPrice  = useSWR('filecoinPrice', FilecoinPrice);
    const fiatPrice = FiatPrice(dataset.price, filecoinPrice.data);

    useEffect(() => {
        const pullDataset = async (datasetId) => {
            const datasetUrl = "/api/v1/dataset/" + datasetId;
            await instance.get(datasetUrl, {withCredentials: true})
                .then((data)=>{
                    const dataset = data.data;
                    setDataset(data.data);
                    setContentID(data.data.contentID);
                    setDatasetImageUrl(datasetImageUrl+dataset.imageFilename);

                    setPrice(Number.parseFloat(dataset.price).toFixed(8).toString().replace(/\.?0+$/,""));

                    setFileType(dataset.fileType);
                    setFileSize(HumanFileSize(dataset.fileSize, true));
                    setUsername(dataset.username);
                    setTimestamp(timeAgo.format(Date.parse(dataset.createdAt)));

                    const getPublisher = async () => {
                        instance.get("/api/v1/user/" + dataset.userID)
                            .then((publish) => {
                                const avatar = publish.data.avatar;
                                publish.data.avatarFilename = (avatar === "") ? defaultAvatar : "/api/v1/image/"+publish.data.avatar;

                                // Convert country code to name
                                const countryObject = Countries.find(c => c.value === publish.data.country);
                                publish.data.countryName = countryObject.label;
                                setPublisher(publish.data);
                                setPublisherAddress(publish.data.filecoinAddress);

                                const checkPublisher = (publish.data.email === localStorage.getItem("email")) ? "orange-button-disable" : "";
                                setDisableBuy(checkPublisher);
                            })

                    }
                    getPublisher();

                    const checkIfPurchased = async () => {
                        instance.get("/api/v1/purchased/" + dataset.id)
                            .then((result) => {
                                setDisableBuy("orange-button-disable");
                            })

                    }
                    checkIfPurchased();

                    const grabBalance = async () => {
                        GetWalletBalance().then((balance)=>{setBalance(balance)});
                    }
                    grabBalance();


                })
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
                <Helmet>
                    <title>Filehive | {decode(dataset.title)}</title>
                </Helmet>

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
                                    <div className="dataset-publisher-name">{publisher.name}</div>
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
                                {dataset.price > 0 &&
                                    <div>
                                        <div className="dataset-metadata-price">{dataset.price} FIL <span
                                            className="tiny-price">({fiatPrice})</span></div>
                                        <div className="dataset-metadata-description">Your payment helps support the dataset creator and Filecoin miners.</div>
                                    </div>
                                }

                                <div className="dataset-metadata-button">
                                    {dataset.price > 0 &&
                                    <button className={"orange-button raise " + disableBuy}
                                            onClick={HandleBuyButton}>Buy Now</button>
                                    }
                                    {dataset.price == 0 &&
                                        <div>
                                            <a href={"/api/v1/download/"+dataset.id} download><button className="orange-button raise" >Download</button></a>
                                            <a href={"https://gateway.ipfs.io/ipfs/"+contentID} target="_blank" className="orange-link"><button className="normal-button raise">View on IPFS Gateway</button></a>
                                        </div>
                                    }
                                </div>

                                <div className="donate-container">
                                    <h3>Donate</h3>
                                    <div className="wd-description">Send FIL to the address below to donate to this dataset creator.</div>
                                    <div className="qr-code-deposit" >
                                        <QRCode value={publisherAddress} size="99" imageSettings={qrSettings} />
                                    </div>
                                    <div className="center">{publisherAddress}</div>
                                    <div className="center"><a className="orange-link" onClick={() =>  navigator.clipboard.writeText(publisherAddress)}>copy</a></div>
                                </div>

                            </div>
                        </div>
                    </div>
                </div>
                <Footer/>
            </div>
    )
}