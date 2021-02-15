import React, {useEffect, useState} from 'react'
import {Table, Tbody, Td, Th, Thead, Tr} from 'react-super-responsive-table';
import 'react-super-responsive-table/dist/SuperResponsiveTableStyle.css';
import {getAxiosInstance} from "../Auth";
import TimeAgo from "javascript-time-ago";
import useSWR from "swr";
import {FiatPrice, FilecoinPrice, truncStringPortion} from "../utilities/images";
import { decode } from "html-entities";

function SalesRows(props) {
    return props.sales.map((sale) => {
        return <SalesRow key={sale.id} metadata={sale} price={props.price}/>;
    });
}

function SalesRow(props) {
    const timeAgo = new TimeAgo('en-US');
    const sale = props.metadata;
    const created = timeAgo.format(Date.parse(props.metadata.timestamp));

    return (
        <Tr>
            <Td>{truncStringPortion(sale.sellerID, 8, 8, 3)}</Td>
            <Td>{truncStringPortion(sale.userID, 8, 8, 3)}</Td>
            <Td><a href={"/dataset/"+sale.datasetID}>{decode(sale.title)}</a></Td>
            <Td>{created}</Td>
            <Td>{sale.price} FIL</Td>
            <Td>{FiatPrice(sale.price, props.price)}</Td> 
        </Tr>
    )
}

export default function AdminSales() {

    const [sales, setSales] = useState([]);
    const filecoinPrice  = useSWR('filecoinPrice', FilecoinPrice);

    const refreshSales = ()=>{
        const instance = getAxiosInstance();
        instance.get('/api/v1/admin/sales')
            .then((res)=>{
                console.log(res);
                const salesitems = res.data.sales;
                setSales(salesitems);
            })
    }

    useEffect(() => {
        refreshSales();
    }, []);

    return (
        <div className="margins-30 bottom-30">
            <h2>Sales</h2>

            <div>
                <Table className="sales-table font-12">
                    <Thead>
                        <Tr>
                            <Th>Seller</Th>
                            <Th>Buyer</Th>
                            <Th>Dataset</Th>
                            <Th>Date</Th>
                            <Th>FIL</Th>
                            <Th>USD</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        <SalesRows sales={sales} price={filecoinPrice.data}/>
                    </Tbody>
                </Table>
            </div>

        </div>
    );
}