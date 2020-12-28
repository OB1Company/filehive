import React from 'react'
import {Link} from "react-router-dom";

function TabbedLinks(props) {
    let links = props.linkNames.map((link)=> {
        return <li class="active"><Link to={link.link}>{link.name}</Link></li>;
    });

    return (
        <ul class="tabbed-links">
            {links}
        </ul>
    )
}

export default TabbedLinks