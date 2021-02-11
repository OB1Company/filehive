import React, {useEffect, useState} from 'react'
import { Table, Thead, Tbody, Tr, Th, Td } from 'react-super-responsive-table';
import 'react-super-responsive-table/dist/SuperResponsiveTableStyle.css';
import {getAxiosInstance} from "../Auth";
import {useHistory} from "react-router-dom";
import TimeAgo from "javascript-time-ago";
import useSWR from "swr";
import {FiatPrice, FilecoinPrice} from "../utilities/images";

function UsersRows(props) {
    let rows = props.users.map((user)=> {
        return <UserRow key={user.id} metadata={user}/>;
    });

    return rows;
}

function UserRow(props) {
    const timeAgo = new TimeAgo('en-US');
    const user = props.metadata;
    const created = timeAgo.format(Date.parse(props.metadata.CreatedAt));

    return (
        <Tr>
            <Td>{user.Name}</Td>
            <Td>{user.Email}</Td>
            <Td>{created} </Td>
            <Td>{user.Admin ? "✅" : ""}</Td>
            <Td>{user.PowergateToken}</Td>
            <Td>{user.PowergateID}</Td>
        </Tr>
    )
}

export default function AdminUsers() {

    const [users, setUsers] = useState([]);

    useEffect(() => {
        const instance = getAxiosInstance();
        instance.get('/api/v1/users')
            .then((res)=>{
                const users = res.data.users;
                setUsers(users);
            })
    }, []);

    return (
        <div className="margins-30">
            <h2>Users 👻</h2>
            <br/>
            <div>
                <Table className="sales-table">
                    <Thead>
                        <Tr>
                            <Th>Name</Th>
                            <Th>Email</Th>
                            <Th>Created</Th>
                            <Th>Admin</Th>
                            <Th>Token</Th>
                            <Th>ID</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        <UsersRows users={users}/>
                    </Tbody>
                </Table>
            </div>

        </div>
    );
}