import React from 'react'
import {Link} from "react-router-dom";
import DataSetsRow from "../components/DataSetsRow";

function DataSetsRows(props) {
    let rows = props.datasets.map((dataset)=> {
        return <DataSetsRow metadata={dataset} rowType={props.rowType}/>;
    });

    return (
        <div class="datasets-rows">
            {rows}
        </div>
    )
}

export default DataSetsRows