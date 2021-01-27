import React, { useState } from 'react'
import DataSetsRows from "../DataSetsRows";
import {Link} from "react-router-dom";

export default function Purchases() {

    const [purchases, setPurchases] = useState([]);

    const trendingURL = "/datasets/trending";

    return (

        <div className="maincontent margins-30">
            <h2>Purchases</h2>

            { purchases.length == 0 &&
            <div>
                <p className="mini-description dashboard-p">You have not made any purchases yet. Check out some of our <a href={trendingURL} className="orange-link">trending datasets</a>!</p>
            </div>
            }

        </div>

    );
}