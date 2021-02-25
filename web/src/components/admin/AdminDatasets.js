import React, {useEffect, useState} from 'react'
import TimeAgo from "javascript-time-ago";
import {Table, Tbody, Td, Th, Thead, Tr} from "react-super-responsive-table";
import {FiatPrice, FilecoinPrice, truncStringPortion} from "../utilities/images";
import {decode} from "html-entities";
import useSWR from "swr";
import {getAxiosInstance} from "../Auth";

function AdminDatasetRows(props) {
    return props.sets.map((set) => {
        return <AdminDatasetRow key={set.id} metadata={set} price={props.price} selectHandler={props.selectHandler}/>;
    });
}

function AdminDatasetRow(props) {
    const set = props.metadata;

    return (
        <Tr>
            <Td><input type="checkbox" name="set" value={set.id} onChange={props.selectHandler}/></Td>
            <Td>{truncStringPortion(set.userID, 8, 8, 3)}</Td>
            <Td>{set.title}</Td>
            <Td>{set.delisted ? "🚫" : "✅"}</Td>
            <Td>{set.totalViews}</Td>
            <Td>0</Td>
            <Td>{set.price} FIL</Td>
            <Td>{FiatPrice(set.price, props.price)}</Td>
        </Tr>
    )
}

export default function AdminDatasets() {

    const [sets, setSets] = useState([]);
    const filecoinPrice  = useSWR('filecoinPrice', FilecoinPrice);
    const [selectedCount, setSelectedCount] = useState(0);
    const [selectedDatasets, setSelectedDatasets] = useState([])

    const refreshSets = ()=>{
        const instance = getAxiosInstance();
        instance.get('/api/v1/admin/datasets')
            .then((res)=>{
                console.log(res);
                const adminSets = res.data.datasets;
                setSets(adminSets);
            })
    }

    useEffect(() => {
        refreshSets();
    }, []);

    const HandleDelist = (e)=>{
        e.preventDefault();
        const instance = getAxiosInstance();
        instance.post("/api/v1/delist", {"datasets":selectedDatasets})
            .then((result)=>{
                refreshSets();
            });

    }

    const HandleRelist = (e)=>{
        e.preventDefault();
        const instance = getAxiosInstance();
        instance.post("/api/v1/relist", {"datasets":selectedDatasets})
            .then((result)=>{
                refreshSets();
            });

    }

    const HandleSelection = (e)=>{
        const checked = e.target.checked;
        if(checked) {
            setSelectedCount(selectedCount+1);
            let sets = selectedDatasets;
            sets.push(e.target.value);
            setSelectedDatasets(sets);
        } else {
            setSelectedCount(selectedCount-1);
            let sets = selectedDatasets;
            const index = sets.indexOf(e.target.value);
            if (index > -1) {
                sets.splice(index, 1);
            }
            setSelectedDatasets(sets);
        }

    }

    return (
        <div className="margins-30 bottom-30">
            <h2>Datasets</h2>
            <br/>
            <div className="admin-toolbar">
                <div className="bold">{selectedCount} Selected</div>
                <div><a href="" className="orange-link2" onClick={HandleDelist}>Delist</a></div>
                <div><a href="" className="orange-link2" onClick={HandleRelist}>Relist</a></div>
            </div>

            <div>
                <Table className="sales-table font-12">
                    <Thead>
                        <Tr>
                            <Th></Th>
                            <Th>Owner</Th>
                            <Th>Title</Th>
                            <Th>Status</Th>
                            <Th>Views</Th>
                            <Th>Purchases</Th>
                            <Th>Price</Th>
                            <Th>Gross Sales</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        <AdminDatasetRows sets={sets} price={filecoinPrice.data} selectHandler={HandleSelection}/>
                    </Tbody>
                </Table>
            </div>

        </div>
    );
}