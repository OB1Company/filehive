import React from 'react'
import {Link} from "react-router-dom";

function TabbedLinks(props) {
    let links = props.linkNames.map((link)=> {
        const active = (link.link === props.activeLink);
        return <li key={link.link} className={active ? 'active' : ''}><Link to={link.link}>{link.name}</Link></li>;
    });

    return (
        <ul className="tabbed-links">
            {links}
        </ul>
    )
}

export default TabbedLinks