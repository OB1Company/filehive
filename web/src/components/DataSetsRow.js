import React, {useState} from 'react'
import {Link} from "react-router-dom";
import {getAxiosInstance} from "./Auth";

function DataSetsRow(props) {
    const imageFilename = props.metadata.imageFilename;
    const title = props.metadata.title;
    const shortDescription = props.metadata.shortDescription;

    const datasetImage = "/api/v1/image/" + imageFilename;

    return (
        <div class="datasets-row">
            <div><img className="datasets-image" src={datasetImage}/></div>
            <div>
                {title}<br/>
                {shortDescription}
            </div>
            <div class="datasets-details">
                <div><button>Details</button></div>
                <div>1.23 FIL</div>
            </div>
        </div>
    )
}

export default DataSetsRow