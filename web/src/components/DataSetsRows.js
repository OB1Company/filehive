import React from 'react'
import DataSetsRow from "../components/DataSetsRow";

function DataSetsRows(props) {
    let rows = props.datasets.map((dataset)=> {
        return <DataSetsRow key={dataset.id} metadata={dataset} rowType={props.rowType}/>;
    });

    return (
        <div className="datasets-rows">
            {rows}
        </div>
    )
}

export default DataSetsRows