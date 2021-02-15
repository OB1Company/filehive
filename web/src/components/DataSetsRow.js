import React, {useState} from 'react'
import {useHistory} from "react-router-dom";
import {FiatPrice, FilecoinPrice, HumanFileSize} from "./utilities/images";
import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
import { decode } from 'html-entities';
import useSWR from "swr";

TimeAgo.addDefaultLocale(en);

function DataSetsRow(props) {
    const history = useHistory();

    const timeAgo = new TimeAgo('en-US')

    const imageFilename = props.metadata.imageFilename;
    const title = props.metadata.title;
    const shortDescription = props.metadata.shortDescription;
    const price = Number.parseFloat(props.metadata.price).toFixed(8).toString().replace(/\.?0+$/,"");

    const filecoinPrice  = useSWR('filecoinPrice', FilecoinPrice);

    const fiatPrice = FiatPrice(props.metadata.price, filecoinPrice.data);

    const fileType = props.metadata.fileType;
    const fileSize = HumanFileSize(props.metadata.fileSize, true);
    const username = props.metadata.username;
    const timestamp = timeAgo.format(Date.parse(props.metadata.createdAt));

    const buttonText = (props.rowType === "edit") ? "Edit" : "Details";
    const gotoPage = (props.rowType === "edit") ? '/dashboard/datasets/'+props.metadata.id : '/dataset/'+props.metadata.id;

    const datasetImage = "/api/v1/image/" + imageFilename;

    const handleClickDatasetRow = (e) => {
        history.push(gotoPage);
    }

    return (
        <div className="datasets-row">
            <div className="datasets-row-image">
                <a href={"/dataset/"+props.metadata.id}><img className="datasets-image" src={datasetImage} alt={decode(title)}/></a>
            </div>
            <div className="datasets-row-info">
                <div className="mini-bold-title"><a href={"/dataset/"+props.metadata.id}>{decode(title)}</a></div>
                <div className="mini-description">{decode(shortDescription)}</div>
                <div className="mini-light-description tag-container">
                    <div>{fileType}</div>
                    <div>{fileSize}</div>
                    <div>{timestamp}</div>
                    <div>{username}</div>
                </div>
            </div>
            <div className="datasets-details">
                <div><button className="normal-button raise" onClick={handleClickDatasetRow}>{buttonText}</button></div>
                <div className="small-orange-text dataset-row-price">{price} FIL</div>
                <div className="mini-light-description">{fiatPrice}</div>
            </div>
        </div>
    )
}

export default DataSetsRow