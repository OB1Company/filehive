import React, { useState }  from 'react'
import { useHistory } from "react-router-dom";
import axios from "axios";
import { ConvertImageToString } from "./utilities/images";

function Create() {

  const [title, setTitle] = useState("");
  const [shortDescription, setShortDescription] = useState("");
  const [fullDescription, setFullDescription] = useState("");
  const [imageFile, setImageFile] = useState("");
  const [fileType, setFileType] = useState("");
  const [price, setPrice] = useState(0);
  const [dataset, setDataset] = useState("");
  const [isError, setIsError] = useState(false);
  const [error, setError] = useState(false);

  const history = useHistory();

  const HandleFormSubmit = async (e) => {
    e.preventDefault();

    // Convert image file to base64 string
    const fileString = await ConvertImageToString(imageFile);

    const data = {
      title: title,
      shortDescription: shortDescription,
      fullDescription: fullDescription,
      image: fileString,
      fileType: fileType,
      price: Number(price),
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
            history.push('/dataset/'+data.data.datasetID);
          });

    } catch(e) {

      setError(e);
      return false;
    }

  };

  const HandleDatasetImage = (e) => {
    setImageFile(e.target.files[0]);
  }

  const HandleDataset = (e) => {

    // Determine file type extension
    setFileType(e.target.files[0].type);

    setDataset(e.target.files[0]);
  }

  return (
      <div className="CreateDataset">
        <h2>Create dataset</h2>
        <div>
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
              <input type="text" name="shortDescription" placeholder="(100 char max)"
                     onChange={e => setShortDescription(e.target.value)}/>
              <span>Explain your dataset in 50 characters or less.</span>
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
            <div>
              <input type="text" name="price" placeholder="5.23" onChange={e => setPrice(e.target.value)}/>
              <span>Set your price in Filecoin (FIL).</span>
            </div>
          </label>

          <label>
            Dataset*
            <div>
              <input type="file" name="dataset" onChange={HandleDataset}/>
              <span>Finally, attach your dataset file(s) in zip format.</span>
            </div>
          </label>

          <div className="form-divider"></div>

          <div>
            <input type="submit" value="Submit" className="orange-button"/>
          </div>

        </form>
        </div>

      </div>
  )
}

export default Create