import React from 'react'
import { Link } from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import TabbedLinks from "../components/TabbedLinks";
import DataSetsRows from "../components/DataSetsRows";

export default function DashboardPage() {

    const linkNames = [
        { name: 'Datasets', link: '/dashboard/datasets' },
        { name: 'Purchases', link: '/datasets/purchases' },
        { name: 'Wallet', link: '/dashboard/wallet' },
        { name: 'Settings', link: '/datasets/settings' }
    ];

    return (
        <div className="container">
            <Header/>
            <TabbedLinks linkNames={linkNames} />
            <DataSetsRows sortby="trending"/>
            <Footer/>
        </div>
    )
}