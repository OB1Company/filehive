import React, {useEffect, useState} from 'react'
import {getAxiosInstance} from "../Auth";
import DataSetsRow from "../DataSetsRow";
import {useHistory} from "react-router-dom";
import TimeAgo from "javascript-time-ago";
import {FiatPrice, FilecoinPrice, HumanFileSize} from "../utilities/images";
import {decode} from "html-entities";
import useSWR from "swr";


function PurchasesRows(props) {

    const [filecoinPrice, setFilecoinPrice] = useState("");
    const retrievePrice = async ()=> {
        const price = await FilecoinPrice();
        setFilecoinPrice(price);
    }
    retrievePrice();

    const reversePurchases = props.purchases.reverse();

    let rows = reversePurchases.map((purchase)=> {
        return <PurchaseRow key={purchase.id} metadata={purchase} filecoinPrice={filecoinPrice}/>;
    });

    return (
        <div className="datasets-rows">
            {rows}
        </div>
    )
}


function PurchaseRow(props) {
    const history = useHistory();

    const timeAgo = new TimeAgo('en-US')

    const imageFilename = props.metadata.imageFilename;
    const title = props.metadata.title;
    const shortDescription = props.metadata.shortDescription;
    const price = Number.parseFloat(props.metadata.price).toFixed(8).toString().replace(/\.?0+$/,"");
    const fiatPrice = FiatPrice(props.metadata.price, props.filecoinPrice);
    const fileType = props.metadata.fileType;
    const username = props.metadata.username;
    const timestamp = timeAgo.format(Date.parse(props.metadata.timestamp));

    const buttonText = "Download";
    const gotoPage = '/dataset/'+props.metadata.datasetID;

    const datasetImage = "/api/v1/image/" + imageFilename;

    const handleClickDatasetRow = (e) => {
        history.push(gotoPage);
    }

    const datasetUrl = "/api/v1/download/"+props.metadata.datasetID;

    return (
        <div className="datasets-row">
            <div className="datasets-row-image">
                <img className="datasets-image" src={datasetImage} alt={decode(title)}/>
            </div>
            <div className="datasets-row-info">
                <div className="mini-bold-title" onClick={handleClickDatasetRow}>{decode(title)}</div>
                <div className="mini-description">{decode(shortDescription)}</div>
                <div className="mini-light-description tag-container">
                    <div>{fileType}</div>
                    <div>{timestamp}</div>
                    <div>{username}</div>
                </div>
            </div>
            <div className="datasets-details">
                <div><a href={datasetUrl} download><button className="normal-button">{buttonText}</button></a></div>
                <div className="small-orange-text dataset-row-price">{price} FIL</div>
                <div className="mini-light-description">{fiatPrice}</div>
            </div>
        </div>
    )
}


export default function Purchases() {

    const [purchases, setPurchases] = useState([]);

    useEffect(() => {
        const getPurchases = async () => {
            const instance = getAxiosInstance();
            const apiUrl = "/api/v1/purchases";
            const res = await instance.get(
                apiUrl
            );
            setPurchases(res.data.purchases);

        }
        getPurchases();
    }, []);

    return (

        <div className="maincontent">
            <h2 className="margins-30">Purchases</h2>

            <div className="">
                { purchases.length === 0 &&
                <div className="margins-30">
                    <p className="mini-description dashboard-p">You have not made any purchases yet. Check out some of our <a href="/datasets/trending" className="orange-link">trending datasets</a></p>
                </div>
                }
                { purchases.length > 0 &&
                <PurchasesRows purchases={purchases}/>
                }
            </div>


        </div>

    );
}