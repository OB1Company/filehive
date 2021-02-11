import React from 'react'
import { useLocation } from 'react-router-dom'
import Header from '../Header'
import Footer from '../Footer'
import TabbedLinks from "../components/TabbedLinks";
import AdminUsers from "../components/admin/AdminUsers";
import AdminDatasets from "../components/admin/AdminDatasets";
import AdminSales from "../components/admin/AdminSales";
import {Helmet} from "react-helmet";

export default function AdminPage() {

    const linkNames = [
        { name: 'Users', link: '/admin/users' },
        { name: 'Datasets', link: '/admin/datasets' },
        { name: 'Sales', link: '/admin/sales' },
    ];

    const location = useLocation();

    const DashboardPage = () => {
        const tab = location.pathname.substring(location.pathname.lastIndexOf('/')+1);

        switch(tab) {
            case "datasets":
                return <AdminDatasets/>;
            case "users":
                return <AdminUsers/>;
            case "sales":
                return <AdminSales/>;
        }

        return <h2>{tab}</h2>

    }

    return (
        <div className="container">
            <Helmet>
                <title>Filehive | Admin</title>
            </Helmet>
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