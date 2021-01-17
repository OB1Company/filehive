import React from 'react'
import {useHistory} from "react-router-dom";

function DataSetsRow(props) {
    const history = useHistory();

    const imageFilename = props.metadata.imageFilename;
    const title = props.metadata.title;
    const shortDescription = props.metadata.shortDescription;
    const price = props.metadata.price;
    const fileType = props.metadata.fileType;
    const fileSize = "3"; //props.metadata.price;
    const username = props.metadata.username;
    const timestamp = "BOBBY"; //props.metadata.timestamp;

    const datasetImage = "/api/v1/image/" + imageFilename;

    const handleClickDatasetRow = (e) => {
        history.push('/dataset/'+props.metadata.id);
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
                <div><button>Details</button></div>
                <div>{price} FIL</div>
            </div>
        </div>
    )
}

export default DataSetsRow