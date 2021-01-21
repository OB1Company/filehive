import React from 'react'
import { useLocation } from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import TabbedLinks from "../components/TabbedLinks";
import Datasets from "../components/dashboard/Datasets";
import Purchases from "../components/dashboard/Purchases";
import Wallet from "../components/dashboard/Wallet";
import Settings from "../components/dashboard/Settings";

export default function DashboardPage() {

    const linkNames = [
        { name: 'Datasets', link: '/dashboard/datasets' },
        { name: 'Purchases', link: '/dashboard/purchases' },
        { name: 'Wallet', link: '/dashboard/wallet' },
        { name: 'Settings', link: '/dashboard/settings' }
    ];

    const location = useLocation();

    const DashboardPage = () => {
        const tab = location.pathname.substring(location.pathname.lastIndexOf('/')+1);

        switch(tab) {
            case "datasets":
                return <Datasets/>;
            case "purchases":
                return <Purchases/>;
            case "wallet":
                return <Wallet/>;
            case "settings":
                return <Settings/>;
        }

        return <h2>{tab}</h2>

    }

    return (
        <div className="container">
            <Header/>
            <div className="maincontent">
                <TabbedLinks linkNames={linkNames} activeLink={location.pathname} />
                <div className="dashboard-container">
                    <DashboardPage/>
                </div>
            </div>
            <Footer/>
        </div>
    )
}