import React from 'react'
import {Link} from "react-router-dom";

function DataSetsRow(props) {
    const imageFilename = props.metadata.imageFilename;
    const title = props.metadata.title;
    const shortDescription = props.metadata.shortDescription;


    return (
        <div class="datasets-row">
            {imageFilename} {title}<br/>
            {shortDescription}
        </div>
    )
}

export default DataSetsRow