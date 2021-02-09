import React, {useEffect, useState} from 'react'
import {getAxiosInstance} from "../Auth";
import {useHistory} from "react-router-dom";
import TimeAgo from "javascript-time-ago";
import {FiatPrice, HumanFileSize} from "../utilities/images";
import { Table, Thead, Tbody, Tr, Th, Td } from 'react-super-responsive-table';
import 'react-super-responsive-table/dist/SuperResponsiveTableStyle.css';


function SalesRows(props) {
    let rows = props.purchases.map((purchase)=> {
        return <SalesRow key={purchase.id} metadata={purchase}/>;
    });

    return rows;
}


function SalesRow(props) {
    const history = useHistory();

    console.log(props);

    const [fiatPrice, setFiatPrice] = useState("");

    const timeAgo = new TimeAgo('en-US');

    const sale = props.metadata;

    const imageFilename = props.metadata.imageFilename;
    const title = props.metadata.title;
    const shortDescription = props.metadata.shortDescription;
    const price = "1.00"; //Number.parseFloat(props.metadata.price).toFixed(8).toString().replace(/\.?0+$/,"");

    const getFiatPrice = async ()=>{
        setFiatPrice(await FiatPrice("1"));
    }
    getFiatPrice();

    const fileType = props.metadata.fileType;
    const fileSize = HumanFileSize(props.metadata.fileSize, true);
    const username = props.metadata.username;
    const timestamp = timeAgo.format(Date.parse(props.metadata.Timestamp));

    const buttonText = "Download";
    const gotoPage = '/dataset/'+props.metadata.datasetID;

    const datasetImage = "/api/v1/image/" + imageFilename;

    const handleClickDatasetRow = (e) => {
        history.push(gotoPage);
    }

    const datasetUrl = "/api/v1/download/"+props.metadata.datasetID;

    return (

        <Tr>
            <Td>{sale.id}</Td>
            <Td>{sale.username}</Td>
            <Td>{sale.price} FIL</Td>
            <Td>{sale.Timestamp}</Td>
        </Tr>

    )
}


export default function Sales() {

    const [sales, setSales] = useState([]);

    useEffect(() => {
        const getSales = async () => {
            const instance = getAxiosInstance();
            const apiUrl = "/api/v1/sales";
            const res = await instance.get(
                apiUrl
            );
            console.log(res);
            setSales(res.data.sales);

        }
        getSales();
    }, []);

    return (

        <div className="maincontent">
            <h2 className="margins-30">Sales</h2>

            {/*<div>*/}
            {/*    <div>Total transactions: 0</div>*/}
            {/*    <div>Total sales: $0</div>*/}
            {/*    <div>Data sold: 0MB</div>*/}
            {/*</div>*/}

            <div className="sales-table-container">
            <Table className="sales-table">
                <Thead>
                    <Tr>
                        <Th>Sale ID</Th>
                        <Th>Buyer</Th>
                        <Th>Price</Th>
                        <Th>Purchase Date</Th>
                    </Tr>
                </Thead>
                <Tbody>
                { sales.length === 0 &&
                <div className="margins-30">
                    <p className="mini-description dashboard-p">You have not made any sales yet. Good luck! 🤞 </p>
                </div>
                }
                { sales.length > 0 &&
                <SalesRows purchases={sales}/>
                }
                </Tbody>
            </Table>
            </div>

        </div>

    );
}