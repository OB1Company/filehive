import React, {useEffect, useState} from 'react'
import { useHistory, useParams } from "react-router-dom";
import axios from "axios";
import {ConvertImageToString, FilecoinPrice} from "./utilities/images";
import ErrorBox, {SuccessBox} from "./ErrorBox";
import {getAxiosInstance} from "./Auth";
import Helmet from "react-helmet";
import { decode } from 'html-entities';

export default function Settings() {

  // Page params
  let { id } = useParams();

  const history = useHistory();

  const [title, setTitle] = useState("");
  const [shortDescription, setShortDescription] = useState("");
  const [fullDescription, setFullDescription] = useState("");
  const [imageFile, setImageFile] = useState("");
  const [fileType, setFileType] = useState("Unknown");
  const [price, setPrice] = useState(0);
  const [error, setError] = useState("");
  const [cid, setCid] = useState("");
  const [success, setSuccess] = useState("");
  const [datasetPrice, setDatasetPrice] = useState("");
  const [dealStatus, setDealStatus] = useState("");
  const [dealPieceCid, setDealPieceCid] = useState("");
  const [dealProposalCid, setDealProposalCid] = useState("");
  const [dealStatusCode, setDealStatusCode] = useState(0);
  const [active, setActive] = useState(true);
  const [inactive, setInactive] = useState(false);

  useEffect(() => {
    const GetDataset = async (datasetId) => {
      const instance = getAxiosInstance();

      const datasetUrl = "/api/v1/dataset/" + datasetId;
      const response = await instance.get(datasetUrl, {withCredentials: true})
      const dr = response.data;

      setTitle(decode(dr.title));
      setShortDescription(decode(dr.shortDescription));
      setFullDescription(decode(dr.fullDescription));
      setPrice(dr.price);
      setFileType(dr.fileType);
      setCid(dr.contentID);

      if(dr.delisted) {
        setInactive(true);
        setActive(false);
      } else {
        setInactive(false);
        setActive(true);
      }

      instance.get("/api/v1/datasetdeal/"+dr.contentID, {withCredentials: true})
          .then((deal)=>{
            setDealStatusCode(deal.data.status);

            if(deal.data.status !== 3 && deal.data.status !== 4) {
              setDealStatus(deal.data.deal_info[0].state_name);
              setDealPieceCid(deal.data.deal_info[0].piece_cid);
              setDealProposalCid(deal.data.deal_info[0].proposal_cid);
            }
          })


    }

    const fetchData = async() => {
      await GetDataset(id);
    };
    fetchData();
  }, []);

  const handleFormSubmit = (e) => {
    e.preventDefault();

    if(title === "") {
      setError("Please specify a title for your dataset");
      return;
    }
    if(shortDescription === "") {
      setError("Please specify a short description for your dataset");
      return;
    }
    if(fullDescription === "") {
      setError("Please specify a full description for your dataset");
      return;
    }
    if(price < 0 || isNaN(price)) {
      setError("Please provide a valid price for your dataset");
      return;
    }

    const handleForm = async() => {

      console.log(this);
      let fileString = "";

      if(imageFile) {
        // Convert image file to base64 string
        fileString = await ConvertImageToString(imageFile);
      }

      const data = {
        id: id,
        title: title,
        shortDescription: shortDescription,
        fullDescription: fullDescription,
        price: Number(price),
      };

      if(fileString !== "") {
        data.image = fileString;
      }

      const csrftoken = localStorage.getItem('csrf_token');
      const instance = axios.create({
        baseURL: "",
        headers: {
          "x-csrf-token": csrftoken,
          "content-type": "multipart/form-data"
        }
      })

      const url = "/api/v1/dataset";
      try {
        await instance.patch(
            url,
            data
        )
            .then((data) => {
              setSuccess("Dataset has been updated")
              setError("");
            })
            .catch((e) => {
              setError(e.response.data.error);
              setSuccess("");
              return false;
            });
      } catch (e) {
        setError(e);
      }
    };
    handleForm();
  };

  const handleSetPrice = async (e)=>{
    setPrice(e.target.value);

    if(!isNaN(e.target.value)) {
      const filecoinPrice = await FilecoinPrice();
      let formatter = new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
        maximumFractionDigits: 4,
      });

      const priceUSD = formatter.format(filecoinPrice*e.target.value);
      setDatasetPrice(priceUSD);
    }
  }

  const handleDatasetImage = (e) => {
    setImageFile(e.target.files[0]);
  }

  const handleVisible = (e)=>{
    const instance = getAxiosInstance();
    instance.get("/api/v1/relist/"+id)
        .then((data)=>{
          setActive(true);
          setInactive(false);
        })
        .catch((err)=>{
          console.error(err);
        })
  }

  const handleInvisible = (e)=>{
    const instance = getAxiosInstance();
    instance.get("/api/v1/delist/"+id)
        .then((data)=>{
          setActive(false);
          setInactive(true);
        })
        .catch((err)=>{
          console.error(err);
        })
  }

  return (
      <div className="CreateDataset">

        <Helmet>
          <title>Filehive | Edit Dataset</title>
        </Helmet>

        <h2>Edit dataset</h2>
        <div>
          <form onSubmit={handleFormSubmit}>
          <label>
            Title*
            <div>
              <input type="text" placeholder="What's the name?" value={title} onChange={e => setTitle(e.target.value)}/>
              <span>Set a clear description title for your dataset.</span>
            </div>
          </label>

          <label>
            Visibility*
            <div className="visibility-toggle">
              <div>
                <input type="radio" name="visibility_on" checked={active} onChange={handleVisible}/>
                <label>Active</label>
                <input type="radio" name="visibility_off" checked={inactive} onChange={handleInvisible}/>
                <label>Inactive</label>
              </div>
              <span>Set your dataset to inactive to hide it from public view.</span>
            </div>
          </label>

          <label>
            Short description*
            <div>
              <input type="text" name="shortDescription" value={shortDescription} placeholder="(100 char max)"
                     onChange={e => setShortDescription(e.target.value)}/>
              <span>Explain your dataset in 50 characters or less.</span>
            </div>
          </label>

          <label>
            Full description*
            <div>
              <textarea name="fullDescription" value={fullDescription} placeholder="Enter description"
                        onChange={e => setFullDescription(e.target.value)}/>
              <span>Fully describe the contents in the dataset and provide example of how the data is structured. The more information the better.</span>
            </div>
          </label>

          <label>
            Image*
            <div>
              <input type="file" name="imageFile" onChange={handleDatasetImage}/>
              <span>Attach a JPG or PNG cover photo for your dataset.</span>
            </div>
          </label>

          <div className="form-divider"></div>

          <h3>Dataset Information</h3>

          <label>
            File Type: {fileType}
          </label>



          <label>
            <h4>IPFS Information</h4>
            <ul className="form-ul">
              <li><a href={"ipfs://"+cid} target="_blank" className="orange-link">ipfs://{cid}</a></li>
              <li><a href={"https://gateway.ipfs.io/ipfs/"+cid} target="_blank" className="orange-link">https://gateway.ipfs.io/ipfs/{cid}</a></li>
            </ul>
          </label>

            {dealStatusCode !== 3 && dealStatusCode != 4 &&
          <label>
            <h4>Filecoin Deal Information</h4>

            <div className="filecoindeal-container">
              <div>Publish Status: <strong>{dealStatus}</strong></div>
              <div>Piece CID: {dealPieceCid}</div>
              <div>Proposal CID: {dealProposalCid}</div>
            </div>

          </label>
            }

          <label>
            Price*
            <div>
              <input type="text" name="price" value={price} placeholder="5.23" onChange={handleSetPrice}/>
              <span>Set your price in Filecoin (FIL).<br/>Estimated price: <strong>{datasetPrice}</strong></span>
            </div>
          </label>

          <div className="form-divider"></div>

            {error &&
            <ErrorBox message={error}/>
            }
            {success &&
            <SuccessBox message={success}/>
            }

          <div>
            <input type="submit" value="Submit" className="raise orange-button"/>
          </div>

        </form>
        </div>

      </div>
  )
}