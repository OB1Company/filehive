import React, {useEffect, useState} from 'react'
import { useHistory } from "react-router-dom";
import axios from "axios";
import {ConvertImageToString, FilecoinPrice} from "./utilities/images";
import ErrorBox, {SuccessBox} from "./ErrorBox";
import spinner from "../images/spinner.gif";
import useSWR from 'swr'
import {getAxiosInstance} from "./Auth";

function Create() {

  const history = useHistory();

  const [title, setTitle] = useState("");
  const [shortDescription, setShortDescription] = useState("");
  const [fullDescription, setFullDescription] = useState("");
  const [imageFile, setImageFile] = useState("");
  const [datasetFilename, setDatasetFilename] = useState("");
  const [fileType, setFileType] = useState("Unknown");
  const [price, setPrice] = useState(0);
  const [dataset, setDataset] = useState("");
  const [error, setError] = useState("");
  const [success] = useState("");
  const [datasetPrice, setDatasetPrice] = useState("");
  const [isCreating, setIsCreating] = useState(false);
  const [activated, setActivated] = useState(false);

  const filecoinPrice  = useSWR('filecoinPrice', FilecoinPrice);


  useEffect(()=>{
    const instance = getAxiosInstance();
    instance.get("/api/v1/user")
        .then((data) => {
          setActivated(data.data.Activated);
        })
  }, []);

  const HandleFormSubmit = (e) => {

    setIsCreating(true);

    e.preventDefault();

    if(imageFile === "") {
      setError("An image to depict your dataset is required");
      setIsCreating(false);
      return;
    }
    if(title === "") {
      setError("Please specify a title for your dataset");
      setIsCreating(false);
      return;
    }
    if(shortDescription === "") {
      setError("Please specify a short description for your dataset");
      setIsCreating(false);
      return;
    }
    if(fullDescription === "") {
      setError("Please specify a full description for your dataset");
      setIsCreating(false);
      return;
    }
    if(price <= 0) {
      setError("Please provide a price for your dataset");
      setIsCreating(false);
      return;
    }
    if(dataset === "") {
      setError("Please choose a dataset to upload");
      setIsCreating(false);
      return;
    }

    const handleForm = async() => {

      // Convert image file to base64 string
      const fileString = await ConvertImageToString(imageFile);

      const data = {
        title: title,
        shortDescription: shortDescription,
        fullDescription: fullDescription,
        image: fileString,
        fileType: fileType,
        price: Number(price),
        filename: datasetFilename
      };

      const formData = new FormData();
      formData.append('metadata', JSON.stringify(data));
      formData.append('file', dataset);

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
        await instance.post(
            url,
            formData
        )
            .then((data) => {
              setIsCreating(false);
              history.push('/dataset/' + data.data.datasetID);
            })
            .catch((e) => {
              console.log(e.response.data);
              setError(e.response.data);
              setIsCreating(false);
            });
      } catch (e) {
        setError(e.response.data);
        setIsCreating(false);
      }
    };
    handleForm();
  };

  const HandleSetPrice = async (e)=>{
    setPrice(e.target.value);

    var formatter = new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      maximumFractionDigits: 4,
    });

    const priceUSD = formatter.format(filecoinPrice.data*e.target.value);
    setDatasetPrice(priceUSD);
  }

  const HandleDatasetImage = (e) => {
    setImageFile(e.target.files[0]);
  }

  const HandleDataset = (e) => {
    setFileType(e.target.files[0].type);
    setDatasetFilename(e.target.files[0].name);
    setDataset(e.target.files[0]);
  }

  const CreateButton = () => {
    if (!isCreating) {
      return <input type="submit" value="Submit" className="orange-button"/>;
    } else {
      return <div>
        <img src={spinner} width="20" height="20" alt="spinner"/> Creating dataset...
      </div>
    }
  }

  return (
      <div className="CreateDataset">
        <h2>Create dataset</h2>
        <div>
          {activated &&
          <form onSubmit={HandleFormSubmit}>
          <label>
            Title*
            <div>
              <input type="text" name="title" placeholder="Title" onChange={e => setTitle(e.target.value)}/>
              <span>Set a clear description title for your dataset.</span>
            </div>
          </label>

          <label>
            Short description*
            <div>
              <input type="text" name="shortDescription" maxLength="100" placeholder="(100 char max)"
                     onChange={e => setShortDescription(e.target.value)}/>
              <span>Explain your dataset in 100 characters or less.</span>
            </div>
          </label>

          <label>
            Full description*
            <div>
              <textarea name="fullDescription" placeholder="Enter description"
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
            <p className="mini-light-description">Filehive automatically deducts a 5% fee from every sale.</p>
            <div>
              <input type="text" name="price" placeholder="5.23" onChange={HandleSetPrice}/>
              <span>Set your price in Filecoin (FIL).<br/>Estimated price: <strong>{datasetPrice}</strong></span>
            </div>
          </label>

          <label>
            Dataset*
            <div>
              <input type="file" name="dataset" onChange={HandleDataset}/>
              <span>Finally, attach your dataset.</span>
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
            <CreateButton/>
          </div>

        </form>
          }
          {!activated &&
            <p className="mini-description dashboard-p">Check your email 📪 to confirm your account and start creating datasets.</p>
          }
        </div>

      </div>
  )
}

export default Create
