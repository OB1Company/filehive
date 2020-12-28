import React from 'react'
import {Link} from "react-router-dom";
import DataSetsRow from "../components/DataSetsRow";

function DataSetsRows(props) {
    // let links = props.linkNames.map((link)=> {
    //     return <li class="active"><Link to={link.link}>{link.name}</Link></li>;
    // });

    return (
        <div class="datasets-rows">
            <DataSetsRow sortby={props.sortby}/>
        </div>
    )
}

export default DataSetsRows