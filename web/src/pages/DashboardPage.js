import React from 'react'
import { useLocation } from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import TabbedLinks from "../components/TabbedLinks";
import DataSetsRows from "../components/DataSetsRows";

export default function DashboardPage() {

    const linkNames = [
        { name: 'Datasets', link: '/dashboard/datasets' },
        { name: 'Purchases', link: '/dashboard/purchases' },
        { name: 'Wallet', link: '/dashboard/wallet' },
        { name: 'Settings', link: '/dashboard/settings' }
    ];

    const location = useLocation();

    return (
        <div className="container">
            <Header/>
            <TabbedLinks linkNames={linkNames} activeLink={location.pathname} />
            <DataSetsRows sortby="trending"/>
            <Footer/>
        </div>
    )
}