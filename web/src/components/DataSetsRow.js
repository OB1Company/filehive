import React from 'react'
import {useHistory} from "react-router-dom";

function DataSetsRow(props) {
    const history = useHistory();

    const imageFilename = props.metadata.imageFilename;
    const title = props.metadata.title;
    const shortDescription = props.metadata.shortDescription;
    const price = Number.parseFloat(props.metadata.price).toFixed(8).toString().replace(/\.?0+$/,"");
    const fileType = props.metadata.fileType;
    const fileSize = "3"; //props.metadata.price;
    const username = props.metadata.username;
    const timestamp = "BOBBY"; //props.metadata.timestamp;

    const buttonText = (props.rowType === "edit") ? "Edit" : "Details";
    const gotoPage = (props.rowType === "edit") ? '/dataset/'+props.metadata.id+'/edit' : '/dataset/'+props.metadata.id;

    const datasetImage = "/api/v1/image/" + imageFilename;

    const handleClickDatasetRow = (e) => {
        history.push(gotoPage);
    }

    return (
        <div class="datasets-row" onClick={handleClickDatasetRow}>
            <div className="datasets-row-image">
                <img className="datasets-image" src={datasetImage}/>
            </div>
            <div className="datasets-row-info">
                <div className="mini-bold-title">{title}</div>
                <div className="mini-description">{shortDescription}</div>
                <div className="mini-light-description">{fileType} {fileSize} {timestamp} {username}</div>
            </div>
            <div className="datasets-details">
                <div><button className="normal-button">{buttonText}</button></div>
                <div className="small-orange-text">{price} FIL</div>
            </div>
        </div>
    )
}

export default DataSetsRow