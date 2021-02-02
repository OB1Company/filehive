import React, {useEffect, useState} from 'react'
import {getAxiosInstance} from "../Auth";
import DataSetsRow from "../DataSetsRow";
import {useHistory} from "react-router-dom";
import TimeAgo from "javascript-time-ago";
import {FiatPrice, HumanFileSize} from "../utilities/images";
import {decode} from "html-entities";


function PurchasesRows(props) {
    let rows = props.purchases.map((purchase)=> {
        return <PurchaseRow key={purchase.id} metadata={purchase}/>;
    });

    return (
        <div className="datasets-rows">
            {rows}
        </div>
    )
}


function PurchaseRow(props) {
    const history = useHistory();

    const [fiatPrice, setFiatPrice] = useState("");

    const timeAgo = new TimeAgo('en-US')

    const imageFilename = props.metadata.imageFilename;
    const title = props.metadata.title;
    const shortDescription = props.metadata.shortDescription;
    const price = "1.00"; //Number.parseFloat(props.metadata.price).toFixed(8).toString().replace(/\.?0+$/,"");

    const getFiatPrice = async ()=>{
        setFiatPrice(await FiatPrice("1"));
    }
    getFiatPrice();

    const fileType = props.metadata.fileType;
    const fileSize = HumanFileSize(props.metadata.fileSize, true);
    const username = props.metadata.username;
    const timestamp = timeAgo.format(Date.parse(props.metadata.Timestamp));

    const buttonText = (props.rowType === "edit") ? "Edit" : "Details";
    const gotoPage = (props.rowType === "edit") ? '/dashboard/datasets/'+props.metadata.id : '/dataset/'+props.metadata.id;

    const datasetImage = "/api/v1/image/" + imageFilename;

    const handleClickDatasetRow = (e) => {
        history.push(gotoPage);
    }

    return (
        <div className="datasets-row" onClick={handleClickDatasetRow}>
            <div className="datasets-row-image">
                <img className="datasets-image" src={datasetImage} alt={decode(title)}/>
            </div>
            <div className="datasets-row-info">
                <div className="mini-bold-title">{decode(title)}</div>
                <div className="mini-description">{decode(shortDescription)}</div>
                <div className="mini-light-description tag-container">
                    <div>{fileType}</div>
                    <div>{fileSize}</div>
                    <div>{timestamp}</div>
                    <div>{username}</div>
                </div>
            </div>
            <div className="datasets-details">
                <div><button className="normal-button">{buttonText}</button></div>
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
            console.log(res);
            setPurchases(res.data.purchases);

        }
        getPurchases();
    }, []);

    return (

        <div className="maincontent">
            <h2 className="margins-30">Purchases</h2>

            <div className="">
                { purchases.length ==0 &&
                <div>
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