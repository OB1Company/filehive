import React, {useEffect, useState} from 'react'
import { useHistory, useParams } from "react-router-dom";
import axios from "axios";
import {ConvertImageToString, FilecoinPrice} from "./utilities/images";
import ErrorBox, {SuccessBox} from "./ErrorBox";
import {getAxiosInstance} from "./Auth";

function Edit() {

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
  const [success] = useState("");
  const [datasetPrice, setDatasetPrice] = useState("");

  const HandleFormSubmit = (e) => {
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
    if(price <= 0) {
      setError("Please provide a price for your dataset");
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
              history.push('/dashboard/datasets/' + id);
            })
            .catch((e) => {
              setError(e.response.data.error);
              return false;
            });
      } catch (e) {
        setError(e);
      }
    };
    handleForm();
  };

  const HandleSetPrice = async (e)=>{
    setPrice(e.target.value);
    const filecoinPrice = await FilecoinPrice();
    var formatter = new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      maximumFractionDigits: 4,
    });

    const priceUSD = formatter.format(filecoinPrice*e.target.value);
    setDatasetPrice(priceUSD);
  }

  const HandleDatasetImage = (e) => {
    setImageFile(e.target.files[0]);
  }

  const GetDataset = async (datasetId) => {
    const instance = getAxiosInstance();

    const datasetUrl = "/api/v1/dataset/" + datasetId;
    const response = await instance.get(datasetUrl, {withCredentials: true})

    console.log(response.data);

    const dr = response.data;
    setTitle(dr.title);
    setShortDescription(dr.shortDescription);
    setFullDescription(dr.fullDescription);
    setPrice(dr.price);
    setFileType(dr.fileType);

  }
  useEffect(() => {
    const fetchData = async() => {
      await GetDataset(id);
    };
    fetchData();
  });

  return (
      <div className="CreateDataset">
        <h2>Edit dataset</h2>
        <div>
          <form onSubmit={HandleFormSubmit}>
          <label>
            Title*
            <div>
              <input type="text" name="title" placeholder="Title" value={title} onChange={e => setTitle(e.target.value)}/>
              <span>Set a clear description title for your dataset.</span>
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
              <input type="file" name="imageFile" onChange={HandleDatasetImage}/>
              <span>Attach a JPG or PNG cover photo for your dataset.</span>
            </div>
          </label>

          <div className="form-divider"></div>

          <h3>Files</h3>

          <label>
            File Type: {fileType}
          </label>

          <label>
            Price*
            <div>
              <input type="text" name="price" value={price} placeholder="5.23" onChange={HandleSetPrice}/>
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
            <input type="submit" value="Submit" className="orange-button"/>
          </div>

        </form>
        </div>

      </div>
  )
}

export default Edit